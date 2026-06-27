; Costanti CP/M-like RetroNet.
; Copia questo blocco all'inizio di un programma .COM finche' retronet-asm non
; supporta include o macro.

.equ BDOS 0x0005

.equ BDOS_TERM 0
.equ BDOS_CONIN 1
.equ BDOS_CONOUT 2
.equ BDOS_DIRECT_IO 6
.equ BDOS_PRINT 9
.equ BDOS_READLINE 10
.equ BDOS_STATUS 11
.equ BDOS_VERSION 12

.equ BDOS_OPEN 15
.equ BDOS_CLOSE 16
.equ BDOS_DELETE 19
.equ BDOS_READSEQ 20
.equ BDOS_WRITESEQ 21
.equ BDOS_MAKE 22
.equ BDOS_RENAME 23
.equ BDOS_SETDMA 26

.equ DEFAULT_DMA 0x0080
.equ DEFAULT_FCB1 0x005C
.equ DEFAULT_FCB2 0x006C
.equ COMMAND_TAIL 0x0080
