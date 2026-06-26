# Architettura CP/M-like

`retronet-cpm` usa `retronet-8080` come core CPU importabile. La macchina carica
un programma `.COM` nella TPA a `0100h`, prepara una pagina zero didattica e
intercetta l'ingresso BDOS convenzionale `CALL 0005h`.

La scelta predefinita dell'ALU e' `cpu.Native`, perche' un ambiente CP/M-like
esegue programmi piu' lunghi rispetto ai test didattici singoli. `cpu.Gate`
resta selezionabile dalla CLI per mostrare lo stesso programma calcolato tramite
il datapath a porte logiche.
