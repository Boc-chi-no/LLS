package model

// Link This Struct representing the data to be stored
type Link struct {
	ShortHash string   `bson:"_id"`
	URL       string   `bson:"url"`
	Token     string   `bson:"token"`
	Created   int64    `bson:"created"`
	Delete    bool     `bson:"delete"`
}
