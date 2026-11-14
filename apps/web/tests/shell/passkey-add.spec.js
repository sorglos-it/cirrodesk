// "Passkey hinzufügen" (apps/users.js): what a fresh WebAuthn credential turns
// into in the trust list, and which rp a version with it carries. The
// credential is built out of Web Crypto (P-256 is in all three engines), so
// these cases need no authenticator and run everywhere.
import { test, expect } from "../fixtures/page.js";

test("a new passkey becomes a key the trust list accepts", async ({ shell }) => {
  const got = await shell.evaluate(async () => {
    const users = await import("/assets/js/apps/users.js");
    const proof = await import("/assets/js/contracts/proof.js");
    const pair = await crypto.subtle.generateKey({ name: "ECDSA", namedCurve: "P-256" }, true, ["sign"]);
    const spki = await crypto.subtle.exportKey("spki", pair.publicKey);
    const cred = (alg, pub) => ({
      rawId: new Uint8Array([1, 2, 3]).buffer,
      response: { getPublicKeyAlgorithm: () => alg, getPublicKey: () => pub },
    });
    const key = await users.passkeyKey(cred(-7, spki));
    const refused = async (c) => users.passkeyKey(c).then(() => false, () => true);
    return {
      kind: key.kind === proof.KIND_PASSKEY, alg: key.alg, cred: [...key.credId],
      valid: await proof.validateUserKey(key),
      eddsa: await refused(cred(-8, spki)),
      noKey: await refused(cred(-7, null)),
      // ES256 named, but the bytes are no key at all.
      junk: await refused(cred(-7, new Uint8Array(70).buffer)),
    };
  });
  expect(got).toEqual({ kind: true, alg: "ES256", cred: [1, 2, 3], valid: true,
    eddsa: true, noKey: true, junk: true });
});

test("the rp of a passkey version comes from the list, or from this page where none is set",
  async ({ shell }) => {
  const got = await shell.evaluate(async () => {
    const { passkeyRP } = await import("/assets/js/apps/users.js");
    const at = { hostname: "ngc.example.org", origin: "https://ngc.example.org" };
    const set = { id: "example.org", origins: ["https://ngc.example.org"] };
    return [
      passkeyRP({ id: "", origins: [] }, at),
      passkeyRP(set, at),
      passkeyRP(set, { hostname: "evil.example", origin: "https://evil.example" }),
      passkeyRP({ id: "", origins: [] }, { hostname: "localhost", origin: "http://localhost:8080" }),
    ];
  });
  expect(got).toEqual([
    { id: "ngc.example.org", origins: ["https://ngc.example.org"] },
    { id: "example.org", origins: ["https://ngc.example.org"] },
    null,
    null,
  ]);
});
