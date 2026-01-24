package models

type Student struct {
	Id        string `bson:"_id,omitempty"`
	FirstName string `bson:"first_name,omitempty"`
	LastName  string `bson:"last_name,omitempty"`
	Email     string `bson:"email,omitempty"`
	Class     string `bson:"class,omitempty"`
	Subject   string `bson:"subject,omitempty"`
}
