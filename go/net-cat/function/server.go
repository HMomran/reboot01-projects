package function

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	ListenAddr string
	ln         net.Listener
	QuitCh     chan struct{}
	Clients    map[string]net.Conn
	Mutex      sync.Mutex
}

func NewServer(listenAddr string) *Server {
	return &Server{
		ListenAddr: listenAddr,
		QuitCh:     make(chan struct{}),
		Clients:    make(map[string]net.Conn),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	go s.acceptLoop()

	<-s.QuitCh
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}
		fmt.Println("Listening on the port: ", conn.RemoteAddr(), time.Now().Format("2006-01-02 15:04:05"))
		s.Mutex.Lock()
		if len(s.Clients) >= 10 {
			s.Mutex.Unlock()
			conn.Write([]byte("Server is full. Maximum number of clients (10) reached.\n"))
			conn.Close()
			continue
		}
		s.Mutex.Unlock()

		go s.readLoop(conn)
	}
}

func (s *Server) Broadcast(sender string, message string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	History(message)

	for name, conn := range s.Clients {
		if sender == name {
			continue
		}
		conn.Write([]byte("\n" + message + "\n"))
		prompt := fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)
		conn.Write([]byte(prompt))
	}
}
