// package main

// import (
// 	"log"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/websocket/v2"
// )

// func main() {
// 	app := fiber.New()

// 	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
// 		defer c.Close()
// 		for {
// 			messageType, msg, err := c.ReadMessage()
// 			if err != nil {
// 				log.Println("Error reading message:", err)
// 				break
// 			}

// 			log.Printf("Received: %s", msg)
// 			if err := c.WriteMessage(messageType, msg); err != nil {
// 				log.Println("Error writing message:", err)
// 				break
// 			}
// 		}
// 	}))

// 	log.Fatal(app.Listen(":3000"))
// }

//---------->ws://localhost:3000/ws

// package main

// import (
// 	"log"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/websocket/v2"
// )

// var clients = make(map[*websocket.Conn]bool)

// func main() {
// 	app := fiber.New()
// 	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
// 		defer func() {
// 			delete(clients, c)
// 			c.Close()
// 		}()
// 		clients[c] = true
// 		log.Println("New client connected")

// 		for {
// 			messageType, msg, err := c.ReadMessage()
// 			if err != nil {
// 				log.Println("Error reading message:", err)
// 				break
// 			}

// 			log.Printf("Received: %s", msg)

// 			for client := range clients {
// 				if client != c {
// 					if err := client.WriteMessage(messageType, msg); err != nil {
// 						log.Println("Error broadcasting message:", err)
// 					}
// 				}
// 			}
// 		}
// 	}))

// 	log.Fatal(app.Listen(":3000"))
// }

// //----------->ws://localhost:3000/ws
// //----------->ws://localhost:3000/ws

package main

import (
	"log"

	"github.com/gofiber/fiber"
	"github.com/gofiber/websocket"
)

type client struct{}

var clients = make(map[*websocket.Conn]client)
var register = make(chan *websocket.Conn)
var broadcast = make(chan string)
var unregister = make(chan *websocket.Conn)

func runHub() {
	for {
		select {
		case connection := <-register:
			clients[connection] = client{}
			log.Println("connection registered")

		case message := <-broadcast:
			log.Println("message received:", message)

			for connection := range clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("write error:", err)

					unregister <- connection
					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
				}
			}

		case connection := <-unregister:
			delete(clients, connection)

			log.Println("connection unregistered")
		}
	}
}

func main() {
	app := fiber.New()

	app.Static("/", "./home.html")

	app.Use(func(c *fiber.Ctx) {
		if websocket.IsWebSocketUpgrade(c) {
			c.Next()
		}
	})

	go runHub()

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			unregister <- c
			c.Close()
		}()

		register <- c

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return
			}

			if messageType == websocket.TextMessage {
				broadcast <- string(message)
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))

	app.Listen(":3000")
}

// in browse -- localhost:3000/ws
