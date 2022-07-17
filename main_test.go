package main_test

import (
	"go-test/domain/entity"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"math/big"
	"net"
	"testing"
	"time"
)

func TestClient1(t *testing.T) {
	msg := &entity.Message{Text: "Hi", Timestamp: time.Now().Unix(), Transaction: &entity.Transaction{
		From:   []byte("A1Wallet"),
		To:     []byte("A3Wallet"),
		Amount: big.NewInt(30).Bytes(),
		Nonce:  0,
		Sign:   []byte("zxcasd"),
	},
	}
	data, _ := proto.Marshal(msg)
	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		t.Fatal(err)
	}
	length, err := conn.Write(data)
	defer conn.Close()
	t.Logf("Data sent, length %d bytes", length)

	receivedData, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Log(err)
	}
	t.Log("DATA RECEIVED: ", receivedData)
	messagePb := entity.Message{}
	err = proto.Unmarshal(data, &messagePb)

	t.Logf("received message: %s, timestamp: %v, transaction: %v", messagePb.Text, messagePb.Timestamp, messagePb.Transaction)

}
func TestClient2(t *testing.T) {
	msg := &entity.Message{Text: "Hi", Timestamp: time.Now().Unix(), Transaction: &entity.Transaction{
		From:   []byte("A2Wallet"),
		To:     []byte("A3Wallet"),
		Amount: big.NewInt(300).Bytes(),
		Nonce:  0,
		Sign:   []byte("zxcasd"),
	},
	}
	data, _ := proto.Marshal(msg)
	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		t.Fatal(err)
	}
	length, err := conn.Write(data)
	defer conn.Close()
	t.Logf("Data sent, length %d bytes", length)

	receivedData, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Log(err)
	}
	t.Log("DATA RECEIVED: ", receivedData)
	messagePb := entity.Message{}
	err = proto.Unmarshal(data, &messagePb)

	t.Logf("received message: %s, timestamp: %v, transaction: %v", messagePb.Text, messagePb.Timestamp, messagePb.Transaction)

}
