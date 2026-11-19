# Den CirroDesk-Store betreiben

Entscheidung S1 (Thomas, 15.11.2026): **A.** Den globalen Store (`trust.StoreGlobal`, Kennung
`global`) betreibt Thomas als statische Dateien auf seinem eigenen Webspace. Der Signierschlüssel
liegt offline als Datei; später geht er über die Vertrauensliste auf den Stick.

**Die Adresse fehlt noch.** Bis Thomas sie nennt, steht hier überall der Platzhalter
`https://store.example.invalid/cd`. Er passt in die Grammatik `trust.ValidStoreURL`: `https://`,
Host, optional Pfad, **kein** `/` am Ende.

Jeder Schritt unten läuft als Test: `TestFirstCDStore` in `apps/cli/classes/cdstore_test.go`.
Ist er rot, stimmt diese Anleitung nicht mehr.

## 1. Werkzeug bauen

```bash
go build -o cd-store ./apps/cli
```

## 2. Schlüssel ziehen (nur Thomas, nie ein Agent)

Auf dem Rechner, der offline bleibt. Der erste Katalog legt die Schlüsseldatei an, wenn sie fehlt
(Rechte 0600), und nennt nur ihren Pfad und den öffentlichen Teil:

```bash
mkdir -p ~/cd-keys && chmod 700 ~/cd-keys
```

Die Datei heißt im Folgenden `~/cd-keys/cd-store.key`. **Sie geht nie auf den Webspace, in kein
Repo, in keine Nachricht.** Eine Sicherung gehört auf einen zweiten Datenträger, der offline
liegt. Ist der Schlüssel weg, bekommt der Store nie wieder eine App (Befund 2 in
`docs/app-weg.md`): jeder neue Katalog braucht dieselbe Unterschrift.

Den öffentlichen Teil nennt jederzeit wieder:

```bash
./cd-store key ~/cd-keys/cd-store.key
```

Die Zeile `Öffentlicher Schlüssel:` ist der Wert für `storeKey` (`/etc/cd/api.json`) und für
`CD_STORE_KEY_DEFAULT` in ulh (Abschnitt 6).

## 3. Katalog bauen

Ein leerer Ordner ist der Webspace im Kleinen. Jede App kommt als Paket hinein, dann wird
unterschrieben:

```bash
mkdir site
./cd-store pack apps/store-calculator site/org.sorglos.calc-1.0.0.cd
./cd-store catalog site ~/cd-keys/cd-store.key global 1 https://store.example.invalid/cd
./cd-store check -store global site
```

- `global` ist die Kennung des Stores, immer genau so.
- `1` ist die Folgenummer. **Jeder neue Katalog braucht eine höhere** als der letzte;
  `cd-store catalog` liest den alten `catalog.json` im Ordner und weist eine kleinere ab. Darum
  den Ordner `site/` behalten und nicht jedes Mal neu anlegen.
- Die Basisadresse ist die Adresse des Stores; die Download-Links im Katalog werden
  `<Adresse>/<Paketname>`.
- `check` muss `store ok` sagen. Sonst nichts hochladen.

Danach liegen in `site/` genau: `catalog.json`, die Pakete (`*.cd`, bei Kauf-Apps zusätzlich
`*.cd.enc` und `*.price.json`) und `store-key.pem`, der **öffentliche** Teil des Schlüssels.

## 4. Hochladen

Den **Inhalt** von `site/` in das Verzeichnis des Webspace, das unter der Store-Adresse
erreichbar ist, zuerst die Pakete, **zuletzt** `catalog.json`: Ein Katalog, der auf ein Paket
zeigt, das noch nicht oben liegt, lässt jede Installation scheitern, die genau dann holt.

`*.price.json` muss nicht hoch (es steckt im Katalog), schadet aber nicht.

## 5. Webspace einstellen

Die API holt die Dateien selbst, auf dem Server, mit Go `net/http` (`appstore/fetch.go`). Daraus
folgt genau das:

| Was | Warum |
|---|---|
| **HTTPS** mit einem Zertifikat, dem Debian/Ubuntu ab Werk trauen | `trust.ValidStoreURL` nimmt nur `https://`. |
| **Antwort 200, keine Weiterleitung** | Die API folgt keiner 3xx (`CheckRedirect`). Auch `http`→`https` oder `www`-Umleitungen der Adresse selbst brechen ab. Die Adresse genau so eintragen, wie der Webspace sie ohne Umleitung ausliefert. |
| **Keine Anmeldeseite, kein Bot-Schutz, kein Cookie-Banner davor** | Die API ist kein Browser; sie bekommt sonst HTML statt JSON. |
| `catalog.json`: `Content-Type: application/json`, `Cache-Control: no-cache` | Ein neuer Katalog soll beim nächsten Holen ankommen, nicht erst nach einem Proxy-Ablauf. |
| `*.cd`, `*.enc`: `Content-Type: application/octet-stream` | Kein Umschreiben, keine Zeichensatz-Umwandlung. Lange cachebar: ein Paketname mit Fassung ändert sich nie. |
| `store-key.pem`: `Content-Type: text/plain` | Nur zum Nachsehen für Menschen. |
| Kompression darf an, Bytes dürfen sich nicht ändern | Prüfsumme und Größe im Katalog gelten für die ausgelieferten Bytes. |
| CORS: **nicht nötig** | Die Oberfläche zeigt die Adresse nur an (`setup.js`) und holt nichts vom Store. |

Für Apache (`.htaccess` im Store-Verzeichnis, falls der Webspace das zulässt):

```apache
AddType application/json .json
AddType application/octet-stream .cd .enc
AddType text/plain .pem
<Files "catalog.json">
  Header set Cache-Control "no-cache"
</Files>
Options -Indexes
```

Prüfen von außen, ohne Umleitung:

```bash
curl -sS -o /dev/null -w '%{http_code} %{content_type}\n' https://store.example.invalid/cd/catalog.json
```

Erwartet: `200 application/json`.

## 6. ulh: der Platzhalter (nur Thomas schiebt dort)

In ulh steht `CD_STORE_KEY_DEFAULT="REPLACE-THIS-LINE-WITH-THE-CD-STORE-KEY"` in
`apps/cli-ulh/assets/scripts/cd-api.sh:97` und `apps/cli-ulh/assets/scripts/cd-all.sh:177`. Wenn Adresse und Schlüssel feststehen,
in **beiden** Dateien:

```sh
CD_STORE_KEY_DEFAULT="<Wert der Zeile 'Öffentlicher Schlüssel:' von cd-store key>"
CD_STORE_URL_DEFAULT="<Store-Adresse, ohne / am Ende>"
```

Steht für die Adresse noch keine `…_DEFAULT`-Zeile neben der des Schlüssels, kommt sie direkt
darunter dazu, und das Skript nimmt sie wie den Schlüssel als Vorgabe für `store_url`. Den
Wortlaut der Zeilen um Zeile 97 und 177 konnte der Store-Partner nicht lesen (das ulh-Repo war von
hier aus nicht erreichbar); die Namen oben sind die aus `docs/was-fehlt.md`, S1.

Ein Push nach ulh-`main` erreicht per Auto-Update sofort alle ulh-Nutzer (A20).

## 7. Eine App dazu, eine neue Fassung

Paket in `site/` legen, `catalog` mit der nächsten Folgenummer, `check`, hochladen wie in 4.
Eine alte Fassung bleibt liegen, solange der Katalog sie nennt.

## Was danach kommt

- Der Schlüssel geht auf den Stick (S3/S4): über die Vertrauensliste, mit einer Änderung in
  `shared/`. Bis dahin unterschreibt die Datei.
