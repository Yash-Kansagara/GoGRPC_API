package models

type Exec struct {
	Id                        string `bson:"_id,omitempty"`
	FirstName                 string `bson:"first_name"`
	LastName                  string `bson:"last_name"`
	Email                     string `bson:"email"`
	Username                  string `bson:"username"`
	Password                  string `bson:"password"`
	PasswordChangedAt         string `bson:"password_changed_at"`
	UserCreatedAt             string `bson:"user_created_at"`
	PasswordResetToken        string `bson:"password_reset_token"`
	PasswordResetTokenExpires string `bson:"password_reset_token_expires"`
	Role                      string `bson:"role"`
	InactiveStatus            int32  `bson:"inactive_status"`
}
