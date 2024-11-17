package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

type Ticket struct {
	ID         int    `json:"id"`
	AssignedTo string `json:"assigned_to"`
}

type Client struct {
	ID   string
	Conn *websocket.Conn
}

var (
	clients      = make(map[*websocket.Conn]Client)
	clientsMutex sync.Mutex
	tickets      []Ticket
	ticketsMutex sync.Mutex
	upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("Server gestartet")
	fmt.Println("keine Tickets")
	go processTickets()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clientID := r.URL.Query().Get("clientID")
	client := Client{
		ID:   clientID,
		Conn: conn,
	}

	clientsMutex.Lock()
	clients[conn] = client
	clientsMutex.Unlock()

	for {
		var message map[string]interface{}
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("Client '%s' wurde geschlossen.\n", client.ID)
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			return
		}

		messageType, ok := message["message_type"].(string)
		if !ok {
			continue
		}

		switch messageType {
		case "client_id":
			clientID, ok := message["client_name"].(string)
			if !ok {
				continue
			}
			clientsMutex.Lock()
			client.ID = clientID
			clientsMutex.Unlock()
			fmt.Printf("Neuer Client: %s\n", clientID)
			printTickets()
			sendTicketInfo(conn)

		case "assign_ticket":
			ticketID, ok := message["ticket_id"].(float64)
			if !ok {
				continue
			}
			assignTicket(conn, client.ID, int(ticketID))
			printTickets()
			// case "client_ready":
			// 	sendTicketInfo(conn)
		}
	}
}

func printTickets() {
	if len(tickets) == 0 {
		fmt.Println("keine Tickets")
	} else {
		fmt.Println("Tickets:")
		for _, ticket := range tickets {
			fmt.Printf("id: %d, assigned_to: '%s'\n", ticket.ID, ticket.AssignedTo)
		}
	}
	fmt.Printf("n: neues Ticket, q: quit\n\n")
}

func processTickets() {
	for {
		var input string
		fmt.Print("n: neues Ticket, q: quit\n\n")
		fmt.Scanln(&input)

		switch input {
		case "n", "N":
			ticketsMutex.Lock()
			newTicket := Ticket{ID: len(tickets) + 1}
			tickets = append(tickets, newTicket)
			ticketsMutex.Unlock()

			fmt.Println("Tickets:")
			for _, ticket := range tickets {
				fmt.Printf("id: %d, assigned_to: '%s'\n", ticket.ID, ticket.AssignedTo)
			}

			sendTicketInfoToAllClients()

		case "q", "Q":
			fmt.Println("Server wird beendet.")
			os.Exit(0)
		}
	}
}

func sendTicketInfoToAllClients() {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for _, client := range clients {
		sendTicketInfo(client.Conn)
	}
}

func sendTicketInfo(conn *websocket.Conn) {
	ticketsMutex.Lock()
	defer ticketsMutex.Unlock()

	message := map[string]interface{}{
		"message_type": "ticket_info",
		"tickets":      tickets,
	}

	err := conn.WriteJSON(message)
	if err != nil {
		log.Println(err)
	}
}

func assignTicket(conn *websocket.Conn, clientID string, ticketID int) {
	ticketsMutex.Lock()
	var assignedTicket *Ticket
	for i := range tickets {
		if tickets[i].ID == ticketID && tickets[i].AssignedTo == "" {
			assignedTicket = &tickets[i]
			break
		}
	}
	ticketsMutex.Unlock()
	if assignedTicket != nil {
		ticketsMutex.Lock()
		assignedTicket.AssignedTo = clientID
		ticketsMutex.Unlock()
		sendTicketInfoToAllClients()
	} else {
		message := map[string]interface{}{
			"message_type": "assign_error",
			"error":        "Dieser Ticket ist bereit zugewiesen.",
		}
		err := conn.WriteJSON(message)
		if err != nil {
			log.Println(err)
		}
	}
}
