# Prossimi Passi

Questo e' il piano tecnico consigliato dopo `retronet-cpm v0.2.0`.

## 1. Origin Logico .COM In `retronet-asm`

Problema attuale: gli esempi `.COM` devono usare indirizzi assoluti calcolati per
runtime `0100h`, ma l'assembler calcola le label da `0`.

Obiettivo: aggiungere una direttiva che cambi la base logica delle label senza
emettere padding nel file.

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

Criteri di successo:

- `label:` nel sorgente vale `0100h + offset`.
- il file `.COM` emesso inizia dal primo byte reale, senza 256 byte iniziali.
- gli esempi CP/M non hanno piu' indirizzi manuali nei commenti.
- test end-to-end assembler -> CP/M aggiornato.

## 2. BDOS File In Scrittura

Estendere il subset file oltre la lettura:

- `22`: make file
- `21`: write sequential
- `19`: delete file
- `23`: rename file

Policy iniziale consigliata:

- default ancora read-only per sicurezza.
- flag CLI futuro `-write-disk` per abilitare scrittura host.
- test con directory temporanee, mai su path reali dell'utente.

Criteri di successo:

- un programma `.COM` crea un file, scrive record da 128 byte e lo rilegge.
- path traversal resta impossibile.
- errori BDOS sono espliciti e documentati.

## 3. Libreria Assembly CP/M-like

Creare esempi o include documentali con costanti comuni:

```asm
BDOS  = 0x0005
PRINT = 9
TERM  = 0
DMA   = 0x0080
```

Finche' `retronet-asm` non supporta include/macro, la libreria puo' vivere come
documentazione copiabile in `examples/README.md` o in `examples/lib/`.

Criteri di successo:

- esempi piu' corti e leggibili.
- meno byte literal nei programmi didattici.
- pattern standard per print, input, open/read e terminate.

## 4. `retronet-terminal`

Separare il terminale dal runtime CP/M:

- buffer output
- input queue
- modalita' raw/testuale
- adattatore CLI
- futuro adattatore websocket

Perche' viene prima del web lab: CP/M, BBS ed emulatori avranno tutti bisogno
dello stesso concetto di terminale.

Criteri di successo:

- `retronet-cpm` puo' usare `retronet-terminal` come console.
- test di input/output indipendenti dalle CPU.
- API pronta per websocket.

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

Dopo i passi sopra:

- supportare command tail a `0080h`.
- inizializzare FCB default a `005Ch` e `006Ch` dalla shell `RUN`.
- aggiungere funzioni BDOS comuni usate da programmi piccoli.
- valutare diagnostiche CP/M reali solo se licenza/provenienza sono gestite
  fuori dal repo.

Regola guida: ogni funzione CP/M va aggiunta con test sintetico, documentazione
e dichiarazione chiara dei limiti.
