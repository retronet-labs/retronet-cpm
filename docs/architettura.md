# Architettura CP/M-like

`retronet-cpm` usa `retronet-8080` come core CPU importabile. La macchina carica
un programma `.COM` nella TPA a `0100h`, prepara una pagina zero didattica e
intercetta l'ingresso BDOS convenzionale `CALL 0005h`.

La scelta predefinita dell'ALU e' `cpu.Native`, perche' un ambiente CP/M-like
esegue programmi piu' lunghi rispetto ai test didattici singoli. `cpu.Gate`
resta selezionabile dalla CLI per mostrare lo stesso programma calcolato tramite
il datapath a porte logiche.

## Mappa base

- `0000h`: warm boot didattico. Se il programma salta qui, il run termina.
- `0005h`: vettore BDOS. Il core intercetta `PC==0005h` e gestisce la funzione
  indicata da `C`.
- `0100h`: inizio TPA e indirizzo di caricamento `.COM`.
- `F000h`: trap interno BDOS, usato solo come indirizzo documentale nella pagina
  zero.
- `EFFEh`: stack iniziale scelto da RetroNet per lasciare spazio ai programmi.

Il runner non esegue un BDOS storico in memoria: quando la CPU arriva al vettore
BDOS, il package `bdos` legge registri e memoria, applica la funzione supportata
e simula il `RET` prelevando l'indirizzo dallo stack 8080.

## ALU

`retronet-8080` mantiene `cpu.Gate` come default del core per coerenza didattica.
`retronet-cpm` invece crea la CPU con `cpu.Native`, salvo richiesta esplicita
`-alu gate`. I due backend sono equivalenti bit-per-bit secondo il test
differenziale del repo 8080; qui la scelta di default privilegia velocita' e
usabilita' della shell.
