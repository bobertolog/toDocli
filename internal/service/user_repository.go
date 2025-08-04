package service

type UserRepository interface {
	Register(username, password string) error
	Login(username, password string) error
}
