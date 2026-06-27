param(
    [string]$Drive = ""
)

$ErrorActionPreference = "Stop"

$root = Resolve-Path (Join-Path $PSScriptRoot "..")
$asmRoot = Resolve-Path (Join-Path $root "..\retronet-asm")

if ($Drive -eq "") {
    $Drive = Join-Path $root ".gocache\demo-drive"
}

New-Item -ItemType Directory -Force $Drive | Out-Null
Get-ChildItem -Path $Drive -Filter "*.COM" | Remove-Item
Set-Content -NoNewline -Encoding ascii (Join-Path $Drive "DOLLAR.TXT") 'CIAO DA A:$'

$examples = @{
    "HELLO.COM" = "hello-bdos.asm"
    "ECHO.COM"  = "echo-input.asm"
    "MINI.COM"  = "mini-menu.asm"
    "TYPE.COM"  = "type-dollar.asm"
}

Push-Location $asmRoot
try {
    foreach ($outName in $examples.Keys) {
        $source = Join-Path $root ("examples\" + $examples[$outName])
        $output = Join-Path $Drive $outName
        go run .\cmd\retronet-asm build $source -o $output
    }
}
finally {
    Pop-Location
}

$input = "DIR`nTYPE DOLLAR.TXT`nRUN HELLO`nRUN MINI`n2RUN TYPE DOLLAR.TXT`nEXIT`n"

Push-Location $root
try {
    go run .\cmd\retronet-cpm -disk $Drive -input $input
}
finally {
    Pop-Location
}
