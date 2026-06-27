# Sessioni Programmatiche

Il package `session` e' il punto di aggancio pensato per `retronet-api`. Permette
di usare shell, terminale e drive CP/M-like senza avviare la CLI.

## Perche' Serve

La CLI e' comoda per una persona, ma un server websocket ha bisogno di un
contratto importabile:

- creare una sessione
- accodare input
- eseguire un comando shell
- leggere i byte prodotti
- leggere uno snapshot dello schermo

`session.Session` offre proprio questi passaggi sopra `retronet-terminal`.

## Esempio Minimo

```go
drive, cleanup, err := disk.NewTemporaryHostDrive(
    "retronet-cpm-web-",
    disk.HostDriveOptions{
        Writable:    true,
        MaxFileSize: 64 * 1024,
        MaxFiles:    64,
    },
)
if err != nil {
    return err
}
defer cleanup()

sess, err := session.New(session.Config{Drive: drive})
if err != nil {
    return err
}

_ = sess.Prompt()
_ = sess.RunCommand("DIR")

delta, _ := sess.DrainOutput()
snapshot, _ := sess.Snapshot()

_ = delta
_ = snapshot
```

## Input Per Programmi Interattivi

Se un programma `.COM` legge dalla console BDOS, l'input va accodato nel
terminale della sessione:

```go
_ = sess.Input([]byte("Z"))
_ = sess.RunCommand("RUN ECHO")
```

La shell e il programma condividono il terminale, ma la sessione non espone path
host o reader arbitrari al client.

## Output E Snapshot

- `DrainOutput()` restituisce i nuovi byte raw e svuota il buffer.
- `Snapshot()` restituisce dimensioni, righe, cursore, input pendente e byte raw
  in attesa.

Un websocket puo' inviare `DrainOutput()` come evento incrementale e usare
`Snapshot()` per riallineare il browser dopo riconnessioni o perdita eventi.

## Limiti

`session` non apre socket, non gestisce autenticazione e non crea policy multi
utente. Queste responsabilita' restano a `retronet-api`.
