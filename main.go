package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	crpto "go-test/crypto"
	"go-test/domain/entity"
	"go-test/infrastructure/persistence"
	"go-test/interfaces"
	"google.golang.org/protobuf/proto"
	"log"
	"math/big"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

func init() {
	//To load our environmental variables.
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {
	dbAddress := os.Getenv("DB_ADDRESS")
	host := os.Getenv("CONN_HOST")
	p1 := os.Getenv("CONN_PORT_S1")
	t := os.Getenv("TYPE")
	base := ""
	if len(os.Args) > 1 {
		base = os.Args[1]
	}
	services, err := persistence.NewRepositories(dbAddress)
	if err != nil {
		panic(err)
	}
	users := interfaces.NewUsers(services.User)
	go startServer(strings.ToLower(t), host, p1, users)

	switch base {
	case "init_user":
		log.Println("Init user")
		go initDB(users)
		go initDB(users)
		go initDB(users)
	}
	//for testing
	//user, _ := users.GetUser([]byte("0x35137192CE684dc06cdaA536c9bB7e717608f61C"))
	//user, _ := users.GetUser([]byte("0x35137192CE684dc06cdaA536c9bB7e717608f61C"))
	//fmt.Println(string(user.Address))
	//fmt.Println(user.Amount)
	//fmt.Println(user.Nonce)
	//fmt.Println(user.PrivateKey)
	//p, _ := crpto.ToECDSA(user.PrivateKey, true)
	//fmt.Println(">>>>>>>> PRIVATE KEY", p)
	sigHandler := func() {

		signChan := make(chan os.Signal, 1)
		signal.Notify(signChan, os.Interrupt)
		sig := <-signChan
		log.Print("Cleanup processes started by ", sig, " signal")
		services.Close()
		//os.Remove("test.db")
		//os.Remove("test.db.lock")
		os.Exit(1)

	}
	sigHandler()
}

func startServer(t, host, port string, us *interfaces.Users) {
	log.Println("starting tcp server on ", host+":"+port)
	listener, err := net.Listen(t, host+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		if conn, err := listener.Accept(); err == nil {
			handleConn(conn, us)
		}
	}
}

func handleConn(conn net.Conn, us *interfaces.Users) {
	log.Println("client connected")

	defer conn.Close()
	buffer := make([]byte, 1024)

	data, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DATA RECEIVED: ", data)
	messagePb := entity.Message{}
	err = proto.Unmarshal(buffer, &messagePb)
	log.Printf("received message: %s, timestamp: %v, transaction: %v", messagePb.Text, messagePb.Timestamp, messagePb.Transaction)
	res := &entity.Message{}
	if err = us.DoTransaction(messagePb.Transaction); err != nil {
		res = &entity.Message{Text: err.Error(), Timestamp: time.Now().Unix()}
		d, _ := proto.Marshal(res)
		conn.Write(d)
		log.Println(err)
	} else {
		res = &entity.Message{Text: "success", Timestamp: time.Now().Unix()}
		d, _ := proto.Marshal(res)
		conn.Write(d)

	}

}

func initDB(us *interfaces.Users) error {
	privateKey := crpto.GeneratePrivateKey()

	if err := us.SaveUser(&entity.User{
		Address:    crpto.GenerateAddress(privateKey),
		Amount:     big.NewInt(100).Bytes(),
		Nonce:      0,
		PrivateKey: crypto.FromECDSA(privateKey),
	}); err != nil {
		return err
	}
	return nil
}
