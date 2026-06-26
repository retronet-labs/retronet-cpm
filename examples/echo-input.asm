.arch i8080
.com

; Programma .COM CP/M-like: stampa un prompt, legge un carattere con BDOS 1
; e lo lascia ecoare dalla console.

.equ BDOS 0x0005

        LXI D, prompt
        MVI C, 9
        CALL BDOS
        MVI C, 1
        CALL BDOS
        LXI D, crlf
        MVI C, 9
        CALL BDOS
        MVI C, 0
        CALL BDOS

prompt: .byte 0x4B, 0x45, 0x59, 0x3F, 0x20, 0x24
crlf:   .byte 0x0D, 0x0A, 0x24
