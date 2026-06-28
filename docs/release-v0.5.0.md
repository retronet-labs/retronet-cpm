# Release v0.5.0

Release dedicata al terminale live CP/M-like.

Novita':

- nuovo comando `go run ./cmd/retronet-cpm-live`
- shell `A>` locale sopra `retronet-terminal/live`
- flag `-disk`, `-steps`, `-alu`, `-write-disk`, `-width`, `-height`, `-line`,
  `-script`
- modalita' scriptata per test e CI
- output della shell normalizzato a CR/LF per schermi terminale
- Dockerfile esteso: builda anche `/app/retronet-cpm-live`
- dipendenza aggiornata a `retronet-terminal v0.3.0`

Uso:

```powershell
go run ./cmd/retronet-cpm-live -disk .
go run ./cmd/retronet-cpm-live -disk . -script "DIR`rHELP`rEXIT`r"
```

Limite noto:

- `RUN` resta sincrono. I programmi `.COM` che richiedono input live durante
  l'esecuzione saranno affrontati con il run loop asincrono previsto per
  `retronet-api`.

Licenza e provenienza:

- nessuna ROM, BDOS, BIOS o immagine disco storica inclusa
- nessun manuale o asset storico copiato
- test e demo restano sintetici/originali

Verifica:

```powershell
go vet ./...
go test -count=1 ./...
go run ./cmd/retronet-cpm -conformance
go run ./cmd/retronet-cpm-live -width 80 -height 12 -script "HELP`rEXIT`r"
git diff --check
```
