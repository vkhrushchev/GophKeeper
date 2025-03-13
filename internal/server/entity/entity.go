package entity

type User struct {
	ID           string `bson:"_id"`
	Username     string `bson:"username"`
	PasswordHash string `bson:"password_hash"`
	PasswordSalt string `bson:"password_salt"`
	Email        string `bson:"email"`
}
