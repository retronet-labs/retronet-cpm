# retronet-cpm - Ambiente CP/M-like per RetroNet

`retronet-cpm` e' il livello didattico che esegue programmi `.COM` CP/M-like
sopra l'emulatore Intel 8080 di RetroNet.

La v0.1 non e' una distribuzione CP/M storica: non include BDOS, BIOS, dischi o
ROM originali. Offre invece una macchina compatibile con le convenzioni minime
utili per esempi, diagnostiche e programmi console: caricamento a `0100h`,
trap `CALL 0005h`, shell `A>` e un drive `A:` mappato a una directory host.

## Stato

- CPU 8080 importata da `github.com/retronet-labs/retronet-8080`.
- Default operativo: ALU `native`, piu' veloce per programmi CP/M-like lunghi.
- ALU `gate` selezionabile per dimostrazioni didattiche.
- Loader `.COM` a `0100h` e pagina zero con vettore BDOS `0005h`.
- BDOS subset: console `0`, `1`, `2`, `6`, `9`, `10`, `11`, `12`; file
  read-only `15`, `16`, `20`, `26`.
- Drive host `A:` read-only con nomi CP/M 8.3.
- Shell `A>` con `DIR`, `TYPE`, `RUN`, `HELP`, `EXIT`.
- Documentazione pubblica in italiano.

## Quick start

```powershell
go test ./...
go run ./cmd/retronet-cpm -conformance
go run ./cmd/retronet-cpm -disk .
go run ./cmd/retronet-cpm -run HELLO.COM
```

La shell parte se non passi `-run`:

```text
A>DIR
A>TYPE README.TXT
A>RUN HELLO
A>EXIT
```

## CLI

- `-disk <dir>`: directory host mappata come drive `A:`.
- `-run <file>`: esegue un `.COM` e termina.
- `-steps <n>`: limite di istruzioni 8080.
- `-alu native|gate`: backend ALU; default `native`.
- `-input <testo>`: input non interattivo per shell o programma.
- `-trace`: trace testuale.
- `-trace-json <file>`: trace JSON Lines.
- `-conformance`: suite sintetica integrata.

`RETRONET_CPM_ALU` puo' impostare il default della CLI. `RETRONET_8080_ALU`
resta intenzionalmente fuori da questo repo: serve alle diagnostiche del core
8080, non al runtime CP/M-like.

## Limiti v0.1

- Nessuna ROM o componente CP/M storico redistribuito.
- BDOS subset, non BDOS completo.
- Drive `A:` host read-only, non immagine disco storica.
- Nessuna emulazione BIOS o periferiche S-100 reali.
