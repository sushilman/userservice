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
$ go mod download
$ DB_URI=mongodb://localhost:27017/users go run main.go
```

## Avilable endpoints
- GET /healthz – a health check endpoint
- POST /v1/users – create a new user
- GET /v1/users – get all users (with filters, the filter parameters are case-sensitive)
- GET /v1/users/:userId – get a user by ID
- PUT /v1/users/:userId - update a user
- DELETE /v1/users/:userId - delete a user

### Validations for POST /v1/users and PUT /v1/users
A very basic field validation exists for POST and PUT operations
The following are the required fields while creating the user
- first_name
- last_name
- email
- password

### Allowed filter parameters for GET /v1/users
- first_name
- last_name
- nickname
- email
- offset – pagination parameters
- limit – 

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
- Password field is hashed before storing (during creation and updating)
- Password is not included in the response to GET requests
- Application loggings are done using `fmt.Printf()`, it is better to use a proper loggers like zerolog

## Possible Extension/Improvements
- Validations:
    - Payload field validations
    - Verfication if a user already exists - eg: a user with duplicate email should not be allowed
- If there are 3rd party external services that should be notified, then notifying them via webhook would be sensible. For that, we would also need a subscription management system, where the interested parties/clients can subscribe to the events that we expose. In the subscription, the subscriber should provide their REST API uri (along with credentials) that should be invoked when an event occurs.
- Make use of a proper logger framework like - zerolog
- Page numbering system can be improved (using cursor instead of pagesize and number)

## About the application structure:
- *apiusers* – is the API layer
- *db* – is the Storage(Database) layer
- *events* – simply holds the event structs for any changes on user entity (create, update, delete)
- *messagebroker* – dummy implementation of a message broker
- *mocks* – for unit/mock tests 
- *services* – service layer implementations
- *usererrors* - custom errors for the service
- *utils* – common utilities like: hashing a password