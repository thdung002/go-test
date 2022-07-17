package interfaces

import (
	"go-test/application"
	"go-test/domain/entity"
)

type Users struct {
	us application.UserAppInterface
}

//Users constructor
func NewUsers(us application.UserAppInterface) *Users {
	return &Users{
		us: us,
	}
}

func (s *Users) DoTransaction(transaction *entity.Transaction) error {
	return s.us.DoTransaction(transaction)
}

func (s *Users) SaveUser(user *entity.User) error {
	return s.us.SaveUser(user)
}

func (s *Users) GetUser(addr []byte) (*entity.User, error) {
	return s.us.GetUser(addr)
}
