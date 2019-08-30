package kafka

import (
	"fmt"
	"log"
	"testing"
)

// go test -v ./app/libs/kafka -test.run Test* -args "-c=${PWD}/app.toml"

func TestProduce(t *testing.T) {
	k := Ins()
	k.Produce("my_test", "hello,world")
}

func TestConsume(t *testing.T) {
	var k = Ins()
	consumer, err := k.NewConsumer("my_test", "zjl")
	if err == nil {
		err := consumer.Consume(func(bytes []byte) {
			fmt.Println(string(bytes))
		})
		if err != nil {
			log.Println(err.Error())
		}
		log.Println("done")
	}
}
