# store-clock – Uhr

Store-App `org.sorglos.clock`: Uhrzeit und Datum, analog (SVG) und digital, in der Zeitzone des
Browsers. Nur Oberfläche, kein Dienst. Titel in de, en, fr, it, es; Uhrzeit und Datum formatiert
`Intl.DateTimeFormat` in der gewählten Sprache. Ausgaben nur mit `textContent`.

Tests: `go test ./apps/store-clock/` (pack, Katalog mit Testschlüssel, `check -store`,
`store.Walk`, node-Tests). Bauen: `ngc-store pack apps/store-clock org.sorglos.clock-1.0.0.ngc`.
