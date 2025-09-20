package publishEngine

import (
	"encoding/json"
	"log"
	"time"

	"github.com/tiago123456789/tqueue/internal/queue"
	"github.com/tiago123456789/tqueue/internal/tcp"
)

func PublishEngine(tcpManager *tcp.TcpManager, queueManager queue.IQueueManager) {
	for {
		if tcpManager.GetTotalConsumer() > 0 {
			for _, client := range tcpManager.GetConsumers() {
				queueName, _ := queueManager.GetQueueConsumerConnected(client.RemoteAddr().String())
				queue, _ := queueManager.GetQueue(queueName)
				if queue == nil {
					continue
				}

				if queue.TotalMessages() > 0 {
					message := queue.Pop()
					if message.Message == "" && len(message.Message) == 0 {
						continue
					}

					jsonMessage, err := json.Marshal(message)
					if err != nil {
						log.Println(err)
						continue
					}

					if tcpManager.IsConsumersClosed(client.RemoteAddr().String()) {
						continue
					}

					client.Write([]byte(string(jsonMessage) + "\n"))
					continue
				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}
