package model

// Link This Struct representing the data to be stored
type Link struct {
	ShortHash string `bson:"_id"`
	URL       string `bson:"url"`
	Password  string `bson:"password"`
	Token     string `bson:"token"`
	Created   int64  `bson:"created"`
	Expire    int64  `bson:"expire"`
	Memo      string `bson:"memo"`
	Delete    bool   `bson:"delete"`
}
