package tcp

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/tiago123456789/tqueue/internal/queue"
	instruction "github.com/tiago123456789/tqueue/pkg/instruction"
	packagetcp "github.com/tiago123456789/tqueue/pkg/packageTcp"
)

type ITcpManager interface {
	AddProducer(address string, conn net.Conn)
	AddConsumer(address string, conn net.Conn)
	RemoveProducer(address string)
	RemoveConsumer(address string)
	GetProducer(address string) (net.Conn, error)
	GetConsumer(address string) (net.Conn, error)
	GetTotalConsumer() int
	GetConsumers() map[string]net.Conn
	StartServer()
	IsConsumersClosed(address string) bool
	handleConnection(conn net.Conn)
}

type TcpManager struct {
	mu                             sync.Mutex
	connectionAlreadyAuthenticated map[string]bool
	producers                      map[string]net.Conn
	consumers                      map[string]net.Conn
	listener                       net.Listener
	queueManager                   queue.IQueueManager
	publishEngine                  func(*TcpManager, queue.IQueueManager)
	consumersClosed                map[string]bool
}

func NewTcpManager(
	queueManager queue.IQueueManager,
	publishEngine func(*TcpManager, queue.IQueueManager),
) *TcpManager {
	return &TcpManager{
		mu:                             sync.Mutex{},
		connectionAlreadyAuthenticated: make(map[string]bool),
		producers:                      make(map[string]net.Conn),
		consumers:                      make(map[string]net.Conn),
		queueManager:                   queueManager,
		publishEngine:                  publishEngine,
		consumersClosed:                make(map[string]bool),
	}
}

func (t *TcpManager) IsConsumersClosed(address string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.consumersClosed[address]
}

func (t *TcpManager) AddProducer(address string, conn net.Conn) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.producers[address] = conn
}

func (t *TcpManager) AddConsumer(address string, conn net.Conn) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.consumers[address] = conn
}

func (t *TcpManager) RemoveProducer(address string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.consumersClosed[address] = true
	delete(t.producers, address)
}

func (t *TcpManager) RemoveConsumer(address string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.consumers, address)
}

func (t *TcpManager) GetProducer(address string) (net.Conn, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.producers[address] != nil {
		return t.producers[address], nil
	}
	return nil, errors.New("Producer not found")
}

func (t *TcpManager) GetConsumer(address string) (net.Conn, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.consumers[address] != nil {
		return t.consumers[address], nil
	}
	return nil, errors.New("Consumer not found")
}

func (t *TcpManager) GetTotalConsumer() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.consumers)
}

func (t *TcpManager) GetConsumers() map[string]net.Conn {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.consumers
}

func (t *TcpManager) handleConnection(conn net.Conn) {
	defer conn.Close()
	partMessage := ""
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err == io.EOF {
			producer, _ := t.GetProducer(conn.RemoteAddr().String())
			consumer, _ := t.GetConsumer(conn.RemoteAddr().String())
			log.Println("producer: ", producer)
			log.Println("consumer: ", consumer)
			if producer != nil {
				log.Println("Client producer disconnected:", conn.RemoteAddr())
				t.queueManager.RemoveQueueProducerConnected(conn.RemoteAddr().String())
				t.RemoveProducer(conn.RemoteAddr().String())
			}

			if consumer != nil {
				log.Println("Client consumer disconnected:", conn.RemoteAddr())
				t.queueManager.RemoveQueueConsumerConnected(conn.RemoteAddr().String())
				t.RemoveConsumer(conn.RemoteAddr().String())
			}

			log.Println("Client disconnected:", conn.RemoteAddr())
			return
		}

		if err != nil {
			log.Println(err)
			return
		}

		items, code := packagetcp.ParseMessage([]byte(partMessage), buf)
		if code == packagetcp.INCOMPLETE {
			partMessage += string(buf)
			continue
		}

		message := items[0]
		isAuthentication := t.connectionAlreadyAuthenticated[conn.RemoteAddr().String()] == false &&
			string(message)[0:4] == instruction.AUTH
		if isAuthentication {
			username := strings.Split(string(message), "|")[1]
			password := strings.Split(string(message), "|")[2]
			if username != os.Getenv("USER_ADMIN") ||
				strings.Trim(password, "\n") != os.Getenv("PASSWORD") {
				conn.Write([]byte(instruction.RESPONSE_NOT_AUTHENTICATED + "\n"))
				log.Println("Client authentication failed:", conn.RemoteAddr().String())
				conn.Close()
				return
			}
			conn.Write([]byte(instruction.RESPONSE_AUTHENTICATED + "\n"))
			log.Println("Client authenticated:", conn.RemoteAddr().String())
			t.connectionAlreadyAuthenticated[conn.RemoteAddr().String()] = true
			continue
		}

		if t.connectionAlreadyAuthenticated[conn.RemoteAddr().String()] == false {
			t.RemoveConsumer(conn.RemoteAddr().String())
			t.RemoveProducer(conn.RemoteAddr().String())
			conn.Write([]byte(instruction.RESPONSE_NOT_AUTHENTICATED + "\n"))
			conn.Close()
			return
		}

		connProducer, err := t.GetProducer(conn.RemoteAddr().String())
		isProducerConnection := connProducer == nil && string(message)[0] == instruction.PRODUCER
		if isProducerConnection {
			t.AddProducer(conn.RemoteAddr().String(), conn)
			queueName := string(message)[2:]
			t.queueManager.CreateQueue(queueName)
			t.queueManager.SetQueueProducerConnected(conn.RemoteAddr().String(), queueName)
			log.Println("Client producer connected:", conn.RemoteAddr())
			conn.Write([]byte(instruction.RESPONSE_OK + "\n"))
			continue
		}

		connConsumer, err := t.GetConsumer(conn.RemoteAddr().String())
		isConsumerConnection := connConsumer == nil && string(message)[0] == instruction.CONSUMER
		if isConsumerConnection {
			t.AddConsumer(conn.RemoteAddr().String(), conn)
			queueName := string(message)[2:]
			t.queueManager.CreateQueue(queueName)
			t.queueManager.SetQueueConsumerConnected(conn.RemoteAddr().String(), queueName)
			log.Println("Client consumer connected:", conn.RemoteAddr())
			conn.Write([]byte(instruction.RESPONSE_OK + "\n"))
			continue
		}

		if code == packagetcp.COMPLETE {
			if string(message)[0:1] == "D" {
				queueName, _ := t.queueManager.GetQueueConsumerConnected(conn.RemoteAddr().String())
				messageId := string(message)[2:]
				t.queueManager.RemoveAvailableMessageById(queueName, strings.Trim(messageId, "\n"))
				conn.Write([]byte(instruction.RESPONSE_OK + "\n"))
				continue
			}

			queueName, _ := t.queueManager.GetQueueProducerConnected(conn.RemoteAddr().String())
			for _, item := range items {
				t.queueManager.Push(queueName, item)
				conn.Write([]byte(instruction.RESPONSE_OK + "\n"))
			}
			partMessage = ""
		}
	}
}
func (t *TcpManager) StartServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	t.listener = listener
	defer listener.Close()

	log.Println("TCP server is listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go t.handleConnection(conn)
		go t.queueManager.RequeueUnavailableMessages()
		go t.publishEngine(t, t.queueManager)
	}
}
