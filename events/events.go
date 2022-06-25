package events

//

import "github.com/sushilman/userservice/models"

const (
	USER_CREATED_TOPIC = "user-created"
	USER_UPDATED_TOPIC = "user-updated"
	USER_DELETED_TOPIC = "user-deleted"
)

type UserCreatedEvent models.User

type UserUpdatedEvent models.User

type UserDeletedEvent struct {
	Id string `json:"id"`
}
