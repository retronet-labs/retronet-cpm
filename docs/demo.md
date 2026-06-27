# Demo Ripetibile

Questa demo usa solo sorgenti assembly originali del repo, file testuali generati
localmente e il drive host temporaneo. Non include ROM, BDOS, BIOS o dischi
storici.

## Esecuzione

Da `C:\work\source\retronet-cpm`:

```powershell
.\scripts\demo-cpm.ps1
```

Lo script:

- crea `.gocache\demo-drive`
- genera `DOLLAR.TXT`
- assembla `HELLO.COM`, `ECHO.COM`, `MINI.COM`, `TYPE.COM` usando il repo sibling
  `retronet-asm`
- lancia una sessione shell con input predefinito

## Transcript Atteso

Il numero di byte puo' cambiare se gli esempi vengono estesi, ma la forma resta:

La shell non ecoa i comandi ricevuti da `-input`, quindi il transcript mostra
prompt e output:

```text
A>DOLLAR.TXT       11
ECHO.COM         ...
HELLO.COM        ...
MINI.COM         ...
TYPE.COM         ...
A>CIAO DA A:$
A>HI
[bdos-terminate ...]
A>1) HELLO
2) BYE
? 2BYE
[bdos-terminate ...]
A>CIAO DA A:
[bdos-terminate ...]
A>
```

`RUN TYPE DOLLAR.TXT` dimostra command tail e FCB default sintetici: la shell
scrive l'argomento nella pagina zero e inizializza il primo FCB a `005Ch`.
