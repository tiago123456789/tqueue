package client

import (
	"fmt"
	"log"
	"net"

	packagetcp "github.com/tiago123456789/tqueue/pkg/packageTcp"
)

type ProducerOptions struct {
	Address  string
	User     string
	Password string
	Queue    string
}

type IProducer interface {
	Connect() (*net.Conn, error)
	Disconnect()
	Send(message string) error
}

type Producer struct {
	conn    *net.Conn
	options *ProducerOptions
}

func NewProducer(options *ProducerOptions) *Producer {
	return &Producer{
		conn:    nil,
		options: options,
	}
}

func (p *Producer) Disconnect() {
	(*p.conn).Close()
}

func (p *Producer) Send(message string) error {
	if p.conn == nil {
		return fmt.Errorf("connection is nil. Hint: call Connect() first")
	}

	_, err := (*p.conn).Write([]byte(message + "\n"))
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	_, err = (*p.conn).Read(buf)
	if err != nil {
		return err
	}

	err = packagetcp.ParseResponse(string(buf))
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) Connect() (*net.Conn, error) {
	conn, err := net.Dial("tcp", p.options.Address)
	if err != nil {
		return nil, err
	}

	p.conn = &conn

	authMessage := fmt.Sprintf("AUTH|%s|%s", p.options.User, p.options.Password)
	_, err = (*p.conn).Write([]byte(authMessage + "\n"))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}

	err = packagetcp.ParseResponse(string(buf))
	if err != nil {
		return nil, err
	}

	log.Println("Client authentication success")

	producerMessage := fmt.Sprintf("P|%s", p.options.Queue)
	_, err = (*p.conn).Write([]byte(producerMessage + "\n"))
	if err != nil {
		return nil, err
	}

	buf = make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}

	err = packagetcp.ParseResponse(string(buf))
	if err != nil {
		return nil, err
	}

	log.Println("Client producer success")

	return &conn, nil
}
