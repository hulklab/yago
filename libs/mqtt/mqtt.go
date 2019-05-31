package mqtt

import (
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/hulklab/yago"
	"log"
	"time"
)

func Ins(id ...string) *MqttConn {

	var name string

	if len(id) == 0 {
		name = "mqtt"
	} else if len(id) > 0 {
		name = id[0]
	}

	v := yago.Component.Ins(name, func() interface{} {
		conf := yago.Config.GetStringMap(name)
		broker := conf["broker"].(string)
		username := conf["username"].(string)
		password := conf["password"].(string)

		val := NewMqttConn(broker, username, password)

		if val == nil {
			log.Println("new default mqtt connect failed")
		}

		return val
	})

	return v.(*MqttConn)
}

//var clientId = "vision-demo-20180927185438"

type MqttConn struct {
	client MQTT.Client
}

type Log struct {
}

func (l Log) Println(v ...interface{}) {
	fmt.Println(v...)
}

func (l Log) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func NewMqttConn(broker, username, password string) *MqttConn {
	var clientId = username
	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetAutoReconnect(true)

	MQTT.DEBUG = Log{}
	log.Println("mqtt.conn", broker, clientId, username, password)

	client := MQTT.NewClient(opts)

	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}

	if token.Error() != nil {
		log.Fatal("mqtt", "connect err", token.Error(), "token:", token)
		return nil
	}

	if client == nil {
		log.Fatal("mqtt", "connect err", "client is nil")
		return nil
	}

	mqttConn := &MqttConn{
		client: client,
	}
	return mqttConn

}

func (m *MqttConn) Pub(topic string, value string, qos int) error {
	if m.client == nil {
		return errors.New("mqtt connect err")
	}

	// 打印底层的调用日志
	//MQTT.DEBUG = Log{}
	fmt.Println(m.client.IsConnected())
	token := m.client.Publish(topic, byte(qos), false, value)

	token.Wait()

	if token.Error() != nil {
		log.Println("mqtt", "pub err", token.Error())
	}

	return token.Error()
}

func (m *MqttConn) Sub(topic string, qos int, cb func(value string)) error {
	if m.client == nil {
		return errors.New("mqtt connect err")
	}

	// 打印底层的调用日志
	//MQTT.DEBUG = Log{}
	if token := m.client.Subscribe(topic, byte(qos), func(client MQTT.Client, message MQTT.Message) {
		value := string(message.Payload())
		log.Println(value)
		cb(value)
	}); token.Wait() && token.Error() != nil {
		log.Println("mqtt", "sub err", token.Error())
		return token.Error()
	}

	return nil
}
