# Release v0.3.0

Novita':

- esempi `.COM` assemblati con origine logica `.com`/`.orgbase`, senza
  indirizzi runtime manuali a `0100h`
- funzioni BDOS file mutanti con opt-in esplicito `-write-disk`: `delete`,
  `write sequential`, `make file`, `rename`
- scrittura file piu' prudente: il drive host resta read-only di default
- libreria didattica copiabile `examples/lib/cpm-bdos.asm`
- esempio `write-file.asm` per creare `OUT.TXT` da un programma `.COM`
- comportamento EOF e ritorni registri BDOS resi piu' espliciti

Limiti:

- resta un ambiente CP/M-like, non una distribuzione CP/M storica
- nessuna ROM, BDOS, BIOS o immagine disco storica inclusa
- la scrittura modifica solo directory host scelte dall'utente e solo con
  `-write-disk`
- niente user area, command tail o FCB default completi in questa release

Verifica release:

```powershell
go vet ./...
go test -count=1 ./...
go run ./cmd/retronet-cpm -conformance
git diff --check
```
