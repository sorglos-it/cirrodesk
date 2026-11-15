# store-files – Dateien

Store-App `org.sorglos.files`: eigene Dateien je Benutzer ablegen, Ordner anlegen, umbenennen,
löschen, hochladen; Texte im selben Fenster bearbeiten, Bilder ansehen (Thomas, F1 = A).
Herunterladen über die Shell (`ngc.save`, F2 = A); fehlt sie, zeigt die App „Download nicht
verfügbar“.

- `main.go`, `classes/`: der Dienst `ngc-files`. `serve` ist ein Aufruf (Vertrag V3,
  `shared/service`), `hold` hält den Container. Dateien liegen unter `/data/users/<Benutzer>`;
  jeder Zugriff geht durch `os.Root`, kein Pfad und kein Link führt hinaus. Namen prüft
  `ValidPath`; Parameter `enc.DecodeStrict`.
- Aktionen: `list`, `mkdir`, `rename`, `remove` antworten JSON (`Answer`), `get` und `put` sind
  Kanäle mit rohen Bytes. Hochladen schreibt erst unter `.upload-<Name>` und benennt am Ende um.
- `ui/`: Oberfläche, Texte in de, en, fr, it, es. Kein `prompt`/`confirm`: der Rahmen hat kein
  `allow-modals`.
- `Containerfile`: Basis und Digest wie beim Terminal (`tools/pin-base.sh`).

Tests: `go test ./apps/store-files/...` (Dienst, Fuzz `FuzzValidPath`/`FuzzParams`, Paket mit
Platzhalter-Abbild, node-Tests der Oberfläche).

Bauen auf dem Baurechner (Linux mit podman), je Architektur:

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.Version=$(cat apps/store-files/VERSION)" -o build/ngc-files ./apps/store-files
podman build --arch amd64 -f apps/store-files/Containerfile -t ngc/files:1.0.0 build/
mkdir -p pkg/services && cp -r apps/store-files/manifest.json apps/store-files/ui pkg/
podman save -o pkg/services/amd64.tar ngc/files:1.0.0
ngc-store pack pkg org.sorglos.files-1.0.0.ngc
```
