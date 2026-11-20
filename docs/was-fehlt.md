# Was fehlt, bis das jemand anderes benutzen kann

**Geschrieben am 30.10.2026.** Diese Datei baut nichts und entscheidet nichts. Sie ist eine
Landkarte.

## Stand am 20.11.2026 (nachgemessen)

Jede Zeile der Abschnitte 1–3 ist am 20.11.2026 neu gemessen: am CirroDesk-Code auf `main`
(`b7a21def`) und an ulh, dem Installationswerkzeug. Was danach auf `main` kam (#247 bis #252), ändert
keinen dieser Einträge; alle Zeilenangaben gelten für `d359b96a`. Dazu kommen zwei neue Zeilen aus der
Sichtung der API-Zweige (A29, A30). Die Tabellen ab „Stand am 15.11.2026" bleiben als Geschichte
stehen.

**Das Wichtigste**

- **69 Zeilen** (16 sperren, 35 halten auf, 18 stören uns; A29 und A30 sind neu, die Zahlen 14 und 32
  unten zählten anders): **22 zu, 21 teilweise, 26 offen.**
- **„Zu" heißt: im Code.** Beim Nutzer ist davon fast nichts. Das letzte Release ist 0.2.2 vom
  23.10.2026 (`gh release list`), ulh-`main` steht seit dem 07.11. still. Was danach gebaut wurde, kommt
  erst mit dem **ersten cd-Release** und dem **Push von `cd-umbenennung`** an — nur zusammen, beides
  macht Thomas. Schon beim Nutzer sind S13 (seit 0.2.0), A20 und U17 (auf ulh-`main`).
- **Abschnitt 5 ist zum Teil überholt.** Thomas hat am 24.10.2026 geschrieben: „keine ahnung entscheide
  selbständig nach den grundsätzen stabilität vor sicherheit vor performance"
  (`docs/thomas-entscheidet.md:5-6`). Danach sind dort 28 von 32 Fragen entschieden (hier stand 27), darunter P10,
  `noexec` (A15), der Digest (A16) und die Rohtasten im Terminal (A17). **Im Entscheidungsprotokoll des
  Entwurfs fehlten diese Entscheidungen**; am 20.11.2026 sind sie nachgetragen. Die
  stehende Zustimmung (S8) ist am 15.11. entschieden, der Push für A20 am 30.10. geschehen.
- **Bei Thomas bleiben fünf Dinge**; die sieben offenen Entscheidungen sind am 20.11.2026 gefallen. Ohne
  ihn baubar ist der Rest,
  meist in kleinen Stücken. Beide Listen stehen am Ende dieses Abschnitts.

**Abkürzungen in den Belegen.** `api/`, `ctl/`, `db/` = `apps/server-api/classes/`,
`apps/server-controller/classes/`, `apps/server-db/classes/`; `web/` = `apps/web/public/assets/js/`.
`TE` = `docs/thomas-entscheidet.md`, `Buch` = `docs/vertraege-v1.2.md`, `Entwurf` =
`docs/superpowers/specs/2026-09-19-ngc-kern-design.md`. ulh-Skripte liegen unter
`apps/cli-ulh/assets/scripts/`: `ulh-main:` = ulh-`main` `2181c9b` (das holt jede ulh-Installation),
`cdu:` = ulh-Zweig `cd-umbenennung` `7c1e554`.

### ulh

- **ulh-`main`** (`2181c9b`, 07.11.2026) ist unverändert (`git ls-remote`). Die Skripte heißen dort
  `ngc-*.sh` und erwarten Releases mit `ngc-*`-Dateien. Pull Requests im ulh-Repo: keine
  (`gh pr list -R sorglos-it/unknown-helper --state all` listet 0).
- **`cd-umbenennung`** (`7c1e554`) liegt nur lokal in `Z:\GitHub\unknown-helper`: 9 Commits, 19
  Dateien. Darin die Umbenennung auf `cd-*.sh`, die Balancer-Datei aus dem Vertrag (A4),
  `RuntimeDirectoryPreserve=yes` für die API, das Entpacken der Oberfläche aus `ui/` (S7) und ein
  Unit-Kommentar zum Beitritt.
- **Prüfsatz heute gemessen** (WSL Debian, OpenSSL 3.5.7, Baum per `git archive`, Sonde `(exit 7)` =
  7): `7c1e554` **864 bestanden, 0 rot**; `2181c9b` **849 bestanden, 0 rot**; übersprungen 0. Nicht
  gemessen: ulhs Prüfung auf GitHub, ein echtes Debian außerhalb von WSL.
- Auf GitHub liegt noch `ulh-s3a-kommentar` (`2094cac`, alte Namen): die Berichtigung S3a, die
  `cd-umbenennung` fehlt.

### 1. Sperrt — heute

| # | Stand am 20.11.2026 | Beleg | Wer |
|---|---|---|---|
| S1 | **offen.** Entschieden (Variante A, 15.11.), Anleitung da. Es fehlen Store-Schlüssel, Adresse und Katalog; ohne eigenen Store hält die Installation an. | Entwurf :1060; `docs/cd-store-betreiben.md` (PR 214); `cdu:cd-api.sh:96-97` (Adresse `https://store.cirrodesk.org`, Schlüssel ein Platzhalter), Abbruch `:413-416`; `ulh-main` gleich mit ngc-Namen | Thomas: Schlüssel, Webspace, Adresse, Push |
| S2 | **teilweise.** Drei echte Apps gebaut (Rechner, Uhr, Dateien), keine veröffentlicht, weil es keinen Laden gibt (S1). Das alte Beispiel `calculator` liegt entgegen der Entscheidung noch im Web-Teil und geht mit jedem Release raus. | PR 202, 215, 212 (`apps/store-*`); Entwurf :1069; `tools/release/entrypoints-web.txt:13` | Veröffentlichen: Thomas (über S1). Beispiel entfernen: baubar, klein · Web, tools |
| S3 | **offen.** Alle vier Skripte tragen den Dateischlüssel `pN6Z6a9NR0aG6aEJIJ5EHg`, nicht den Stick. | `cdu:cd-api.sh:88`, `cd-all.sh:169`, `cd-controller.sh:89`, `cd-web.sh:105` (`ulh-main` gleich); `docs/release-schluessel.md:17-21`; TE :17 (Fragen 2, 3 bleiben bei Thomas) | Thomas: Entscheidung, Schlüssel |
| S3a | **offen.** Der falsche Satz steht auf beiden ulh-Ständen; die Berichtigung liegt nur auf `ulh-s3a-kommentar`, mit alten Namen. | `cdu:cd-all.sh:167`, `cd-api.sh:85`, `cd-controller.sh:87`, `cd-web.sh:102` | baubar, klein · ulh-Zweig |
| S4 | **offen.** `Sign` unterschreibt weiter die rohe Nachricht; der Stick kann das nicht. | `shared/sign/sign.go:120` | Thomas (folgt aus S3); danach groß · shared/, Web, tools, ulh |
| S5 | **offen.** Entschieden ist ein Unterbefehl `test-key <Datei>` (Frage 6), gebaut ist er nicht; der Schlüssel bleibt eine nackte Zuweisung. | TE :81; `test-key` 0 Treffer auf `ulh-main`, `cdu`, `ulh-s3a-kommentar`; `cdu:cd-api.sh:88` | baubar, mittel · ulh-Zweig |
| S6 | **zu.** Einen eigenen Passkey hinzufügen geht. | PR 203; `web/apps/users.js:689` (`credentials.create`); `apps/web/tests/shell/passkey-add.spec.js` | — |
| S7 | **zu.** Das Web-Paket legt die Oberfläche unter `ui/`, der Controller weist anderes ab, ulh entpackt die `ui/`-Form. | PR 211, 213; `tools/pack-release.sh:178-224`; `TestInstallWebRefusesAPackageWithoutUI` (`ctl/update/install_test.go:1503`); `cdu` `32ef215`, `cd-web.sh:626-656`, ulh-Fall „an archive written by cd-store web goes in" | Release und Push, nur zusammen (Thomas) |
| S8 | **zu.** Ohne gespeicherte Fassung fragt der Controller das Programm oder liest `ui/VERSION`. Ein Lauf an der Anlage fehlt. | PR 210; Entwurf :1062; `ctl/update/install.go:396-455`; `TestInstallWithoutAWrittenVersionAsksTheProgram` (`install_test.go:1532`) | Lauf an der Anlage (Thomas) |
| S9 | **teilweise.** `cd-db backup/restore` gibt es und wird gebaut, kommt aber auf keinen Datenbankknoten: kein Release-Teil, kein ulh-Skript. Nie gegen eine echte MariaDB gelaufen. | PR 198, 231; `tools/build-release.sh:202` (`api_progs="cd-api cd-store"`); `TestThereIsNoDatabasePart` (`tools/release/release_test.go:56`); in ulh installiert kein Skript `cd-db` (beide Stände) | Weg wählen: Koordinator; dann baubar, mittel · tools, cd-store, ulh. Lauf: Anlage |
| S10 | **offen.** Entschieden ist ein Lauf mit echtem Namen vor der Veröffentlichung (Frage 32); gemacht ist er nicht. | TE :384 | Thomas: öffentlicher Name, Anlage |
| S11 | **offen.** Nur VM 301 trägt die Rolle `web`. | `docs/proxmox-sechs.md:40`; kein Eintrag im Protokoll | entschieden 20.11.: zweiter Knoten, gleicher Name, gleiches Zertifikat; danach mittel · ulh, tools, docs |
| S12 | **teilweise, fast zu.** Alle fünf Sprachen haben je 459 Schlüssel in 11 Dateien. Rund 15 Texte sind noch deutsch, etwa „im Monat" und „je Benutzer" an jedem Preis und die Terminal-App. | PR 230, 233, 238–241; `apps/web/public/assets/lang/` (mit node gezählt); `web/apps/store.js:1067,1069`; `apps/web/assets/apps/terminal/app.js:56-59` | baubar, klein · Web |
| S13 | **zu, seit 0.2.0 auch beim Nutzer.** Eine Lizenz mit Ausstellungszeit in der Zukunft wird abgewiesen. War schon am 30.10. erledigt. | `api/appstore/license.go:101`; `TestALicenseIssuedInTheFutureIsRefused` (`license_test.go:736`); `a9592aa6` und `d1f0f175` sind Vorfahren von `v0.2.0` | — |
| S14 | **zu.** Die Oberfläche öffnet die Vertrauensliste bei der neuesten Fassung. Auf `main` erst seit dem 08.11. (PR 164), nicht in 0.2.2. | `589324b8` ist Vorfahr von `main`, nicht von `v0.2.2`; `apps/web/tests/shell/trustgrowth.spec.js:47` | — |
| S15 | **zu im Code.** `cd-ein` wurde nie eingerichtet (Punkt 100). Eine solche Anlage antwortet jetzt `E_NOT_SET_UP` und schickt zur Einrichtung (wie A3). | Buch :8235; `api/login/start.go:79-86`, `finish.go:102-108`; `TestALoginBeforeTheSetupSaysNotSetUp` (`notsetup_test.go:20`) | Maschine geht mit dem Neuaufsetzen (Thomas) |

### 2. Hält auf — heute

| # | Stand am 20.11.2026 | Beleg | Wer |
|---|---|---|---|
| A1 | **teilweise.** Der Controller schreibt jetzt, warum er eine Aktivierung abweist. Die API schreibt bei ihren eigenen Abweisungen (`E_VERIFY`) nur „Der Rechner hat die Aufgabe abgelehnt" ins Audit, keine eigene Zeile. Die Ursache aus R5 ist nie gefunden. | PR 216; `ctl/install/activate.go:193-202`; `TestTheJournalSaysWhyAnActivationWasRefused` (`reason_test.go:232`); `api/appstore/activate.go:199,382`; `api/coreapp/store.go:333-339` | API-Zeile: baubar, klein · API. Ursache: Anlage |
| A2 | **zu.** Fehlt die Datenbank, startet der Controller nicht und sagt es. | PR 225; `ctl/node/run.go:582-585`; `TestAControllerWithoutItsDatabaseDoesNotStart` (`nodb_test.go:17`) | — |
| A3 | **zu.** Eigener Code `E_NOT_SET_UP`; die Anmeldung schickt zur Einrichtung. | PR 209, 220, 234; `shared/wire/wire.go:470-473`; `notsetup_test.go:20`; `apps/web/tests/shell/busy.spec.js:296` | — |
| A4 | **teilweise.** Die Balancer-Vorlage trägt `connectiontimeout=2`, `Retry-After` und `LogLevel proxy:info`; in ulh nur auf `cdu`. Die nackte `500` bleibt mit Absicht. Warum api-2 sechs `500` gab: offen. | `ctl/apache/apache.go:109-140`; `TestTheBalancerAnswersAMissingMemberHonestly` (`apache_test.go:685`); `cdu` `5146658`; Buch :6959 | Push, Release, Anlage (Thomas) |
| A5 | **zu.** Eine Zeile, wenn ein Laden ausfällt, und eine, wenn er wiederkommt. | PR 204; `api/appstore/fetch.go:195,209-227`; `TestJournalSaysAChangeOnce` (`journal_test.go:14`) | — |
| A6 | **zu.** Die Kanalnummer wird frei, sobald beide CLOSE da sind; das nächste OPEN geht nicht mehr verloren. | PR 205, 236; `ctl/link/channel.go:571,713`; `TestOpenWhileOwnCloseIsStillBeingWritten` (`channel_test.go:1424`) | — |
| A7 | **teilweise.** Es war ein halber Vorgang; der Rest ist gebaut. Die Oberfläche fragt die Liste nach „Aufnehmen" nicht neu ab. Nicht nachgemessen. | `api/coreapp/users.go:266-339`; `TestAListedControllerLeavesTheWaitingList` (`users_test.go:1069`); `web/apps/users.js:1671-1676` | Liste neu abfragen: baubar, klein · Web. Messung: Anlage |
| A8 | **offen, ohne Anlage lösbar.** Die „fünfte" Einstellung ist die HSTS-Zeile des ganzen Hosts, dieselbe wie die vierte. ulh zählt Zeilen statt Einstellungen und warnt darum ohne Grund. | `apps/web/templates/apache-vhost.conf:404`, `:640`; `ulh-main:ngc-web.sh:1084-1087` | baubar, klein · ulh-Zweig, docs |
| A8a | **teilweise.** Fehlt ein Controller, steht nach zwei Runden eine Warnung im Journal. Die Oberfläche meldet es nicht, und fehlende API- oder Datenbankknoten prüft niemand. | PR 208; `api/server/missing.go:31-50`; `TestAbsenceSaysOnceGoneAndOnceBack` (`missing_test.go:12`); das Feld `online` gibt es (`shared/core/results.go:142`) | baubar, mittel · Web |
| A8b | **teilweise.** Bei jedem Start meldet CirroDesk eine abweichende Spalte im Journal. Es startet trotzdem, und so bleibt es (entschieden 20.11.). | PR 159; `db/sqlstore/drift.go:25,110` (kam mit PR 198); `TestCheckColumnsFindsTheBentColumn` (`drift_test.go:15`); `docs/schema-pruefung-beim-start.md:330-339` | entschieden 20.11.: Weg A wird nicht gebaut, das Blatt ist nachgezogen; der Fall für `OncePeek` (A19a) bleibt · Datenbank |
| A8c | **teilweise.** Jede API-Anfrage hat jetzt eine Zeile mit dem Knoten. Apaches Protokoll nennt den Arbeiter des Lastverteilers weiter nicht. | PR 204; `api/server/reqlog.go:22-45`; `TestRequestLogNamesNodePathAndStatusOnly` (`reqlog_test.go:18`); `apache-vhost.conf:319`, `BALANCER_WORKER` 0 Treffer | baubar, klein · Web |
| A9 | **zu.** Abgelaufene Lizenz: kein Container, `E_LICENSE`. | PR 200; `ctl/serve/app.go:245-259`; `TestAnAppWithoutALicenceReachesNoContainer` (`reason_test.go:400`) | — |
| A10 | **zu im Code.** Jede App bekommt eine eigene cgroup mit `pids.max`, `memory.max` und `cpu.max`, und der Controller liest jede Zahl zurück. „Delegate oder linger" entfällt. | PR 155; Buch Punkt 101; `ctl/runtime/cgroup.go:199-339`; `cgroup_test.go:249,292` | Anlage (Thomas) |
| A11 | **zu.** Eine Grenze, die nicht greift, fällt auf; ein Container außerhalb seiner cgroup wird angehalten. Rest: beim `create` geht podmans stderr verloren. | PR 155; `cgroup.go:317-339`; `TestEnsureCgroupRefusesWhenTheKernelDoesNotHoldTheLimit` (`cgroup_test.go:292`); `ctl/runtime/podman.go:350` | Rest: baubar, klein · Controller |
| A12 | **zu im Code.** Was eine App in `/tmp` schreibt, zählt gegen ihre Speichergrenze (Vorgabe 256 MiB); gemessen: 400 MiB führen zum Abbruch. | `ctl/runtime/podman.go:270-287`; `cgroup.go:279-286`; Buch :8470-8482 | Anlage (Thomas) |
| A13 | **teilweise.** Der Beitritt prüft die MAC jetzt vor dem Verbrauch. Die Kopplung von Controllern verbraucht weiter zuerst, und den Zähler aus Frage 23 gibt es nirgends. | PR 206, 221; `api/join/join.go:122,163`; `api/controllers/pair.go:89`; `TestPairSpendsTheSecretBeforeTheMAC` (`pair_test.go:123`); TE :277 | baubar: Zähler mittel · shared/, API; danach Kopplung klein · API |
| A14 | **zu.** Eine zweite Anfrage für einen wartenden Namen wird abgewiesen. Rest: der Text sagt noch „Der Beitrittscode ist verbraucht", seit PR 221 falsch. | PR 201; `api/join/join.go:146-157`; `TestASecondRequestForAWaitingNameIsRefused` (`join_test.go:478`); `apps/server-api/main.go:246` | Rest: baubar, klein · API |
| A15 | **offen, aber entschieden.** `noexec` auf `/data` ist seit 24.10. beschlossen (Frage 25, „vor der Veröffentlichung"), nicht gebaut. | TE :305; `ctl/runtime/podman.go:304` (`rw,nosuid,nodev`) | baubar, klein · Controller, docs |
| A16 | **teilweise.** Der Digest ist seit 15.10. festgenagelt, und wer ihn anhebt, ist geregelt. Es fehlt die Prüfung, die meldet, wenn er alt wird (Frage 26). | `apps/server-controller/templates/terminal/Containerfile:30`; `tools/pin-base.sh:11-24`; `pin-base` sonst 0 Treffer in `tools/` und `.github/` | baubar, klein · tools |
| A17 | **offen, aber entschieden.** Rohtasten (Frage 27); gebaut ist weiter das Zeilenfeld. | TE :327; `apps/web/assets/apps/terminal/app.js:118-119,377-403` | baubar, mittel bis groß · Web |
| A18 | **offen, nur an der Anlage.** Der Code schreibt `innodb_encrypt_log=ON`; das `OFF` auf db-3 kam aus einer Handdatei. | `cdu:cd-db.sh:517` (`ulh-main` gleich); `apps/server-db/templates/galera.cnf:62`; Buch Punkt 93 | Neuaufsetzen (Thomas) |
| A19 | **offen.** `wsrep_dirty_reads` ist nirgends gesetzt (entschieden 20.11.: bleibt aus). | 0 Treffer in `apps/` beider Repos; Buch :6954 | eine Zeile `wsrep_dirty_reads=OFF` mit Grund in `apps/server-db/templates/galera.cnf`: klein · Datenbank |
| A19a | **zu.** Der Ablauf eines Einmalwerts wird zusätzlich in Go geprüft. | PR 181; `db/sqlstore/once.go:163-169`; `TestOnceUseRefusesAnExpiredValueOnABentColumn` (`expiredtype_test.go:73`) | — |
| A20 | **zu, auf ulh-`main`.** Die API darf `/etc/ngc` schreiben (seit 30.10., Merge `a61e3fc`); die `/etc/cd`-Form steht auf `cdu`. | `ulh-main:ngc-api.sh:57,1085`, `ngc-all.sh:1443`; `cdu:cd-api.sh:1094` | Push der Umbenennung (Thomas) |
| A21 | **offen.** Kein Skript legt `/var/lib/cd/staged` an. | gelesen in `cdu:cd-api.sh:646-656`; `mkdir` oder `ensure_dir` darauf: 0 Treffer auf beiden ulh-Ständen | baubar, klein · ulh-Zweig |
| A22 | **teilweise.** Sicherung und Rücksicherung gibt es als Befehle; `cd-db` kommt nicht zum Nutzer (S9). In der Oberfläche gibt es beides nicht. | `apps/server-api/main.go:73`; `db/dbcli/dbcli.go:47,112`; „backup" und „sicherung": 0 Treffer in `apps/web/public/`, `apps/web/assets/` | baubar, groß · shared/, API, Controller, Web |
| A23 | **offen.** Das Wiederherstellen legt Benutzer und Datenordner einer App nicht an; die Form ist am 20.11. entschieden: ein eigener Schritt, nicht `node.Restore`. | `ctl/node/restore.go:139,216`; `docs/sicherung-und-weg-zurueck.md:993` | dbKey über die alte `controller.json` vorgeben: Anleitung und ulh-Fall, klein · docs, ulh-Zweig |
| A24 | **teilweise.** Gebaut wie am 15.11. entschieden: ein Kanalfehler schließt nur den Kanal. Ob das den gemessenen Fall trifft (Punkt 77, `ErrLink` im Controller), ist offen. | PR 224; `api/controllers/tunnel.go:203-225`; `TestAContractBreakOnOneChannelLeavesTheOthersAlone` (`tunnel_test.go:249`) | Nachmessung: Anlage (Thomas) |
| A25 | **teilweise.** `cd.store/refresh` gibt es in Vertrag und API, höchstens einmal je Minute. Die Oberfläche ruft es nicht; der Knopf „Jetzt nachsehen" baut nur den Schirm neu. | PR 227; `shared/core/core.go:56`; `TestRefreshNowFetchesAtMostOnceAMinute` (`refresh_test.go:8`); `ACT_REFRESH` nur in `web/contracts/core.js:66,142`; `web/shell/stale.js:51` | baubar, klein · Web |
| A26 | **teilweise.** Der Kopplungscode gilt 24 Stunden, und „Kopplungscode zeigen" gibt jederzeit einen neuen aus. Entschieden ist eine Stunde (Frage 18). | `shared/tunnel/tunnel.go:35` (`PairLifetime = 86400`); `web/apps/users.js:1530`; TE :215 | baubar, klein · shared/ (Koordinator) |
| A27 | **offen.** Die Anleitung nennt vier WSL-Fallen, nicht das gemeinsame Netz. | `docs/installation-kleines-system.md:780-802` | baubar, klein · docs |
| A28 | **offen.** Die Anleitung beginnt weiter mit „zwei Entscheidungen von Thomas" und beschreibt das Web-Paket in alter Form. | `docs/installation-kleines-system.md:57,94-96,281` | baubar, klein · docs (endgültig mit dem cd-Release) |
| A28a | **offen.** Für `/assets/` setzt die Vorlage keinen `Cache-Control`; jede Anmeldung lädt ihr Rechenwerk neu. | `apache-vhost.conf:594` (nur Kommentar), `:636-642`; `web/contracts/proof.js:254` | baubar, mittel · Web |
| A29 | **offen, neu.** Den Vertrauensanker außerhalb der Datenbank legt nur die Kopplung eines Controllers an; die API rückt ihn nur vor, wo die Datei schon liegt. Ein API-Server ohne gekoppelten Controller hat diesen Schutz nicht, auch ein beigetretener nicht: der Beitrittscode trägt den Anker, die API zeigt ihn nur an. Im Labor am 24.10. gemessen: api-1 und api-2 ohne Datei, die Liste auf Fassung 3. | `ctl/node/main.go:287`; `api/state/trust.go:228,239` (`advanceAnchor`), `:287-294` (`ResealAnchor`: „An installation without the file never gets one here"); `api/join/run.go:201`; `docs/proxmox-sechs.md:309-317` | entschieden 20.11.: zuerst ein Rückweg für root (E), dann anlegen bei Einrichtung und Beitritt und vorrücken beim Auffrischen (B–D) in einem Zug; E mittel, B–D klein · API, docs |
| A30 | **offen, neu; die Behebung liegt in PR #257 (offen).** Ein App-Paket lässt sich auf einer eingerichteten Anlage nicht hochladen. Hochladen ist immer „kritisch", der erste Verwalter steht auf Stufe `critical`, aber die Oberfläche schickt den Upload ohne unterschriebene Aufgabe und ohne Bestätigung; die API antwortet `E_CONFIRM`. Die Test-Attrappe stand auf Stufe `off`, darum fiel es in keinem Fall auf. Rest laut #257: `cd.store/api/set-update` geht weiter ohne Aufgabe. | `shared/core/core.go:310-318`; `web/auth/setup.js:181`; `web/apps/store.js:547` (`this.ask(A.upload, null, …)`), `:884-886`; `shared/task/task.go:175,183`; `api/rights/check.go:146,165-166`; `web/shell/fake-link.js:597`; PR #257 | #257 prüfen und zusammenführen (Koordinator) |

### 3. Stört uns — heute

| # | Stand am 20.11.2026 | Beleg | Wer |
|---|---|---|---|
| U1 | **teilweise.** Alle sechs Grenzen der Vertrauensliste haben seit 18.10. einen Fall. Es fehlt ein Fall, der `session.go:481` gegen eine ungeprüfte Liste hält; `vertragsdeckung.md` sagt noch „kein Fall". | `TestTheSizeCapsOfAListAreValidated` (`shared/trust/trust_test.go:866`), `TestTheBudgetAndTheCapsAgree` (`budget_test.go:19`); `shared/session/session.go:481` | baubar, klein · shared/ (nur Fälle) |
| U2 | **zu.** Alle vier Prüffälle sind in Ordnung. Nebenbei: E15 überspringt sich auf GitHub, weil dort kein busybox-Abbild liegt. | `ctl/install/activate_test.go:1193`; `ctl/node/run_test.go:525` (PR 217); `ctl/terminal/terminal_test.go:540`; `ctl/serve/attack_test.go:400-446` | — |
| U3 | **teilweise.** `TestReleaseScript` läuft bei jedem Push mit echtem `jq`. `tools/release-gate.sh` ist nie gelaufen (Schritte 4–6), `govulncheck` läuft in keiner Prüfung. | `gh run view 44636049436` (`b7a21def` grün); `tools/release-gate.sh:11`; `release-gate` und `govulncheck` 0 Treffer in `.github/workflows/` | baubar, klein · tools; gegen ein echtes Release erst mit dem cd-Release |
| U4 | **offen, aber entschieden.** `cd-api migrate` lädt weiter die ganze `api.json`, und niemand startet es außer von Hand. Entschieden ist P7 (a). | `apps/server-api/main.go:71,136-137,165-166`; `LoadMigrate` 0 Treffer; TE :227 | baubar, groß · shared/, API, Controller, ulh |
| U5 | **teilweise.** P10 ist am 24.10. entschieden: ein unbekannter Code wird angenommen. Gebaut ist es nicht, und im Protokoll steht es nicht. | TE :144; `shared/wire/wire.go:532-535` weist ab, gehalten von `wire_test.go:436`; `web/contracts/wire.js:570` | Protokoll: Koordinator. Bau: wenige Zeilen, aber Vertragsänderung · shared/, Web |
| U6 | **offen, aber entschieden.** P8, eigener Schließgrund für den Leerlauf, nach P10. | TE :154; `E_IDLE` 0 Treffer in `shared/`, `apps/` | baubar, groß · shared/, Controller, Web |
| U7 | **offen, aber entschieden.** P9 (`E_BUSY`), nach P8. „Datenbank besetzt" wird weiter ein interner Fehler. | TE :164; `api/httpapi/fault.go:83-84`; `E_BUSY` nur in Kommentaren | baubar, groß · shared/, API, DB, Web |
| U8 | **offen, aber entschieden.** P1 (a). `store.ValidDownloadURL` gibt es nicht; Pfad und Query einer Download-Adresse prüft niemand. | TE :174; `shared/store/catalog.go:208`; `api/appstore/fetch.go:433-446` | baubar, klein im Umfang (Vertragsergänzung) · shared/, API |
| U9 | **offen, aber entschieden.** P11, P12, P5, P14 und P15; keiner gebaut. | TE :185, :195, :205, :251, :262; `shared/dbapi/dbapi.go:83-99` (kein Verweisfeld) | P11, P12, P14: klein · shared/. P5: groß · shared/, DB, Controller. P15: groß, nach dem Schnitt |
| U10 | **teilweise.** Entschieden ist eine eigene Release-Nummer (Frage 31). Gebaut ist eine Nummer, die in alle Programme gestempelt wird; im Vertrag steht es nicht. | TE :374; `tools/build-release.sh:5-6,223`; `shared/store/catalog.go:146` | Satz: klein · shared/. Ob das Gebaute zur Entscheidung passt: Koordinator |
| U11 | **offen.** ulh legt `/var/lib/cd` mit `0755` und `release.seq` mit `0644` an. Entschieden 20.11.: `0711` und `0640`. | `cdu:cd-api.sh:734,736` (ebenso die anderen drei Skripte, `ulh-main` gleich); Buch :7293 | alle vier Skripte gleich ändern: klein · ulh-Zweig |
| U12 | **zu.** Der Kommentar an `--read-only` sagt jetzt, was stimmt. | PR 155; `ctl/runtime/podman.go:270-286` | — |
| U13 | **teilweise.** Die Herkunft ist geklärt: die Warnzeilen von rootless podman. Sie zählen weiter als stderr der App. | `ctl/runtime/cgroup.go:43-51`; `ctl/serve/call.go:207-210` | baubar, klein · Controller |
| U14 | **teilweise.** Die Vorrichtung ist behoben (Punkt 85: danach 30 von 30 grün). Warum der Runner häufiger stückelt, ist ungemessen. | `ctl/terminal/terminal_test.go:683,712` (`7a946603`); Buch :6159 | baubar, klein · docs |
| U15 | **zu.** `STAND.md` ist vom 19.11.2026, mit Nachträgen vom 20.11. | `docs/STAND.md:3`; PR 244, 246, 248 | — |
| U16 | **zu, war kein Mangel.** Die Vertrauensliste kennt nur `api` und `controller`; `web` ist die Rolle eines Controllers. | `shared/trust/trust.go:50-53` | — |
| U17 | **zu.** Einheit und Programm sagen beide `0711`. | `ulh-main:ngc-controller.sh:1244` (seit 25.10.); `cdu:cd-controller.sh:1244`; `ctl/runtime/run.go:65-67` | — |
| U18 | **offen.** Die Oberfläche schickt den Anmeldebeweis weiter, bevor sie prüft, mit wem sie spricht; am 20.11. als gewollt entschieden (siehe U19). | `web/auth/login.js:177,195,201`; `api/login/trust.go:33` | nichts zu bauen |

### Was nur Thomas kann

1. **Release-Schlüssel** (S3, danach S4): Stick oder Datei (TE, Fragen 2 und 3).
2. **Laden** (S1, damit S2): Store-Schlüssel ziehen, Webspace und Adresse festlegen (ulh nimmt
   `https://store.cirrodesk.org` an), Katalog hochladen, Schlüssel in ulh eintragen (TE, Frage 4).
3. **Erstes cd-Release und Push von `cd-umbenennung` nach ulh-`main`**, nur zusammen (`docs/STAND.md`,
   „Auslieferungskette"). Erst damit kommt alles, was oben „zu" heißt, beim Nutzer an.
4. **Anlage neu aufsetzen** mit diesem Release. Damit erledigen sich A18 und S15. Erst danach messbar:
   die Ursache von A1, die `500` von api-2 (A4), die Nachmessungen A7 und A24, der erste Selbstlauf
   (S8) und `cd-db` gegen eine echte MariaDB (S9).
5. **Ein öffentlicher Name** für den einen Lauf mit echtem Zertifikat (S10, TE Frage 32).

**Die sieben Entscheidungen, die nicht auf dem Blatt vom 24.10. standen, sind am 20.11.2026 gefallen**
(Nachtrag): Thomas hat bestätigt, dass seine Vorgabe vom 24.10. auch für sie gilt, und der Koordinator hat
nach der Rangfolge entschieden. `wsrep_dirty_reads` (A19) bleibt aus; die Form des Wegs zurück (A23) ist
festgelegt; bei abweichender Spalte (A8b) startet CirroDesk weiter und meldet es; die Rolle `web` bekommt
einen zweiten Knoten (S11); `/var/lib/cd` wird `0711` (U11); „erst beweisen, dann prüfen, wem" (U18) ist
gewollt; ein API-Server legt seinen Anker selbst an, sobald es einen Rückweg gibt (A29). Das Warum steht
je Zeile im Protokoll des Entwurfs, Abschnitt 6. Bei Thomas bleiben nur 1–5.

**Neu beim Entscheiden gefunden (20.11.2026), noch nicht eingeordnet:**
- **U19: Jede Sitzung kann die Passwort-Prüfwerte aller Konten lesen.** Die Vertrauensliste trägt für
  jedes Passwortkonto den abgeleiteten öffentlichen Schlüssel und die KDF-Parameter
  (`shared/proof/proof.go`, `UserKey`: `pub`, `kdf`), und `api/login/trust.go` gibt die Kette jeder
  Sitzung. Wer angemeldet ist, kann damit das Passwort jedes anderen Kontos offline raten (Argon2id mit
  den echten Parametern). Passkey-Konten betrifft es nicht. Entscheidung fehlt; zuerst abwägen.
- **Vermutet, nicht gemessen:** Der Testknoten `nodeNonPrimary` lehnt schon das `SET` beim Verbinden ab
  (`db/mariadb/fakenode_test.go:51,243`); eine echte MariaDB erlaubt `SET` auf einem Knoten ohne
  Mehrheit. Dann greift das Ausweichen auf den nächsten Knoten (`db/mariadb/nodes.go:107-144`) im Feld
  vermutlich nicht. Messen an der Anlage nach dem Neuaufsetzen.

### Ohne Thomas baubar (Vorschlag)

Nach Teilen. Größe grob: klein = eine Stelle, mittel = mehrere Dateien oder zwei Teile, groß = über
mehrere Teile oder eine Vertragsänderung.

- **Koordinator (`shared/`, Protokoll):** ~~die Entscheidungen vom 24.10. ins Protokoll~~ (am 20.11.
  nachgetragen) · ~~PR #257 prüfen und zusammenführen, A30~~ (am 20.11. zusammengeführt) ·
  P10 bauen, U5 (wenige Zeilen, aber Vertragsänderung; mit Web) · `PairLifetime` auf eine Stunde, A26
  (klein) · Sätze zu P11, P12, P14 und zur Release-Nummer, U9 und U10 (klein) · Fälle für U1 (klein) ·
  `store.ValidDownloadURL`, U8 (klein, mit API) · Zähler für falsche MACs, A13 (mittel, mit API) · Weg
  für `cd-db` auf die Datenbankknoten, S9 (mittel, mit tools, cd-store, ulh) · danach in der Reihenfolge
  von TE: P8 (U6), P9 (U7), P7 (U4), P5 (U9), je groß.
- **Web:** Resttexte übersetzen, S12 · Knopf für `cd.store/refresh`, A25 · Liste nach „Aufnehmen" neu
  abfragen, A7 · Arbeiter des Lastverteilers ins Apache-Protokoll, A8c · Beispiel `calculator`
  entfernen, S2 · ein Fall, der einen Passkey anlegt und sich damit anmeldet, S6 · den Kommentar
  „nothing is ticked" in `users.js` an den Code angleichen (Systemrechte, 20.11.) — alle klein. Meldung
  „mir fehlt ein Teil", A8a (mittel) · Zwischenspeicher für `/assets/`, A28a (mittel) · Terminal mit
  Rohtasten, A17 (mittel bis groß) · Sicherung und Rücksicherung in der Oberfläche, A22 (groß, mit
  `shared/`, API, Controller).
- **API:** eigene Zeile beim Abweisen einer Aktivierung, A1 · Text „Der Beitrittscode ist verbraucht"
  berichtigen, A14 · die Kopplung prüft die MAC vor dem Verbrauch, A13 (nach dem Zähler) — alle klein.
  Vertrauensanker, A29: zuerst der Rückweg für root (E, mittel; vorher `docs/recovery-modus.md`
  lesen), danach anlegen bei Einrichtung und Beitritt und vorrücken beim Auffrischen (B–D, klein) in
  einem Zug.
- **Controller:** `noexec` auf `/data`, A15 · podmans Warnzeilen nicht als stderr der App zählen, U13 ·
  stderr beim `create` behalten, A11 — alle klein.
- **Datenbank:** ein Fall für `OncePeek` mit verbogener Spalte, A19a · `wsrep_dirty_reads=OFF` mit
  Grund in `templates/galera.cnf`, A19 — beide klein.
- **ulh-Zweig `cd-umbenennung`** (den Push macht Thomas): Berichtigung S3a nachziehen ·
  `/var/lib/cd/staged` anlegen, A21 · die HSTS-Zeile nicht doppelt zählen, A8 · `/var/lib/cd` auf
  `0711` und `release.seq` auf `0640`, in allen vier Skripten, U11 · ein Fall für den dbKey aus der
  alten `controller.json`, A23 — alle klein. Unterbefehl `test-key <Datei>`, S5 · ein vorhandenes
  Web-Zertifikat übernehmen, statt eines zu erzeugen, S11 — beide mittel.
- **tools:** Altersprüfung des Digests, A16 · `govulncheck` und `release-gate.sh` Schritte 5–6 gegen
  den Testvektor in die Prüfung, U3 · den veralteten Kommentar in `tools/sicherung.sh:14-15`
  berichtigen, A23 — klein. Zweiter Web-Knoten in `tools/proxmox-sechs*.sh`, S11 — mittel.
- **docs:** WSL-Netz, A27 · Anleitung Abschnitt 1, A28 · Buch zu A8
  und U14 · `vertragsdeckung.md`, U1 · Zeilenangabe in `docs/cd-store-betreiben.md`, S1 ·
  `docs/zweiter-benutzer-schritte.md:21-28` nennt PR 118 als nicht verschmolzen (verschmolzen seit
  24.10., ausgerollt erst mit dem cd-Release) · Anleitung: der dbKey bei einer Neuinstallation über die
  alte `controller.json`, A23 · `docs/proxmox-sechs.md` für den zweiten Web-Knoten, S11 — alle klein.

**So gemessen.** Vier Lesedurchgänge (S, A1–A14, A15–A28a, U) haben jede Zeile am Code nachgesehen und
nichts geändert; über 50 ihrer Belege sind danach einzeln nachgeprüft (`git grep`, `sed -n`, `gh`),
alle stimmten. A29 und A30 hat der Partner „Zweige API" am 20.11. gefunden; beide sind am Code von
`d359b96a` nachgesehen. CirroDesk-Tests wurden nicht gestartet; die Linux-Prüfung auf GitHub ist für `b7a21def`
grün (`gh run view 44636049436`). **Nicht gemessen:** alles an der Anlage (sie fährt noch NGC 0.2.x),
der Inhalt der Release-Archive und ob `store.cirrodesk.org` erreichbar ist.

## Stand am 15.11.2026 (Koordinator)

Die Tabellen unten sind der Stand vom 30.10. Seitdem, alles in `main` oder als PR offen:

| Punkt | Stand |
|---|---|
| S1 | entschieden (A, eigener Webspace); Anleitung PR 214. Offen: `storeUrl`, Schlüssel und ulh-Push (Thomas) |
| S6 | gebaut: Passkey hinzufügen, PR 203 |
| S7 | entschieden (Paket unter `ui/`); Controller-Seite PR 213, Packen PR 211 mit `tools/pack-release.sh`. Offen: ulh `cd-web.sh` (Thomas) |
| S8 | entschieden; Controller fragt das Programm nach seiner Fassung (PR 210); cd-web über `ui/VERSION`. `core.Unattended` bleibt |
| S9 | gebaut: `cd-db backup/restore`, PR 198. Offen: Lauf gegen eine echte MariaDB, `cd-db` in den Bau-Skripten |
| S12 | entschieden: nach dem WhiteSur-Look |
| S14 | war am 31.10. schon behoben (589324b8) |
| A1 | Controller-Hälfte PR 216 |
| A2 | entschieden (Dienst startet ohne Datenbank nicht), Bau beim Controller |
| A3 | entschieden (`E_NOT_SET_UP`); Journal PR 209, Vertrag und Antwort PR 220 |
| A5, A8c | PR 204 |
| A8a | PR 208 |
| A9 | PR 200 |
| A11 | war mit #155 gebaut |
| A13 | `OncePeek` PR 206, Aufruf im Beitritt PR 221 |
| A14 | PR 201 |
| A19a | mit den Riegeln in Go erledigt |
| A24 | entschieden (nur der Kanal); README stimmte schon, die API baut um |
| A25 | `cd.store/refresh`, Zweig `shared-store-refresh` |
| E32 | PR 218 |
| U2 | PR 217 |

Entscheidungen dazu: Entscheidungsprotokoll im Entwurf, Einträge vom 15.11.2026.

**Die Frage:** CirroDesk läuft. Acht Maschinen, zwei Apps, drei Datenbankknoten, Paketfilter, 95
Punkte im Entscheidungsbuch, drei veröffentlichte Fassungen. Trotzdem fehlt eine Antwort auf:
**was hält einen Fremden auf?**

**Woher das kommt.** Gesammelt wurden vier Quellen, dann gezählt:

| Quelle | Stellen |
|---|---:|
| `docs/thomas-entscheidet.md` — 32 Fragen + 3 Messungen | **35** |
| `docs/bis-zur-veroeffentlichung.md` — H1–H10 + H4a, R1–R17, B1–B8 | **36** |
| `docs/messungen-m1-m3.md` — offene Lücken | **3** |
| `docs/vertraege-v1.2.md`, Punkte 63–95 — „zu bauen" / „offen" / „nicht gemessen" | **24** |
| Meldungen anderer Stränge, alle am 30.10.2026 gemessen (A8a–A8c, A19a, A28a, S14, S15, U16–U18) | **10** |
| **gefunden, mit Doppelungen** | **108** |

Nach Abzug der Doppelungen (dieselbe Sache steht oft an drei Stellen) bleiben **64 eigene
Einträge**: **14 sperren**, **32 halten auf**, **18 stören nur uns**.

**Was sich seit dem 18.10.2026 geschlossen hat** — am 30.10.2026 selbst nachgesehen, nicht
übernommen:

| Alt | Heute |
|---|---|
| H1 Repo ist privat | **zu.** `gh repo view` sagt `PUBLIC`. |
| H2 kein Release | **zu.** Drei: 0.2.0, 0.2.1, 0.2.2 (`gh release list`). |
| H3 kein `cd-*`-Skript auf ulh-`main` | **zu.** Alle fünf liegen dort, und `catalog.yaml` kennt CirroDesk. |
| H7 `Updater: nil` | **zu.** `run.go:667` verdrahtet den Updater. |
| H8 eine Installation hat einen Benutzer | **zu.** `users.js:783` schreibt in `next.users` (Punkt 73). |
| H9 vier Behebungen auf unverbundenen Zweigen | **zu** (Punkt 63 ff.). |

Das ist der halbe Weg. Was bleibt, steht unten.

---

## 1. Sperrt — ohne das kann ein Fremder CirroDesk nicht benutzen

| # | Worum | Wo gemessen | Was es braucht |
|---|---|---|---|
| **S1** | **Es gibt keinen Laden.** `CD_STORE_KEY_DEFAULT="REPLACE-THIS-LINE-WITH-THE-CD-STORE-KEY"` steht am 30.10.2026 unverändert auf ulh-`main` (`cd-api.sh:97`). Jede Installation muss ihren eigenen Store nennen — es gibt keine Apps von CirroDesk. | Punkt 56, 57 · Frage 4 | **Entscheidung** (Thomas), danach bauen |
| **S2** | **Zwei Apps, und beide sind Beispiele.** Im Repo liegen `calculator` und `terminal`. Auf der Anlage läuft als zweite eine **Uhr**. CirroDesk verspricht einen Schreibtisch mit Apps in Containern. | Punkt 95 · B3 | bauen |
| **S3** | **Der Release-Schlüssel ist eine Datei auf einer Testmaschine.** Der öffentliche Schlüssel in den vier Skripten ist der mit der Kennung `pN6Z6a9NR0aG6aEJIJ5EHg` — der dateibasierte, nicht der vom Stick. Drei veröffentlichte Fassungen sind damit unterschrieben. Wer auf dieser einen Maschine root ist, unterschreibt CirroDesk für alle. | Punkt 57 · H4 · Frage 2/3 | **Entscheidung** (Thomas) |
| **S3a** | **Der Kommentar darüber sagt das Gegenteil.** „The placeholder below is not a key, so every install stops at the check" — darunter steht ein echter, gültiger Schlüssel. Es stoppt nichts, und das Repo ist jetzt öffentlich. | H4a, am 30.10. auf ulh-`main` nachgesehen | drei Zeilen |
| **S4** | **Punkt 19 ist nicht gebaut.** `shared/sign/sign.go:112` unterschreibt die rohe Nachricht, nicht einen Digest. Darum **kann** der Stick nicht unterschreiben, selbst wenn Thomas sich für ihn entscheidet. Änderung in `shared/` — nur im v1.3-Schnitt. | Punkt 19 (gilt) | bauen, nach P10 |
| **S5** | **Ein Fremder kann die Kette nicht mit eigenem Schlüssel prüfen.** `CD_RELEASE_KEY="…"` ist eine nackte Zuweisung, anders als alles daneben. Nur mit `sed` in vier Skripten zu ändern — und wer die Skripte verbiegt, prüft nicht mehr die Skripte. | Punkt 62, Befund 2 · Frage 6 | bauen |
| **S6** | **„Mit Passkey anmelden" kann nicht gehen.** `credentials.create` hat auf `main` **null** Treffer in `apps/web/public/` (am 30.10. gezählt). Die Oberfläche empfiehlt Passkeys selbst. 22 Vertragsposten leben nur im Prüfstand. | R2 | bauen |
| **S7** | **Ein `cd-web`-Update kann nicht gehen.** Die Skripte bauen Namen im Wurzelverzeichnis, der Controller packt nur `ui/` aus, `store.WalkUI` weist das Wurzel-tar ab. Ein Fremder installiert und aktualisiert die Oberfläche nie. | P3 · Punkt 56 (als falsch gemessen) · Frage 7 | **Entscheidung**, danach bauen |
| **S8** | **Die Selbstaktualisierung endet am ersten Knoten.** Vier Läufe in Folge `cd-api 0.2.2 db-1 E_INTERNAL`. Ursache: ein unbeaufsichtigter Lauf darf kein CirroDesk-Teil einsetzen — die stehende Zustimmung ist gebaut und wirkungslos. | Punkt 82 | **Entscheidung** (a oder b), danach bauen |
| **S9** | **Benutzerdaten werden nicht gesichert.** Zwei Stellen im Code schicken den Admin an `cd-db backup`. **Das Programm existiert nicht** — `apps/server-db/` hat kein `main.go`, auf keiner Maschine liegt ein `cd-db`. Gesichert wird nur, was die Anlage ohnehin wiederbekäme. | Punkt 81 | bauen |
| **S10** | **Für `.lan` stellt keine Zertifizierungsstelle aus.** Alles seit Punkt 58 ist mit selbst unterschriebenem Zertifikat gemessen, in einem Browser, der es für den Lauf annimmt. Die Setup-Seite weigert sich über `http://`. Was ein Fremder mit `letsencrypt` erlebt, weiß niemand. | Frage 32 · R10 | **Messung** auf einer Maschine mit echtem Namen |
| **S11** | **Die Oberfläche liegt auf einer Maschine.** Die Rolle `web` hat nur VM 301 (`cd-web-api-1`). Fällt sie, ist CirroDesk weg — auch wenn die anderen fünf laufen. | `docs/proxmox-sechs.md:33` | **Entscheidung** + bauen |
| **S12** | **Die Oberfläche ist deutsch.** Die fünf Sprachdateien tragen je rund 21 Schlüssel, alle aus Anmeldung und Einrichtung; rund 233 Zeichenketten stehen als deutscher Text im Code. `lang.js` nennt fünf Sprachen und `FALLBACK = "en"` — ein englischer Nutzer bekommt trotzdem Deutsch. Die README ist englisch. | B1, Zählung am 30.10. bestätigt | bauen |
| **S13** | **Eine Lizenz mit Ausstellungszeit in der Zukunft sperrt jede spätere aus** — dauerhaft, und der Verwalter liest „gültig". Die Zeile liegt in eingefrorenem `shared/store/catalog.go:321`. | R1 (S1) · Punkt 63 | **Entscheidung** P10 zuerst |
| **S14** | **Die erste Anmeldung sagt „Anmeldung fehlgeschlagen" — bei jedem.** Die Oberfläche baut ihre Vertrauensmaske aus der **ältesten** Fassung der Liste; dort steht nur api-1. Antwortet api-2, findet `verifyAPI` keinen Aussteller, wirft `session: malformed certificate` ohne Code, und die Oberfläche zeichnet „Passwort falsch". Beleg: acht Fassungen aus `cd_blob`, Lastverteiler-Parität 3 + 5 = 8. | am 30.10.2026 von einem anderen Strang gemessen | bauen (ein Strang entwirft die Behebung) |
| **S15** | **`cd-ein` hat keine Vertrauensliste** und ist damit unbenutzbar: 302 Klagen „state: no trust list yet" in sieben Tagen, `login/finish` antwortet dort jedem `E_UNREACHABLE`. Auf dieser Maschine kann sich niemand anmelden, nie. | am 30.10.2026 von einem anderen Strang gemessen | bauen (niemand sitzt daran) |

**Das sind 14 Sachen** (S3a gehört zu S3).

---

## 2. Hält auf — es geht, kostet ihn aber eine Stunde oder einen Befund, den wir kennen

### Er sieht nicht, was los ist

| # | Worum | Wo | Was es braucht |
|---|---|---|---|
| **A1** | `E_VERIFY` beim Aktivieren ist nicht erklärt. Weder API noch Controller schreibt eine Zeile über die Abweisung. | R5 | messen |
| **A2** | Ist die Datenbank des Controllers weg, kommt der Dienst `active` hoch und hat einfach keine Apps. Bewusst so gelassen. | R14 | bauen |
| **A3** | Eine migrierte, nicht eingerichtete Anlage sagt bei **jeder** Anmeldung „die Datenbank konnte nicht erreicht werden" — über eine Datenbank, die einwandfrei antwortet. 12 288 solcher Zeilen in 3,4 s gemessen. | B7, hängt an P9 | **Entscheidung** |
| **A4** | Ein fehlendes Mitglied des Lastverteilers wird dem Benutzer als `HTTP 500` zugestellt, und nichts schreibt es auf. Warum api-2 sechs `500` gab und api-1 keine: offen. | Punkt 90, 91, 92 | messen |
| **A5** | Die Anlage fragt stündlich nach ihrem Laden, bekommt keine Antwort und sagt es niemandem. | Punkt 79 | bauen |
| **A6** | Eine Kanalnummer gilt als belegt, nachdem die API die Freigabe schon gelesen hat — das nächste `OPEN` darauf verschwindet lautlos. | Punkt 86 | bauen |
| **A7** | Nach „Aufnehmen" zeigt die Liste den Server weiter, der Anker bleibt bei `seq: 1`. Ob ausbleibende Auffrischung oder halber Vorgang: **nicht gemessen**. | B6 · Punkt 62 | messen, dann bauen |
| **A8** | Die fünfte Schutzeinstellung der Oberfläche ist nicht gemessen — die Oberfläche sagt über sich selbst, dass sie eine nicht kennt. | B8 · Punkt 56 | messen |
| **A8a** | **Ein Teil kann wegbrechen, ohne dass es jemand sieht.** Controller auf svc-1 aus: alle App-Container weg — die Anlage meldete trotzdem 40 von 40 guten Anmeldungen. CirroDesk hat keine Stelle, die sagt „mir fehlt ein Teil". | am 30.10.2026 von einem anderen Strang gemessen | bauen |
| **A8b** | **CirroDesk startet ohne eine Klage, wenn ein Spaltentyp vom Erwarteten abweicht** — und liest danach eine falsche Reihenfolge. Gemessen an `cd_blob`, auf SQLite und auf MariaDB. | am 30.10.2026 von einem anderen Strang gemessen | bauen |
| **A8c** | **Kein Protokoll je Anfrage.** Weder api-1 noch api-2 schreibt eine Zeile für eine einzelne Anfrage (33 bzw. 75 Journalzeilen, alles Start, Tunnel, Runden), und das Protokollformat des Stellvertreters trägt den Namen des Lastverteiler-Arbeiters nicht. Bei einem Fehler über zwei Knoten ist **nicht feststellbar, welcher geantwortet hat**. Hat heute eine Ursachensuche fast zum Stillstand gebracht. | am 30.10.2026 von einem anderen Strang gemessen | bauen (niemand sitzt daran) |

### Riegel, die gebaut sind und nicht greifen

| # | Worum | Wo | Was es braucht |
|---|---|---|---|
| **A9** | `E_LICENSE` verlässt das Controller-Paket nicht. Ein App-Kanal läuft auch für eine App mit abgelaufener Lizenz. | R6 | bauen |
| **A10** | `--pids-limit 512` und `--memory 256m` liegen **nicht an** — gemessen im Verhalten: 600 Aufgaben über einer Decke von 512, keine schlug fehl. Ursache: kein App-Container hat eine eigene cgroup (`Delegate=no` auf allen drei gemessenen Knoten). | Punkt 72 (Posten 1) · Punkt 95 | **Entscheidung** (Delegate oder linger), dann bauen |
| **A11** | Der Controller bemerkt nicht, dass podman einen Schalter verwirft. Ein still verworfener Schalter ist schlimmer als ein fehlender: der Prüffall bleibt grün. | Punkt 72 (Posten 2) | bauen |
| **A12** | `--read-only` lässt vier beschreibbare Ordner offen (`/tmp`, `/run`, `/var/tmp` je 3,9 GiB, `/dev/shm`). Zusammen mit `memory.max = max` ist das ein Weg, den Speicher des Wirts zu füllen, ohne eine Grenze zu verletzen, die irgendwo steht. | Punkt 95 | **Entscheidung** (cgroup oder Größe der tmpfs) |
| **A13** | Der Einmal-Code wird vor der MAC-Prüfung verbraucht: wer die Adresse kennt, verbrennt jeden wartenden Beitrittscode. **Blockiert** — keine Methode in `shared/dbapi` liest einen Einmal-Wert, ohne ihn zu verbrauchen. | Punkt 26 · Frage 23 | **Entscheidung** P10, dann bauen |
| **A14** | Eine zweite Beitrittsanfrage für denselben Namen überschreibt die erste, und der Admin erfährt nichts davon. | Punkt 27 · Frage 24 | bauen |
| **A15** | `noexec` auf `/data` fehlt: eine App darf dort abgelegte Skripte ausführen. Verhaltensänderung für jede Anlage, die das heute tut. | Punkt 51 · Frage 25 | **Entscheidung** |
| **A16** | Das Terminal-Abbild hängt an `debian:13-slim`, nicht an einem Digest — wer den Tag beherrscht, beherrscht einen Container mit `needs:["root"]`. Und: **wer hebt den Digest an?** | Punkt 39 · Frage 26 | **Entscheidung** |
| **A17** | Terminal: Zeilenfeld statt Rohtasten — ein Admin tippt irgendwann sein `sudo`-Passwort in ein Feld, das es anzeigt. | Punkt 41 · Frage 27 | bauen |
| **A18** | Auf db-3 steht die Protokollverschlüsselung weiter auf `OFF`. | Punkt 93 | bauen |
| **A19** | Bei verlorener Mehrheit ist auch das **Lesen** zu (`ERROR 1047` beim `SELECT`). Ob `wsrep_dirty_reads` an soll: offen — veraltete Daten gegen gar keine. | Punkt 92 | **Entscheidung** |
| **A19a** | **Ein abgelaufener Einmalwert wird auf SQLite angenommen** und gibt sein Geheimnis heraus, wenn der Spaltentyp abweicht — der Riegel steht in SQL und kippt dort um. Gehört zu A8b und zu A13. | am 30.10.2026 von einem anderen Strang gemessen | bauen (ein Strang baut zwei Riegel in Go) |

### Die Anlage von Hand nachbessern

| # | Worum | Wo | Was es braucht |
|---|---|---|---|
| **A20** | `ReadWritePaths` nennt `/etc/cd` nicht (`cd-api.sh:1065` auf ulh-`main`, am 30.10. nachgesehen). Die API kann den Vertrauensanker nicht hochziehen — die Vertrauensliste kann nicht wachsen. | R9 · Punkt 62 · Frage 5 | zwei Zeilen; Push nur durch Thomas |
| **A21** | `/var/lib/cd/staged` legt kein Skript an. Wer ohne Release-Adresse installiert, legt den Ordner selbst. | Punkt 62, Befund 3 | eine Zeile |
| **A22** | Eine Sicherung geht nur über die Befehlszeile, eine Rücksicherung gar nicht. „backup" kommt in `apps/web/public/` **nirgends** vor (am 30.10. gegengeprüft). Das ist das Erste, was ein DSM-Nutzer sucht. | B2 | bauen |
| **A23** | Der Weg zurück geht — er bleibt an einem Systembenutzer hängen, den keine Sicherung kennt. | Punkt 87 | bauen |
| **A24** | Eine fehlgeschlagene Einrichtung reißt eine fremde, laufende App aus der Leitung. Die Trennung, die ein Container leistet, endet am gemeinsamen Tunnel. | Punkt 77 | bauen |
| **A25** | Eine frische Anlage hat eine Stunde lang keinen Laden; der Knopf „Jetzt nachsehen" fehlt. | Punkt 66 · Frage 29 | bauen |
| **A26** | Der Kopplungscode gilt unbefristet, und es gibt keinen sichtbaren Weg, einen neuen auszugeben. | Frage 18 (P6) | bauen |
| **A27** | Zwei CirroDesk-Anlagen können auf einem WSL-Rechner nicht nebeneinander laufen (alle Instanzen teilen ein Netz). Fehlt nur in der Anleitung — und genau dort probiert ein Fremder CirroDesk zuerst. | R8 · Punkt 62 | ein Absatz |
| **A28** | `docs/installation-kleines-system.md` nennt als Erstes „zwei Entscheidungen von Thomas" samt Selberbauen und Selberunterschreiben. Mit drei veröffentlichten Fassungen ist die Anleitung an dieser Stelle überholt. | eigener Vergleich 30.10. | schreiben |
| **A28a** | **Kein Zwischenspeicher für die Rechendateien.** Thomas' Browser holte in einer Woche `argon2.js` 55-mal und `_md.js`/`_u64.js`/`utils.js` je 83-mal: die Anmeldung lädt ihr Rechenwerk bei jedem Versuch neu. Kostet Zeit und Datenmenge — und deshalb misst die Frist auf die Herleitung heute Laden **plus** Rechnen. | am 30.10.2026 von einem anderen Strang gemessen | bauen (niemand sitzt daran) |

**Das sind 32 Sachen.**

---

## 3. Stört uns — technische Schuld, die ein Fremder nie sieht

| # | Worum | Wo |
|---|---|---|
| **U1** | **Elf von siebzehn ungemessenen Vertragsposten sind Sicherheitsriegel.** Das Schwerste ist ein einziges `if` (`trust.go:489`): `MaxNodes`, `MaxStores`, `MaxAuto` haben **keinen** Fall, ebenso `MaxKeysPerUser`, `MaxRoot`, `MaxEntries`. `session.go:455` indiziert ein Feld fester Länge aus einer Zahl, die aus der verteilten Vertrauensliste kommt. | R3 · `vertragsdeckung.md` |
| **U2** | Vier Prüffälle in Ordnung bringen: einer lässt einen Systembenutzer auf der Maschine zurück, einer startet den echten Controller und hängt 10 Minuten, E15 maß das **falsche** Dateisystem. | R7 |
| **U3** | `tools/release-gate.sh` Schritte 4–6 sind nie ausgeführt; `govulncheck` nie auf echtem Debian; `TestReleaseScript` nie mit echtem `jq`. Das Tor vor dem ersten Release ist selbst ungemessen. | R12 |
| **U4** | Der Weg der Migrationen ist nie gelaufen. Dazu verlangt `cd-api migrate` eine ganze `api.json`, die ein reiner db-Knoten nicht hat. | R13 · P7 · Frage 19 |
| **U5** | **P10 ist das Tor:** darf ein Vertrag von `shared/` überhaupt etwas dazubekommen? Ohne P10 bleibt jede Ergänzung für immer ein Bruch — S4, S13, U6, U7 und A13 hängen daran. | Frage 11 · Punkt 54 |
| **U6** | P8 — eigener Schließgrund für die Leerlaufgrenze. Ein Leerlauf bleibt von einem Absturz nicht zu unterscheiden. | Frage 12 · Punkt 35 |
| **U7** | P9 — eigener Code `E_BUSY`. Heute liest der Benutzer bei drei verschiedenen Lagen immer „Dienst nicht erreichbar". | Frage 13 |
| **U8** | P1 — `store.ValidDownloadURL` wird nicht ausgeführt; Pfad und Query einer Download-Adresse prüft dauerhaft niemand. | Frage 14 · Punkt 25 |
| **U9** | P11, P12, P5, P14, P15 — vier Sätze in den Vertrag und ein Verweisfeld in `dbapi.App` (sonst bleibt ein Abbild ohne Verweis für immer liegen und zählt als Benutzerdaten). | Fragen 15–17, 21, 22 |
| **U10** | Welche der fünf Nummern ist die Fassung eines Releases? Heute gibt es fünf Programmnummern und keine Release-Nummer. | Frage 31 (D4) |
| **U11** | `/var/lib/cd` ist `0755` und von jedem App-Benutzer auflistbar, `release.seq` ist `0644`. Kein Loch, aber die App erfährt den Aufbau des Servers und die laufende Release-Nummer. | Punkt 95 |
| **U12** | Der Kommentar an `--read-only` („the only place it may write") ist falsch. Genau diese Sorte Satz wird sonst als Zusicherung weitergelesen. | Punkt 95 |
| **U13** | Nicht geklärt: die 477 Byte stderr, die das Protokoll einer App zuschreibt, die sie nicht geschrieben hat. | Punkt 78 |
| **U14** | Nicht geklärt: warum ein Prüffall auf einem GitHub-Runner häufiger flattert als auf dem eigenen Rechner. | Punkt 85 |
| **U15** | `docs/STAND.md` ist vom **10.10.2026** und beschreibt eine Lage, die es nicht mehr gibt („nichts davon ist je gelaufen"). Wer es zuerst liest, liest 20 Tage Rückstand. | eigener Vergleich 30.10. |
| **U16** | `web-api-1` steht in der Vertrauensliste als `controller` statt als `web`. | am 30.10.2026 von einem anderen Strang gemessen |
| **U17** | Die Einheit des Controllers verlangt `RuntimeDirectoryMode=0700`, danach steht der Ordner auf 711 — erklärt (der Controller setzt ihn selbst), heißt aber: **aus einem Dateimodus lässt sich nicht auf die aufgespielte Fassung schließen**. | am 30.10.2026 von einem anderen Strang gemessen |
| **U18** | **Die Oberfläche prüft erst, mit wem sie spricht, nachdem sie den Anmeldebeweis geschickt hat.** Der Kettenpfad für die Vertrauensliste verlangt eine Sitzung (`login/trust.go:34` über `h.Live`), also läuft `verifyAPI` erst nach `loginFinish`. **Kein Passwort geht hinaus** (berichtigt am 20.11.2026; hier stand „Das ist kein Leck", hinaus geht aber ein Prüfwert zum Offline-Raten, siehe U19 im Stand vom 20.11.): das Passwort verlässt den Browser nie — `contracts/proof.js` leitet im Worker einen Schlüssel ab und schickt eine Unterschrift über die Bytes der Client-Bescheinigung, an Nonce und Zweck gebunden; ein abgefangener Beweis gilt für diese eine Bescheinigung und diesen einen Nonce. Die Argon2-Parameter der ungeprüften Gegenstelle hält `proof.validateKDF` gegen den Vertrag, bevor abgeleitet wird. Bleibt eine **Frage an den Entwurf**: ist „erst beweisen, dann prüfen, wem" so gewollt? | am 30.10.2026 von einem anderen Strang gemessen |

**Das sind 18 Sachen.**

---

## 4. Nie gemessen — die blinden Flecken

Zehn Tage, dieselben Leute, dieselbe Hardware. Das Folgende ist **nicht** offen im Sinne von
„noch nicht gebaut" — es ist nie angesehen worden. Keine Zahl, keine Vermutung.

- **Ein Fremder.** Niemand außer Thomas und den Sitzungen hat CirroDesk je installiert. Jede
  Anleitung ist von denen geprüft, die wussten, was gemeint war.
- **Mehrere Benutzer gleichzeitig.** Ein zweiter Mensch ist seit Punkt 73 anlegbar. Dass ein an
  api-1 angelegter Benutzer sich an api-2 anmelden kann, steht dort ausdrücklich als **nicht
  gemessen**. Zwei Menschen gleichzeitig am Schreibtisch: nie.
- **Datenmengen.** Die echte Anlage hat **einen** Benutzer mit **einem** Schlüssel, **null**
  Stores und zwei Apps. Die Vertragsgrenzen liegen drei Größenordnungen darüber (U1). Eine
  große Datei, ein voller App-Ordner, ein Katalog mit hundert Einträgen: nie.
- **Lange Laufzeit.** Die verteilte Anlage läuft seit dem 23.10.2026. Kein Dauerlauf, keine
  Messung des Speicherverbrauchs über Tage, kein Wachstum von Protokoll oder Platte über Zeit.
- **Ein Neustart der ganzen Anlage.** Gemessen sind Teilausfälle: eine Maschine aus, ein Dienst
  gestoppt, zwei von drei Datenbankknoten weg. **Alle sechs gleichzeitig kalt hochfahren** —
  die Lage, in der Galera über `gvwstate.dat` entscheidet, ob es von selbst wiederkommt — ist
  nie gefahren.
- **Eine volle Platte.** `ErrDisk` ist gebaut und im Ernstfall nie gelaufen. Was ein Verwalter
  sieht, wenn `/var/lib/cd` voll ist, weiß niemand.
- **Ein Upgrade von Fassung zu Fassung.** Es gibt 0.2.0, 0.2.1 und 0.2.2. Dass eine laufende
  Anlage von einer auf die nächste kommt, ist nie durchgelaufen — die Runde endet am ersten
  Knoten (S8).
- **Eine zweite echte App im Repo.** Auf der Anlage läuft eine Uhr. Was eine *andere* App
  braucht, ist nicht gemessen, weil es keine gibt; `cd/portversuch` startet `sleep infinity`.
- **MAC-Fehlversuche.** CirroDesk zeichnet sie **nirgends** auf. „Null Fehlversuche" heißt „nichts
  gesehen", nicht „nichts passiert". Die Schwelle 3 in Frage 23 steht darum unter Vorbehalt.
- **Ein Knoten mit `Delegate=yes`.** Ob `PIDsLimit = 512` im Betrieb trägt, ist erst dort
  messbar. Folge laut Punkt 72: 512 reißen bei rund 53 Terminal-Sitzungen — unter den 64, die
  ein Browser öffnen darf.
- **Andere Hardware, andere Distribution.** Alles läuft auf Debian 13 auf einem Proxmox-Wirt.
  Ein anderer Wirt, ein Blech ohne Virtualisierung, ein anderes podman mit anderer
  `containers.conf` (dort ist der Prozessdeckel eine andere Zahl): nie.
- **Der Weg über ein echtes, öffentliches Zertifikat.** Siehe S10. Es ist der erste Schirm,
  den ein Fremder sieht.
- **Die Einzel-VM 208 ist absichtlich abgeräumt.** Die Messungen der Punkte 72 und 95, die auf
  ihr entstanden, sind auf dieser Maschine **nicht wiederholbar**. Wer sie nachfahren will,
  braucht eine neue Maschine.

---

## 5. Was nur Thomas entscheiden kann

Hier steht nichts Neues — nur, welche der Zeilen oben keinem Agenten gehören.

- **Der Release-Schlüssel** (S3): Stick oder Datei. Wird es der Stick, kommt S4 dazu, und das
  ist eine Änderung in `shared/`.
- **Der Laden** (S1): eigener Store-Schlüssel oder ein CirroDesk-Store, und wer ihn betreibt.
- **Ein Push nach ulh-`main`** (A20): erreicht per Auto-Update sofort alle ulh-Nutzer.
- **P10** (U5): das Tor vor jeder Vertragsergänzung — und damit vor S4, S13, U6, U7 und A13.
- **Verhaltensänderungen für fremde Anlagen:** `noexec` auf `/data` (A15), Digest statt Tag und
  wer ihn anhebt (A16), `wsrep_dirty_reads` (A19), die stehende Zustimmung (S8).
- **Ein öffentlicher Name** für die Messung aus S10.

---

## 6. Vermutet, nicht geprüft

- **Die Gruppen sind ein Urteil**, keine Messung. „Sperrt" heißt: ich sehe keinen Weg, wie ein
  Fremder CirroDesk damit sinnvoll benutzt. Wer anders urteilt, verschiebt eine Zeile.
- **Keine Zahl in Stunden oder Tagen.** Die Spalte „Was es braucht" sagt nur: bauen, messen
  oder entscheiden. Schätzungen wären erfunden.
- **Zehn Einträge habe ich nicht selbst gemessen** (A8a, A8b, A8c, A19a, A28a, S14, S15,
  U16, U17, U18). Sie kommen als Meldungen anderer Stränge vom 30.10.2026, übernommen und
  nicht nachgeprüft.
- **Eine „Container-Verdoppelung" wurde gemeldet, ohne Einzelheiten.** Sie steht hier
  darum nicht als eigener Eintrag.
- **Die 108 enthalten Doppelungen**, und die Zahl der eigenen Einträge (64) ist von Hand
  zusammengeführt. Wer anders schneidet, kommt auf eine andere Zahl.
- **Die Übernahmen aus `bis-zur-veroeffentlichung.md` sind vom 18.10.2026.** Was ich am
  30.10. selbst nachgesehen habe, steht in der Quellenspalte mit Datum; der Rest ist
  übernommen und kann seitdem gefallen sein — so wie H1, H2, H3, H7 und H8 oben.
- **Die Zählung der Sprachschlüssel** ist am 30.10. über die fünf Dateien gezählt; die 233
  deutschen Zeichenketten sind eine Musterzählung vom 18.10., kein Abgleich.
- **Den Code hinter den meisten Einträgen der Gruppe 3 habe ich nicht angesehen.** Sie stehen
  dort, weil ihre Quelle sie so hinterlässt.
