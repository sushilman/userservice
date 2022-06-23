# FACEIT User Service
Demo user service

## How to start the application
- docker-compose up

## Avilable endpoints
- GET /healthz
- POST /v1/users
- GET /v1/users
- GET /v1/users/:userId
- PUT /v1/users/:userId
- DELETE /v1/users/:userId

## Run the unit tests
- `go test ./...`

## Assumptions Made
- Assuming that the other services that are interested in the changes in the user entities are internal, message broker system like Kafka/RabbitMQ/SQS can be used to notify them

## Possible Extension/Improvements
- If there are 3rd party external services that should be notified, then notifying them via webhook would be sensible. For that, we would also need a subscription management system, where the interested parties/clients can subscribe to the events that we expose. In the subscription, the subscriber should provide their REST API uri (along with credentials) that should be invoked when an event occurs.
- Make use of a proper logger framework like - zerolog
- Page numbering system can be improved (using cursor instead of pagesize and number)
