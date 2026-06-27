# Release v0.4.0

Novita':

- integrazione di `github.com/retronet-labs/retronet-terminal v0.1.0` come
  adattatore console BDOS per la CLI `-run`
- esempi assembly aggiornati a `.include "lib/cpm-bdos.asm"`
- `TYPE.COM` didattico basato sul default FCB a `005Ch`
- command tail sintetico a `0080h`
- FCB default sintetici a `005Ch` e `006Ch` dalla shell `RUN`
- demo ripetibile `scripts/demo-cpm.ps1` con transcript documentato

Licenza e provenienza:

- nessuna ROM, BDOS, BIOS, font, terminfo o disco storico incluso
- tutti i programmi `.COM` della demo sono assemblati da sorgenti originali del
  repo
- i file del drive demo sono generati localmente dallo script

Limiti:

- non e' CP/M completo
- niente user area
- niente wildcard/CCP storico completo
- il websocket bridge resta futuro lavoro di `retronet-api`

Verifica release:

```powershell
go vet ./...
go test -count=1 ./...
go run ./cmd/retronet-cpm -conformance
.\scripts\demo-cpm.ps1
git diff --check
```
