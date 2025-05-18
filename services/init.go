package services

var UserService *userService

func init() {
	UserService = &userService{}
}
