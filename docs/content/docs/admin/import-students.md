---
title: "Leerlingen importeren"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: true
---

Leerlingen kunnen zowel handmatig als automatisch worden toegevoegd.

# Handmatig toevoegen
1. Ga naar Leerlingen.
2. Klik op `Toevoegen`.
    ![Leerling importeren](/new-student.png)
3. Geef de zermelo-gebruiker een naam, bijvoorbeeld een leerlingcode.
    Dit kan door op de naam te klikken. Wanneer naast het veld geklikt wordt, slaat Timeterm het op.
    Voor wijzigingen geldt hetzelfde.
4. Klik op `Toewijzen`, vul de pascode in en klik nogmaals op `Toewijzen`.
    ![Pascode toewijzen](/assign-student-passcode.png)
5. Het overzicht ziet er nu als volgt uit:
    ![Leerlingenoverzicht](/students-overview.png)

# Automatisch toevoegen
Timeterm ondersteund hiervoor zelf geen mogelijkheden.
Om dit toch te doen moet het systeembeheer een script schrijven om 
de leerlingen zelf toe te voegen. Dit kan door te communiceren met de 
Timeterm API. De API-documentatie van Timeterm is [hier](https://api.timeterm.nl/reference) te vinden.