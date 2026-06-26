.arch i8080

; Mini menu CP/M-like. Legge un tasto con BDOS 1:
;   1 -> stampa HELLO
;   2 -> stampa BYE
;   altro -> stampa ?
;
; Gli indirizzi sono assoluti per .COM caricato a 0100h.

.equ BDOS 0x0005
.equ ONE 0x011D
.equ TWO 0x0123
.equ PRINT 0x0129
.equ MENU 0x0133
.equ MSG1 0x0148
.equ MSG2 0x0150
.equ BAD 0x0156

        LXI D, MENU
        MVI C, 9
        CALL BDOS
        MVI C, 1
        CALL BDOS
        CPI 0x31
        JZ ONE
        CPI 0x32
        JZ TWO
        LXI D, BAD
        JMP PRINT
        LXI D, MSG1
        JMP PRINT
        LXI D, MSG2
        JMP PRINT
        MVI C, 9
        CALL BDOS
        MVI C, 0
        CALL BDOS

        .byte 0x31, 0x29, 0x20, 0x48, 0x45, 0x4C, 0x4C, 0x4F
        .byte 0x0D, 0x0A
        .byte 0x32, 0x29, 0x20, 0x42, 0x59, 0x45
        .byte 0x0D, 0x0A
        .byte 0x3F, 0x20, 0x24
        .byte 0x48, 0x45, 0x4C, 0x4C, 0x4F, 0x0D, 0x0A, 0x24
        .byte 0x42, 0x59, 0x45, 0x0D, 0x0A, 0x24
        .byte 0x3F, 0x0D, 0x0A, 0x24
