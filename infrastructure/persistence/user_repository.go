package persistence

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	_const "go-test/const"
	"go-test/crypto"
	"go-test/domain/entity"
	"go-test/domain/repository"
	"log"
	"math/big"
)

type UserRepo struct {
	db *bolt.DB
}

func NewUserRepository(db *bolt.DB) *UserRepo {
	return &UserRepo{db}
}

var _ repository.UserRepository = &UserRepo{}

func (u UserRepo) SaveUser(user *entity.User) error {

	return u.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}

		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("users"))

		// Marshal user data into bytes.
		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put(user.Address, buf)
	})
}

func (u UserRepo) GetUser(address []byte) (*entity.User, error) {
	user := &entity.User{}
	//get data on username is a key
	err := u.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		v := b.Get(address)
		err := json.Unmarshal(v, user)
		if err != nil {
			return err
		}
		return nil
	})
	return user, err
}

func (u UserRepo) DoTransaction(transaction *entity.Transaction) error {
	fromUser, err := u.GetUser(transaction.From)
	if err != nil {
		return err
	} else if err == nil && fromUser == nil {
		return _const.ERR_USER_NOT_FOUND
	}
	toUser, err := u.GetUser(transaction.To)
	if err != nil {
		return err
	} else if err == nil && fromUser == nil {
		return _const.ERR_USER_NOT_FOUND
	}
	//Convert private key from []byte to ecdsa type
	privK, err := crypto.ToECDSA(fromUser.PrivateKey, true)
	if transaction.Nonce != fromUser.Nonce+1 {
		return _const.ERR_NONCE
	}
	//verify signature
	if !crypto.VerifySignature(privK, transaction.Sign) {
		return _const.ERR_VERIFY_SIGNATURE
	}
	//convert amount
	fromAmount := new(big.Int).SetBytes(fromUser.Amount)
	transactionAmount := new(big.Int).SetBytes(transaction.Amount)

	//verify sender amount and transaction amount
	if fromAmount.Cmp(transactionAmount) < 0 {
		return _const.ERR_INSUFFICIENT
	}
	if len(transaction.From) != 42 || len(transaction.To) != 42 {
		return _const.ERR_ADDRESS_LENGTH
	}

	//calculate for uint256 from amount
	newFromAmount := new(big.Int).Sub(fromAmount, transactionAmount).Bytes()
	err = u.SaveUser(&entity.User{
		Address:    fromUser.Address,
		Amount:     newFromAmount,
		Nonce:      fromUser.Nonce + 1,
		PrivateKey: fromUser.PrivateKey,
	})
	log.Println("Sender amount", newFromAmount)
	//calculate for uint256 to amount
	toAmount := new(big.Int).SetBytes(toUser.Amount)
	newToAmount := new(big.Int).Add(toAmount, transactionAmount).Bytes()
	err = u.SaveUser(&entity.User{
		Address:    toUser.Address,
		Amount:     newToAmount,
		Nonce:      0,
		PrivateKey: toUser.PrivateKey,
	})
	log.Println("Receiver amount", newToAmount)

	return nil
}
