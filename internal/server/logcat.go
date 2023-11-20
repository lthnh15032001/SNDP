package chserver

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan string)
var mutex = &sync.Mutex{}

func runAdbCommand(args ...string) (*exec.Cmd, *bufio.Scanner, error) {
	cmd := exec.Command("adb", args...)
	fmt.Printf("cmd %s \n", cmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	return cmd, bufio.NewScanner(stdout), nil
}

func connectAdb(svr string) error {
	cmd, _, err := runAdbCommand("connect", svr)
	if err != nil {
		return err
	}

	return cmd.Run()
}
func disconnectAdb(svr string) error {
	cmd, _, err := runAdbCommand("disconnect", svr)
	if err != nil {
		return err
	}

	return cmd.Run()
}
func handleLogcatOutput(conn *websocket.Conn, svr string) {
	// Connect to the Android device
	if err := connectAdb(svr); err != nil {
		log.Println("Error connecting to ADB:", err)
		return
	}

	cmd, scanner, err := runAdbCommand("logcat")
	if err != nil {
		log.Println("Error getting command stdout:", err)
		return
	}
	defer func() {
		if err := disconnectAdb(svr); err != nil {
			log.Println("Error disconnecting from ADB:", err)
		}
	}()

	cmd.Start()

	for scanner.Scan() {
		mutex.Lock()
		message := scanner.Text()
		mutex.Unlock()

		broadcast <- message
	}
}

func (s *Server) handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {

	host := strings.Split(r.Host, ":")[0]
	port := r.URL.Query().Get("port")
	hostport := fmt.Sprintf("%s:%s", host, port)
	fmt.Println("r.URL", hostport)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	mutex.Lock()
	clients[conn] = hostport
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		conn.Close()

		if err := disconnectAdb(hostport); err != nil {
			log.Println("Error disconnecting from ADB:", err)
		}
	}()

	// Send existing logcat output to the new client
	cmd, scanner, err := runAdbCommand("-s", hostport, "logcat")
	if err == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Connecting to ADB...\n"))
		if err := connectAdb(hostport); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Error connecting to ADB\n"))
		}
		conn.WriteMessage(websocket.TextMessage, []byte("\n"))
		cmd.Start()
		// scanner := bufio.NewScanner(cmd.Stdout)
		for scanner.Scan() {
			message := scanner.Text()
			conn.WriteMessage(websocket.TextMessage, []byte(message+"\n"))
		}
	}

	go handleLogcatOutput(conn, hostport)
}

func (s *Server) broadcastLogcatOutput() {
	for {
		message := <-broadcast
		mutex.Lock()
		for client, deviceInfo := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(message+"\n"))
			if err != nil {
				log.Println("Error writing message to client:", err)
				client.Close()
				delete(clients, client)

				if err := disconnectAdb(deviceInfo); err != nil {
					log.Println("Error disconnecting from ADB:", err)
				}
			}
		}
		mutex.Unlock()
	}
}

// func main() {
// 	go broadcastLogcatOutput()

// 	http.HandleFunc("/ws", handleWebSocketConnection)

// 	// Serve static files (you may need to adjust this depending on your project structure)
// 	http.Handle("/", http.FileServer(http.Dir(".")))

// 	log.Println("Server listening on :3001")
// 	err := http.ListenAndServe(":3001", nil)
// 	if err != nil {
// 		log.Fatal("Error starting server:", err)
// 	}
// }
