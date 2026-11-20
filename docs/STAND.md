# Stand: CirroDesk

**Zuletzt aktualisiert:** 19.11.2026 · gemessen an `origin/main` **`b2262020`** (Merge von PR #242,
1612 Commits)

> **Die Fassung davor beschrieb den 10.10.2026.** Sie nannte die Pull Requests #26–#30 als offen und
> api-Aufgabe 24 („`main.go` verdrahten") als nächstes Ziel. Beides ist erledigt: #26–#30 sind
> zwischen dem 12. und 18.10.2026 zusammengeführt, `apps/server-api/main.go` gibt es (316 Zeilen),
> und drei Fassungen sind veröffentlicht. Ein Versuch vom 10.11.2026, diese Datei nachzuziehen
> (Zweig `stand-aktuell`), wurde nie zusammengeführt; was davon noch gilt, steht jetzt hier.
> **Bei jeder Zahl steht, woran sie gemessen ist.**

## Wo wir stehen

CirroDesk ist ein Arbeitsplatz im Browser mit App-Store. Bis zum 15.11.2026 hieß es **NGC**
(NextGenCloud). Der Code ist öffentlich, drei Fassungen sind veröffentlicht. **Für Fremde benutzbar
ist es noch nicht** — was fehlt, steht in [`was-fehlt.md`](was-fehlt.md).

| Was | Stand am 19.11.2026 | Gemessen mit |
|---|---|---|
| `main` | `b2262020`, 1612 Commits | `git rev-list --count origin/main` |
| Offene Pull Requests | außer #244, der diese Fassung bringt, **einer:** #243 `labor-cd` (Labornamen auf `cd`, Reste der ulh-Umbenennung) | `gh pr list -R sorglos-it/cirrodesk` |
| Zuletzt zusammengeführt | #242 am 19.11., davor #203–#241 am 14.–16.11. | `gh pr list --state merged` |
| Prüfung auf GitHub | `b2262020` grün. Go: 3765 Fälle bestanden, 0 rot, 27 übersprungen. Web: chromium 612, firefox 609, webkit 604 bestanden, je 0 rot, 0 flatternd | `gh run view 44573010049 --log` |
| Veröffentlicht | 0.2.0 (19.10.), 0.2.1 und **0.2.2** (beide 23.10.) | `gh release list` |
| Entscheidungsbuch | bis **Punkt 102** | [`vertraege-v1.2.md`](vertraege-v1.2.md) |
| Playwright (Web-Prüfsatz) | **1.65.0** seit #192 (09.11.2026). Der Sprung von 1.55.1 machte aus rund 290 übersprungenen Fällen je Engine gemessene. Für Firefox belegt, warum: Firefox 141 (zu 1.55.1) warf bei `exportKey("jwk")` eines Ed25519-Schlüssels `OperationError`, Firefox 159 (zu 1.65.0) nicht | `apps/web/tests/package-lock.json`; Nachricht von `681d4880` (#192) |

### Die Teile

| Teil | Ordner | Fassung |
|---|---|---|
| Webinterface | `apps/web/` | 0.2.1 |
| API | `apps/server-api/` | 0.2.0 |
| Controller | `apps/server-controller/` | 1.0.0 |
| Datenbank | `apps/server-db/` | 0.1.0 |
| Werkzeug `cd-store` | `apps/cli/` | 0.2.0 |
| App „Uhr" (`cd-uhr`) | `apps/uhr/` | 1.0.1 |
| Store-Apps Rechner, Uhr, Dateien | `apps/store-calculator/`, `apps/store-clock/`, `apps/store-files/` | je 1.0.0 |

Fassung aus `VERSION`, bei den Store-Apps aus `manifest.json`. **Der Recovery-Modus liegt noch nicht
auf `main`:** `apps/cli-recovery/`, `apps/cli-stick/`, `apps/server-recovery/` und `lib/rust/` aus
[`recovery-modus.md`](recovery-modus.md) gibt es dort nicht.

**Der Updater ist verdrahtet:** `apps/server-controller/classes/node/run.go:718`
(`Updater: updater{…}`), seit `5ff4bc89` vom 24.10.2026 („controller: wire the updater, and serve
the stage channel"). Frage 1 aus [`thomas-entscheidet.md`](thomas-entscheidet.md) ist damit gebaut;
das Blatt gibt den Stand vom 24.10. wieder und nennt dort noch `Updater: nil`.

## Neuer Name: NGC → CirroDesk

- **15.11.2026:** NGC heißt CirroDesk, Kürzel `cd` — im Code, in den Programmen und in den lebenden
  Dokumenten (PR #237, letzte Reste #242). Datierte Berichte behalten die alten Namen und sagen das
  in ihrem Kopf (z. B. `CHANGELOG.md`).
- **19.11.2026:** Das GitHub-Repo heißt **`sorglos-it/cirrodesk`** (öffentlich). Die alte Adresse
  `sorglos-it/nextgencloud` leitet weiter (`gh api repos/sorglos-it/nextgencloud` antwortet mit
  `sorglos-it/cirrodesk`).
- **Das ulh-Repo ist `sorglos-it/unknown-helper`** (öffentlich). In seinem `main` (`2181c9bb` vom
  07.11.2026) heißen die Skripte noch `apps/cli-ulh/assets/scripts/ngc-*.sh`. Der Zweig
  **`cd-umbenennung`** (Partner-Thread ulh-cd) ist fertig: 6 Commits vor ulh-`main`, zuletzt
  `5146658` (cd-web schreibt die Balancer-Datei aus dem Vertrag) und `010cfdc` (CHANGELOG). Dort heißen
  die Skripte `apps/cli-ulh/assets/scripts/cd-*.sh`, und er bringt die Balancer-Zeilen aus
  `shared/ulh.md` mit, die ulh fehlten (`connectiontimeout=2`, `Retry-After`, `LogLevel`). **Er liegt
  nur lokal** in `Z:\GitHub\unknown-helper`, nicht auf GitHub. **Nachtrag 20.11.2026:** jetzt 9 Commits
  auf ulh-`main` `2181c9b`, dazugekommen `6747f07` (`RuntimeDirectoryPreserve=yes` auch für die API),
  `32ef215` (viertes Stück, siehe „Auslieferungskette") und `7c1e554` (Unit-Kommentar nennt nur
  `join-secret`); weiter nur lokal.
- **Labor, laufende Arbeit (Entscheidung Thomas, 19.11.2026):** Auch die Laborbezeichnungen werden
  umbenannt. Die Proxmox-VMs 300–309 heißen seit dem 19.11.2026 `cd-vorlage`, `cd-web-api-1`,
  `cd-api-2`, `cd-db-1` bis `cd-db-3`, `cd-svc-1`, `cd-svc-2`, `cd-store` und `cd-ein` (Angabe des
  Koordinators). PR #243 (`labor-cd`, offen) stellt `tools/` und die lebenden Dokumente auf
  `cd.mein.lan` und `cd-*` um. **In den laufenden Gästen gelten Hostnamen, Zertifikate und
  `/etc/hosts` noch für `ngc.mein.lan`** (dort läuft NGC 0.2.x), bis die Anlage mit dem ersten
  `cd`-Release neu aufgesetzt wird. Historische Angaben (etwa `ssh ngc`, 10.2.3.208, datierte
  Berichte) bleiben unverändert.
  **Nachtrag 20.11.2026:** api-2 ist am 24.10.2026 beigetreten
  ([`proxmox-sechs.md`](proxmox-sechs.md)), also mit dem alten `-join` von vor PR #247, der ein
  eigenes Schlüsselpaar meldete statt der Schlüssel aus `api.json`. api-2 steht darum vermutlich
  unter Schlüsseln in der Vertrauensliste, die er nicht hat — **nicht gemessen**. Unbemerkt, weil
  nichts zu api-2 durchkommt ([`ausfall-je-maschine.md`](ausfall-je-maschine.md)). Das bleibt so,
  bis die Anlage neu aufgesetzt wird.
- **Die drei Releases tragen `ngc-`Namen** (`ngc-api-*.tar`, `ngc-controller-*.tar`,
  `ngc-web-all.tar`; `gh release view v0.2.2`). Darum gehen das erste Release mit `cd`-Namen und
  `cd-umbenennung` nur **zusammen** raus — siehe Nachtrag unter „Auslieferungskette".

## Nächster Schritt

1. **Umbenennung abschließen:** Der ulh-Zweig `cd-umbenennung` ist fertig und liegt nur lokal; er und
   das erste `cd`-Release gehen zusammen raus — **beides macht Thomas**. Dazu die Labornamen (PR #243);
   die Anlage selbst wird mit dem ersten `cd`-Release neu aufgesetzt.
2. **Die Sperren aus [`was-fehlt.md`](was-fehlt.md)**, Abschnitt „Stand am 15.11.2026". Danach
   zusammengeführt, laut PR-Titel: A24 (#224), A2 (#225), A25 (#227), A3 (#234) und S12, die
   Oberfläche in fünf Sprachen (#230, #233, #238–#241). **Erst nachsehen, dann anfangen** — eine
   Sperre kann längst zu sein.
3. **Was nur Thomas kann**, steht in [`was-fehlt.md`](was-fehlt.md), „Stand am 20.11.2026",
   Abschnitt „Was nur Thomas kann": Release-Schlüssel (S3), Laden (S1), das erste cd-Release mit dem
   Push nach ulh-`main`, die Anlage neu aufsetzen, ein öffentlicher Name für S10. **Alles andere
   entscheidet der Koordinator** nach Thomas' Vorgabe vom 24.10.2026 (Stabilität vor Sicherheit vor
   Leistung); seit dem 20.11.2026 gilt sie ausdrücklich auch für neue Fragen. Jede Entscheidung steht
   im Protokoll des Entwurfs, Abschnitt 6.
4. ~~**Die 31 Zweige ohne Pull Request**: je Zweig entscheiden — Pull Request oder löschen.~~
   **Gesichtet am 20.11.2026** (nächster Abschnitt): 30 können weg, einer bleibt bis zu seinem Bau.
   Die drei Fragen an Thomas sind am selben Tag entschieden (Protokoll).
5. **Ohne Thomas baubar:** die Liste am Ende von [`was-fehlt.md`](was-fehlt.md), Abschnitt „Stand am
   20.11.2026" (PR #254), nach Teilen geordnet.

## Zweige: was wirklich offen ist (gemessen am 19.11.2026)

**Gesichtet am 20.11.2026** von vier Partner-Threads, je Zweig am Code und gegen `main`; die Belege
stehen in ihren Notizen unter `D:\claude\side-session-notes\local_09105105-84d8-4761-814a-f7842c285333\`
(`zweige-web.md`, `zweige-api.md`, `zweige-sichten-2026-11-20.md`, `zweige-doku-2026-11-20.md`).
Gelöscht ist noch nichts; das entscheidet Thomas.

| Ergebnis | Zweige |
|---|---|
| Rest als PR auf `main`, Zweig kann weg | `fristen-hintergrund` (#256), `ngc-rootrechte` (#259), `pool-schreibkanal` (#255), `ngc-ersteapp` (#260), `ngc-abmelden` (#258), `ngc-upload-bestaetigung` und `web-pruef-main` (#257), `web-reif` (#251, #257), `ngc-knotenfassung` (#251), `zeit-vm` (#249), `punkt54-paar-echt` und `messung-punkt53-paar` (#250), `docs-apps-bauen` (#252), `stand-aktuell` (#244, #246, #253), `ngc-join-keys` (#247) |
| schon in anderer Form auf `main`, kann weg | `p-web-s12-store` (#238), `endgeraet-echo-merge`, `term-echo`, `web-lang` (S12 baute anders), `ngc-knotenkennung`, `ctl-aktivierung-kanalfall`, `ngc-fscrypt-echt`, `ngc-appweg`, `riegel-ed25519-messen`, `release-packer`, `pruefauftraege`, `pruefergebnisse` |
| bleibt, bis gebaut | `api-knoten-kennung`: Abschnitt 9 nennt je Stelle die Zeile für A8c (Apache-Protokoll) |
| entschieden am 20.11.2026 (Protokoll), kann weg | `veroeffentlichungsliste` (es gilt die Sprachschicht auf `main`), `ngc-zweiter-mensch` (Weg (b): ein Konto für alle Sitzungen; das Blatt ist `c83da1d1`), `claude/elegant-thompson-6m2un0` (Plan Standard-Apps vom 25.09. nicht übernommen) |

200 Zweige liegen auf GitHub außer `main`:

| Anzahl | Was |
|---:|---|
| 142 | ganz in `main` enthalten |
| 1 | inhaltlich leer: jede Zeile steht schon in `main` (`ngc-uhr-fassung-pruefung`) |
| 20 | Pull Request zusammengeführt, danach kein Commit mehr auf dem Zweig — der Unterschied zu `main` ist also auf `main` entstanden |
| 5 | Pull Request zusammengeführt, **danach** noch Commits: `ausfall` (7), `gui-gestaltung` (2), `zweite-app`, `docs/portierbarkeit`, `design-marken` (je 1) |
| 1 | offener Pull Request: #243 `labor-cd` |
| 31 | **tragen Zeilen, die `main` nicht hat, und hatten nie einen Pull Request** |

Die drei, die schon am 10.11.2026 als dringend galten, tragen ihre Arbeit weiterhin:

- ~~**`ngc-join-keys`** (06.11.): „join: register the keys this server holds, and write none to
  `/run`" — kein privater Schlüssel mehr unter `/run`, also Sicherheit. 5 Dateien, +258 Zeilen.~~
  **Erledigt 20.11.2026:** als Zweig `join-keys` über PR #247 auf `main` (Merge `35565572`) —
  Cherry-Pick, Konflikte nur an Namen, dazu ein Fall, der auf `b7a21def` 5 von 5 rot war.
  `ngc-join-keys` selbst trägt nur noch die alten Namen und kann weg.
- **`fristen-hintergrund`** (16.10.): 13 Dateien, +1719 Zeilen; ein Commit darin nennt ein
  Datenrennen, das ein `-race`-Lauf gefunden hat.
- **`ngc-rootrechte`** (19.10.): das Rechtemodell über Web, API und Controller, 30 Dateien, +1358
  Zeilen.

Keiner der 31 ändert `shared/`. Ob ihre Arbeit noch gebraucht wird oder längst anders gebaut ist,
sagt die Zahl nicht — unter den 31 ist auch `stand-aktuell`, den diese Fassung ablöst.

<details><summary>Alle 31, neueste zuerst</summary>

| Zweig | Letzter Commit | Commits vor `main` | Zeilen, die `main` nicht hat |
|---|---|---:|---:|
| `p-web-s12-store` | 15.11.2026 | 1 | 66 |
| `stand-aktuell` | 13.11.2026 | 2 | 170 |
| ~~`ngc-join-keys`~~ über #247 auf `main` (20.11.) | 06.11.2026 | 1 | 333 |
| `docs-apps-bauen` | 04.11.2026 | 1 | 359 |
| `api-knoten-kennung` | 03.11.2026 | 4 | 463 |
| `ngc-zweiter-mensch` | 29.10.2026 | 3 | 200 |
| `ngc-rootrechte` | 19.10.2026 | 6 | 2248 |
| `veroeffentlichungsliste` | 19.10.2026 | 35 | 4965 |
| `ngc-fscrypt-echt` | 19.10.2026 | 1 | 34 |
| `ngc-knotenkennung` | 18.10.2026 | 5 | 1358 |
| `ngc-abmelden` | 18.10.2026 | 4 | 535 |
| `web-lang` | 18.10.2026 | 2 | 3250 |
| `ngc-appweg` | 18.10.2026 | 2 | 146 |
| `ctl-aktivierung-kanalfall` | 18.10.2026 | 2 | 1115 |
| `web-reif` | 18.10.2026 | 14 | 4287 |
| `zeit-vm` | 18.10.2026 | 1 | 91 |
| `web-pruef-main` | 18.10.2026 | 1 | 890 |
| `term-echo` | 18.10.2026 | 1 | 51 |
| `ngc-ersteapp` | 18.10.2026 | 1 | 303 |
| `ngc-knotenfassung` | 17.10.2026 | 3 | 662 |
| `ngc-upload-bestaetigung` | 17.10.2026 | 2 | 669 |
| `endgeraet-echo-merge` | 17.10.2026 | 2 | 34 |
| `fristen-hintergrund` | 16.10.2026 | 5 | 1849 |
| `punkt54-paar-echt` | 16.10.2026 | 4 | 361 |
| `pool-schreibkanal` | 16.10.2026 | 3 | 1034 |
| `riegel-ed25519-messen` | 15.10.2026 | 1 | 238 |
| `messung-punkt53-paar` | 13.10.2026 | 1 | 307 |
| `release-packer` | 02.10.2026 | 7 | 208 |
| `pruefauftraege` | 25.09.2026 | 4 | 368 |
| `claude/elegant-thompson-6m2un0` | 25.09.2026 | 4 | 343 |
| `pruefergebnisse` | 25.09.2026 | 1 | 286 |

</details>

**So gemessen**, nach `git fetch --prune`, je Zweig:

    git merge-base --is-ancestor origin/<zweig> origin/main                  # ganz enthalten?
    git diff --numstat origin/main origin/<zweig> -- <dateien des zweigs>    # Zeilen, die main fehlen
    gh pr list -R sorglos-it/cirrodesk --state all                           # gab es je einen PR?

`<dateien des zweigs>` nennt `git diff --name-only $(git merge-base origin/main origin/<zweig>)
origin/<zweig>`. „Danach noch Commits" heißt: Commit-Zeit nach der Merge-Zeit des Pull Requests.
**Drei Punkte (`origin/main...<zweig>`) taugen nicht:** das vergleicht gegen den Abzweigpunkt und wird
nie leer.

## Ledger

Die sechs Ledger der ersten Pläne (`.superpowers/sdd/<plan>/progress.md`, nur im Hauptcheckout
`Z:\GitHub\cirrodesk`, nicht im Repo) enden zwischen dem 22.09. und 13.10.2026 (Dateidatum). Danach
lief die Arbeit über Pull Requests und [`was-fehlt.md`](was-fehlt.md).

**Wer einen Ledger liest: jeder hat sein eigenes Wort für „fertig".** Gezählt am 19.11.2026: Plan 0
schreibt `complete` (35-mal). Web, API, Controller und Datenbank mischen `complete`, `COMPLETE` und
`ABGESCHLOSSEN`, drei davon auch `abgeschlossen`. Recovery schreibt vor allem `ABGESCHLOSSEN` — sein
eigener Eintrag sagt, dass `grep 'Task N: complete'` dort nur Aufgabe 1 fand, obwohl 1 bis 6 fertig
waren. `DONE` und `DONE_WITH_CONCERNS` heißen dort „umgesetzt", nicht „geprüft und fertig". **Nach
allen Wörtern suchen, nicht nach einem.**

## Im ulh-Repo (öffentlich — dort schiebt nur Thomas, nie ein Agent)

Gemessen am 19.11.2026 an ulh-`main` (`2181c9bb`):

- **Kein offener Pull Request.** Neben `main` liegen 18 Zweige unter `archiv/ulh/` und
  `ulh-s3a-kommentar` (1 Commit vor `main`, ohne Pull Request — die Berichtigung S3a aus
  [`was-fehlt.md`](was-fehlt.md)). Die am 10.10.2026 genannten Stände `75c63f4` und `33dce46` liegen
  nur noch auf `archiv/ulh/ngc-web-pruefung` und `archiv/ulh/ngc-db`, nicht in `main`.
- **Bei der API fehlt `RuntimeDirectoryPreserve=yes`.** API und Controller nennen beide
  `RuntimeDirectory=ngc` (`ngc-api.sh:1050`, `ngc-controller.sh:1231`), aber nur der Controller
  behält das Verzeichnis (`ngc-controller.sh:1256`, `ngc-all.sh:1536`). Ohne die Zeile löscht systemd
  das Verzeichnis, wenn die Unit stoppt (`man systemd.exec`) — stoppt die API, ist `/run/ngc` also
  auch für den Controller weg. Ein Nachtrag liegt auf `archiv/ulh/api-runtimedir-preserve`. Auf
  `cd-umbenennung` ist es genauso: die Zeile steht nur in `cd-controller.sh:1256` und `cd-all.sh:1536`.
  **Nachtrag 20.11.2026:** auf `cd-umbenennung` behoben mit `6747f07` — die Zeile steht jetzt auch in
  der API-Unit (`cd-api.sh:1072`, `cd-all.sh:1430`); in ulh-`main` fehlt sie weiter.
- **Die Vorlage für den virtuellen Host wird gelesen:** `ensure_vhost_template` in `ngc-web.sh`.
- **Berichtigt gegenüber `stand-aktuell`:** Dass `Delegate=yes` in ulh-`main` fehlt, ist **kein
  Loch, sondern entschieden** — Punkt 101, Abschnitt „Warum nicht `Delegate=yes`": der Controller
  legt die cgroup selbst an.
- **Nachtrag 20.11.2026, Beitritt:** Seit PR #247 schreibt `cd-api -join` weder `join-noise` noch
  `join-sign`, und ulh hat beide nie gelesen — `join_installation` liest nur `join-secret` und behält
  `noiseKey`/`signKey` aus `api.json`. Die Unit-Kommentare, die die zwei Dateien nannten, sind auf
  `cd-umbenennung` berichtigt: `7c1e554` vom Koordinator (`cd-api.sh`, `cd-controller.sh`, `cd-all.sh`
  und die Bytevergleiche der drei `*.test.sh`; ulh-Tests 864 bestanden, 0 rot, laut Koordinator).

## Schon entschieden — nicht neu aufmachen

- **Go, nicht Rust:** Punkt 67 im Entscheidungsbuch, gemessen am 20.10.2026. Der Versuch liegt unter
  `versuch-rust/` als Beleg, nicht als Bestandteil. Für den Recovery-Modus ist Rust dagegen gesetzt:
  alles, was Schlüssel berührt ([`recovery-modus.md`](recovery-modus.md)).

## Fallen dieser Umgebung (am 10.10.2026 gemessen, alle bezahlt)

Vollständig mit Messungen: `Z:\GitHub\claude\docs\programmierung.md`, Abschnitt B7 „Prüfen und
Messen" (gilt für jedes Projekt), und `docs/arbeitsweise.md` (CirroDesk-eigene Lehren).
**Berichtigt 19.11.2026:** hier stand `Z:\GitHub\claude-md\messen.md` — die Datei gibt es nicht mehr.
Die, die am meisten gekostet haben:

- **Arbeitsbäume werden während der Arbeit hohl.** Mehrmals gemessen: ein Baum stand nach einer
  Wiederherstellung bei 0 fehlenden Dateien und kurz darauf wieder bei **629**; parallel 145, 123,
  159 und einmal **alle 1097**. Ursache **nicht gemessen**. Darum: `git diff --name-only
  --diff-filter=D HEAD | wc -l` muss **0** sein, und zwar **vor jeder gemeldeten Zahl im selben
  Aufruf**. Fehlt Fremdcode, scheitert der Bau sichtbar — **fehlt eine Testdatei, läuft der Satz
  mit weniger Fällen und sieht grün aus.** Wiederherstellen nur die **gelöschten** Pfade
  (`git restore --source=HEAD -- <pfade>`), nie `git restore .`.
- **`tools/check.sh` war im Arbeitsbaum kaputt** und ist repariert: `git ls-files` starb an
  „dubious ownership" (Code **128**), und der go-Werkzeugkasten an derselben Frage beim
  VCS-Stempel. **Der Ausstiegscode einer Pipe gehört dem letzten Glied** — mit `| tail` kam **0**,
  während das Skript tot war. **Der Beleg ist die gedruckte Zeile `check: ok`, nicht die Null.**
- **`$?` ist in `wsl -- sh -c '…'` schon tot**, auch in einfachen Anführungszeichen (Sonde:
  im Befehlstext **0**, dieselben Zeilen als **Datei** 7 und 1). Messskript als Datei, Code **in**
  WSL schreiben, Sonde `(exit 7)` vor der ersten echten Messung.
- **Zeitfälle flattern unter Last, und bisher war jedes Mal die Vorrichtung schuld, nicht der
  Code.** Am 12.10.2026 zwei Fälle gemessen, nachgestellt (je einzeln in **16 gleichzeitigen
  Prozessen** — die billigere Nachstellung des vollen Satzes; reine Rechenlast warf keinen um) und
  behoben:
  - `TestThePoolRetiresAConnectionAtItsLifetime`: **25 rot von 1600 → 0 von 1600**. Ursache
    gemessen: `database/sql` zieht eine überfällige Verbindung an **zwei** Stellen ein — beim
    Herausnehmen **und** beim Zurückgeben. Bei einer Frist von 1 ms greift die zweite unter Last
    immer. Die erwartete Zahl war **abgeleitet** und hing daran, dass ein
    Unter-Millisekunden-Fenster genau einmal getroffen wird. Behoben, indem die **Uhr aus der
    Zusicherung** genommen wurde (die Ersatzverbindung wird gehalten, also kann sie niemand
    zurückgeben) — die Zusicherung bleibt exakt `== 1`.
  - `TestCloseEndsTheLinkInOrder`: **5 rot von 3200 → 0 von 3200**. Ursache gemessen: die Attrappe
    setzte `SO_LINGER 0`, ihr Close schickte also **RST** statt FIN, und **auf Windows verwirft ein
    eintreffendes RST den noch ungelesenen Empfangspuffer** — damit verschluckte es das Echo des
    Close-Rahmens, das schon abgeschickt war. Behoben, indem die **Ursache** wegfiel (Gegenseite
    ohne `SO_LINGER 0`); die Fälle, die einen Reset brauchen, behalten ihn.

  **Beide Male war die erste Frage die richtige: ist es ein echter Fehler?** Beim Reset wurde das
  ausdrücklich geprüft, bevor dem Fall Luft gegeben wurde — der Rahmen ging raus, mit dem richtigen
  Status, und die Gegenseite sah ihn. **Einem Fall Luft zu geben, der einen echten Fehler verdeckt,
  ist schlimmer als das Flattern.**

  **Darum für jede Meldung:** einen roten Zeitfall **allein nachfahren** und die Rate nennen, und
  eine Behebung mit **Rate gegen Rate** belegen (gleiche Zahl Läufe vorher und nachher), nicht mit
  „jetzt ging es". Und je Fall einen Mutanten: die übliche falsche Behebung ist, die Zusicherung so
  weit zu lockern, dass sie immer stimmt.

  **Offen, unter dreifach überzeichneter Last gefunden:** `TestAWrongSubprotocolIsRefused` fiel in
  1 von 10 Paketläufen um („not at once but after 2.4 s") — dieselbe Sorte, ein „sofort" an der Uhr
  gemessen. **Ursache nicht gemessen**, und ob die Last der Lage im vollen Satz entspricht, ist
  offen.
- **`ok` ohne `-v` sagt nichts über die Zahl der gelaufenen Fälle.** Am 12.10.2026
  gemessen: `go test` druckte `ok`, unter `-v` waren **10 von 11** Fällen `SKIP` — die
  Attrappe baute nicht (ein Heredoc hatte die Escapes gefressen), und **`t.Skipf` schreibt
  seine Zeile nur unter `-v`**. Ein Satz, der sich stumm selbst abschaltet, sieht grün aus.
  Darum: **die Fallzahl nennen, nicht „ok"** — `-v` fahren und zählen, und ein Nachbau, der
  nicht baut, ruft `t.Fatalf`, nicht `t.Skipf`. Ebenso: `go list` stirbt im Arbeitsbaum an
  derselben „dubious ownership" wie `check.sh` (`error obtaining VCS status: exit status
  128`) und meldet dabei **`packages=0`** — eine Null, die nach einer Zahl aussieht;
  Abhilfe `-buildvcs=false`.
- **`go test` über alle Controller-Pakete stirbt in 3 von 5 Läufen**, und zwar nicht an
  einem Testfall: `fatal error: unknown caller pc`, immer in `scanstack`/`markroot` des
  Müllsammlers der Go-Laufzeit (`go1.27.1`), immer im Paket `classes/install`. **Allein
  ist dasselbe Paket grün**, auch unter künstlicher Last. Am 12.10.2026 gemessen, 3/5 —
  und ebenso 3/5 auf dem Stand **davor** (`01f8a6e`), also nicht von einer Änderung
  verursacht. Daneben gemessen: `go test -race ./apps/server-api/...` scheitert
  reproduzierbar mit `open …httpapi.test.exe: Zugriff verweigert`, was auf den
  Virenschutz (Bitdefender) deutet — das ist aber eine **andere** Erscheinung (sie bricht den
  **Bau** ab, der Absturz passiert **im Lauf**), und in sechs roten Ausgaben kommt sie nicht vor.
  **Die Ursache ist gefunden und liegt nicht bei uns: golang/go #81731**, ein Speicherfehler der
  Go-Laufzeit 1.27.0/1.27.1 auf windows/amd64 — `LockOSThread` um einen Systemaufruf herum
  überschreibt Nachbarspeicher auf dem Goroutine-Stack. Der abstürzende Stack führt über
  `os/user.Lookup` → `syscall.LookupAccountName`, und `os/user/lookup_windows.go:309` ruft genau
  dieses `LockOSThread`; nur drei Controller-Pakete hängen an `os/user`, und `install` ist das,
  dessen Fälle den Weg wirklich gehen. Eigener Code ist ausgeschlossen (kein `go:linkname`, kein
  Assembler, kein cgo, kein `unsafe`, kein `LockOSThread` in `apps/`, `shared/`, `vendor/`).
  **Behoben in go1.27.2** (Backport #81829), darum steht in `go.mod` jetzt diese Fassung.
  **Die Rate schwankt mit dem Maschinenzustand, nicht mit dem Stand** — morgens 3/5 rot, mittags
  0 von 102 Satzläufen, auch mit 16 gleichzeitigen Läufen und kaltem Cache. **Warum: offen.**
  **Die Trennlinie für jede Meldung:** ein roter Lauf mit `unknown caller pc`,
  `unexpected return pc`, `found pointer to free object`, `sweep increased allocation count`,
  `fatal error: fault` oder `Exception 0x` ist **kein Befund** — verwerfen und wiederholen.
  Ein roter Lauf mit `--- FAIL` **ist** einer. ~~Und: **der Hauptcheckout kann den Controller-Satz
  gar nicht messen** — `apps/server-controller/classes` gibt es auf `main` nicht, `go list` meldet
  dort 1 Paket „no test files" und Ausstiegscode 0, was wie grün aussieht.~~ **Berichtigt
  19.11.2026:** `apps/server-controller/classes/` gibt es inzwischen auf `main`.
  **Daraus folgt für jede Meldung: ein einzelnes grünes `go test ./...` ist hier keine
  Aussage.** Das eigene Paket einzeln messen, und wenn der Gesamtsatz genannt wird, ihn
  mehrfach fahren und die Zahl der roten Läufe mitnennen.
- **`git` stirbt in der Bash mit „Out of memory", während 98 GB frei sind.** Gemessen am
  10.10.2026 bei vielen gleichzeitigen Strängen: `git commit` brach mit „malloc failed (tried to
  allocate 1048576 bytes)" ab, `git push` mit „calloc failed", und derselbe Befehl lief
  **unverändert** über PowerShell (`git -C Z:\GitHub\cirrodesk push`) sauber durch. Es ist nicht
  der Arbeitsspeicher, sondern die Fork-Schicht von Git-für-Windows (dieselbe Messung zeigte
  `post-commit: fork: Resource temporarily unavailable`). **Also: bricht git in der Bash mit
  „Out of memory" ab, erst den Arbeitsspeicher nachsehen, dann denselben Befehl über PowerShell
  wiederholen — und vorher prüfen, ob der Commit schon liegt (`git log --oneline -1`), sonst
  entsteht er zweimal.**
- **Der Scratchpad ist zwischen den Strängen geteilt**, nicht nur `/tmp`: eine Messdatei wurde
  mitten im Lauf von einem anderen Strang überschrieben. Je Strang ein eigener Unterordner, und
  eine Zahl aus einer überschriebenen Datei wird **neu gemessen**, nicht geborgen.

### Nachträge bis 20.11.2026

- **Der Hauptcheckout stand nicht immer auf `main`.** Am 10.11.2026 stand er auf
  `feldzug-echter-server`, 912 Commits hinter `main` (gemessen im Zweig `stand-aktuell`); am
  19.11.2026 stand er auf `main`. Jede Datei dort kann also alt sein, auch `CLAUDE.md` und diese.
  Vor dem Lesen:

  ```powershell
  git -C Z:\GitHub\cirrodesk rev-parse --abbrev-ref HEAD
  git -C Z:\GitHub\cirrodesk rev-list --count HEAD..origin/main
  ```

  Steht da nicht `main` und `0`: mit `git show origin/main:<datei>` lesen oder einen eigenen
  Arbeitsbaum anlegen, unter `D:\claude\working\cirrodesk-<name>`, **nie auf `Z:`**.
- **Der git-Index eines Checkouts ist zwischen allen Sitzungen darin geteilt:** `git commit --
  <pfade>`, nie `git add -A` (`docs/arbeitsweise.md`). Und **nie `git config` in einem
  Arbeitsbaum** — alle Arbeitsbäume eines Repos teilen eine Konfiguration.
- **Ein zweiter Prüflauf für denselben Zweig bricht den ersten ab**, auch Push und Pull Request
  desselben Commits (`.github/workflows/linux.yml`, `concurrency` mit `cancel-in-progress`).
  `cancelled` ist kein `failure`: je Lauf die Schlussfolgerung lesen (`gh run list`). Auf `main`
  zuletzt am 16.11.2026: `85244a0d` `cancelled`, gleich danach `2e7bc052` `success`.
- **Flatternd zählt rot.** `tools/ci-count-web.mjs` wertet einen Fall, der erst im zweiten Versuch
  grün wird, als rot.
- **Einen Flatterer zuerst am Protokoll zuordnen:** eine gefallene Zusicherung heißt Code — im
  Quelltext suchen; nur Verbindungsfehler oder eine wandernde Zahl übersprungener Fälle heißen
  Umgebung. Gemessen in [`arbeitsweise.md`](arbeitsweise.md), „Sechs Messungen gleichzeitig messen
  die Maschine, nicht den Code": zwei rote Läufe mit 9 und 15 Fällen, keine einzige gefallene
  Zusicherung, nur „Could not connect to server" — und ein Serverausfall macht aus einem
  übersprungenen Fall einen laufenden.
- **Die Testanlage startet jede Nacht durch.** Der Sicherungsauftrag des Proxmox-Wirts läuft mit
  `mode stop` ab 03:00 (`/etc/pve/jobs.cfg`) und hält dabei jede Maschine an; die CirroDesk-VMs
  trifft es gegen 04:40 Ortszeit des Wirts. Gemessen am 25.10.2026,
  [`sicherung-und-weg-zurueck.md`](sicherung-und-weg-zurueck.md), Abschnitt 12. Der Zustand „läuft
  seit Tagen" kommt dort also nicht vor, und ein Fehler, der sich über Stunden anhäuft, ist morgens
  weggewischt.

### Ein dritter Zeitflattrer: `TestStartTakesTheSameTimeForAnUnknownName`

`apps/server-api/classes/login`. Er verlangt, dass ein unbekannter Name nicht schneller abgelehnt
wird als ein bekannter, und misst dafür einen Faktor zwei. Am 13.10.2026 einmal rot im vollen
`tools/check.sh`-Lauf (`263,877µs` gegen `118,197µs`, „this machine cannot measure a factor of two
today"), **allein 0 von 10 rot**, und der zweite vollständige Lauf ebenfalls grün mit `check: ok`.

Gleiche Familie wie die zwei oben: eine abgeleitete Erwartung gegen ein Fenster im
Mikrosekundenbereich, auf einer Maschine mit anderen Lasten. **Ursache nicht gemessen** — er liegt
auf keinem Weg, der an diesem Tag geändert wurde (der Fall ruft den Handler direkt, nicht über den
Router). Wer ihn rot sieht: erst allein wiederholen, und mit `GOMAXPROCS=1` prüfen, bevor er zum
Befund wird.

### `tools/check.sh` scheitert im Arbeitsbaum mit Code 128 — UNC-Pfad, nicht Eigentum

Nicht „dubious ownership": `git worktree add` schreibt in die `.git`-Datei einen **UNC-Pfad**
(`//nas.mein.lan/freigaben/…`), den Git-Bash über seine Pfadumschreibung zerlegt —
`fatal: not a git repository: /mnt/c/<baum>/Z:/GitHub/…`. **Abhilfe gemessen:**
`MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*'`, dann läuft derselbe Befehl. `Z:/…` in die
`.git`-Datei zu schreiben hilft **nicht** (`fatal: not a git repository: (NULL)`).

### `git stash push` ohne `-u` lässt neue Testdateien stehen

Der „Grundstand"-Lauf misst dann **alte Quelle gegen neue Fälle**. Einmal war das nützlich (es
wurde die Rot-vorher-Messung des ganzen Satzes), aber es wäre leicht als Grundstand durchgegangen.
Mit `-u` dagegen verschwindet ein eigener Messordner im Arbeitsbaum mitten im Lauf —
**Messdateien gehören in den Scratchpad, nicht in den Baum.**

### Ein vierter Zeitflattrer: `TestFinishTakesTheSameTimeForAKeyNobodyHolds`

`apps/server-api/classes/login`, 1 von 2 `check.sh`-Läufen rot (174,57 s), allein 0 von 3. Der
Null-Kontrollzweig des Falls feuert selbst („this machine cannot measure a factor of two today").
Gleiche Familie wie die drei darüber. **Freigemessen** mit zehn verschachtelten Paaren
(Grundstand und Änderung abwechselnd, gleicher Maschinenzustand je Paar): **0/10 gegen 0/10**.

**Und eine Lehre dazu:** `GOMAXPROCS=1` ist bei einem Fall, der **Uhrzeiten** vergleicht, der
falsche Hebel — `messen.md` nennt es für ein Flattern an einer Reihenfolge zwischen Goroutinen.
Das Verengen des Planers macht die Uhr nur verrauschter; der Lauf war nach über 30 Minuten nicht
durch. Bei einem Zeitvergleich trennt der **gepaarte Lauf gegen den Grundstand**, nicht der
verengte Planer.

## Wie hier gearbeitet wird (Head-Officer)

- **Jeder Auftrag beginnt mit „sieh nach, ob es diese Aufgabe noch gibt."** Am 10.10. sind zwei
  Stränge auf fertige Arbeit geschickt worden, weil der Koordinator einen Ledger falsch gelesen
  hat; nur weil der Strang nachsah, blieb ein 540-Zeilen-Bericht unversehrt.
- **Greift der Koordinator in einen laufenden Arbeitsbaum, sagt er es dem Strang im selben
  Moment** — mit Datei, Zeit und Grund. Zwei Stränge haben sonst eine Geistersuche gemacht.
- **Eine Planberichtigung hat drei Orte:** Plan, Hauptcheckout, neu erzeugter Auftrag — und
  **innerhalb** des Aufgabenstücks noch Kopf, Rumpf und **Codeblöcke**. Dieselbe Berichtigung
  brauchte am 03.10. drei Anläufe.
- **Dokumente wandern mit dem Zweig.** Ein Strangzweig kann Dutzende Commits hinter `main` liegen;
  ein Verweis auf „Punkt 20" ist dort etwas anderes. **Die Überschrift zitieren, nicht die Nummer**,
  und vor einer Prüfung, die Dokumente heranzieht, `origin/main` in den Zweig ziehen.

## Vor dem Bauen zu prüfen

- Sobald die Server per Remote Control verbunden sind: ob sich die Sitzungen gegenseitig
  Nachrichten schicken können (Entwurf 4.7 setzt das voraus).
- Was Thomas bereitstellt: siehe Entwurf 4.7, „Voraussetzungen“ (sieben Server, ein Pi, je ein
  YubiKey 5 NFC und Nitrokey 3 NFC).
- ~~Auf dem ersten Linux-Server mit echtem `jq`: `TestReleaseScript` (31 Fälle, auch der mit einem
  Null-Byte) muss laufen und grün sein; bisher lief es nur mit einem Ersatz für `jq`.~~
  **Berichtigt 19.11.2026:** Er läuft bei jedem Push in der Linux-Prüfung auf GitHub. Ohne `jq`
  überspringt er sich selbst; im Lauf 44573010049 (`b2262020`) steht er in keiner der drei
  SKIPPED-Listen, und rot war nichts.

## Auslieferungskette ulh ↔ Release (zwingende Reihenfolge, Stand 03.10.2026)

**Viertes Stück für das erste cd-Release (19.11.2026, aus dem Code gelesen, nicht ausprobiert):** Laut
Entscheidung S7 legt das Web-Paket die Oberfläche unter `ui/` (`apps/cli/classes/contents.go` verlangt
`ui/public/index.php`). ~~ulhs `unpack_release` in `cd-web.sh` entpackt das Archiv aber unverändert und
verlangt `public/index.php` an der Wurzel, auch auf `cd-umbenennung`. Ohne Änderung bricht die
Installation des Web-Teils mit dem ersten cd-Release ab.~~ **Erledigt 20.11.2026 auf `cd-umbenennung`**
(`32ef215`): `unpack_release` nimmt nur Namen unter `ui/` an (auch kein `ui//…`, das nach dem
Abstreifen ein absoluter Pfad wäre) und entpackt mit `--strip-components=1`; `public/index.php` und
`templates/apache-vhost.conf` liegen wie bisher. Laut Koordinator ulh-Tests 864 bestanden, 0 rot, und
die Gegenprobe ohne `--strip-components` 13 rot. Liegt wie der ganze Zweig nur lokal und geht mit ihm
raus.

**Nachtrag 19.11.2026, Umbenennung:** ulhs `main` holt `releases/latest` und erwartet `ngc-*`-Dateien
und `ngc/`-Signaturen, wie sie `v0.2.2` trägt. Das erste Release mit `cd`-Namen bricht deshalb jede
ulh-Installation, bis der ulh-Zweig `cd-umbenennung` in ulhs `main` liegt. **Release und dieser Zweig
gehen nur zusammen raus**, beides macht Thomas. Alte Adressen von `sorglos-it/nextgencloud` leiten seit
der Umbenennung des GitHub-Repos weiter (gemessen: ulhs `curl --location` kommt mit 206 an).

**Berichtigt 19.11.2026:** Von den vier Stücken unten liegen drei nachweislich — die Vorlage auf
`main` (`apps/web/templates/apache-vhost.conf`), der Packer auf `main` (`tools/pack-release.sh`, Liste
`tools/release/extra-web.txt`), und ulh-`main` liest die Vorlage (`ensure_vhost_template` in
`ngc-web.sh`). Veröffentlicht sind drei Releases. **0.2.2 trägt die Vorlage** (nachgesehen am
19.11.2026: `ngc-web-all.tar` aus `gh release download v0.2.2`, 93 Einträge, darunter
`templates/apache-vhost.conf` und `templates/apache-apis.conf.example`). Die ulh-Namen unten (`cd-web-kopfzeilen`, `cd-web`, `cd-web.sh`) hat die Umbenennung vom
15.11.2026 mitgeändert; damals hießen sie `ngc-…`, und das Skript heißt in ulh-`main` noch heute
`ngc-web.sh` (erst auf `cd-umbenennung` `cd-web.sh`). `75c63f4` liegt nur noch auf
`archiv/ulh/ngc-web-pruefung`. Der Kommentar aus dem
Nachtrag vom 10.10.2026 (unten) steht
**noch unberichtigt** in `extra-web.txt` (Zeilen 27 und 33: „cd-web.sh does NOT read it").

Vier Stücke hängen aneinander, und **in der falschen Reihenfolge macht jedes
einzelne die Installation kaputt**. Gemessen, nicht geschätzt:

1. **`apps/web/templates/apache-vhost.conf`** (Zweig `web`, `c5c3591`) — die
   Vorlage, die der Server ausliefert. Sie trägt jetzt die host-weite HSTS, die
   verengte Cipher-Liste, eigene Protokolle und `FallbackResource disabled` im
   assets-Block.
2. **Der Release-Packer** (Zweig `release-packer`, `61e21f0`) — er nimmt diese
   Vorlage namentlich auf (`tools/release/extra-web.txt:30`), sie landet beim
   Nutzer unter `/var/www/cd/templates/apache-vhost.conf`.
3. **Ein gebautes und unterschriebenes Release**, das dieses Stück enthält.
4. **Der ulh-Zweig `cd-web-kopfzeilen`** (`75c63f4`) — `cd-web.sh` **liest**
   die Vorlage aus dem Release, statt sich einen zweiten virtuellen Host zu
   schreiben.

**Warum die Reihenfolge Pflicht ist:** ohne 1-3 findet Schritt 4 die Datei
nicht. Das Skript bricht dann ausdrücklich ab („Update CirroDesk to a release that
carries the file"), ändert nichts — **und damit kann niemand mehr einrichten.**
Dieser ulh-Zweig allein ist also **schlechter** als der Stand davor
(`ad5ad0f`), obwohl er ein Sicherheitsloch schließt.

**Also: 1 und 2 nach `main`, dann ein Release, erst dann darf 4 zusammengeführt
werden.** Das ulh-Repo ist öffentlich und jeder Push nach `main` erreicht per
Auto-Update sofort alle Nutzer — Schritt 4 macht **Thomas**, nicht der Agent.

**Was auf der ersten Maschine mit einem echten Apache noch zu messen ist**,
steht als Befehlsliste in `apps/web/README.md` (neun Zeilen) und als A-D im
Kopf der Vorlage selbst. **Befehl D zählt am meisten:** er prüft, dass eine
fehlende Datei unter `/assets/` von Apache und nicht von PHP beantwortet wird —
sonst stehen dort zwei `Cross-Origin-Resource-Policy`, und die Wand wirkt nicht.

### Nachtrag 10.10.2026: ein Kommentar des Packers ist überholt — vor dem Release berichtigen

Der Packer nimmt `templates/apache-vhost.conf` mit, und der Kommentar daneben
(`tools/release/extra-web.txt`) sagt: *„cd-web.sh does NOT read it — it writes its own virtual
host, lines 978-1028 — and the task that makes it read this one cannot start before the file is on
the server."*

**Das war am 02.10. richtig und ist es nicht mehr.** Gemessen am 10.10.2026:

| Zweig | Treffer auf `ensure_vhost_template` / `templates/apache-vhost.conf` |
|---|---|
| `cd-web` (`ad5ad0f`, der alte Stand) | **0** |
| `cd-web-kopfzeilen` (`75c63f4`, der neue) | **12** |

Das Skript **liest** die Vorlage jetzt, bricht sogar ab, wenn sie fehlt. Der Satz „die Aufgabe, die
es lesen lässt, kann nicht beginnen" beschreibt eine Aufgabe, **die inzwischen erledigt ist**.

**Also, in dieser Reihenfolge:** wenn `cd-web-kopfzeilen` zusammengeführt wird, im selben Zug den
Kommentar in `extra-web.txt` berichtigen. Er ist kein Beiwerk: er ist der Grund, warum die Datei
überhaupt auf der Liste steht, und der Nächste liest ihn als Begründung.

**Was dagegen gemessen stimmt** (Packer-Nachmessung, 10.10.2026): die **gepackte** Vorlage ist
**bitgleich** zur neuen (`git hash-object` des aus dem Archiv gezogenen Files gleich
`origin/web:apps/web/templates/apache-vhost.conf`, 40 374 Bytes), zwei Packläufe geben dieselbe
Prüfsumme, **66 Dateien, 0 unter `tests/`**, und alles, was die Vorlage nennt und beim Nutzer
liegen muss, liegt im Archiv. Offen bleiben zwei Kleinigkeiten: die mitgelieferte `README.md` und
die Vorlage verweisen auf `docs/STAND.md`, **das nicht im Archiv liegt** (für einen Nutzer ohne
Repo ein Verweis ins Leere), und eine Zusicherung der Dichtheitsprüfung **überspringt weiter laut**,
weil `jq` auf diesem Rechner fehlt.

## Offene Entscheidungen

- ~~**Eine Anlage kann keinen zweiten Benutzer bekommen.**~~ **Berichtigt 19.11.2026: gebaut** am
  24.10.2026 (`72ceb7e7` „web: let an admin create a role and a second user"). Was dabei nur ein
  Mensch tun kann: [`zweiter-benutzer-schritte.md`](zweiter-benutzer-schritte.md). Entwurf und
  Abwägung: [`superpowers/specs/2026-10-23-zweiter-benutzer.md`](superpowers/specs/2026-10-23-zweiter-benutzer.md).
- ~~**ES GIBT KEINEN RELEASE-PACKER**~~ **Berichtigt 19.11.2026: es gibt ihn**
  (`tools/pack-release.sh`, `tools/build-release.sh`), und drei Releases sind veröffentlicht. Offen
  bleibt der Release-Schlüssel (S3 und S4 in [`was-fehlt.md`](was-fehlt.md)).
- **Vertragsfassung 1.2:** Was der nächste Schnitt an `shared/` tragen muss,
  steht gesammelt in [`vertraege-v1.2.md`](vertraege-v1.2.md) — **83 Einträge**
  aus allen vier Strängen, jeder mit Herkunft und Messung (bis 18.10.2026 stand
  hier „18 Punkte"). Mehrere Module verweisen im Code darauf; die Datei ist die
  Liste, die dort gemeint ist. **Berichtigt 19.11.2026:** die Datei reicht
  inzwischen bis Punkt 102; ein Vorschlag für den Schnitt 1.3 liegt in
  [`vertraege-v1-3-vorschlag.md`](vertraege-v1-3-vorschlag.md).
  - Seit dem 18.10.2026 hat sie **ein Inhaltsverzeichnis mit Stand je Punkt**
    (gilt · überholt durch Punkt N · berichtigt Punkt N · offen) und **die
    offenen Entscheidungen gesammelt am Ende**. Wer eine Entscheidung sucht,
    liest diesen einen Abschnitt, nicht 83 Einträge.
  - **Die Nummern 19 bis 35 sind dort doppelt vergeben:** einmal als Punkt der
    Sammelliste am Anfang, einmal als Abschnitt `## Punkt N` darunter. „Punkt
    30" meint also zwei verschiedene Dinge. Das ist der zweite Grund für die
    Regel oben unter „Wie hier gearbeitet wird": **die Überschrift zitieren,
    nicht die Nummer.**

- **Release-Stick:** Mit welchem Werkzeug der Hardware-Stick das Release signiert, wird vor dem
  ersten Release entschieden. `cd-store` signiert nur Test-Releases.
  Stand 01.10.2026: PIV auf Firmware 5.4.3 kann kein Ed25519, und der Vertrag verlangt reines
  Ed25519 (`shared/sign/sign.go` prüft `len(sig) != ed25519.SignatureSize`). Unterschrieben wird
  `"cd/" + purpose + NUL + payload`, gemessen **193 Byte** im Beispiel. Eine Probe, ob der
  OpenPGP-Teil des Sticks diese Bytes **direkt** unterschreibt, liegt bereit; der normale
  `gpg`-Weg unterschreibt einen Hash **plus OpenPGP-Anhang** und verifiziert daher nie.
- **Für den Sicherheitsdurchgang (Plan 7) vorgemerkt:** eine feste Bindung der Oberfläche einer
  Store-App an ihren Store, die einen Wechsel des Admin-Schlüssels übersteht; `enc.Budget` soll
  Mitglieder in Rohwerten zählen; eine gemeinsame Prüfung der Form von URLs für alle Verträge.
- **Aus Aufgabe 27 (Web), je mit Messung im Aufgabenbericht:**
  - **Zurückgedrehte Vertrauensliste:** Liefert die API auf `from=N+1` eine ältere, gültig
    unterschriebene Fassung, wirft die Shell heute ihre geprüfte Liste weg und schlägt Alarm.
    Entwurf 4.6 will das Gegenteil: ablehnen, ohne Alarm, bei der neueren bleiben. Sonst kann eine
    übernommene API jede Sitzung abmelden. Entscheiden, dann `shell/trustlist.js` nachziehen.
  - **`wire.MaxCallResult` auf Ende-zu-Ende-Kanälen:** `net/channel.js` setzt die Grenze nur, wenn
    der Aufrufer sie nennt. `channel.call` nennt sie, `channel.openE2E` nicht – richtig für eine
    Aktion mit `stream`, offen für eine ohne. Festlegen, wer sie dort setzt, bevor die erste
    solche App kommt.
  - ~~**Sicherheitskopfzeilen für `assets/`**~~ – erledigt mit Aufgabe 29 (Web):
    `apps/web/templates/apache-vhost.conf` ist die Vorlage für den virtuellen Host und setzt für
    `assets/` vier Kopfzeilen: `nosniff`, `Cross-Origin-Resource-Policy`, HSTS und eine **eigene**
    Richtlinie, nicht die der Seite. Je Kopfzeile steht die Begründung in der Datei, auch für die
    bewusst fehlenden. `tests/shell/csp.spec.js` prüft es an einer `.js`, einer `.css`, einer Datei
    ohne Endung und einer, die es nicht gibt; die Mock-API liest die Werte aus der Vorlage und
    prüft dabei, dass der Block unter dem DocumentRoot liegt – ein Block woanders bricht den Lauf.
  - **Offene Reste davon, drei:**
    1. **`apps/cli-ulh/assets/scripts/cd-web.sh` im ulh-Repo (Aufgabe 26) nachziehen.** Es schreibt einen eigenen
       virtuellen Host unter demselben Namen und ist eine Aufgabe älter. Vier Unterschiede: *für*
       die Vorlage, dass das Skript für `assets/` **keine** Kopfzeile außer HSTS setzt (das Loch)
       und dass es `/assets/` aus seiner Rewrite-Regel nimmt, sodass eine fehlende Datei dort von
       Apache selbst beantwortet wird statt von `index.php`. **Gegen** die Vorlage, dass das Skript
       die Zertifikatserneuerung mit mod_md löst – ein bewegliches Teil weniger als certbot in der
       Vorlage – und dass es `cd-web.conf` bei jedem `install`/`configure` **überschreibt**: ein
       von Hand gehärteter Server verliert den Asset-Block ohne Warnung, bis das Skript ihn selbst
       mitbringt. **Berichtigt 19.11.2026:** Pfad berichtigt (hier stand `scripts/cd-web.sh`; in
       ulh-`main` heißt die Datei bis `cd-umbenennung` noch `ngc-web.sh`). Das Skript liest die
       Vorlage inzwischen (`ensure_vhost_template`); ob die zwei Punkte „gegen die Vorlage" damit
       erledigt sind, ist nicht nachgesehen.
    2. **Kein Apache im Repo.** Fünf `curl`-Befehle in `apps/web/README.md` („What no test here can
       measure") machen aus „ungemessen" „gemessen"; der wichtigste klärt, ob
       `Header always set` auf dem Rückfall-404 **ergänzt** statt zu ersetzen. Ergebnis hierher.
    3. **Nur `public/assets/` ist gedeckt.** Eine Datei direkt in `public/` hätte keine der vier
       Kopfzeilen. Der Riegel ist ein Fall, nicht die Konfiguration: direkt im Webroot liegt nur
       `index.php`, alles Neue gehört unter `assets/`.
  - **`HEAD` auf `/`** antwortet 405 statt wie `GET`, und ein 405 trägt keine `Allow`-Kopfzeile.
    Kein Loch, aber nicht nach HTTP. Lassen oder ändern.

Alle weiteren Entscheidungen stehen im Entwurf, Kapitel 6.

## Hinweise

- Aufbau des Repos: `CLAUDE.md` im Stamm, Abschnitt „Aufbau“. `vendor/` bleibt im Stamm, weil Go es
  neben `go.mod` verlangt.
- Fremdcode: alle Module aus Entwurf, Kapitel 7 liegen seit dem 22.09.2026 in `vendor/`, aufgezählt
  in `THIRD-PARTY.md`. Was noch kein Teil einbindet, hält `tools/modules.go` fest; wer ein Modul
  zuerst benutzt, nimmt seine Zeile dort heraus.
- **`-race` IST auf dieser Maschine messbar** — berichtigt am 12.10.2026, und zwar gegen einen Satz,
  der vorher hier stand und falsch war. Der C-Übersetzer liegt unter
  ~~`C:\esp\tools\llvm-mingw\bin\gcc.exe`~~ **`D:\claude\sdk\esp-tools\llvm-mingw\bin\gcc.exe`**
  (clang 23.1.1, `x86_64-w64-windows-gnu`), also in der Werkzeugkette der ESP-Arbeit — darum hat ihn
  keine Suche an den „üblichen Orten" gefunden. **Berichtigt 19.11.2026:** `C:\esp\tools\` gibt es
  nicht mehr. **Mit dem neuen Pfad neu gemessen am 19.11.2026** über den ganzen Baum (`main`
  `c3fc2139`, go1.27.4, `go test -race -count=1 -buildvcs=false -timeout 30m -json ./...`, 40 min):
  3612 Fälle bestanden, 1 rot, 102 übersprungen, `DATA RACE` 0; 67 Pakete grün, rot nur
  `shared/trust` mit `TestSignatureBeforeDecoding` (`Open: 65632 bytes allocated`, Grund unten).
  Auf GitHub läuft `-race` bei
  jedem Push (Linux, `shared/trust` getrennt, `.github/workflows/linux.yml`).
  So lief es am 12.10.2026, selbst gemessen (`relay` und `link`: 0 DATA RACE, Code 0):

  ```powershell
  $env:CGO_ENABLED="1"; $env:CC="D:\claude\sdk\esp-tools\llvm-mingw\bin\gcc.exe"
  go test -count=1 -buildvcs=false -race ./apps/server-api/...
  ```

  **Ein Paket scheitert dabei, und zwar aus einem gemessenen Grund:** `shared/trust` wird unter
  `-race` rot, ohne ein einziges Datenrennen (`DATA RACE` 0 über alle 30 Pakete).
  `TestSignatureBeforeDecoding` misst **Speicher-Obergrenzen**, und die Instrumentierung des
  Race-Detektors belegt selbst mehr: `Open: 65736 bytes allocated` gegen eine Grenze um 65 544.
  Also: **`-race` über alles fahren und `shared/trust` getrennt nennen** — der Fall kann unter
  `-race` nicht halten, solange er die Instrumentierung nicht mitrechnet, und die Zeile liegt in
  eingefrorenem `shared/`. Über den Rest des Baums: **30 Pakete, 0 DATA RACE** (12.10.2026).

  **Wie der falsche Satz entstand, weil das die eigentliche Lehre ist:** sechs Stränge haben
  unabhängig nach einem Übersetzer gesucht und keinen gefunden, einer davon rekursiv in vierzehn
  Orten. Daraus habe ich „es gibt keinen" gemacht und es als Regel aufgeschrieben — und danach
  sechs weiteren Strängen mitgegeben, **nicht** zu suchen. **Dass niemand etwas gefunden hat, ist
  keine Messung, dass es nichts gibt.** Eine Abwesenheit belegt man mit der Liste der Orte, an
  denen man gesucht hat, nicht mit der Zahl der Sucher.
- Ein Go-Update mit neuen Unicode-Tabellen ist eine Vertragsänderung (`TestValidText` hält es an).
- Der claude-mem-Dienst (`bun … worker-service.cjs`) leckt Speicher. Brechen Builds oder
  Fuzz-Läufe ab, zuerst den freien Speicher prüfen.

## Skills

- ~~`superpowers:brainstorming` und `superpowers:writing-plans`: Wiederherstellung, dann die Pläne 1–4~~
- ~~`superpowers:subagent-driven-development`: die Pläne ausführen~~

**Berichtigt 19.11.2026:** `superpowers` ist installiert, aber seit dem 09.10.2026 abgeschaltet
(`Z:\GitHub\claude\docs\programmierung.md`, Abschnitt B10); diese Skills stehen in den Sitzungen
nicht zur Verfügung.

## Artefakte

- Pläne der vier Teile (23.09.2026): [Webinterface](superpowers/plans/2026-09-23-ngc-plan-1-web.md),
  [API](superpowers/plans/2026-09-23-ngc-plan-2-api.md),
  [Controller](superpowers/plans/2026-09-23-ngc-plan-3-controller.md),
  [Datenbank](superpowers/plans/2026-09-23-ngc-plan-4-db.md); Recovery:
  [Paket 0](superpowers/plans/2026-09-23-recovery-paket-0-testumgebung.md)
- Testumgebung des Recovery-Strangs: VM `ngc-recovery` (Debian 13, `ssh ngc-recovery`); Pakete
  eingerichtet, `sudo` ohne Passwort, **kein** `/dev/kvm` (Proxmox reicht keine Virtualisierung
  durch), QEMU läuft dort emuliert
- Entwurf (freigegeben, mit Entscheidungsprotokoll): [`docs/superpowers/specs/2026-09-19-ngc-kern-design.md`](superpowers/specs/2026-09-19-ngc-kern-design.md)
- Plan 0: [`docs/superpowers/plans/2026-09-19-ngc-plan-0-vertraege-grundgeruest.md`](superpowers/plans/2026-09-19-ngc-plan-0-vertraege-grundgeruest.md)
  (historisch; maßgeblich sind `shared/` und `shared/README.md`)
- Fixpläne, Berichte, Prüfungen, Protokoll: `.superpowers/sdd/2026-09-19-ngc-plan-0-vertraege-grundgeruest/`
  (nicht im Repo)
- Verträge: [`shared/README.md`](../shared/README.md), [`shared/ulh.md`](../shared/ulh.md)
- Prototyp: [`docs/prototype.html`](prototype.html), Original unter `W:\cd\index.html`
- Repo: <https://github.com/sorglos-it/cirrodesk> (öffentlich)
- Facebook-Post: ~~zurückgestellt, solange das Repo privat ist; Vermerk in
  `Z:\GitHub\facebook\README.md`, Abschnitt „Bestand“~~ **Berichtigt 19.11.2026:** das Repo ist
  öffentlich, und `Z:\GitHub\facebook\` gibt es nicht mehr. Ob ein Post kommt, entscheidet Thomas.
