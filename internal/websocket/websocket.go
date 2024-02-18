// websocketmanager/websocket_manager.go
package websocketmanager

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var (
	instance *Manager
	once     sync.Once
)

// Manager manages WebSocket connections and messages.
type Manager struct {
	connections map[*websocket.Conn]bool
	lock        sync.Mutex
}

// GetInstance returns the singleton instance of the WebSocket Manager.
func GetInstance() *Manager {
	once.Do(func() {
		instance = &Manager{
			connections: make(map[*websocket.Conn]bool),
		}
	})
	return instance
}

// AddConnection adds a new WebSocket connection to the manager.
func (m *Manager) AddConnection(conn *websocket.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.connections[conn] = true
}

// RemoveConnection removes a WebSocket connection from the manager.
func (m *Manager) RemoveConnection(conn *websocket.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.connections[conn]; ok {
		delete(m.connections, conn)
	}
}

// BroadcastMessage sends a message to all managed WebSocket connections.
func (m *Manager) BroadcastMessage(message string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for conn := range m.connections {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Printf("Error broadcasting message: %v", err)
			conn.Close()
			delete(m.connections, conn)
		}
	}
}
