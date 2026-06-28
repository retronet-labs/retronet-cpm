# Terminale Live CP/M-like

`retronet-cpm-live` collega la sessione CP/M-like a `retronet-terminal/live`.
E' il primo modo locale per usare la shell `A>` come un terminale, invece che
come semplice stdin/stdout a righe.

## Avvio

```powershell
cd C:\work\source\retronet-cpm
go run ./cmd/retronet-cpm-live -disk .
```

Con dimensioni esplicite:

```powershell
go run ./cmd/retronet-cpm-live -disk . -width 100 -height 30
```

Modalita' scriptata, utile per test e demo ripetibili:

```powershell
go run ./cmd/retronet-cpm-live -disk . -script "DIR`rHELP`rEXIT`r"
```

## Uso

Nel terminale:

```text
A>DIR
A>TYPE README.TXT
A>RUN HELLO
A>HELP
A>EXIT
```

Tasti speciali:

- Backspace: cancella il carattere precedente nella riga corrente
- Invio: esegue la riga
- `Ctrl+L`: pulisce lo schermo e ristampa `A>`
- `Ctrl+Q`, `Ctrl+C`, `Ctrl+D`: esce

Se raw mode non e' disponibile:

```powershell
go run ./cmd/retronet-cpm-live -disk . -line
```

## Architettura

```text
tasti utente
    |
    v
retronet-terminal/live
    |
    v
cmd/retronet-cpm-live handler
    |
    v
session.Session -> shell -> BDOS -> retronet-terminal
```

Il package `live` non conosce CP/M. Sa solo leggere byte, gestire raw mode,
ridisegnare lo snapshot e inviare output a delta. Il comando CP/M fornisce
l'handler che accumula una riga e chiama `session.RunCommand`.

## Limiti

- e' CP/M-like, non CP/M storico completo
- non include BDOS, BIOS, ROM o dischi storici
- `RUN` e' sincrono: i programmi `.COM` che chiedono input durante l'esecuzione
  non hanno ancora un canale live pienamente asincrono
- per input programmato durante `RUN`, usare per ora la CLI classica con `-input`
  o il package `session`

Questo limite e' intenzionale per v0.5.0: il prossimo passo naturale e' portare
lo stesso modello in `retronet-api`, dove websocket e goroutine renderanno
possibile alimentare input anche mentre un programma e' in esecuzione.
