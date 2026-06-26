# Contesto operativo per agenti

## Obiettivo

Implementare e mantenere `retronet-cpm`: ambiente CP/M-like didattico sopra
`retronet-8080`, testato, importabile e documentato in italiano.

## Decisioni da preservare

- Programmi `.COM` caricati a `0x0100`.
- `CALL 0005h` intercettato come ingresso BDOS.
- `PC=0x0000` trattato come warm boot/terminazione programma.
- Default operativo su `cpu.Native`, non sul default implicito `cpu.Gate`.
- `cpu.Gate` resta disponibile per demo e confronto didattico.
- Nessuna ROM, BDOS, BIOS o disco storico viene incluso nel repo.
- Il progetto e' CP/M-like: non dichiarare compatibilita' CP/M completa.

## Verifica

```powershell
go test -count=1 ./...
go vet ./...
go run ./cmd/retronet-cpm -conformance
git diff --check
```
