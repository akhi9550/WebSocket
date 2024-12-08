package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer c.Close()
		for {
			messageType, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}

			log.Printf("Received: %s", msg)
			if err := c.WriteMessage(messageType, msg); err != nil {
				log.Println("Error writing message:", err)
				break
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))
}

//---------->ws://localhost:3000/ws
