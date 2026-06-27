.arch i8080
.include "lib/cpm-bdos.asm"
.com

; Mini menu CP/M-like. Legge un tasto con BDOS 1:
;   1 -> stampa HELLO
;   2 -> stampa BYE
;   altro -> stampa ?

        LXI D, menu
        MVI C, BDOS_PRINT
        CALL BDOS
        MVI C, BDOS_CONIN
        CALL BDOS
        CPI 0x31
        JZ one
        CPI 0x32
        JZ two
        LXI D, bad
        JMP print
one:    LXI D, msg1
        JMP print
two:    LXI D, msg2
        JMP print
print:  MVI C, BDOS_PRINT
        CALL BDOS
        MVI C, BDOS_TERM
        CALL BDOS

menu:   .byte 0x31, 0x29, 0x20, 0x48, 0x45, 0x4C, 0x4C, 0x4F
        .byte 0x0D, 0x0A
        .byte 0x32, 0x29, 0x20, 0x42, 0x59, 0x45
        .byte 0x0D, 0x0A
        .byte 0x3F, 0x20, 0x24
msg1:   .byte 0x48, 0x45, 0x4C, 0x4C, 0x4F, 0x0D, 0x0A, 0x24
msg2:   .byte 0x42, 0x59, 0x45, 0x0D, 0x0A, 0x24
bad:    .byte 0x3F, 0x0D, 0x0A, 0x24
