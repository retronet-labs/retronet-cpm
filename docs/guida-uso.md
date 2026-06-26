# Guida D'Uso

Questa guida parte da un clone locale del repo.

```powershell
cd C:\work\source\retronet-cpm
```

## Verificare Che Tutto Funzioni

```powershell
$env:GOCACHE='C:\work\source\retronet-cpm\.gocache'
go test -count=1 ./...
go run ./cmd/retronet-cpm -conformance
```

Output atteso dalla conformance:

```text
PASS bdos-print-string ...
PASS bdos-console-output ...
PASS bdos-direct-console-input ...
PASS warm-boot ...
PASS unsupported-bdos ...
conformance passed=5 failed=0
```

## Avviare La Shell A>

La shell usa la directory indicata da `-disk` come drive `A:`.

```powershell
mkdir C:\tmp\cpm
Set-Content -NoNewline C:\tmp\cpm\README.TXT "CIAO DA A:"
go run ./cmd/retronet-cpm -disk C:\tmp\cpm
```

Dentro la shell:

```text
A>DIR
A>TYPE README.TXT
A>HELP
A>EXIT
```

`DIR` mostra solo file compatibili con nomi CP/M 8.3.

## Eseguire Un Programma .COM

Se hai gia' un file `.COM`:

```powershell
go run ./cmd/retronet-cpm -disk C:\tmp\cpm -run HELLO.COM
```

Se il file e' nella directory corrente puoi indicare anche un percorso:

```powershell
go run ./cmd/retronet-cpm -run .\HELLO.COM
```

Se il nome non ha estensione e non e' un path, la CLI prova `.COM`:

```powershell
go run ./cmd/retronet-cpm -disk C:\tmp\cpm -run HELLO
```

## Scegliere L'ALU

Default:

```powershell
go run ./cmd/retronet-cpm -run HELLO.COM -alu native
```

Dimostrazione didattica con ALU a porte:

```powershell
go run ./cmd/retronet-cpm -run HELLO.COM -alu gate
```

Default da variabile d'ambiente:

```powershell
$env:RETRONET_CPM_ALU='native'
```

`RETRONET_8080_ALU` non viene letto da `retronet-cpm`: e' riservato alle
diagnostiche del repo `retronet-8080`.

## Usare Input Non Interattivo

Per programmi che leggono dalla console:

```powershell
go run ./cmd/retronet-cpm -run ECHO.COM -input "X"
```

Lo stesso vale per la shell:

```powershell
go run ./cmd/retronet-cpm -disk C:\tmp\cpm -input "DIR`nEXIT`n"
```

## Trace

Trace testuale:

```powershell
go run ./cmd/retronet-cpm -run HELLO.COM -trace
```

Trace JSON Lines:

```powershell
go run ./cmd/retronet-cpm -run HELLO.COM -trace-json trace.jsonl
```

Il trace include istruzioni 8080 e chiamate BDOS.

## Scrivere Sul Drive Host

Per sicurezza il drive `A:` e' read-only. I programmi che usano BDOS `make`,
`write sequential`, `delete` o `rename` falliscono con `A=FFh` se non abiliti la
scrittura:

```powershell
go run ./cmd/retronet-cpm -disk C:\tmp\cpm -run WRITE.COM -write-disk
```

Usa `-write-disk` solo su directory temporanee o preparate apposta.

## Assemblare Gli Esempi

Con `retronet-asm` come repo sibling:

```powershell
cd C:\work\source\retronet-asm
go run ./cmd/retronet-asm build C:\work\source\retronet-cpm\examples\hello-bdos.asm -o C:\tmp\cpm\HELLO.COM
```

Poi:

```powershell
cd C:\work\source\retronet-cpm
go run ./cmd/retronet-cpm -disk C:\tmp\cpm -run HELLO
```

Output atteso:

```text
HI
program=HELLO.COM stop=bdos-terminate ...
```

## Provare `type-dollar.asm`

Assembla:

```powershell
cd C:\work\source\retronet-asm
go run ./cmd/retronet-asm build C:\work\source\retronet-cpm\examples\type-dollar.asm -o C:\tmp\cpm\TYPE.COM
```

Prepara il file letto via FCB:

```powershell
Set-Content -NoNewline C:\tmp\cpm\DOLLAR.TXT "CIAO DA A:$"
```

Esegui:

```powershell
cd C:\work\source\retronet-cpm
go run ./cmd/retronet-cpm -disk C:\tmp\cpm -run TYPE
```

Output atteso:

```text
CIAO DA A:
```

Il `$` e' terminatore BDOS `9`, quindi non viene stampato.
