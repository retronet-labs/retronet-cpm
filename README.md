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
  `15`, `16`, `19`, `20`, `21`, `22`, `23`, `26`.
- Drive host `A:` read-only con nomi CP/M 8.3.
- Shell `A>` con `DIR`, `TYPE`, `RUN`, `HELP`, `EXIT`.
- Console `.COM` adattata a `retronet-terminal` per input/output condivisi con
  futuri websocket.
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
- `-write-disk`: abilita le funzioni BDOS mutanti sulla directory host.
- `-conformance`: suite sintetica integrata.

`RETRONET_CPM_ALU` puo' impostare il default della CLI. `RETRONET_8080_ALU`
resta intenzionalmente fuori da questo repo: serve alle diagnostiche del core
8080, non al runtime CP/M-like.

## Documentazione

- [Stato attuale](docs/stato-attuale.md): cosa e' stato implementato finora e
  quali repo sono coinvolti.
- [Guida d'uso](docs/guida-uso.md): comandi pratici per conformance, shell,
  programmi `.COM`, esempi assembly e trace.
- [Architettura](docs/architettura.md): macchina CP/M-like, pagina zero, BDOS e
  scelta ALU.
- [BDOS](docs/bdos.md): funzioni console e file read-only supportate.
- [Shell](docs/shell.md): comandi `A>`.
- [Demo ripetibile](docs/demo.md): transcript locale senza ROM o dischi storici.
- [Prossimi passi](docs/prossimi-passi.md): piano tecnico dopo v0.2.
- [Release v0.3.0](docs/release-v0.3.0.md): origine `.COM`, BDOS write opt-in
  e libreria assembly didattica.

## Limiti attuali

- Nessuna ROM o componente CP/M storico redistribuito.
- BDOS subset, non BDOS completo.
- Drive `A:` host read-only di default; scrittura solo con `-write-disk`.
- Nessuna emulazione BIOS o periferiche S-100 reali.
