# Release v0.2.0

Novita':

- test end-to-end locale `retronet-asm -> .COM -> retronet-cpm`
- subset BDOS file read-only: `open`, `close`, `read sequential`, `set DMA`
- esempi didattici CP/M-like: echo input, mini menu, TYPE-like
- documentazione BDOS aggiornata con FCB e DMA

Limiti:

- solo lettura sequenziale, niente write/delete/rename
- nessuna immagine disco storica
- gli esempi `.COM` usano ancora indirizzi assoluti calcolati per `0100h`
