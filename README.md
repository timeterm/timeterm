# Timeterm

Timeterm is een embedded roostersysteem voor scholen, wat op het moment alleen streeft Zermelo aan te spreken.
Het is een PWS-project wat nog in ontwikkeling is, het is nog niet klaar.

Timeterm is open-source, maar niet alle delen vallen onder dezelfde licentie. Zie de verschillende mappen voor de verschillende licenties die gebruikt zijn. De root van de repository valt onder Apache 2.

© 2020 Kees Blok, Robert van der Maas en Rutger Broekhoff

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
	 nats-manager en het backend laden automatisch omgevingsvariabelen het bestand `.env`.
	 Er is een bestand `.env.example` in de map van nats-manager toegevoegd waarin voorbeeldwaarden voor de vereiste
	 omgevingsvariabelen staan.
7. (in een andere terminalsessie/venster) `timeterm $ cd backend`
8. `backend $ docker-compose up -d nats`
9. `backend $ ./backend`  
   Voor deze stap geldt ook dat de omgevingsvariabelen juist ingesteld moeten zijn. Dit kan hier ook met `.env` gedaan worden.

