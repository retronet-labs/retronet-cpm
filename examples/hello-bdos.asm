.arch i8080
.com

        LXI D, msg
        MVI C, 9         ; BDOS print string, terminata da '$'
        CALL 0x0005
        MVI C, 0         ; BDOS terminate
        CALL 0x0005
msg:    .byte 0x48, 0x49, 0x24
