package storageDriver

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tiago123456789/tqueue/pkg/types"
)

type SqliteStorageDriver struct {
	db                *sql.DB
	queueName         string
	mu                sync.Mutex
	messagesUnavaible *LinkedList
	messages          *LinkedList
}

func NewSqliteStorageDriver(db *sql.DB, queueName string) *SqliteStorageDriver {
	return &SqliteStorageDriver{
		db:        db,
		queueName: queueName,
	}
}

func (i *SqliteStorageDriver) getItem() *types.QueueItem {
	rows, err := i.db.Query(`
	UPDATE queue SET available_at = ?, is_available = false 
	WHERE id in (
		SELECT id FROM queue WHERE is_available = true AND queue_name = ? 
		ORDER BY created_at ASC LIMIT 1
	)
	RETURNING id, message, available_at
	`, time.Now().Add(time.Second*30), i.queueName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []types.QueueItem
	for rows.Next() {
		var u types.QueueItem
		var availableAt string
		err := rows.Scan(&u.Id, &u.Message, &availableAt)
		if err != nil {
			log.Fatal(err)
		}
		layout := "2006-01-02 15:04:05.999999999-07:00"
		u.AvailableAt, err = time.Parse(layout, availableAt)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, u)
	}

	if len(items) == 0 {
		return nil
	}

	return &items[0]
}

func (i *SqliteStorageDriver) GetByIdFromUnavaible(id string) *types.QueueItem {
	stmt, err := i.db.Prepare("DELETE FROM queue WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (i *SqliteStorageDriver) PushToUnavaible(message string) {
	// PS: when call the pop method, the message is setted as unavaible
}

func (i *SqliteStorageDriver) Push(message string) {
	stmt, err := i.db.Prepare("INSERT INTO queue(queue_name, id, message, available_at) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		i.queueName,
		uuid.New().String(),
		message,
		time.Now(),
	)
	if err != nil {
		log.Fatal(err)
	}

}

func (i *SqliteStorageDriver) Pop() *types.QueueItem {
	item := i.getItem()
	if item == nil {
		return nil
	}
	return item
}

func (i *SqliteStorageDriver) TotalMessages() int {
	total := 0
	row := i.db.QueryRow("SELECT COUNT(*) FROM queue where is_available = true AND queue_name = ?", i.queueName)
	row.Scan(&total)
	return total
}

func (i *SqliteStorageDriver) totalUnavaibleMessages() int {
	total := 0
	row := i.db.QueryRow("SELECT COUNT(*) FROM queue where is_available = false AND queue_name = ?", i.queueName)
	row.Scan(&total)
	return total
}

func (i *SqliteStorageDriver) RequeueUnavailableMessages() {
	for i.totalUnavaibleMessages() > 0 {
		_, err := i.db.Exec(`
		UPDATE queue SET is_available = true 
		WHERE is_available = false AND available_at < ?
		`, time.Now())
		if err != nil {
			log.Fatal(err)
		}
	}
}
