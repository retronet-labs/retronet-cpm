.arch i8080
.include "lib/cpm-bdos.asm"
.com

        LXI D, msg
        MVI C, BDOS_PRINT
        CALL BDOS
        MVI C, BDOS_TERM
        CALL BDOS
msg:    .byte 0x48, 0x49, 0x24
