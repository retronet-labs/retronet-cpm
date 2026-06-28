# CLI

```powershell
go run ./cmd/retronet-cpm -conformance
go run ./cmd/retronet-cpm -disk C:\tmp\cpm
go run ./cmd/retronet-cpm -disk C:\tmp\cpm -run HELLO.COM
go run ./cmd/retronet-cpm -run HELLO.COM -alu gate -trace
go run ./cmd/retronet-cpm -disk C:\tmp\cpm -run WRITE.COM -write-disk
go run ./cmd/retronet-cpm-live -disk C:\tmp\cpm
```

Flag principali:

- `-disk`: directory host per il drive `A:`.
- `-run`: file `.COM` da eseguire; se manca l'estensione, la CLI prova `.COM`.
- `-steps`: limite massimo di istruzioni 8080.
- `-alu`: `native` per uso normale, `gate` per demo a porte logiche.
- `-input`: input testuale non interattivo.
- `-trace`: trace leggibile.
- `-trace-json`: trace JSON Lines.
- `-write-disk`: abilita funzioni BDOS che modificano la directory host.
- `-conformance`: suite sintetica.

Variabile d'ambiente:

```powershell
$env:RETRONET_CPM_ALU = "native"
```

`RETRONET_8080_ALU` non viene letto: quel nome resta riservato alle diagnostiche
del repo `retronet-8080`.

## Terminale Live

`retronet-cpm-live` usa lo stesso drive e lo stesso backend ALU della CLI
classica, ma presenta una shell `A>` su `retronet-terminal/live`:

```powershell
go run ./cmd/retronet-cpm-live -disk C:\tmp\cpm
go run ./cmd/retronet-cpm-live -disk C:\tmp\cpm -width 100 -height 30
go run ./cmd/retronet-cpm-live -disk C:\tmp\cpm -script "DIR`rTYPE README.TXT`rEXIT`r"
```

Flag aggiuntivi:

- `-width`: larghezza dello schermo live.
- `-height`: altezza dello schermo live.
- `-line`: fallback manuale a input di riga se raw mode non e' disponibile.
- `-script`: comandi live per test e CI.

Il live e' pensato per comandi shell (`DIR`, `TYPE`, `RUN`, `HELP`, `EXIT`).
I programmi `.COM` che richiedono input durante l'esecuzione restano un limite
del primo adattatore live: saranno gestiti meglio dal run loop asincrono previsto
per `retronet-api`.
