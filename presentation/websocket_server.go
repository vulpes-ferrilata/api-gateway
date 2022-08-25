package presentation

import (
	"log"

	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
)

var events = neffos.Namespaces{
	"catan": {
		neffos.OnNamespaceConnected: func(c *neffos.NSConn, msg neffos.Message) error {
			log.Printf("[%s] connected to namespace [%s].", c, msg.Namespace)
			return nil
		},
		neffos.OnNamespaceDisconnect: func(c *neffos.NSConn, msg neffos.Message) error {
			log.Printf("[%s] disconnected from namespace [%s].", c, msg.Namespace)
			return nil
		},
		neffos.OnRoomJoined: func(c *neffos.NSConn, msg neffos.Message) error {
			log.Printf("[%s] connected to room [%s].", c, msg.Room)
			return nil
		},
		neffos.OnRoomLeft: func(c *neffos.NSConn, msg neffos.Message) error {
			log.Printf("[%s] left from room [%s].", c, msg.Room)
			return nil
		},
	},
}

func NewWebsocketServer() (*neffos.Server, error) {
	websocketServer := neffos.New(websocket.DefaultGobwasUpgrader, events)

	return websocketServer, nil
}
