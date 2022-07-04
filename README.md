# User Service
Demo user service.
The service serves both REST API server and gRPC server.
- HTTP API is served in the port `:8080`
- gRPC is served in the port `50051`

## How to start the application
The application can be launched in two different ways:

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
## Available gRPC methods
 - `Check()` and `Watch()` – for healthcheck according to the grpc google.golang.org/grpc/health/grpc_health_v1 healthcheck standard
 - `CreateUser()` – create a new user
 - `GetUsers()` - fetch all users, filterable
 - `GetUserById()` – fetch a specific user by ID
 - `UpdateUser()` – update a specific user (identified by ID)
 - `DeleteUser()` – delete a specific user (identified by ID)

## Available endpoints REST API
 - GET /healthz – a health check endpoint
 - POST /v1/users – create a new user
 - GET /v1/users – get all users (with filters, the filter parameters are case-sensitive)
 - GET /v1/users/:userId – get a user by ID
 - PUT /v1/users/:userId - update a user
 - DELETE /v1/users/:userId - delete a user

### Validations for POST /v1/users and PUT /v1/users
A very basic field validation exists for POST and PUT operations
The following are the `required` fields while creating the user
 - first_name
 - last_name
 - email
 - password

### Allowed filter parameters for GET /v1/users
The filter parameters are case sensitive
 - first_name
 - last_name
 - nickname
 - email
 - offset – pagination parameters
 - limit – maximum number of items per request

## Run the unit tests
```
$ go test ./...
```

## Integration tests
- Integration tests have been skipped for this demo as suggested

## Trying out the Endpoints
- You can import the Postman collection - `userservice.postman_collection.json` and try out the endpoints

## Assumptions Made
- Assumed that the fields `first_name`, `last_name`, `email`, and `password` are required
- Assuming that the other services that are interested in the changes in the user entities are internal, message broker system like Kafka/RabbitMQ/SQS can be used to notify them
- For this demo implementation, boilerplace code for kafka like message broker has been skipped, instead the message broker is emulated with a dummy _Publish()_ function
- Also skipped designing the OpenAPI/Swagger documentation
- Password field is hashed before storing (during creation and updating)
- Password is not included in the response to GET requests

## Possible Extension/Improvements
### Basic Improvements
 - Validations:
    - Payload field validations
    - Verfication if a user already exists - eg: a user with duplicate email should not be allowed
 - If there are 3rd party external services that should be notified, then notifying them via webhook would be sensible. For that, we would also need a subscription management system, where the interested parties/clients can subscribe to the events that we expose. In the subscription, the subscriber should provide their REST API uri (along with credentials) that should be invoked when an event occurs.
 - Make use of a better logger framework like - zerolog
 - Page numbering system can be improved (using cursor instead of pagesize and number)

### Production Deployment
 - It is better to use even lean docker image for production deployment / image - instead of the `golang-alpine` image, we can just use the `alpine` image to run the already built application
 - The default parameters used in the application like PORT number, TIMEOUT, DB_URI etc should be passed as environment variable to the service
 - Make use of migration image and script to prepare the database schema.
 - Also make use DB indexes. It would be smart to create indexes on the fields that are allowed the query parameters (filterables).
 - Proper API documentation - like OpenAPI/Swagger with payload schema for each of the endpoints should be considered
 - Extensive integration tests with different flows and end-to-end scenarios should be executed on deployment

### Scaling
 - There can be several ways to scale the service depending on the usage.
 - The straight forward approach would be horizontal scaling – simply deploy more instances (increase the maximum number of pods in kubernetes)
 - Although, more refined approach would be to determine which operations are more in demand
 - If the READ operations (`GET /v1/users` and `GET /v1/users/:userId`) are more frequent, then it might make sense to implement caching layers
 - In case of too frequent WRITE operations – one approach would be to use the pure event driven architecure with the CQRS pattern, where the incoming requests (write operations) are only written to kafka / of similarly high performant message broker, and the response would be `202` Accepted instead of `201` Created. Updating of the database takes place asynchronously after the message is consumed from kafka. This approach also enables the service to  independently scale the READ componnet and WRITE component.

## The folder structure:
 - *apiusers* – the API layer
 - *db* – the Storage(Database) layer
 - *events* – simply holds the event structs for any changes on user entity (create, update, delete)
 - *grpc* – the gRPC layer
 - *messagebroker* – a dummy message broker
 - *mocks* – for unit/mock tests 
 - *services* – service layer implementations
 - *usererrors* - custom errors for the service
 - *utils* – common utilities like: hashing a password