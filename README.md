# Timeterm

[![Pipeline Status](https://gitlab.com/timeterm/timeterm/badges/master/pipeline.svg)](https://gitlab.com/timeterm/timeterm/-/pipelines)

Timeterm is een embedded roostersysteem voor scholen, wat op het moment alleen streeft het roosterysteem Zermelo te ondersteunen.

Timeterm is open-source, maar niet alle delen vallen onder dezelfde licentie. Zie de verschillende mappen voor de verschillende licenties die gebruikt zijn. De root van de repository valt onder Apache 2.

© 2020 Kees Blok, Robert van der Maas en Rutger Broekhoff

## Oplevering

Het project is woensdag 16 december 2020 opgeleverd voor BM3a. Er is toen een tag met de naam 'bm3a' aangemaakt.
Klik [hier](https://gitlab.com/timeterm/timeterm/-/tree/bm3a) om de broncode op dat moment in GitLab te bekijken.

Het project is donderdag 11 februari 2020 opgeleberd voor BM3b (12 februari 2020). Er is toen een tag met de naam 'bm3b' aangemaakt.
Klik [hier](https://gitlab.com/timeterm/timeterm/-/tree/bm3b) om de broncode op dat moment in GitLab te bekijken.

## Projectstructuur

`/`:
- `api/` - API-documentatie en verdere informatie (over bijv. de Zermelo API)
  - `reference/` - OpenAPI-documentatie van de Timeterm API, gehost op https://api.timeterm.nl/reference/
- `backend/` - Broncode voor het backend van Timeterm (Go), gehost op https://api.timeterm.nl/
- `design/` - Designbestanden van de UI en het apparaat
- `docs/` - Praktische documentatie van Timeterm, gehost op https://docs.timeterm.nl/
- `frontend-admin-web/` - Broncode voor het administratorsfrontend (React, TypeScript & JavaScript), gehost op https://admin.timeterm.nl/
- `frontend-embedded-devtools/` - Broncode voor ontwikkelingstools voor frontend-embedded (nepkaartlezer om lokaal in te loggen) (Qt, C++ & JavaScript)
- `frontend-embedded/` - Broncode voor het frontend wat op een embedded Timeterm-apparaat draait (Qt, C++ & JavaScript)
- `mfrc522/` - Broncode voor het gebruik van een MFRC522 op een Raspberry Pi 4 (maakt gebruik van het Linux GPIO character device via [libgpiod](https://git.kernel.org/pub/scm/libs/libgpiod/libgpiod.git/)) (C++)
- `nats-manager/` - Broncode voor het programma wat NATS (JetStream) streams en AuthZ beheert (en een NATS account server implementeert) met JWTs (Go)
- `oci-images/` - Docker images (build-time en overigen)
- `os/` - Broncode voor Timeterm OS (wat draait op het embedded apparaat), gebaseerd op Boot2Qt (BitBake)
- `proto/` - Beschrijvingen van berichten die gebruik maken van (Google) Protocol Buffers

Het project maakt gebruik van GitLab CI. Alleen de gewijzigde componenten worden gebouwd.

## Het backend opstarten (development)

De eerste keer opstarten vereist wat meer stappen dan de keren daarop.
We gaan er vanuit dat je toegang hebt tot [Docker Compose](https://docs.docker.com/compose/) en de [Vault](https://www.vaultproject.io/) CLI.

1. `timeterm $ cd backend`
2. `backend $ docker-compose up -d vault postgres`
3. `backend $ vault operator init -address http://localhost:8300`  
   Deze stap kan weggelaten worden wanneer Vault al eens eerder is opgestart. 
	 Sla de initial root token en de unseal keys op een veilige plek op (en zorg vooral dat je ze niet vergeet).
4. `backend $ vault operator unseal -address http://localhost:8300`  
   In het geval van deze setup moet dit commando 3x uitgevoerd worden (de helft + 1 (quorum) van de aangemaakte unseal keys moet geleverd worden). 
5. `backend $ cd ../nats-manager`
6. `nats-manager $ ./nats-manager`  
   Voor deze stap moet je nats-manager al gebouwd hebben en je omgevingsvariabelen moeten ook juist ingesteld zijn.
	 nats-manager en het backend laden automatisch omgevingsvariabelen uit het bestand `.env`.
	 Er is een bestand `.env.example` in de map van nats-manager toegevoegd waarin voorbeeldwaarden voor de vereiste
	 omgevingsvariabelen staan.
7. (in een andere terminalsessie/venster) `timeterm $ cd backend`
8. `backend $ docker-compose up -d nats`
9. `backend $ ./backend`  
   Voor deze stap geldt ook dat de omgevingsvariabelen juist ingesteld moeten zijn. Dit kan hier ook met een `.env`-bestand gedaan worden.

