# Roadmap retronet-cpm

## v0.1

- caricamento `.COM` a `0100h`
- trap BDOS `CALL 0005h`
- subset BDOS console
- drive host `A:` read-only
- shell `A>` minimale
- CLI e conformance sintetica

## v0.2

- funzioni BDOS file read-only: open, close, read sequential, set DMA
- esempi didattici: echo, mini menu, TYPE-like su file `$`-terminated
- test end-to-end assembler `i8080` -> `.COM` -> `retronet-cpm`

## v0.3

- origine logica `.COM` via `.com`/`.orgbase` in `retronet-asm`
- BDOS write opt-in con `-write-disk`
- libreria assembly didattica `examples/lib/cpm-bdos.asm`

## v0.4

- `retronet-terminal` come console condivisa
- `.include` negli esempi assembly
- command tail e FCB default sintetici
- demo ripetibile locale

## v0.4.1

- terminale condiviso anche nella shell `RUN`
- conformance estesa a terminale, tail, FCB e BDOS write
- documentazione didattica di compatibilita' CP/M-like

## v0.4.2

- package `session` per uso programmatico API-ready
- opzioni di sicurezza sul drive host: limiti file, limiti dimensione e drive
  temporaneo
- documentazione ponte verso `retronet-api`

## Dopo v0.4.2

- immagini disco solo con provenienza e licenza chiare
- profili BIOS didattici, non storici inventati
