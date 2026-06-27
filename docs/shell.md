# Shell A>

La shell e' un front-end minimo per il drive `A:`. I programmi avviati con
`RUN` usano la console condivisa basata su `retronet-terminal`.

- `DIR`: elenca i file host rappresentabili come nomi CP/M 8.3.
- `TYPE <file>`: stampa un file testuale dal drive.
- `RUN <programma[.COM]> [argomenti]`: carica il programma a `0100h`, prepara
  command tail/FCB default e lo esegue.
- `HELP`: mostra i comandi disponibili.
- `EXIT`: chiude la shell.

Il drive e' read-only. I nomi sono normalizzati in maiuscolo 8.3 e i percorsi
con separatori o `..` sono rifiutati per evitare traversal fuori dalla directory
host.

Esempio:

```text
A>DIR
HELLO.COM        16
README.TXT      120
A>RUN HELLO
HI
[bdos-terminate steps=4 bdos=2]
```

Con argomenti:

```text
A>RUN TYPE DOLLAR.TXT
```

La tail viene scritta a `0080h` e i primi due argomenti compatibili 8.3
inizializzano gli FCB default a `005Ch` e `006Ch`. Wildcard, user area e parsing
CP/M completo restano fuori scope.
