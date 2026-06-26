# BDOS Console Subset

La v0.1 implementa solo le funzioni console necessarie per esempi e programmi
diagnostici piccoli. Non implementa FCB, file BDOS, user area o dischi CP/M
storici.

| C | Nome | Comportamento |
|---|------|---------------|
| `0` | terminate | termina il programma |
| `1` | console input | legge un byte, lo restituisce in `A/L` e lo ecoa |
| `2` | console output | scrive il byte in `E` |
| `6` | direct console I/O | `E=FFh` legge se disponibile, altrimenti `A=0`; altri valori scrivono `E` |
| `9` | print string | stampa da `DE` fino a `$` |
| `10` | buffered input | usa il buffer CP/M-like `max,count,data...` |
| `11` | console status | `A/L=FFh` se input disponibile, altrimenti `0` |
| `12` | version | restituisce `0022h` in `HL` come stub CP/M 2.2 didattico |
| `15` | open file | apre un FCB read-only dal drive `A:` |
| `16` | close file | chiude l'FCB aperto |
| `19` | delete file | cancella un file, solo se il drive e' scrivibile |
| `20` | read sequential | legge record da 128 byte nel DMA |
| `21` | write sequential | scrive record da 128 byte dal DMA, solo se il drive e' scrivibile |
| `22` | make file | crea un file e apre l'FCB in scrittura |
| `23` | rename file | rinomina un file usando il nuovo nome a offset `16` dell'FCB |
| `26` | set DMA | imposta l'indirizzo DMA da `DE` |

Le funzioni non supportate restituiscono `ErrUnsupportedFunction`. Questa scelta
e' intenzionale: e' meglio fallire in modo chiaro che fingere un BDOS completo.

## FCB Read-Only

Il subset file usa l'FCB CP/M standard: drive a offset `0`, nome a `1..8`,
estensione a `9..11` e record corrente `CR` a offset `32`.

`open` carica il file in memoria host del runtime, `set DMA` imposta il buffer
di trasferimento e `read sequential` copia record da 128 byte nel DMA corrente.
L'ultimo record viene riempito con `1Ah`, il carattere EOF testuale usato spesso
in ambiente CP/M.

## Scrittura Controllata

Le funzioni mutanti (`delete`, `write sequential`, `make`, `rename`) funzionano
solo se il drive implementa `disk.MutableDrive`. La CLI lo abilita solo con:

```powershell
go run ./cmd/retronet-cpm -disk C:\tmp\cpm -write-disk
```

Senza `-write-disk`, le funzioni BDOS mutanti restituiscono fallimento (`A=FFh`)
e non modificano il filesystem host.
