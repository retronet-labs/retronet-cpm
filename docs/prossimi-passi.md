# Prossimi Passi

Questo e' il piano tecnico consigliato dopo `retronet-cpm v0.2.0`.

## 1. Origin Logico .COM In `retronet-asm`

Stato: implementato con `.orgbase <addr>` e alias `.com`.

La direttiva cambia la base logica delle label senza emettere padding nel file.

Possibili forme:

```asm
.arch i8080
.orgbase 0x0100
```

oppure:

```asm
.arch i8080
.com
```

Risultato:

- `label:` nel sorgente vale `0100h + offset`.
- il file `.COM` emesso inizia dal primo byte reale, senza 256 byte iniziali.
- gli esempi CP/M usano label al posto di indirizzi manuali.
- test end-to-end assembler -> CP/M aggiornato.

## 2. BDOS File In Scrittura

Stato: implementato con opt-in `-write-disk`.

- `22`: make file
- `21`: write sequential
- `19`: delete file
- `23`: rename file

Policy attuale:

- default ancora read-only per sicurezza.
- flag CLI `-write-disk` per abilitare scrittura host.
- test con directory temporanee, mai su path reali dell'utente.

Criteri di successo:

- un programma `.COM` crea un file, scrive record da 128 byte e lo rilegge.
- path traversal resta impossibile.
- errori BDOS sono espliciti e documentati.

## 3. Libreria Assembly CP/M-like

Stato: aggiunto `examples/lib/cpm-bdos.asm` come blocco copiabile con costanti
comuni:

```asm
BDOS  = 0x0005
PRINT = 9
TERM  = 0
DMA   = 0x0080
```

`retronet-asm` supporta `.include`, quindi gli esempi possono includere il file:

```asm
.include "lib/cpm-bdos.asm"
```

Criteri di successo:

- esempi piu' corti e leggibili.
- meno byte literal nei programmi didattici.
- pattern standard per print, input, open/read e terminate.

## 4. `retronet-terminal`

Stato: implementato come repo separato `retronet-terminal` e integrato nella CLI
`-run` tramite adattatore BDOS.

Responsabilita' separate dal runtime CP/M:

- buffer output
- input queue
- schermo testuale derivato
- adattatore CLI
- futuro adattatore websocket

Perche' viene prima del web lab: CP/M, BBS ed emulatori avranno tutti bisogno
dello stesso concetto di terminale.

Criteri di successo:

- `retronet-cpm` usa `retronet-terminal` come console per `-run`.
- test di input/output indipendenti dalle CPU nel repo terminale.
- API byte-oriented pronta per un futuro websocket.

## 5. Demo Unica Documentata

Preparare un transcript ripetibile:

```text
A>DIR
A>TYPE DOLLAR.TXT
A>RUN HELLO
A>RUN MINI
A>RUN TYPE
A>EXIT
```

Artefatti:

- script PowerShell che assembla gli esempi in una directory temporanea.
- file `DOLLAR.TXT` generato automaticamente.
- transcript in `docs/demo.md`.

Criteri di successo:

- clone dei repo sibling.
- un comando prepara i `.COM`.
- un comando lancia la demo.
- output documentato e stabile.

## 6. Compatibilita' CP/M Incrementale

Stato: primo incremento implementato in forma sintetica e documentata.

- command tail a `0080h`.
- FCB default a `005Ch` e `006Ch` dalla shell `RUN`.

Restano aperti:

- aggiungere funzioni BDOS comuni usate da programmi piccoli.
- valutare diagnostiche CP/M reali solo se licenza/provenienza sono gestite
  fuori dal repo.

Regola guida: ogni funzione CP/M va aggiunta con test sintetico, documentazione
e dichiarazione chiara dei limiti.
