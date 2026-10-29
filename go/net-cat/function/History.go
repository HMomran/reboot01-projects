package function

import (
	"net"
	"sync"
)

var (
	history      []string
	historyMutex sync.Mutex
)

func History(message string) {
	historyMutex.Lock()
	defer historyMutex.Unlock()
	history = append(history, message)
}

func HistoryLoad(conn net.Conn) {
	historyMutex.Lock()
	defer historyMutex.Unlock()
	for _, msg := range history {
		conn.Write([]byte(msg + "\n"))
	}
}
