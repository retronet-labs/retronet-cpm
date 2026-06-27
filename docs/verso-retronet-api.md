# Verso `retronet-api`

Questo documento definisce il ponte tecnico prima di creare `retronet-api`.
L'obiettivo e' evitare che il futuro server HTTP/websocket reinventi shell,
terminale o accesso disco.

## Package Da Usare

`retronet-api` dovrebbe importare:

- `github.com/retronet-labs/retronet-cpm/session`
- `github.com/retronet-labs/retronet-cpm/disk`
- `github.com/retronet-labs/retronet-terminal`

Non dovrebbe:

- eseguire la CLI `retronet-cpm`
- leggere/scrivere path host forniti dal client
- includere ROM, BDOS, BIOS o immagini disco storiche
- duplicare parser terminale o logica ANSI

## Flusso Sessione

```text
HTTP create session
    |
    v
disk.NewTemporaryHostDrive(...)
    |
    v
session.New(...)
    |
    +--> websocket input  -> Session.Input(...)
    +--> websocket command -> Session.RunCommand(...)
    +--> websocket output -> Session.DrainOutput()
    +--> websocket sync   -> Session.Snapshot()
```

## Eventi Websocket Suggeriti

Input da browser:

```json
{"type":"input","data":"A"}
```

Comando shell:

```json
{"type":"command","line":"DIR"}
```

Output incrementale:

```json
{"type":"output","data":"A>DIR\r\nHELLO.COM\r\n"}
```

Snapshot:

```json
{
  "type": "snapshot",
  "width": 80,
  "height": 24,
  "cursorRow": 2,
  "cursorCol": 0,
  "rows": ["A>DIR                                                                           "]
}
```

Il formato definitivo vivra' in `retronet-api`, ma questi campi riflettono il
contratto gia' esposto da `retronet-terminal`.

## Sicurezza

- una directory temporanea per sessione
- `MaxFileSize` e `MaxFiles` sempre impostati
- nessun path assoluto dal client
- scrittura opt-in e limitata
- cleanup esplicito
- timeout/step limit sempre presenti

## Cosa Resta Fuori

- autenticazione utenti
- persistenza sessioni
- upload/download file via HTTP
- multiplexing multi-emulatore
- xterm.js o UI browser

Queste parti appartengono a `retronet-api` e `retronet-ui`, non a
`retronet-cpm`.
