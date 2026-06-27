.arch i8080
.include "lib/cpm-bdos.asm"
.com

; TYPE-like minimale: usa l'FCB default a 005Ch inizializzato dalla shell
; con RUN TYPE DOLLAR.TXT, legge un record nel DMA 0200h e lo stampa come
; stringa BDOS 9. Il file deve contenere un '$' nel primo record.

.equ DMA 0x0200

        LXI D, DEFAULT_FCB1
        MVI C, BDOS_OPEN
        CALL BDOS
        LXI D, DMA
        MVI C, BDOS_SETDMA
        CALL BDOS
        LXI D, DEFAULT_FCB1
        MVI C, BDOS_READSEQ
        CALL BDOS
        LXI D, DMA
        MVI C, BDOS_PRINT
        CALL BDOS
        MVI C, BDOS_TERM
        CALL BDOS
