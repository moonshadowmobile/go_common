package mq

import (
	"fmt"
	"crypto/tls"
	"crypto/x509"
	"github.com/streadway/amqp"
	"io/ioutil"
	"strings"
)

type RabbitMQManager struct {
	Connection *amqp.Connection
}
type RabbitMQConfig struct {
	Protocol string
	Addr string
	Port string
	Username string
	Password string
	TLS *RabbitMQTLSConfig
}
type RabbitMQTLSConfig struct {
	Ca string
	Cert string
	Key string
}

func CreateRabbitMQManager(rmqc *RabbitMQConfig) (*RabbitMQManager, error) {
	var conn *amqp.Connection
	var err error

	rmq := new(RabbitMQManager)
	u := strings.Join([]string{
		rmqc.Protocol, "://", rmqc.Username, ":", rmqc.Password, "@",
		rmqc.Addr, ":", rmqc.Port,
	}, "")
	rmqtls := rmqc.TLS
	if rmqtls != nil {
		cfg := new(tls.Config)
		cfg.MaxVersion = tls.VersionTLS10

		cfg.RootCAs = x509.NewCertPool()
		ca, err := ioutil.ReadFile(rmqtls.Ca)
		if err != nil {
			fmt.Println("can't read ca file")
			return nil, err
		}
		cfg.RootCAs.AppendCertsFromPEM(ca)
		cert, err := tls.LoadX509KeyPair(rmqtls.Cert, rmqtls.Key)
		if err != nil {
			fmt.Println("cant read cert or key")
			return nil, err
		}
		cfg.Certificates = append(cfg.Certificates, cert)
		conn, err = amqp.DialTLS(u, cfg)
		if err != nil {
			return nil, err
		}
	} else {
		conn, err = amqp.Dial(u)
		if err != nil {
			return nil, err
		}
	}

	rmq.Connection = conn

	return rmq, nil
}

