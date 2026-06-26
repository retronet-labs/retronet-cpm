.arch i8080

        LXI D, 0x010D    ; indirizzo runtime del messaggio nel .COM caricato a 0100h
        MVI C, 9         ; BDOS print string, terminata da '$'
        CALL 0x0005
        MVI C, 0         ; BDOS terminate
        CALL 0x0005
        .byte 0x48, 0x49, 0x24
