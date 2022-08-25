package websocket

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/kataras/iris/v12"
)

type Server interface {
	Serve(ctx iris.Context)
}

type server struct {
}

func (s *server) Serve(ctx iris.Context) {
	conn, _, _, err := ws.UpgradeHTTP(ctx.Request(), ctx.ResponseWriter())
	if err != nil {
		panic(err)
	}

	err = wsutil.WriteServerMessage(conn, ws.OpText, []byte("hello"))
	if err != nil {
		panic(err)
	}

	go func() {
		defer conn.Close()

		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				panic(err)
			}

			err = wsutil.WriteServerMessage(conn, op, msg)
			if err != nil {
				panic(err)
			}
		}
	}()
}
