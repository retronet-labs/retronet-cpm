.arch i8080

; Programma .COM CP/M-like: stampa un prompt, legge un carattere con BDOS 1
; e lo lascia ecoare dalla console.
;
; Indirizzi calcolati per caricamento .COM a 0100h:
;   PROMPT = 011Ah
;   CRLF   = 0120h

.equ BDOS 0x0005
.equ PROMPT 0x011A
.equ CRLF 0x0120

        LXI D, PROMPT
        MVI C, 9
        CALL BDOS
        MVI C, 1
        CALL BDOS
        LXI D, CRLF
        MVI C, 9
        CALL BDOS
        MVI C, 0
        CALL BDOS

        .byte 0x4B, 0x45, 0x59, 0x3F, 0x20, 0x24
        .byte 0x0D, 0x0A, 0x24
