package frontpage

import (
	"encoding/json"
	"time"

	"github.com/StevenZack/fasthttp"
	"github.com/StevenZack/frontpage/logx"
	"github.com/StevenZack/pubsub"
	"github.com/StevenZack/websocket"
)

type WsServer struct {
	ChanID                         string
	ClientNum                      int
	DisableExitWhenAllClientClosed bool
	shutdown                       func() error
}

var (
	upgrader = websocket.FastHTTPUpgrader{}
)

func NewWsServer(shutdown func() error) *WsServer {
	return &WsServer{
		ChanID:   "ws",
		shutdown: shutdown,
	}
}

func (w *WsServer) ws(cx *fasthttp.RequestCtx) {
	e := upgrader.Upgrade(cx, func(c *websocket.Conn) {
		defer c.Close()
		ps := pubsub.NewPubSub()
		go func() {
			defer ps.UnSub()
			defer func() {
				w.ClientNum--
				time.Sleep(time.Second * 3)
				w.shutdownIfNeed()
			}()
			w.ClientNum++

			// read
			for {
				_, b, e := c.ReadMessage()
				if e != nil {
					logx.Error(e)
					return
				}
				w.handleMsg(b)
			}
		}()

		// write
		ps.Sub(w.ChanID, func(v interface{}) {
			if s, ok := v.(string); ok {
				c.WriteMessage(websocket.TextMessage, []byte(s))
				return
			}

			b, e := json.Marshal(v)
			if e != nil {
				logx.Error(e)
				return
			}
			c.WriteMessage(websocket.TextMessage, b)
		})
	})
	if e != nil {
		cx.Error(e.Error(), fasthttp.StatusBadRequest)
		return
	}
}

// shutdownIfNeed shutdown server if no client connected
func (w *WsServer) shutdownIfNeed() {
	if w.ClientNum < 1 && !w.DisableExitWhenAllClientClosed {
		w.shutdown()
	}
}

func (w *WsServer) handleMsg(b []byte) {

}
