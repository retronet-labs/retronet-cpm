.arch i8080
.com

; Crea OUT.TXT e scrive un record sequenziale dal buffer "SAVED$".
; Richiede retronet-cpm -write-disk.

.equ BDOS 0x0005
.equ BDOS_TERM 0
.equ BDOS_MAKE 22
.equ BDOS_WRITESEQ 21
.equ BDOS_CLOSE 16
.equ BDOS_SETDMA 26

        LXI D, fcb
        MVI C, BDOS_MAKE
        CALL BDOS
        LXI D, record
        MVI C, BDOS_SETDMA
        CALL BDOS
        LXI D, fcb
        MVI C, BDOS_WRITESEQ
        CALL BDOS
        LXI D, fcb
        MVI C, BDOS_CLOSE
        CALL BDOS
        MVI C, BDOS_TERM
        CALL BDOS

fcb:    .byte 0x00
        .byte 0x4F, 0x55, 0x54, 0x20, 0x20, 0x20, 0x20, 0x20
        .byte 0x54, 0x58, 0x54
        .byte 0x00, 0x00, 0x00, 0x00
        .byte 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00
        .byte 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00
        .byte 0x00

record: .byte 0x53, 0x41, 0x56, 0x45, 0x44, 0x24
