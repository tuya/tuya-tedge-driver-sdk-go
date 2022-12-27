package mqttclient

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
)

//同时，username 应该由驱动传入，而不是从专家模式中读取，以及配置文件中的username!!!
const (
	DriverClientPrefix = "driver-"
	DriverTopicPrefix  = "tedge/driver/"
	InnerTopicPrefix   = "tedge/inner/"

	DriverTX = "tx"
	DriverRX = "rx"
)

const (
	Auth = iota + 1 // 连接鉴权
	Sub             // 设备订阅校验
	Pub             // 设备发布校验
	UnSub
	Connected
	Closed
)

type (
	AsyncMsg struct {
		Id   int64
		Type int
		Data json.RawMessage
	}

	AuthCheck struct {
		ClientId string
		Username string
		Password string
		Pass     bool
		Msg      string
	}

	PubTopic struct {
		ClientId string
		Username string
		Topic    string
		QoS      byte
		Retained bool
		Pass     bool
		Msg      string
	}
	// SubTopic 三方设备或服务订阅topic校验
	SubTopic struct {
		Topic string
		QoS   byte
		Pass  bool
		Msg   string
	}
	SubTopics struct {
		ClientId string
		Username string
		Topics   []SubTopic
	}

	// ConnectedNotify 三方设备或服务连接成功后通知对应驱动
	ConnectedNotify struct {
		ClientId string
		Username string
		IP       string
		Port     string
	}

	// ClosedNotify 三方设备或服务断开连接后通知对应驱动
	ClosedNotify struct {
		ClientId string
		Username string
	}

	UnSubNotify struct {
		ClientId string
		Username string
		Topics   []string
	}
)

type MqttClient struct {
	clientId string
	username string
	subTopic string
	pubTopic string
	driver   commons.MqttDriver
	logger   commons.TedgeLogger
	opts     *mqtt.ClientOptions
	client   mqtt.Client
}

func NewMqttClient(server, username string, logger commons.TedgeLogger) (*MqttClient, error) {
	pubT := DriverTopicPrefix + username + DriverTX
	subT := DriverTopicPrefix + username + DriverRX
	clientId := DriverClientPrefix + username
	h := md5.New()
	h.Write([]byte(username))
	rs := h.Sum(nil)
	pwd := hex.EncodeToString(rs)
	opts := mqtt.NewClientOptions().AddBroker(server).SetClientID(clientId).SetUsername(username).SetPassword(pwd[8:24])
	opts = opts.SetAutoReconnect(true).SetCleanSession(true).SetKeepAlive(5 * time.Second).SetMaxReconnectInterval(10 * time.Second).
		SetConnectRetry(true).SetConnectRetryInterval(time.Second)

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logger.Errorf("mqtt connection lost, err: %s", err)
	})

	// add mqtt internal log
	mqtt.ERROR = log.New(os.Stdout, "", log.LstdFlags)
	return &MqttClient{
		clientId: clientId,
		username: username,
		subTopic: subT,
		pubTopic: pubT,
		logger:   logger,
		opts:     opts,
	}, nil
}

func (mc *MqttClient) SetDriver(driver commons.MqttDriver) {
	mc.driver = driver
}

func (mc *MqttClient) SetClient(client mqtt.Client) {
	mc.client = client
}

func (mc *MqttClient) GetOpts() *mqtt.ClientOptions {
	return mc.opts
}

func (mc *MqttClient) Connect() error {
	token := mc.client.Connect()
	token.WaitTimeout(time.Second * 5)
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt token error: %w", err)
	}
	return nil
}
func (mc *MqttClient) Disconnect() {
	mc.client.Disconnect(1000)
}

func (mc *MqttClient) OnConnectHandler(handler mqtt.OnConnectHandler) mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		mc.logger.Infof("mqtt client connect success")
		if token := mc.client.Subscribe(mc.subTopic, byte(1), mc.onMessageReceived); token.Wait() && token.Error() != nil {
			mc.logger.Errorf("failed to sub topic:%s, err: %s", mc.subTopic, token.Error())
			return
		}
		mc.logger.Infof("subscribe success,topic: %s", mc.subTopic)
		// sub driver topic
		if handler != nil {
			handler(client)
		}
	}

}

func (mc *MqttClient) onMessageReceived(client mqtt.Client, message mqtt.Message) {
	mc.logger.Infof("Received message: %s", message.Payload())
	var req AsyncMsg
	decoder := json.NewDecoder(bytes.NewReader(message.Payload()))
	decoder.UseNumber()
	if err := decoder.Decode(&req); err != nil {
		mc.logger.Errorf("decode error: %s", err)
		return
	}
	var aack *AsyncMsg
	switch req.Type {
	case Auth:
		var auth AuthCheck
		if err := json.Unmarshal(req.Data, &auth); err != nil {
			mc.logger.Errorf("unmarshal auth msg error: %s", err)
			return
		}
		mc.logger.Infof("client auth, clientId: %s", auth.ClientId)

		// call driver interface
		if pass, err := mc.driver.Auth(auth.ClientId, auth.Username, auth.Password); err != nil {
			auth.Pass = false
			auth.Msg = err.Error()
		} else {
			auth.Pass = pass
		}

		payload, err := json.Marshal(auth)
		if err != nil {
			mc.logger.Errorf("marshal auth ack msg error: %s", err)
			return
		}
		aack = &AsyncMsg{
			Id:   req.Id,
			Type: Auth,
			Data: payload,
		}
	case Sub:
		var sub SubTopics
		if err := json.Unmarshal(req.Data, &sub); err != nil {
			mc.logger.Errorf("unmarshal sub topics msg error: %s", err)
			return
		}
		for i := range sub.Topics {
			mc.logger.Infof("sub: client: %s sub topic: %s, qos: %d", sub.ClientId, sub.Topics[i].Topic, sub.Topics[i].QoS)
			// call
			if pass, err := mc.driver.Sub(sub.ClientId, sub.Username, sub.Topics[i].Topic, sub.Topics[i].QoS); err != nil {
				sub.Topics[i].Pass = false
				sub.Topics[i].Msg = err.Error()
			} else {
				sub.Topics[i].Pass = pass
			}
		}
		payload, err := json.Marshal(sub)
		if err != nil {
			mc.logger.Errorf("marshal error: %s", err)
			return
		}
		aack = &AsyncMsg{
			Id:   req.Id,
			Type: Sub,
			Data: payload,
		}
	case Pub:
		var pub PubTopic
		if err := json.Unmarshal(req.Data, &pub); err != nil {
			mc.logger.Errorf("unmarshal pub topic msg error: %s", err)
		}
		mc.logger.Infof("pub: client: %s publish msg to topic: %s", pub.ClientId, pub.Topic)
		if pass, err := mc.driver.Pub(pub.ClientId, pub.Username, pub.Topic, pub.QoS, pub.Retained); err != nil {
			pub.Pass = false
			pub.Msg = err.Error()
		} else {
			pub.Pass = pass
		}

		payload, err := json.Marshal(pub)
		if err != nil {
			mc.logger.Errorf("marshal pub topic ack msg error: %s", err)
			return
		}
		aack = &AsyncMsg{
			Id:   req.Id,
			Type: Pub,
			Data: payload,
		}
	case UnSub:
		var unsub UnSubNotify
		if err := json.Unmarshal(req.Data, &unsub); err != nil {
			mc.logger.Errorf("unmarshal unsub notify msg error: %s", err)
			return
		}
		mc.logger.Infof("unsub: client: %s unsub topic: %+v", unsub.ClientId, unsub.Topics)
		mc.driver.UnSub(unsub.ClientId, unsub.Username, unsub.Topics)
		return
	case Connected:
		var connected ConnectedNotify
		if err := json.Unmarshal(req.Data, &connected); err != nil {
			mc.logger.Errorf("unmarshal connected notify msg error: %s", err)
			return
		}
		mc.logger.Infof("connected: client: %s connected", connected.ClientId)
		mc.driver.Connected(connected.ClientId, connected.Username, connected.IP, connected.Port)
		return
	case Closed:
		var closed ClosedNotify
		if err := json.Unmarshal(req.Data, &closed); err != nil {
			mc.logger.Errorf("unmarshal closed  notify msg error: %s", err)
			return
		}
		mc.logger.Infof("closed: client: %s closed", closed.ClientId)
		mc.driver.Closed(closed.ClientId, closed.Username)
		return
	default:
		mc.logger.Errorf("Received message on topic: %s,Message: %s", message.Topic(), message.Payload())
		return
	}

	buff, err := json.Marshal(aack)
	if err != nil {
		mc.logger.Errorf("marshal error: %s", err)
		return
	}

	if err := mc.publish(mc.pubTopic, byte(1), false, buff); err != nil {
		mc.logger.Errorf("publish msg to topic(%s) error: %s", mc.pubTopic, err)
	}
}

func (mc *MqttClient) Publish(topic string, qos byte, retained bool, message []byte) error {
	if !strings.HasPrefix(topic, InnerTopicPrefix) {
		return mc.publish(topic, qos, retained, message)
	}
	return errors.New("TopicNameInvalid")
}

func (mc *MqttClient) publish(topic string, qos byte, retained bool, message []byte) error {
	token := mc.client.Publish(topic, qos, retained, message)
	token.WaitTimeout(5 * time.Second)

	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (mc *MqttClient) Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error {
	if !strings.HasPrefix(topic, InnerTopicPrefix) {
		mc.logger.Infof("driver sub topic: %s...", topic)
		if token := mc.client.Subscribe(topic, qos, handler); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		mc.logger.Infof("driver sub topic: %s succeed", topic)
		return nil
	}
	return errors.New("TopicNameInvalid")
}

func (mc *MqttClient) UnSubscribe(topic string) error {
	mc.logger.Infof("driver unsub topic: %s...", topic)
	if token := mc.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	mc.logger.Infof("driver unsub topic: %s succeed", topic)
	return nil
}
