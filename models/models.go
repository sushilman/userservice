package models

// used for the POST request - to accept the payload
type UserCreation struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Country   string `json:"country"`
}

// used for the GET request - as the response payload
type User struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Country   string `json:"country"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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
	Country   string `form:"country"`
	FirstName string `form:"last_name"`
	LastName  string `form:"last_name"`
	Offset    uint   `form:"offset"`
	Limit     uint   `form:"limit"`
}
