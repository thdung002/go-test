package main

import (
	"go-test/crypto"
	"go-test/domain/entity"
	"google.golang.org/protobuf/proto"
	"log"
	"math/big"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	client := ""
	if len(os.Args) > 1 {
		client = os.Args[1]
	}
	switch client {
	case "clients1":
		clients1()
	case "clients2":
		clients2()
	default:
		clients1()
		clients2()
	}
	sigHandler := func() {

		signChan := make(chan os.Signal, 1)
		signal.Notify(signChan, os.Interrupt)
		sig := <-signChan
		log.Print("Cleanup processes started by ", sig, " signal")
		os.Exit(1)

	}
	sigHandler()

}

func clients1() {
	p, _ := crypto.ToECDSA([]byte{202, 63, 51, 51, 151, 11, 17, 110, 201, 63, 99, 186, 177, 150, 14, 197, 89, 50, 206, 164, 102, 158, 130, 52, 235, 55, 169, 207, 55, 193, 28, 41}, true)
	sign := crypto.CreateSignature(p, []byte("0x35137192CE684dc06cdaA536c9bB7e717608f61C|sender"))
	msg := &entity.Message{Text: "Hi", Timestamp: time.Now().Unix(), Transaction: &entity.Transaction{
		From:   []byte("0x35137192CE684dc06cdaA536c9bB7e717608f61C"),
		To:     []byte("0xD642dA57BAe683F18E2bADcEE22741BBbe85006e"),
		Amount: big.NewInt(30).Bytes(),
		Nonce:  1,
		Sign:   sign,
	},
	}
	data, _ := proto.Marshal(msg)
	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		log.Fatal(err)
	}
	length, err := conn.Write(data)
	defer conn.Close()
	log.Printf("Data sent, length %d bytes", length)
	buffer := make([]byte, 1024)

	receivedData, err := conn.Read(buffer)
	if err != nil {
		log.Println(err)
	}
	log.Println("DATA RECEIVED: ", receivedData)
	messagePb := entity.Message{}
	err = proto.Unmarshal(buffer, &messagePb)

	log.Printf("received message: %s, timestamp: %v, transaction: %v", messagePb.Text, messagePb.Timestamp, messagePb.Transaction)

}

func clients2() {
	p, _ := crypto.ToECDSA([]byte{60, 245, 151, 92, 65, 169, 243, 169, 113, 234, 62, 31, 244, 62, 50, 57, 220, 162, 193, 124, 198, 87, 30, 86, 67, 185, 168, 213, 42, 155, 54, 129}, true)
	sign := crypto.CreateSignature(p, []byte("0x904b774bD1891269B0D01C7Bf92a611c623fd848|sender"))

	msg := &entity.Message{Text: "Hi", Timestamp: time.Now().Unix(), Transaction: &entity.Transaction{
		From:   []byte("0x904b774bD1891269B0D01C7Bf92a611c623fd848"),
		To:     []byte("0xD642dA57BAe683F18E2bADcEE22741BBbe85006e"),
		Amount: big.NewInt(300).Bytes(),
		Nonce:  1,
		Sign:   sign,
	},
	}
	data, _ := proto.Marshal(msg)
	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		log.Fatal(err)
	}
	length, err := conn.Write(data)
	defer conn.Close()
	log.Printf("Data sent, length %d bytes", length)
	buffer := make([]byte, 1024)

	receivedData, err := conn.Read(buffer)
	if err != nil {
		log.Println(err)
	}
	log.Println("DATA RECEIVED: ", receivedData)
	messagePb := entity.Message{}
	err = proto.Unmarshal(buffer, &messagePb)

	log.Printf("received message: %s, timestamp: %v, transaction: %v", messagePb.Text, messagePb.Timestamp, messagePb.Transaction)

}
