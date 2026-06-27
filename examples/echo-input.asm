.arch i8080
.include "lib/cpm-bdos.asm"
.com

; Programma .COM CP/M-like: stampa un prompt, legge un carattere con BDOS 1
; e lo lascia ecoare dalla console.

        LXI D, prompt
        MVI C, BDOS_PRINT
        CALL BDOS
        MVI C, BDOS_CONIN
        CALL BDOS
        LXI D, crlf
        MVI C, BDOS_PRINT
        CALL BDOS
        MVI C, BDOS_TERM
        CALL BDOS

prompt: .byte 0x4B, 0x45, 0x59, 0x3F, 0x20, 0x24
crlf:   .byte 0x0D, 0x0A, 0x24
