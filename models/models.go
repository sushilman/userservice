package models

// used for the POST and PUT request - to accept the payload
type UserCreation struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Country   string `json:"country"`
}

// used for the GET request - as the response payload
// also reused as the DB model (to keep things simple)
type User struct {
	Id        string `json:"id" bson:"id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Nickname  string `json:"nickname" bson:"nickname"`
	Password  string `json:"-" bson:"password"`
	Email     string `json:"email" bson:"email"`
	Country   string `json:"country" bson:"country"`
	CreatedAt string `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}

// used as the response payload for the POST request
type UserCreationResponse struct {
	Link string `json:"link"`
}

type GetUsersResponse struct {
	Data  []User          `json:"data"`
	Links PaginationLinks `json:"links"`
}

type PaginationLinks struct {
	Prev string `json:"prev,omitempty"`
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
}

type GetUserQueryParams struct {
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
	NickName  string `form:"nickname"`
	Email     string `form:"email"`
	Country   string `form:"country"`
	Offset    uint   `form:"offset"`
	Limit     uint   `form:"limit"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
