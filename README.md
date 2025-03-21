# Docker Kafka Go
Example local Kafka setup using Docker and Golang, along with an integrated database.


## Kafka Configruation
The main Kafka configuration settings are in the following files:
- `docker-compose.yaml`
- `.env.kafka`
- `.env.compose`


## Usage
Start the Service, Kafka instance and the database:
```
make dev-start
```

Watch the services (Consumer) logs:
```
make dev-watch
```

In a seperate terminal window, publish an event using the external go script
```
make dev-produce-event
```

Watch the logs...




## To Do:
- Remove the `.env.kafka` dependancy from the `make` commands (put env vars somewhere else?)
- Pass `.env` file to the external script (`cmd/event-publisher`), but keep all environment variables in the same place (the same file).
- Fix Golang Air server reloading (likely a OSX/M2/Docker issue - see the `Dockerfile.compose` file)
- Setup Kafka TLS