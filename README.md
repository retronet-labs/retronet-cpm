# retronet-cpm - Ambiente CP/M-like per RetroNet

`retronet-cpm` e' il livello didattico che esegue programmi `.COM` CP/M-like
sopra l'emulatore Intel 8080 di RetroNet.

La v0.1 non e' una distribuzione CP/M storica: non include BDOS, BIOS, dischi o
ROM originali. Offre invece una macchina compatibile con le convenzioni minime
utili per esempi, diagnostiche e programmi console: caricamento a `0100h`,
trap `CALL 0005h`, shell `A>` e un drive `A:` mappato a una directory host.

## Stato iniziale

- CPU 8080 importata da `github.com/retronet-labs/retronet-8080`.
- Default operativo: ALU `native`, piu' veloce per programmi CP/M-like lunghi.
- ALU `gate` selezionabile per dimostrazioni didattiche.
- Documentazione pubblica in italiano.

## Quick start

```powershell
go test ./...
go run ./cmd/retronet-cpm -conformance
go run ./cmd/retronet-cpm -disk . 
go run ./cmd/retronet-cpm -run HELLO.COM
```

## Limiti v0.1

- Nessuna ROM o componente CP/M storico redistribuito.
- BDOS console subset, non BDOS completo.
- Drive `A:` host read-only, non immagine disco storica.
- Nessuna emulazione BIOS o periferiche S-100 reali.
