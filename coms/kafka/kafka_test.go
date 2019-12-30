package kafka

import (
	"fmt"
	"testing"
)

//  go test -v ./coms/kafka  -args "-c=${PWD}/example/conf/app.toml"

// go test -v ./coms/kafka -run TestSyncProducer_Produce -args "-c=${PWD}/example/conf/app.toml"
func TestSyncProducer_Produce(t *testing.T) {
	k := Ins()
	p, err := k.SyncProducer()
	if err != nil {
		t.Fatalf("init sync producer error: %s", err)
	}
	defer p.Close()

	_, _, err = p.Produce("i am a sync produced msg")
	if err != nil {
		t.Fatalf("sync produce error: %s", err)
	}
}

// go test -v ./coms/kafka -run BenchmarkSyncProducer_Produce -bench=BenchmarkSyncProducer_Produce -args "-c=${PWD}/example/conf/app.toml"
func BenchmarkSyncProducer_Produce(b *testing.B) {
	k := Ins()
	p, err := k.SyncProducer()
	if err != nil {
		b.Fatalf("init sync producer error: %s", err)
	}
	defer p.Close()

	for i := 0; i < 1000; i++ {
		_, _, err := p.Produce("i am an sync produced msg")
		if err != nil {
			b.Fatalf("sync produce error: %s", err)
		}
	}
}

// go test -v ./coms/kafka -run TestAsyncProducer_Produce -args "-c=${PWD}/example/conf/app.toml"
func TestAsyncProducer_Produce(t *testing.T) {
	k := Ins()
	// 如果没有填写topic，默认采用配置文件中的topic
	p, err := k.AsyncProducer("demo")
	if err != nil {
		t.Fatalf("init async producer error: %s", err)
	}
	defer p.Close()

	p.Produce("i am an async produced msg")
}

// go test -v ./coms/kafka -run BenchmarkAsyncProducer_Produce -bench=BenchmarkAsyncProducer_Produce -args "-c=${PWD}/example/conf/app.toml"
func BenchmarkAsyncProducer_Produce(b *testing.B) {
	k := Ins()
	p, err := k.AsyncProducer()
	if err != nil {
		b.Fatalf("init async producer error: %s", err)
	}
	defer p.Close()

	for i := 0; i < 1000; i++ {
		p.Produce("i am an async produced msg")
	}
}

// go test -v ./coms/kafka -run TestConsumer_Consume -args "-c=${PWD}/example/conf/app.toml"
func TestConsumer_Consume(t *testing.T) {
	var k = Ins()
	consumer, err := k.NewConsumer("zjl")
	if err != nil {
		t.Errorf("new consumer error: %s", err)
	}
	err = consumer.Consume(func(bytes []byte) {
		fmt.Println(string(bytes))
	})
	if err != nil {
		t.Errorf("consume error: %s", err)
	}
}
