package function

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	displayWelcomeMessage(conn)

	name := getUserName(conn, buf, s)
	if name == "" {
		return
	}

	HistoryLoad(conn)
	joinMsg := fmt.Sprintf("%s has joined our chat...", name)
	s.Broadcast(name, joinMsg)

	defer func() {
		leaveMsg := fmt.Sprintf("%s has left our chat...", name)
		s.Broadcast(name, leaveMsg)
		s.Mutex.Lock()
		delete(s.Clients, name)
		s.Mutex.Unlock()
		conn.Close()
	}()

	handleChatMessages(conn, buf, name, s)
}

func displayWelcomeMessage(conn net.Conn) {
	file, err := os.Open("TCP-chat.txt")
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			conn.Write([]byte("\n"))
			conn.Write([]byte(scanner.Text()))
		}
	} else {
		fmt.Println("Error", err)
	}
}

func getUserName(conn net.Conn, buf []byte, s *Server) string {
	var name string
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading name:", err)
			return ""
		}
		nameInput := strings.TrimSpace(string(buf[:n]))

		filteredName := make([]rune, 0, len(nameInput))
		for _, r := range nameInput {
			if r >= 32 && r <= 126 {
				filteredName = append(filteredName, r)
			}
		}
		name = string(filteredName)

		if name == "" || len(name) > 20 {
			conn.Write([]byte("Invalid name. Please use only ASCII characters and keep it under 20 characters.\n[ENTER YOUR NAME]: "))
			continue
		}

		s.Mutex.Lock()
		if _, exists := s.Clients[name]; exists {
			s.Mutex.Unlock()
			conn.Write([]byte("This name is already taken. Please choose a different name.\nEnter your name: "))
			continue
		}

		if len(s.Clients) >= 10 {
			s.Mutex.Unlock()
			conn.Write([]byte("Server is full. Maximum number of clients (10) reached.\n"))
			return ""
		}

		s.Clients[name] = conn
		s.Mutex.Unlock()
		break
	}
	return name
}

func handleChatMessages(conn net.Conn, buf []byte, name string, s *Server) {
	for {
		msg := fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)
		conn.Write([]byte(msg))

		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("client disconnected:", name, conn.RemoteAddr(), time.Now().Format("2006-01-02 15:04:05"))
				break
			}
			fmt.Println("read error:", err)
			break
		}

		text := strings.TrimSpace(string(buf[:n]))
		cleanText := filterASCII(text)

		if cleanText == "" {
			continue
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formatted := fmt.Sprintf("[%s][%s]:%s", timestamp, name, cleanText)
		s.Broadcast(name, formatted)
	}
}

func filterASCII(text string) string {
	filtered := make([]rune, 0, len(text))
	for _, r := range text {
		if r >= 32 && r <= 126 {
			filtered = append(filtered, r)
		}
	}
	return string(filtered)
}
