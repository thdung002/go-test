package persistence

import (
	"fmt"
	"github.com/boltdb/bolt"
	"go-test/domain/repository"
	"time"
)

type Repositories struct {
	User repository.UserRepository
	db   *bolt.DB
}

func NewRepositories(address string) (*Repositories, error) {
	db, err := loadDatabase(address)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		User: NewUserRepository(db),
		db:   db,
	}, nil

}

func loadDatabase(address string) (*bolt.DB, error) {
	var err error
	db, err := bolt.Open(address, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err == nil {
		fmt.Println("connect database successful", address)
		return db, nil
	}
	defer db.Close()
	return nil, err
}

func (s *Repositories) Close() error {
	return s.db.Close()
}
