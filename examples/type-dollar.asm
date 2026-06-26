.arch i8080

; TYPE-like minimale: apre DOLLAR.TXT dal drive A:, legge un record nel DMA
; 0200h e lo stampa come stringa BDOS 9. Il file deve contenere un '$' nel primo
; record, per esempio: CIAO DA A:$.
;
; Indirizzi calcolati per caricamento .COM a 0100h:
;   FCB = 0125h
;   DMA = 0200h

.equ BDOS 0x0005
.equ FCB 0x0125
.equ DMA 0x0200

        LXI D, FCB
        MVI C, 15
        CALL BDOS
        LXI D, DMA
        MVI C, 26
        CALL BDOS
        LXI D, FCB
        MVI C, 20
        CALL BDOS
        LXI D, DMA
        MVI C, 9
        CALL BDOS
        MVI C, 0
        CALL BDOS

        .byte 0x00
        .byte 0x44, 0x4F, 0x4C, 0x4C, 0x41, 0x52, 0x20, 0x20
        .byte 0x54, 0x58, 0x54
        .byte 0x00, 0x00, 0x00, 0x00
        .byte 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00
        .byte 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00
        .byte 0x00
