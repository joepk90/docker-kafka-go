# Event Publisher
The script is used to publish events from outside the `docker-compose` network which then get picked up by the `docker-kafka-go` event handler.

To publish an event, first make sure the `docker-compose` services are running. Then run the following command:
```
make dev-produce-event
```