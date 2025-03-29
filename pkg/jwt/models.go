package jwt

type AuthLevel string

const (
	Unknown    AuthLevel = "unknown"
	Guest      AuthLevel = "guest"
	User       AuthLevel = "user"
	Student    AuthLevel = "student"
	Instructor AuthLevel = "instructor"
	Moderator  AuthLevel = "moderator"
	Admin      AuthLevel = "admin"
)
