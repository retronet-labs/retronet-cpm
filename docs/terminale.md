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

Da v0.5.0 esiste anche un terminale locale interattivo:

```powershell
go run ./cmd/retronet-cpm-live -disk C:\tmp\cpm
```

Questo comando usa il package `retronet-terminal/live`: raw mode, rendering
dello snapshot e output a delta vivono nel repo terminale, mentre `retronet-cpm`
fornisce solo l'handler che traduce i tasti in comandi della sessione CP/M-like.

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
- Il terminale live v0.5.0 esegue comandi shell in modo sincrono; i programmi
  `.COM` che chiedono input mentre sono in esecuzione richiedono il futuro run
  loop asincrono.
- Le sequenze ANSI supportate sono quelle documentate nel repo
  `retronet-terminal`.
