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

Le funzioni non supportate restituiscono `ErrUnsupportedFunction`. Questa scelta
e' intenzionale: e' meglio fallire in modo chiaro che fingere un BDOS completo.
