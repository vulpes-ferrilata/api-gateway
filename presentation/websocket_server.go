package presentation

import (
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
)

var events = neffos.Namespaces{
	"catan": {
		neffos.OnNamespaceConnected: func(c *neffos.NSConn, msg neffos.Message) error {
			websocket.GetContext(c.Conn).Application().Logger().Infof("[%s] connected to namespace [%s]", c, msg.Namespace)

			return nil
		},
		neffos.OnNamespaceDisconnect: func(c *neffos.NSConn, msg neffos.Message) error {
			websocket.GetContext(c.Conn).Application().Logger().Infof("[%s] disconnected from namespace [%s]", c, msg.Namespace)

			return nil
		},
		neffos.OnRoomJoined: func(c *neffos.NSConn, msg neffos.Message) error {
			websocket.GetContext(c.Conn).Application().Logger().Infof("[%s] connected to room [%s]", c, msg.Room)

			return nil
		},
		neffos.OnRoomLeft: func(c *neffos.NSConn, msg neffos.Message) error {
			websocket.GetContext(c.Conn).Application().Logger().Infof("[%s] left from room [%s]", c, msg.Room)

			return nil
		},
	},
	"chat": {
		neffos.OnNamespaceConnected: func(c *neffos.NSConn, msg neffos.Message) error {
			websocket.GetContext(c.Conn).Application().Logger().Infof("[%s] connected to namespace [%s]", c, msg.Namespace)

			return nil
		},
		neffos.OnNamespaceDisconnect: func(c *neffos.NSConn, msg neffos.Message) error {
			websocket.GetContext(c.Conn).Application().Logger().Infof("[%s] disconnected from namespace [%s]", c, msg.Namespace)

			return nil
		},
		neffos.OnRoomJoined: func(c *neffos.NSConn, msg neffos.Message) error {
			websocket.GetContext(c.Conn).Application().Logger().Infof("[%s] connected to room [%s]", c, msg.Room)

			return nil
		},
		neffos.OnRoomLeft: func(c *neffos.NSConn, msg neffos.Message) error {
			websocket.GetContext(c.Conn).Application().Logger().Infof("[%s] left from room [%s]", c, msg.Room)

			return nil
		},
	},
}

func NewWebsocketServer() (*neffos.Server, error) {
	websocketServer := neffos.New(websocket.DefaultGobwasUpgrader, events)

	return websocketServer, nil
}
