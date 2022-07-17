package repository

import "go-test/domain/entity"

type UserRepository interface {
	SaveUser(*entity.User) error
	GetUser([]byte) (*entity.User, error)
	DoTransaction(*entity.Transaction) error
}
