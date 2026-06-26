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

## Dopo v0.1

- immagini disco con provenienza chiara
- profili BIOS didattici
- origin logico `.COM` in `retronet-asm`, per evitare indirizzi manuali `0100h`
