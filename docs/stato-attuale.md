# Stato Attuale

Questo documento riassume cosa esiste oggi in `retronet-cpm` e come si collega
al resto dell'ecosistema RetroNet.

## Versioni Pubblicate

- `retronet-8080 v0.1.1`: emulatore Intel 8080 con ALU `gate`/`native`,
  validazione diagnostica 8080EXM e API importabile.
- `retronet-cpm v0.1.0`: prima macchina CP/M-like pubblica, con loader `.COM`,
  BDOS console subset, shell e drive host read-only.
- `retronet-cpm v0.2.0`: aggiunge integrazione assembler, funzioni BDOS file
  read-only e nuovi esempi didattici.

## Cosa Implementa `retronet-cpm`

- Loader `.COM`: carica il programma a `0100h`.
- Pagina zero didattica:
  - `0000h`: warm boot, usato come fine programma.
  - `0005h`: vettore BDOS, intercettato dal runtime.
  - `F000h`: trap interno documentale del BDOS.
- CPU: usa `github.com/retronet-labs/retronet-8080/cpu`.
- ALU: default `cpu.Native`, piu' veloce per programmi lunghi; `cpu.Gate`
  resta selezionabile con `-alu gate`.
- Stack iniziale: `EFFEh`.
- Command tail: `0080h`, inizializzata da `RUN <programma> [argomenti]`.
- FCB default: `005Ch` e `006Ch` dai primi due argomenti 8.3 della shell.
- Drive `A:`: directory host read-only con nomi CP/M 8.3.
- Shell `A>`: `DIR`, `TYPE`, `RUN`, `HELP`, `EXIT`.
- Terminale: programmi `.COM` e shell `RUN` usano `retronet-terminal` come
  console condivisa.
- Sessioni: package `session` per creare sessioni API-ready senza invocare la
  CLI.
- Live: `cmd/retronet-cpm-live` collega `session` a `retronet-terminal/live` per
  una shell `A>` locale interattiva.
- Trace: testuale con `-trace`, JSON Lines con `-trace-json`.

## BDOS Supportato

Console:

- `0`: terminate
- `1`: console input con echo
- `2`: console output
- `6`: direct console I/O
- `9`: print string terminata da `$`
- `10`: buffered input
- `11`: console status
- `12`: version stub `0022h`

File read-only:

- `15`: open file da FCB
- `16`: close file
- `20`: read sequential, record da 128 byte nel DMA
- `26`: set DMA

File mutanti controllati:

- `19`: delete file
- `21`: write sequential
- `22`: make file
- `23`: rename file

Queste funzioni modificano il filesystem solo quando la CLI apre il drive con
`-write-disk`.

Le funzioni non implementate restituiscono errore esplicito. Il progetto non
finge un BDOS completo.

## Test E Conformance

La suite locale copre:

- core CP/M-like, loader e warm boot
- default ALU `native` e override `gate`
- BDOS console
- FCB read-only e DMA
- drive host 8.3 e blocco path traversal
- shell `DIR`, `TYPE`, `RUN`, `EXIT`
- command tail e FCB default sintetici
- CLI e `RETRONET_CPM_ALU`
- conformance sintetica
- terminale condiviso, command tail, FCB default e write opt-in
- sicurezza drive host con limiti configurabili
- test end-to-end locale `retronet-asm -> .COM -> retronet-cpm`

Gate usati:

```powershell
go test -count=1 ./...
go vet ./...
go run ./cmd/retronet-cpm -conformance
git diff --check
```

## Cosa Non E' Ancora CP/M

- Nessun BDOS/BIOS storico incluso.
- Nessuna immagine disco storica.
- Nessuna scrittura file BDOS senza opt-in esplicito `-write-disk`.
- Nessun supporto user area.
- Nessun parsing completo di wildcard o opzioni CCP storiche.
- Nessuna emulazione BIOS o periferiche S-100.
- Gli esempi assembly usano `.com`/`.orgbase` in `retronet-asm` per label
  logiche a partire da `0100h`, senza indirizzi runtime manuali.

La promessa attuale e' "CP/M-like didattico", non compatibilita' CP/M completa.
