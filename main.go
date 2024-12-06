package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

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

package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool)

func main() {
	app := fiber.New()
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			delete(clients, c)
			c.Close()
		}()
		clients[c] = true
		log.Println("New client connected")

		for {
			messageType, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}

			log.Printf("Received: %s", msg)

			for client := range clients {
				if client != c {
					if err := client.WriteMessage(messageType, msg); err != nil {
						log.Println("Error broadcasting message:", err)
					}
				}
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))
}

//----------->ws://localhost:3000/ws
//----------->ws://localhost:3000/ws