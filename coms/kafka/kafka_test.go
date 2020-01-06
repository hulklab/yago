package kafka

import (
	"fmt"
	"strings"
	"testing"
)

// go test -v ./coms/kafka  -args "-c=${PWD}/example/conf/app.toml"

// go test -v ./coms/kafka -run TestSyncProducer_Produce -args "-c=${PWD}/example/conf/app.toml"
func TestSyncProducer_Produce(t *testing.T) {
	k := Ins()
	p, err := k.SyncProducer()
	if err != nil {
		t.Fatalf("init sync producer error: %s", err)
	}
	defer p.close()

	for i := 0; i < 10; i++ {
		_, _, err = p.Produce("demo", fmt.Sprintf("sync msg: %d", i))
		if err != nil {
			t.Fatalf("sync produce error: %s", err)
		}
	}

}

// go test -v ./coms/kafka -run BenchmarkSyncProducer_Produce -bench=BenchmarkSyncProducer_Produce -args "-c=${PWD}/example/conf/app.toml"
func BenchmarkSyncProducer_Produce(b *testing.B) {
	k := Ins()
	p, err := k.SyncProducer()
	if err != nil {
		b.Fatalf("init sync producer error: %s", err)
	}
	defer p.close()

	for i := 0; i < 1000; i++ {
		_, _, err := p.Produce("demo", fmt.Sprintf("sync msg: %d", i))
		if err != nil {
			b.Fatalf("sync produce error: %s", err)
		}
	}
}

// go test -v ./coms/kafka -run TestAsyncProducer_Produce -args "-c=${PWD}/example/conf/app.toml"
func TestAsyncProducer_Produce(t *testing.T) {
	k := Ins()
	p, err := k.AsyncProducer()
	if err != nil {
		t.Fatalf("init async producer error: %s", err)
	}
	defer p.close()

	for i := 0; i < 10; i++ {
		p.Produce("demo", fmt.Sprintf("async msg: %d", i))
	}
}

// go test -v ./coms/kafka -run BenchmarkAsyncProducer_Produce -bench=BenchmarkAsyncProducer_Produce -args "-c=${PWD}/example/conf/app.toml"
func BenchmarkAsyncProducer_Produce(b *testing.B) {
	k := Ins()
	p, err := k.AsyncProducer()
	if err != nil {
		b.Fatalf("init async producer error: %s", err)
	}
	defer p.close()

	for i := 0; i < 1000; i++ {
		p.Produce("demo", fmt.Sprintf("async msg: %d", i))
	}
}

// go test -v ./coms/kafka -run TestConsumer_Consume -args "-c=${PWD}/example/conf/app.toml"
func TestConsumer_Consume(t *testing.T) {
	var k = Ins()
	consumer, err := k.NewConsumer("zjl", "demo", "demo1")
	if err != nil {
		t.Errorf("new consumer error: %s", err)
	}
	err = consumer.Consume(func(bytes []byte) bool {
		if strings.Contains(string(bytes), "5") {
			return false
		}
		fmt.Println(string(bytes))
		return true
	})
	if err != nil {
		t.Errorf("consume error: %s", err)
	}
}

func TestConsumer_Consume1(t *testing.T) {
	var k = Ins()
	consumer, err := k.NewConsumer("zjl", "demo", "demo1")
	if err != nil {
		t.Errorf("new consumer error: %s", err)
	}
	err = consumer.Consume(func(bytes []byte) bool {
		if strings.Contains(string(bytes), "1") {
			return false
		}
		fmt.Println(string(bytes))
		return true
	})
	if err != nil {
		t.Errorf("consume error: %s", err)
	}
}
