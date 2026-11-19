# Installation über ulh

CirroDesk wird mit [ulh](https://github.com/sorglos-it/unknown-helper) installiert, konfiguriert und deinstalliert.
Die Skripte liegen direkt im öffentlichen ulh-Repo: `apps/cli-ulh/assets/scripts/<skript>.sh`,
eingetragen in dessen `apps/cli-ulh/assets/catalog.yaml` unter der Kategorie „CirroDesk", mit einem
Test in `apps/cli-ulh/tests/<skript>.test.sh`. Aufbau nach einem vorhandenen Skript in
`apps/cli-ulh/assets/scripts/` und ulhs `docs/ulh/guide.md`; `docs/ulh/scripts.md` wird nicht von Hand
gepflegt, sondern mit `bash tools/ulh/gen-scripts-doc.sh` erzeugt, und `CHANGELOG.md` bekommt eine
Zeile. Der Eintrag im Katalog braucht alle fünf Sprachen (de, en, fr, it, es) — eine fehlende
Übersetzung macht ulhs eigenen Test rot; was die Skripte selbst ausgeben, ist englisch.
(Bis 22.09.2026 lagen die Skripte in `scripts/`, der Katalog hieß `config.yaml`, und es gab
`SCRIPTS.md`, `TREE.md` und `scripts/_template.sh`. Alles fünf ist seit dem Umbau von ulh
auf den Projektstruktur-Standard weg; gemessen am 02.10.2026 auf ulhs `main`.) Jedes Teil schreibt das Skript für seinen Server auf einem eigenen Branch im ulh-Repo
(`cd-web`, `cd-api`, `cd-controller`, `cd-db`) und liefert per Pull Request; der Koordinator
hält diese Konventionen fest und führt zusammen. Jede ulh-Installation holt beim Start `main`,
deshalb kommt nur Geprüftes dorthin.

## Skripte

| Skript | Server | Was es tut |
|---|---|---|
| `cd-db` | DB 1 bis 3 | MariaDB mit Galera und TLS einrichten, Verschlüsselung einschalten, Konten `cd` (lesen und schreiben) und `cd_migrate` (Schema) anlegen; zeigt am Ende das Zertifikat der CA der Installation (PEM) |
| `cd-api` | API 1 und 2 | `cd-api` und `cd-store` nach `/usr/lib/cd/`, `api.json` schreiben, eigene Schlüssel und ein selbst signiertes TLS-Zertifikat (100 Jahre gültig) erzeugen, beim Beitritt mit `join_code` nur die gemeinsamen Schlüssel der Installation von einer laufenden API übernehmen (`secretKey` und `challengeKey`: Beide müssen auf allen API-Knoten gleich sein, sonst beantwortet jeder Knoten denselben unbekannten Anmeldenamen anders und verrät damit, welche Benutzer es gibt), `cd-api.service` starten; zeigt am Ende den Einmal-Code (erster API-Server) und den TLS-Fingerabdruck (`pin.sh`), beim Beitritt `tunnel.JoinFingerprint` und den Anker des Beitrittscodes |
| `cd-web` | Web | Apache und PHP einrichten, das Webinterface nach `/var/www/cd/` (Webroot nur dessen `public/`), Zertifikat (Let's Encrypt, selbst signiert oder Import), `web.json` schreiben, Lastverteilung vorbereiten; prüft das Zertifikat jeder API gegen ihren Fingerabdruck und schreibt `config.APITrustFile` |
| `cd-controller` | alle | `cd-controller` nach `/usr/lib/cd/`, Schlüssel erzeugen, mit dem Kopplungscode koppeln, `controller.json` schreiben (mit `roles`, bei Rolle `db` auch `migrate`), `cd-controller.service` starten, Podman installieren, wenn der Server Dienste trägt |
| `cd-all` | Pi | alles auf einer Maschine mit SQLite; `configure` koppelt den Controller nach der Ersteinrichtung |

Jedes Skript kann `install`, `configure` und `uninstall`. Ein zweiter Lauf von `install` ändert
nichts, was schon stimmt.

## answer.yaml

Damit die Skripte ohne Rückfragen laufen. Schlüssel je Skript:

- `cd-db`: `hostname` (Name oder IPv4-Adresse **dieses** Knotens, nach `trust.ValidHost` und
  `host.sh`; er muss in `nodes` vorkommen und wird der Name im Knotenzertifikat),
  `nodes` (Liste der Knoten, je `<host>:<port>` nach `trust.ValidHostPort`),
  `bootstrap` (erster Knoten: ja/nein), `password_cd`, `password_migrate` (zwei verschiedene
  Passwörter, ohne Komma: ulh reiht die Antworten mit Kommas) und `ca_dir` (Ordner für den
  privaten Schlüssel der CA, Vorgabe `/etc/cd/db-ca`, nur beim Bootstrap-Knoten). Der erste Knoten
  erzeugt die CA der Installation und zeigt ihr Zertifikat (PEM); die anderen zeigen ihre Anfrage
  (CSR) und bekommen das Zertifikat mit `cd-db sign` auf dem Bootstrap-Knoten. Der private
  Schlüssel der CA und der des Knotens kommen nie in answer.yaml, in eine Konfigurationsdatei oder
  ins Repo. Das Passwort für die Kopie zu einem neuen Knoten (`wsrep_sst_auth`) steht ebenfalls
  nicht in answer.yaml: `cd-db` erzeugt es je Installation und legt es mit Modus 0400 für `mysql`
  ab.
- `cd-api`: `node` (z. B. `api-1`; auch der Name in der Vertrauensliste und Teil des
  Fingerabdrucks beim Beitritt), `hostname` (Name, unter dem Controller und Apache den Server
  erreichen, nach `trust.ValidHost`; die Skripte prüfen wie `shared/testdata/ulh/host.sh`:
  DNS-Name in Kleinbuchstaben oder IPv4, kein IPv6. Daraus `url` = `wss://<hostname>:7443/tunnel`
  in `api.json` und das TLS-Zertifikat für diesen Namen), `db_driver`,
  `db_nodes` (je `<host>:<port>` nach `trust.ValidHostPort`), `db_name`, `db_user`,
  `db_password`, `db_ca` (MariaDB: Zertifikat der CA von `cd-db`), `store_url` (nach
  `trust.ValidStoreURL`) und `store_key` (der globale Store: Adresse und öffentlicher Schlüssel,
  base64url; nur beide zusammen oder keins, ohne Angabe der offizielle CirroDesk-Store, eine normale
  Installation braucht also keins von beiden; `api.json` nimmt sie als `storeUrl` und
  `storeKey`), `join_code` (leer beim ersten API-Server, sonst der Beitrittscode einer laufenden
  API: `cd-join:` und base64url des JSON, geprüft mit `tunnel.ParseJoinCode`)
- `cd-web`: `hostname` (wie bei `cd-api`, nach `trust.ValidHost` und `host.sh`), `cert_mode`
  (`letsencrypt`, `selfsigned`, `import`), `cert_file`, `key_file`, `apis` (Adressen
  `https://<hostname>:7443` nach `trust.ValidWebURL`), `api_pins` (je API ihr
  TLS-Fingerabdruck, wie `cd-api` ihn am Ende mit `pin.sh` zeigt, in derselben Reihenfolge wie
  `apis`)
- `cd-controller`: `name` (Name des Servers nach `wire.ValidName`: Kleinbuchstaben, Ziffern und
  Bindestrich, beginnt mit einem Buchstaben, höchstens 32 Zeichen), `pair_code` (der Kopplungscode
  aus der Shell: `cd-pair:` und base64url des JSON, geprüft mit `tunnel.ParsePairCode`; nie in der
  Kommandozeile, das Skript reicht ihn mit `printf '%s\n' "$pair_code" | cd-controller pair`
  durch), `roles`
  (`web`, `api`, `db`, `services`); mit Rolle `db` zusätzlich `migrate_nodes` (je
  `<host>:<port>` nach `trust.ValidHostPort`), `migrate_name`, `migrate_user` (`cd_migrate`),
  `migrate_password` und `migrate_ca` (Zertifikat der CA von `cd-db`) für `migrate` in
  `controller.json`
- `cd-all`: `hostname` (wie bei `cd-api`, nach `trust.ValidHost` und `host.sh`), `cert_mode`
  für `install`; `pair_code` für `configure` nach der Ersteinrichtung

## Regeln

- Programmdateien übernimmt ein Skript nur, wenn die Signatur des Release stimmt. Die Vorlage mit
  Testdaten ist `shared/testdata/release/verify.sh`, aufgerufen je Teil als
  `verify.sh <Verzeichnis> <Teil> <Arch>` (etwa `cd-api amd64` oder `cd-web all`). Im
  Verzeichnis liegen `release.json` (ein `sign.Envelope`), `release-key.pem` (der Release-Schlüssel
  als PEM, fest im Skript) und die geladenen Dateien. Werkzeuge aus den Paketquellen: `sh`, `jq`,
  `openssl` ab 3.0, `base64`, `sha256sum`, `od`, `tr`, `cut`, `wc` und `mktemp`. Die Schritte:
  1. `payload` und `sig` aus `release.json` (je ein Text) streng von base64url in Bytes umwandeln:
     nur `A-Z a-z 0-9 - _`, und die Bytes ergeben wieder genau diesen Text.
  2. Die Nachricht ist genau `cd/release/v1`, ein Null-Byte und die Nutzlast;
     `openssl pkeyutl -verify -pubin -inkey release-key.pem -rawin` muss sie mit `sig` bestätigen.
  3. `files` der Nutzlast ist eine Liste mit genau einem Eintrag für Teil und Arch; darin sind `url`
     und `sha256` Text und `size` eine Zahl.
  4. Der Dateiname ist das letzte Stück der `url` nach dem letzten `/` und passt auf
     `^[A-Za-z0-9][A-Za-z0-9._-]*$`. Nur diese Datei wird geprüft und danach installiert; andere
     Dateien im Verzeichnis und Dateien anderer Teile zählen nicht.
  5. `sha256` wie in 1 umwandeln, in Hex (`od -An -tx1 | tr -d ' \n'`, genau 64 Zeichen) mit
     `sha256sum` der Datei vergleichen, ebenso ihre Größe mit `size`.
  6. Nie ein älteres Release (das tut das ulh-Skript, nicht `verify.sh`): `seq` der Nutzlast mit der
     Nummer in `/var/lib/cd/release.seq` vergleichen (fehlt die Datei, ist noch nichts
     installiert). Ist sie kleiner, bricht das Skript ab; dieselbe Nummer darf zur Reparatur noch
     einmal installiert werden. Nach einer erfolgreichen Installation schreibt das Skript die Nummer
     dorthin. Go prüft genauso (`store.VerifyRelease` mit dieser Nummer als `lastSeq`, sonst
     `store.ErrRollback`).

  `verify.sh` legt Arbeitsdateien in einem Verzeichnis von `mktemp -d` ab, das es am Ende löscht,
  und schreibt nichts in das Verzeichnis der Dateien. Geht etwas schief, meldet es
  `release wrong: <Grund>` und endet mit 1, bei falschem Aufruf mit 2; nur wenn alles stimmt, gibt es
  `release ok` aus und endet mit 0. `TestReleaseScript` in `shared/store` prüft das in 31 Fällen;
  auf den Linux-Servern muss er mit dem echten `jq` bestehen. `release.json` bleibt im Format
  `release/v1`; spätere Formate kommen in eigene Dateien (Regel „Signierte Dokumente, die bleiben“
  in `shared/README.md`).
- Release-Dateien benennen Softwareteile: `cd-web` (`all`), `cd-api` und `cd-controller` (`amd64`,
  `arm64`). Die Migrationen stecken in `cd-api`.
- Rollen eines Servers (`roles`): `web`, `api`, `db`, `services`. Welche Softwareteile dazu gehören,
  sagt `tunnel.Components`: immer `cd-controller`, dazu `cd-web` bzw. `cd-api` je nach Rolle.
- Mit Rolle `db` kann der Controller beim Update die Migrationen einspielen. Genau einen dieser
  Controller bestimmt die API, die das Update steuert (sie hält `dbapi.LockUpdate`): unter den
  verbundenen Controllern mit Rolle `db` den mit der kleinsten Node-ID, scheitert er, den nächsten.
  Dafür steht der Zugang des Kontos `cd_migrate` (auf dem Pi: SQLite-Datei und Schlüssel) in
  `controller.json` unter `migrate`. Andere Server bekommen diesen Eintrag nicht.
- Konfigurationsdateien sind striktes JSON (`shared/README.md`, „Striktes JSON“): nie `null`,
  ein leeres Feld fehlt ganz. Es gelten die Regeln von `config`: in `controller.json` ein bis vier
  verschiedene `roles`, `migrate` genau mit Rolle `db`, `webConfig` mit Rolle `web`; in `web.json`
  `rp.id` und `rp.origins` beide oder keins, `rp.id` ein Hostname in Kleinbuchstaben. Kein Passwort
  darf wörtlich `[secret]` lauten.
- **Hosts und Adressen:** Ein Skript übernimmt einen Host nur nach `host.sh`
  (`shared/testdata/ulh/host.sh`: DNS-Name in Kleinbuchstaben oder IPv4, kein IPv6 – ein
  IPv6-Netz erreicht man über den Hostnamen) und schreibt Adressen nur in der Form der Regel
  „Adressen“ in `shared/README.md`. Wie `host.sh` setzt es vorher `LC_ALL=C`, sonst nehmen
  `[a-z]` und `grep` in einer anderen Sprachumgebung weitere Buchstaben an.
- Schlüssel erzeugt immer der Server selbst; nichts Privates geht über die Leitung, außer beim
  Beitritt weiterer API-Server (`join_code`), und dort verschlüsselt und erst nach der Unterschrift
  eines Admins.
- Konfigurationsdateien mit Schlüsseln gehören root und haben Modus 0600 (`config`).
- `uninstall` fragt, was mit den Daten passiert: löschen oder aufheben. Aufgehobene Daten bleiben
  unter `/var/lib/cd/` liegen, damit eine Neuinstallation sie wiederfindet.
- Kopplungs- und Beitrittscodes zeigt die Shell erst, wenn sie zu ihrer neuesten geprüften Fassung
  der Vertrauensliste passen (`tunnel.PairCode.Check`, `tunnel.JoinCode.Check`).
- Am Ende zeigt `cd-api` beziehungsweise `cd-all` den Einmal-Code für das erste Admin-Konto an
  (`session.ValidSetupCode`); ihn liest nur root, und er kommt in kein Log. `cd-controller`
  zeigt seinen Fingerabdruck (`tunnel.Fingerprint` seines Schlüssels) und den Anker aus dem
  Kopplungscode (Nummer und `tunnel.Fingerprint` des Hashes); der Admin vergleicht beides mit der
  Anzeige der Shell. Beim Beitritt zeigt `cd-api` ebenso den Anker aus dem Beitrittscode (Nummer
  und `tunnel.Fingerprint` des Hashes) zum selben Vergleich.
- **TLS intern, nie ungeprüft:** Kein Skript und kein Programm schaltet eine Zertifikatsprüfung
  ab. Den Fingerabdruck eines TLS-Zertifikats rechnet `shared/testdata/tls/pin.sh` (SHA-256 des
  öffentlichen Schlüssels, base64url; Testdaten `tls/api.pem`, `tunnel.json` `tlsPin`, geprüft von
  `TestTLSPinScript`). Kann `pin.sh` die Datei nicht lesen, gibt es keinen Fingerabdruck aus und
  endet mit 1, bei falschem Aufruf mit 2; die Skripte prüfen den Exit-Status und brechen sonst ab.
  Vertrauen trägt nur dieser Fingerabdruck des Schlüssels, nicht das Zertifikat: `cd-api` erzeugt
  sein Zertifikat selbst, selbst signiert für `hostname` und seinen Schlüssel, 100 Jahre gültig, und
  erneuert es nur mit demselben Schlüssel; einen neuen Schlüssel gibt es nur zusammen mit einem
  neuen Fingerabdruck in der Vertrauensliste. `cd-web` holt das Zertifikat jeder API
  (`openssl s_client -connect <host>:7443 -servername <host>`), vergleicht den Fingerabdruck mit
  `api_pins`, schreibt passende Zertifikate nach `/etc/cd/tls/api-trust.pem`
  (`config.APITrustFile`) und `web.json` (`apis` mit `url` und `tls`) und richtet Apache mit
  `SSLProxyEngine on`, `SSLProxyVerify require`, `SSLProxyCACertificateFile` = diese Datei,
  `SSLProxyCheckPeerName on`, `SSLProxyCheckPeerExpire on`, `SSLProxyProtocol TLSv1.3` und
  `SSLProxyCipherSuite TLSv1.3 TLS_AES_256_GCM_SHA384` ein (Entwurf 4.5: nur AES-256-GCM). Ist der
  Controller des Webservers gekoppelt, übernimmt er das aus der Vertrauensliste. Bei jedem Aufbau
  eines Tunnels und einmal am Tag vergleicht er das Zertifikat, das er über `tunnel.PinnedTLS`
  bekommen hat, mit dem in dieser Datei; weicht es ab (derselbe Schlüssel, das beweist der
  Fingerabdruck), schreibt er die Datei neu und lädt Apache neu. So sperrt weder ein erneuertes
  Zertifikat noch ein Datum Apache aus. Den Beitritt und die Kopplung führen die Programme aus
  (`cd-api`, `cd-controller`); sie prüfen TLS mit dem Fingerabdruck aus dem Code. Der
  Kopplungscode steht dabei nie in der Kommandozeile: `cd-controller pair` fragt ihn am Terminal
  ab und nimmt ihn sonst als erste Zeile der Standardeingabe. So reicht das Skript ihn mit einer
  eingebauten Shell-Anweisung durch (`printf '%s\n' "$pair_code" | cd-controller pair`) und legt
  ihn in keine Prozessliste: `/proc/<pid>/cmdline` liest ohne `hidepid` jeder lokale Benutzer.
  `cd-controller pair` legt auch die leere lokale Datenbank an. Das Skript startet den Dienst
  `cd-controller` deshalb erst nach `pair`: Fehlt die Datenbank, startet der Dienst nicht (A2,
  Entscheidung vom 15.11.2026), weil er eine frische Anlage sonst nicht von einer unterscheiden
  könnte, deren Datenbank verloren ging.
- **Lastverteilung in Apache:** `config.BalancerFile` (`/etc/cd/apache-apis.conf`).
  - Es schreiben nur `cd-web` bei `install`/`configure` (aus `apis` und `api_pins`, nach der
    Prüfung der Fingerabdrücke) und danach der Controller mit Rolle `web` (aus `web.json`),
    jeweils atomar, Modus 0644; Apache wird erst nach `apachectl configtest` neu geladen. Sonst
    schreibt niemand hinein.
  - Der virtuelle Host von `cd-web` bindet die Datei ein.
  - Inhalt: je API ein `BalancerMember` mit der Adresse aus `web.json`, die TLS-Einstellungen der
    Regel „TLS intern“ und die Weiterleitung von `/api/` samt WebSocket an den Balancer.
  - Die genaue Form ist diese Vorlage; `cd-web` und der Controller schreiben genau sie,
    Zeile für Zeile. Das Beispiel zeigt zwei APIs; je weiterer API kommt genau eine weitere
    `BalancerMember`-Zeile dazu, mit derselben Einrückung und denselben Werten hinter der
    Adresse. Es gibt keine andere erlaubte Abweichung.

    ```apache
    # /etc/cd/apache-apis.conf
    # Written by cd-web (install, configure) and after that by cd-controller with
    # role web, both from /etc/cd/web.json. Never edit by hand; every write is
    # atomic, mode 0644, and Apache is reloaded only after apachectl configtest.
    SSLProxyEngine on
    SSLProxyVerify require
    SSLProxyCACertificateFile /etc/cd/tls/api-trust.pem
    SSLProxyCheckPeerName on
    SSLProxyCheckPeerExpire on
    SSLProxyProtocol TLSv1.3
    SSLProxyCipherSuite TLSv1.3 TLS_AES_256_GCM_SHA384

    <Proxy "balancer://cd-api">
        BalancerMember "https://api-1.cd.test:7443" retry=5 timeout=65 connectiontimeout=2 keepalive=On upgrade=websocket
        BalancerMember "https://api-2.cd.test:7443" retry=5 timeout=65 connectiontimeout=2 keepalive=On upgrade=websocket
        ProxySet lbmethod=byrequests failontimeout=On
    </Proxy>

    # A 503 of the balancer's own carries the same header a busy database gives.
    Header always set Retry-After "5" "expr=%{REQUEST_STATUS} == 503"
    # A member that cannot be reached belongs in the error log, not only in the
    # access log.
    LogLevel proxy:info proxy_balancer:info

    ProxyPreserveHost Off
    ProxyPass        "/api/" "balancer://cd-api/api/"
    ProxyPassReverse "/api/" "balancer://cd-api/api/"
    ```

    Wer diese Vorlage ändert, ändert einen Vertrag: `TestBalancerTemplateMatchesTheContract`
    in `apps/server-controller/classes/apache` vergleicht sie byteweise mit dem, was
    `apache.Render` erzeugt, und wird rot, sobald eine der beiden Seiten allein wandert.
- **Reihenfolge nach der Ersteinrichtung:** zuerst den Controller des Webservers koppeln. Erst er
  schreibt `trust-anchor.json`; bis dahin kann sich nur der Browser anmelden, der die
  Ersteinrichtung gemacht hat.
- **`cd-all`:** `install` richtet alles ein, startet die API und zeigt den Einmal-Code; `web.json`
  bekommt Adresse und Fingerabdruck der lokalen API ohne Rückfrage. Der Controller läuft erst nach
  `configure` mit `pair_code`, einem Kopplungscode aus der Shell nach der Ersteinrichtung (alle vier
  Rollen, `migrate` auf die SQLite-Datei). Bis dahin kennt nur der Browser der Ersteinrichtung den
  Anker.
- Nach der Installation spielt der Controller die Updates ein (4.8); ulh bleibt der Notweg.
- **Kopplung eines Controllers an die API** (eingetragen 03.10.2026, aus
  Aufgabe 12 des API-Plans): der Controller fragt
  `POST <tunnel.PathPair>?secret=<base64url des SHA-256 ueber den
  Kopplungscode>`. Der Code selbst geht **nie** ueber die Leitung, nur sein
  Streuwert — und er ist einmal brauchbar. `cd-controller` braucht dafuer die
  Antwort `pair_code`, die oben schon steht.
