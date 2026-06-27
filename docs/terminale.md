# Terminale Condiviso

`retronet-cpm` usa `retronet-terminal` come console dei programmi `.COM`.
Questo separa il terminale dal runtime CP/M-like:

```text
programma .COM
    |
    | CALL 0005h
    v
BDOS console
    |
    v
retronet-terminal
```

## Cosa Cambia Per L'Utente

La CLI continua a stampare su stdout:

```powershell
go run ./cmd/retronet-cpm -run HELLO.COM
```

Internamente, pero', ogni byte scritto dal BDOS passa anche dal terminale
condiviso. Questo permette al futuro `retronet-api` di usare lo stesso output
per websocket o UI browser.

## Shell E Programmi Interattivi

La shell `A>` e i programmi lanciati con `RUN` condividono lo stesso reader. Se
uno script passa input non interattivo, un programma puo' leggere un carattere e
poi la shell continua a leggere il comando successivo.

Esempio:

```powershell
go run ./cmd/retronet-cpm -disk C:\tmp\cpm -input "RUN ECHO`nZEXIT`n"
```

Il programma riceve `Z`; dopo la sua terminazione la shell legge `EXIT`.

## Perche' Serve Prima Di `retronet-api`

Un websocket ha bisogno di due oggetti stabili:

- byte prodotti dal programma, da inviare al client
- snapshot dello schermo, per riallineare il client se perde eventi

Questi concetti ora vivono in `retronet-terminal`, non in `retronet-cpm`.

## Limiti

- Non e' un VT100 completo.
- Non include font, ROM o terminfo.
- Non interpreta programmi CP/M storici arbitrari.
- Le sequenze ANSI supportate sono quelle documentate nel repo
  `retronet-terminal`.
