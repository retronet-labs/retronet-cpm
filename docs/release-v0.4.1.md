# Release v0.4.1

Release di consolidamento prima di `retronet-api`.

Novita':

- aggiorna `retronet-terminal` a `v0.1.1`
- usa il terminale condiviso anche dalla shell `RUN`
- amplia la conformance sintetica con:
  - terminal console adapter
  - command tail a `0080h`
  - FCB default a `005Ch` e `006Ch`
  - failure BDOS write su drive read-only
  - BDOS write su drive mutabile in memoria
- aggiunge documentazione didattica:
  - compatibilita' CP/M-like
  - terminale condiviso
  - esempi pratici di command tail e FCB

Licenza e provenienza:

- nessuna ROM, BDOS, BIOS, font, terminfo o disco storico incluso
- tutti i test sono sintetici e scritti nel repo
- i programmi di esempio restano sorgenti originali RetroNet

Verifica release:

```powershell
go vet ./...
go test -count=1 ./...
go run ./cmd/retronet-cpm -conformance
.\scripts\demo-cpm.ps1
git diff --check
```
