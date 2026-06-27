# Sicurezza Del Drive Host

Il drive `A:` di `retronet-cpm` e' una directory host vista con nomi CP/M 8.3.
Questa scelta e' didattica e comoda, ma quando arrivera' `retronet-api` diventa
importante limitare con precisione cosa puo' essere letto o scritto.

## Regole Gia' Applicate

- I nomi sono normalizzati in maiuscolo 8.3.
- Path assoluti, separatori `/` e `\`, `:` e `..` sono rifiutati.
- La scrittura e' disabilitata di default.
- `-write-disk` e `HostDriveOptions{Writable: true}` sono opt-in espliciti.
- `RenameFile` non sovrascrive un file esistente.

## Limiti Configurabili

`disk.HostDriveOptions` permette di impostare:

- `Writable`: abilita operazioni mutanti
- `MaxFileSize`: dimensione massima per read/write, `0` significa nessun limite
- `MaxFiles`: numero massimo di file CP/M visibili, `0` significa nessun limite

Esempio:

```go
drive, err := disk.NewHostDriveWithOptions(root, disk.HostDriveOptions{
    Writable:    true,
    MaxFileSize: 64 * 1024,
    MaxFiles:    64,
})
```

Per sessioni web e test isolati:

```go
drive, cleanup, err := disk.NewTemporaryHostDrive(
    "retronet-cpm-web-",
    disk.HostDriveOptions{Writable: true, MaxFileSize: 64 * 1024, MaxFiles: 64},
)
defer cleanup()
```

## Raccomandazione Per `retronet-api`

Ogni sessione web dovrebbe usare una directory temporanea dedicata, non una
directory scelta liberamente dal client. Il client puo' caricare file tramite una
API controllata, ma non deve mai passare path host arbitrari.

Policy consigliata:

- root temporanea per sessione
- scrittura disabilitata finche' non serve davvero
- limite file e dimensione sempre impostati
- cleanup a fine sessione
- nessuna esposizione di path assoluti nei messaggi websocket

## Licenze E Provenienza

Il drive non contiene ROM o dischi storici. Le demo generano file testuali locali
e assemblano programmi da sorgenti originali del repo.
