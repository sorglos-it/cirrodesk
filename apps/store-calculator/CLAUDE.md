# store-calculator – Taschenrechner mit Klammern

Store-App `org.sorglos.calc`: Grundrechenarten, Punkt vor Strich, geschachtelte Klammern. Nur
Oberfläche, kein Dienst, kein Container.

- `manifest.json`, `ui/` (index.html, app.js, app.css): genau das kommt ins Paket.
- Der Ausdruck wird mit einem eigenen Parser (rekursiver Abstieg) in `ui/app.js` gerechnet,
  nie mit `eval` oder `Function`. Grenzen: 256 Zeichen, Tiefe 64.
- Texte in den 5 Sprachen der Shell (de, en, fr, it, es), Fallback en. Ausgaben nur mit
  `textContent`.
- Tests: `go test ./apps/store-calculator/` packt, unterschreibt einen Katalog mit Testschlüssel,
  prüft ihn mit `ngc-store check -store` und das Paket mit `store.Walk`; dazu laufen die
  node-Tests mit Fuzz. Fuzz länger: `NGC_FUZZ=200000 node --test apps/store-calculator/tests/`,
  einen Fehler nachstellen mit `NGC_FUZZ_SEED=<Saat>`.

Bauen:

```bash
go build -o ngc-store ./apps/cli
./ngc-store pack apps/store-calculator org.sorglos.calc-1.0.0.ngc
```
