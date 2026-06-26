# CLI

```powershell
go run ./cmd/retronet-cpm -conformance
go run ./cmd/retronet-cpm -disk C:\tmp\cpm
go run ./cmd/retronet-cpm -disk C:\tmp\cpm -run HELLO.COM
go run ./cmd/retronet-cpm -run HELLO.COM -alu gate -trace
```

Flag principali:

- `-disk`: directory host per il drive `A:`.
- `-run`: file `.COM` da eseguire; se manca l'estensione, la CLI prova `.COM`.
- `-steps`: limite massimo di istruzioni 8080.
- `-alu`: `native` per uso normale, `gate` per demo a porte logiche.
- `-input`: input testuale non interattivo.
- `-trace`: trace leggibile.
- `-trace-json`: trace JSON Lines.
- `-conformance`: suite sintetica.

Variabile d'ambiente:

```powershell
$env:RETRONET_CPM_ALU = "native"
```

`RETRONET_8080_ALU` non viene letto: quel nome resta riservato alle diagnostiche
del repo `retronet-8080`.
