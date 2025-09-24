package client

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tiago123456789/tqueue/pkg/instruction"
	packagetcp "github.com/tiago123456789/tqueue/pkg/packageTcp"
)

type ConsumerOptions struct {
	Address  string
	User     string
	Password string
	Queue    string
	Handler  func(Message) error
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
			partMessage = ""
			for _, item := range items {
				if item == instruction.RESPONSE_OK {
					continue
				}

				message, err := packagetcp.GetMessage(item)

				if err != nil {
					return
				}

				if message.Id == "" || message.Message == "" {
					continue
				}
				messageToProcess := Message{
					Id:      message.Id,
					Message: message.Message,
				}
				err = p.options.Handler(messageToProcess)
				if err == nil {
					(*p.conn).Write([]byte("D|" + string(message.Id) + "\n"))
				}
			}
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
