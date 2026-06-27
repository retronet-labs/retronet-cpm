# Compatibilita' CP/M-like

`retronet-cpm` non e' CP/M originale e non include componenti storici. E' un
ambiente didattico che riproduce alcune convenzioni utili per capire come un
programma `.COM` interagisce con una macchina 8080, una pagina zero, il BDOS e
un drive testuale.

## Implementato

| Area | Stato |
| --- | --- |
| Caricamento `.COM` | Programma caricato a `0100h` nella TPA. |
| Warm boot | Salto a `0000h` ferma il programma. |
| Vettore BDOS | `CALL 0005h` intercettato dal runtime. |
| Stack iniziale | `SP=EFFEh`. |
| Console BDOS | Funzioni `0`, `1`, `2`, `6`, `9`, `10`, `11`, `12`. |
| File FCB | `open`, `close`, `read sequential`, `set DMA`. |
| File mutanti | `delete`, `write sequential`, `make`, `rename`, solo con `-write-disk`. |
| Command tail | `0080h`, massimo 126 caratteri piu' CR. |
| FCB default | `005Ch` e `006Ch` dai primi due argomenti della shell `RUN`. |
| Shell | `DIR`, `TYPE`, `RUN`, `HELP`, `EXIT`. |
| Terminale | Console `.COM` adattata a `retronet-terminal`. |

## Sintetico, Non Storico

Questi comportamenti sono scelti per essere chiari e testabili, non per fingere
un sistema CP/M completo:

- il BDOS non e' codice 8080 in memoria: e' una trap Go sul vettore `0005h`
- il drive `A:` e' una directory host con nomi 8.3
- il default FCB supporta nomi 8.3 semplici, senza wildcard complete
- il terminale e' ASCII/ANSI generico, non un terminale storico specifico
- `RETRONET_CPM_ALU` sceglie il backend ALU della CPU, ma il default resta
  `native` per velocita'

## Fuori Scope

- ROM, BIOS, BDOS o CCP storici
- immagini disco storiche
- user area
- wildcard CP/M complete
- periferiche S-100 reali
- compatibilita' binaria generale con programmi CP/M arbitrari

## Esempio: Command Tail E FCB

Shell:

```text
A>RUN TYPE DOLLAR.TXT
```

Effetto prima di eseguire `TYPE.COM`:

- `0080h` contiene la lunghezza della tail e poi `DOLLAR.TXT\r`
- `005Ch` contiene il primo FCB: nome `DOLLAR`, estensione `TXT`
- `006Ch` viene preparato vuoto, perche' non c'e' un secondo argomento

Il programma `examples/type-dollar.asm` apre proprio `DEFAULT_FCB1`, quindi il
nome file non e' hardcoded nel binario.

## Come Verificarla

```powershell
go run ./cmd/retronet-cpm -conformance
```

La suite sintetica copre console, terminale, command tail, FCB default, BDOS
write opt-in e failure su drive read-only.
