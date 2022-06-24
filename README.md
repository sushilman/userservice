# FACEIT User Service
Demo user service

## How to start the application
### Docker
```bash
$ docker-compose up
```

### Run application natively
```bash
$ docker run -d -p 27017:27017 mongo:latest
$ DB_URI=mongodb://localhost:27017/users go run main.go
```

## Avilable endpoints
- GET /healthz
- POST /v1/users
- GET /v1/users
- GET /v1/users/:userId
- PUT /v1/users/:userId
- DELETE /v1/users/:userId

## Run the unit tests
```
$ go test ./...
```

## Testing the Endpoints
- You can import the Postman collection - `userservice.postman_collection.json` and try out the endpoints

## Assumptions Made
- Assuming that the other services that are interested in the changes in the user entities are internal, message broker system like Kafka/RabbitMQ/SQS can be used to notify them
- For this demo implementation, boilerplace code for kafka like message broker has been skipped, instead the message broker is emulated with a dummy _Publish()_ function
- Also skipped designing the OpenAPI/Swagger documentation

## Possible Extension/Improvements
- Validations:
    - Payload validations and sanitizations
    - Validations if a user already exists - eg: a user with duplicate email should not be allowed
- If there are 3rd party external services that should be notified, then notifying them via webhook would be sensible. For that, we would also need a subscription management system, where the interested parties/clients can subscribe to the events that we expose. In the subscription, the subscriber should provide their REST API uri (along with credentials) that should be invoked when an event occurs.
- Make use of a proper logger framework like - zerolog
- Page numbering system can be improved (using cursor instead of pagesize and number)

## About the application structure:
- *apiusers* package is the API layer
- *db* package is the Storage(Database) layer
- *events* package is 