# Shell A>

La shell e' un front-end minimo per il drive `A:`:

- `DIR`: elenca i file host rappresentabili come nomi CP/M 8.3.
- `TYPE <file>`: stampa un file testuale dal drive.
- `RUN <programma[.COM]>`: carica il programma a `0100h` e lo esegue.
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
