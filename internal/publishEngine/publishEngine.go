package publishEngine

import (
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
					if message == "" && len(message) == 0 {
						continue
					}

					client.Write([]byte(message + "\n"))
					continue
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
}
