# Release v0.4.2

Release preparatoria per `retronet-api`.

Novita':

- nuovo package `session`
  - `New`
  - `Input`
  - `RunCommand`
  - `Prompt`
  - `DrainOutput`
  - `Snapshot`
- `shell.Config` accetta un terminale condiviso esplicito
- `disk.HostDriveOptions` per:
  - scrittura opt-in
  - limite dimensione file
  - limite numero file
- `disk.NewTemporaryHostDrive` per sessioni isolate
- documentazione API-ready:
  - sessioni programmatiche
  - sicurezza drive host
  - ponte verso `retronet-api`

Licenza e provenienza:

- nessuna ROM, BDOS, BIOS o immagine disco storica inclusa
- nessun path host arbitrario esposto come contratto API
- test e demo restano sintetici/originali

Verifica:

```powershell
go vet ./...
go test -count=1 ./...
go run ./cmd/retronet-cpm -conformance
.\scripts\demo-cpm.ps1
git diff --check
```
