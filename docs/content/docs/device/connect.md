---
title: "Koppelen"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: true
---

Hier staat gedocumenteerd hoe een apparaat gekoppeld kan worden.
{{< hint info >}}
   Deze stap moet uitgevoerd worden **nadat** de netwerken 
   en de andere dingen correct zijn ingesteld!
{{< /hint >}}
1. Ga naar Apparaten.
2. Klik op Configuratie exporteren. Er wordt nu een bestand gedownload met de naam  
   timeterm-config.json. Dit bestand mag **niet** hernoemd worden.
    ![Configuratie exporteren](/export-config.png)
3. Als het bestand niet de naam **timeterm-config.json** heeft, moet de naam hiernaar gewijzigd worden.
   Toevoegingen zoals (1), (2) enzovoorts zijn niet toegestaan! Als de naam niet correct is, zal
   de koppeling niet werken.
4. Formateer een USB-stick naar FAT32-bestandssysteem.
5. Plaats het bestand **timeterm-config.json** op de geformatteerde USB.
6. Steek de USB-stick in een USB-poort van een apparaat.
7. Zet het apparaat aan.
8. Het apparaat stelt zichzelf nu verder in en verbindt vanzelf met het netwerk.
   U kunt het WiFi-icoontje in de gaten houden om te kijken of het apparaat al verbonden is.
9. Het apparaat voegt zichzelf aan de database toe. De naam wijzigen kan door
op het veld naam te klikken. De wijziging wordt opgeslagen als naast het veld geklikt wordt.
