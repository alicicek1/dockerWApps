package categoryType

import "time"

type Category struct {
	Id        string    `json:"_id" bson:"_id,omitempty"`
	Name      string    `json:"name,omitempty" bson:"name"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}

type CategoryPostRequestModel struct {
	Name string `json:"name,omitempty"`
}
