package application

import (
	"go-test/domain/entity"
	"go-test/domain/repository"
)

type userApp struct {
	us repository.UserRepository
}

var _ UserAppInterface = &userApp{}

type UserAppInterface interface {
	DoTransaction(transaction *entity.Transaction) error
	SaveUser(*entity.User) error
	GetUser([]byte) (*entity.User, error)
}

func (u *userApp) DoTransaction(transaction *entity.Transaction) error {
	return u.us.DoTransaction(transaction)
}
func (u *userApp) SaveUser(user *entity.User) error {
	return u.us.SaveUser(user)
}

func (u *userApp) GetUser(addr []byte) (*entity.User, error) {
	return u.us.GetUser(addr)
}
