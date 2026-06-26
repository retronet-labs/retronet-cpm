# Esempi

Questa cartella ospita esempi CP/M-like senza ROM storiche. I programmi `.COM`
generati localmente non sono versionati.

`hello-bdos.asm` stampa `HI` usando `CALL 0005h` e la funzione BDOS `9`.

Esempi disponibili:

- `hello-bdos.asm`: hello world con BDOS `9`.
- `echo-input.asm`: prompt, input BDOS `1`, eco e terminazione.
- `mini-menu.asm`: menu a due scelte con `CPI`/`JZ`.
- `type-dollar.asm`: lettura FCB read-only di `DOLLAR.TXT` e stampa dal DMA.
- `write-file.asm`: crea `OUT.TXT` e scrive un record; richiede `-write-disk`.
- `lib/cpm-bdos.asm`: costanti BDOS copiabili negli esempi.

Per provarli:

```powershell
go run ..\retronet-asm\cmd\retronet-asm build examples\hello-bdos.asm -o HELLO.COM
go run .\cmd\retronet-cpm -run HELLO.COM
```

Per `type-dollar.asm`, prepara un file `DOLLAR.TXT` nel drive:

```text
CIAO DA A:$
```

Per `write-file.asm`, usa esplicitamente il flag di scrittura:

```powershell
go run .\cmd\retronet-cpm -disk C:\tmp\cpm -run WRITE.COM -write-disk
```
