// Package join lets a further API server join a running installation: this
// API hands out a one-time code (cd.users/join-code), takes the new server's
// request under that code's secret and keeps it for an admin to compare and
// sign (cd.users/pending-api). Nothing here puts a node into the trust list;
// only a signed version of the list does that. Once the admin has signed, the
// new server fetches the keys every API shares over PathKeys (keys.go), and
// Run (run.go) is the other side of the same way, the one the new server
// itself walks; no command line reaches it yet -- main.go is still a stub.
package join

import (
	"context"
	"crypto/ecdh"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sorglos-it/cirrodesk/apps/server-api/classes/audit"
	"github.com/sorglos-it/cirrodesk/apps/server-api/classes/dbwait"
	"github.com/sorglos-it/cirrodesk/apps/server-api/classes/httpapi"
	"github.com/sorglos-it/cirrodesk/apps/server-api/classes/reason"
	"github.com/sorglos-it/cirrodesk/apps/server-api/classes/secrets"
	"github.com/sorglos-it/cirrodesk/apps/server-api/classes/state"
	"github.com/sorglos-it/cirrodesk/shared/core"
	"github.com/sorglos-it/cirrodesk/shared/dbapi"
	"github.com/sorglos-it/cirrodesk/shared/enc"
	"github.com/sorglos-it/cirrodesk/shared/trust"
	"github.com/sorglos-it/cirrodesk/shared/tunnel"
	"github.com/sorglos-it/cirrodesk/shared/wire"
)

// Join is the issuing side: the API an admin asks for a join code and the one
// a new server sends its request to. It keeps no state of its own; everything
// lies in the database, so any API of the installation can take over.
type Join struct {
	S *state.State

	// Log takes the one line a join writes (Accept). Package server sets it;
	// a nil one means no line, never a refused join.
	Log *audit.Log
}

// New returns the join side for s. It satisfies coreapp.JoinCodes.
func New(s *state.State) *Join { return &Join{S: s} }

// NewJoinCode makes the one-time code for the ulh script of a new API server:
// it names this API, its TLS pin and the newest version of the trust list, so
// the shell can check the code against what it verified itself before showing
// it. It mirrors controllers.NewPairCode, with dbapi.OnceJoin instead of
// dbapi.OncePair.
//
// The secret travels in the code and nowhere else. What the database keeps is
// its SHA-256 as the one-time value, and the secret itself as that value's
// data, so the API can check a MAC against it exactly once.
func (j *Join) NewJoinCode(ctx context.Context) (tunnel.JoinCode, int64, error) {
	v, err := j.S.Trust(ctx)
	if err != nil {
		return tunnel.JoinCode{}, 0, err
	}
	secret, err := secrets.Bytes(tunnel.SecretSize)
	if err != nil {
		return tunnel.JoinCode{}, 0, err
	}
	pub, err := j.noisePublic()
	if err != nil {
		return tunnel.JoinCode{}, 0, err
	}
	expires := j.S.Now() + tunnel.PairLifetime
	h := sha256.Sum256(secret)
	if err := j.S.DB.OnceCreate(ctx, dbapi.OnceJoin, h[:], enc.B64.EncodeToString(secret), expires); err != nil {
		return tunnel.JoinCode{}, 0, err
	}
	return tunnel.JoinCode{
		URL: j.S.Cfg.URL, Key: pub, TLS: enc.Bytes(j.S.Pin), Secret: enc.Bytes(secret),
		Anchor: trust.Anchor{Seq: v.List.Seq, Hash: enc.Bytes(v.Hash)},
	}, expires, nil
}

// Accept is the handler of tunnel.PathJoin?secret=<base64url of the SHA-256
// of the join secret>. The query names which one-time value to spend; the MAC
// proves that the sender holds the secret itself, which never travels back.
// That the method is POST is the router's business, as it is for
// controllers.Pair: nothing here looks at r.Method.
//
// The order is peek, check, spend (A13, Frage 23): the one-time value is read
// with OncePeek, the MAC is checked against it, and only a request that holds
// the secret spends it with OnceUse. Before, the value was spent FIRST, so
// whoever knew the hash in the query could burn every waiting join code with a
// wrong MAC. Guessing at the MAC instead gains nothing: it is keyed with a
// secret of tunnel.SecretSize random bytes. OnceUse still decides alone which
// of two requests with the right secret wins.
//
// Nothing here puts the new server anywhere near the trust list: the request
// is kept for the admin to compare and sign.
//
// Every attempt that gets as far as reading the one-time value is written to
// the log, as in controllers.Pair, and for the same measurement: twenty
// refused attempts left the log unchanged on 24.10.2026
// (docs/messungen-m1-m3.md, M3). audit.Log.Pairing says why the line hangs on
// the code being spent and not on the request arriving.
func (j *Join) Accept(w http.ResponseWriter, r *http.Request) {
	body, ok := httpapi.Body(w, r)
	if !ok {
		return
	}
	// The deadline of a path somebody is waiting on, and it is set HERE
	// rather than on the route because this path is one of the four that do
	// not go over httpapi.Routes.Mux (server.routes) -- so the deadline point
	// 49 gave the five session routes never reached it. Without one a store
	// call below waits busy_timeout times maxTries, 80,6 s measured (point
	// 48), and then answers the word for a database that is AWAY.
	//
	// After Body and not before it, which is the one thing this does better
	// than the route does: Body may take httpapi.readLimit (15 s) for a body
	// that trickles in, and a deadline started before it would spend itself on
	// the reading and then call a slow client's database busy. Point 49 named
	// that as the price of setting it on the route; here it costs nothing to
	// avoid.
	ctx, cancel := dbwait.User(r.Context())
	defer cancel()
	req, err := tunnel.ParseJoinRequest(body)
	if err != nil {
		httpapi.Fail(w, wire.CodeProtocol)
		return
	}
	hash, err := enc.B64.DecodeString(r.URL.Query().Get("secret"))
	if err != nil || len(hash) != sha256.Size {
		httpapi.Fail(w, wire.CodeProtocol)
		return
	}
	pub, err := j.noisePublic()
	if err != nil {
		httpapi.Fail(w, wire.CodeInternal)
		return
	}
	data, err := j.S.DB.OncePeek(ctx, dbapi.OnceJoin, hash, j.S.Now())
	if err != nil {
		// Unknown, expired and already used are one answer: telling them apart
		// would say whether a hash was ever a secret of this installation.
		//
		// A database that did not come free is NOT one of those three, and
		// answering E_VERIFY for it told a server that is joining correctly
		// that its code was wrong -- measured, 403 while a second program
		// held the file, and an admin would have made a new code for nothing.
		// Busy here opens no oracle: the deadline runs out before there is an
		// answer about any row, so it runs out alike for a hash that exists
		// and one that never did. TestAcceptSaysNothingAboutWhetherACodeExists
		// holds that with the store calls, the headers and the body.
		if dbwait.Busy(err) {
			httpapi.FailCoded(w, wire.CodeUnreachable, err)
			return
		}
		httpapi.Fail(w, wire.CodeVerify)
		return
	}
	secret, err := enc.B64.DecodeString(data)
	if err != nil {
		// Our own row does not read. That is damage,
		// not a wrong proof, and the log says so; the answer outside is the
		// same either way.
		j.note(ctx, wire.CodeVerify, reason.Defect)
		httpapi.Fail(w, wire.CodeVerify)
		return
	}
	if req.Verify(secret, pub) != nil {
		// THE line Frage 23 is about: somebody who knew the hash of a real
		// code and not its secret. The code stays valid for the server that
		// holds it.
		j.note(ctx, wire.CodeVerify, reason.JoinProof)
		httpapi.Fail(w, wire.CodeVerify)
		return
	}
	switch _, err := j.S.DB.BlobGet(ctx, dbapi.BlobJoin, req.Node); {
	case err == nil:
		// Frage 24 (a): a request already waits under this name, and
		// overwriting it would let the admin sign a server other than the one
		// he is waiting for without noticing. The first one stays; the new
		// server is told, so its admin can sort the name out.
		// shortcut: two APIs taking two requests for one name at the same
		// instant can still both pass this read; fine while joins are rare and
		// by hand, a conditional put in dbapi would close it.
		j.note(ctx, wire.CodeForbidden, reason.JoinWaits)
		httpapi.Fail(w, wire.CodeForbidden)
		return
	case !errors.Is(err, dbapi.ErrNotFound):
		j.note(ctx, wire.CodeUnreachable, reason.DB)
		httpapi.FailCoded(w, wire.CodeUnreachable, err)
		return
	}
	if _, err := j.S.DB.OnceUse(ctx, dbapi.OnceJoin, hash, j.S.Now()); err != nil {
		// THE busy arm of this path. OncePeek and BlobGet are reads, and a
		// read passes a second program's write lock (WAL); OnceUse is the first
		// write, so this is where a busy database really stops a server that
		// is joining correctly. Only a request that holds the secret gets here
		// -- the MAC is checked above -- so the answer tells nobody anything
		// about a code he does not already hold.
		if dbwait.Busy(err) {
			httpapi.FailCoded(w, wire.CodeUnreachable, err)
			return
		}
		// Another request with the same secret spent it between the peek and
		// here, or it ran out in that moment: the one answer for a code that
		// is not there.
		httpapi.Fail(w, wire.CodeVerify)
		return
	}
	if err := j.keep(ctx, req); err != nil {
		// The code is gone and nothing was kept: the admin waits for a
		// request that will never show up in cd.users/pending-api.
		j.note(ctx, wire.CodeUnreachable, reason.DB)
		// The code is the same E_UNREACHABLE as before -- a word of its own
		// for a busy database is a new wire.Code, and that list is in frozen
		// shared/wire (BLOCKED, P9 of the v1.3 row). What the error adds is
		// the Retry-After and the reason in the log, which is the whole of
		// what tells busy from away for a client that reads headers and for a
		// human who reads the log.
		httpapi.FailCoded(w, wire.CodeUnreachable, err)
		return
	}
	j.note(ctx, "", reason.None)
	httpapi.NoContent(w)
}

// note writes one line of the log about a join. A log that does not answer
// never costs the join: the attempt is the event, and whoever looks for a hole
// does not need the line to be written, only the event to be missing.
func (j *Join) note(ctx context.Context, code wire.Code, why reason.Reason) {
	if j.Log == nil {
		return
	}
	_ = j.Log.Pairing(ctx, audit.ActionJoin, code, why)
}

// keep stores the request of a new API server as dbapi.BlobJoin. The admin
// compares tunnel.JoinFingerprint of exactly these values with the output of
// the new server's ulh script before listing it, so nothing here is trusted
// on its own.
func (j *Join) keep(ctx context.Context, req tunnel.JoinRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return j.S.DB.BlobPut(ctx, dbapi.Blob{Kind: dbapi.BlobJoin, ID: req.Node, Data: data,
		Updated: j.S.Now()})
}

// Waiting answers cd.users/pending-api.
func (j *Join) Waiting(ctx context.Context) ([]core.PendingAPI, error) {
	blobs, err := j.S.DB.BlobFrom(ctx, dbapi.BlobJoin, 0)
	if err != nil {
		return nil, err
	}
	out := make([]core.PendingAPI, 0, len(blobs))
	for _, b := range blobs {
		// The stored row goes through the contract's parser again: a broken
		// row must never reach the shell as something to sign.
		req, err := tunnel.ParseJoinRequest(b.Data)
		if err != nil {
			continue
		}
		out = append(out, core.PendingAPI{Node: req.Node, Noise: req.Noise, Sign: req.Sign,
			TLS: req.TLS, URL: req.URL, Since: b.Updated})
	}
	return out, nil
}

// noisePublic returns this API's public Noise key, derived from the private
// half in the configuration so that nothing has to keep the two in step.
func (j *Join) noisePublic() (enc.Bytes, error) {
	k, err := ecdh.X25519().NewPrivateKey(j.S.Cfg.NoiseKey.Bytes())
	if err != nil {
		return nil, err
	}
	return enc.Bytes(k.PublicKey().Bytes()), nil
}
