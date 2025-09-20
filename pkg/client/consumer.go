package client

import (
	"fmt"
	"io"
	"log"
	"net"

	packagetcp "github.com/tiago123456789/tqueue/pkg/packageTcp"
)

type ConsumerOptions struct {
	Address  string
	User     string
	Password string
	Queue    string
	Handler  func(string) error
}

type IConsumer interface {
	Connect() (*net.Conn, error)
	Disconnect()
	Start()
}

type Consumer struct {
	conn    *net.Conn
	options *ConsumerOptions
}

func NewConsumer(options *ConsumerOptions) *Consumer {
	return &Consumer{
		conn:    nil,
		options: options,
	}
}

func (p *Consumer) Start() {
	partMessage := ""
	for {
		buf := make([]byte, 1024)
		_, err := (*p.conn).Read(buf)
		if err == io.EOF {
			log.Println("Sever closed connection with the client")
			return
		}

		if err != nil {
			log.Println(err)
			return
		}

		items, code := packagetcp.ParseMessage([]byte(partMessage), buf)
		if code == packagetcp.INCOMPLETE {
			partMessage += string(buf)
		} else if code == packagetcp.COMPLETE {
			for _, item := range items {
				p.options.Handler(item)
			}
			partMessage = ""
		}

	}
}

func (p *Consumer) Disconnect() {
	(*p.conn).Close()
}

func (p *Consumer) Connect() (*net.Conn, error) {
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

	consumerMessage := fmt.Sprintf("C|%s", p.options.Queue)
	_, err = (*p.conn).Write([]byte(consumerMessage + "\n"))
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

	log.Println("Client consumer success")

	return &conn, nil
}
