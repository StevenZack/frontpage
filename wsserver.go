package frontpage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/StevenZack/pubsub"
)

type WsServer struct {
	ChanID                         string
	ClientNum                      int
	DisableExitWhenAllClientClosed bool
	shutdown                       func() error
}

var (
	upgrader = websocket.Upgrader{}
)

func NewWsServer(shutdown func() error) *WsServer {
	return &WsServer{
		ChanID:   "ws",
		shutdown: shutdown,
	}
}

func ws(w http.ResponseWriter, r *http.Request) {

}
func (server *WsServer) ws(w http.ResponseWriter, r *http.Request) {
	c, e := upgrader.Upgrade(w, r, nil)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	defer c.Close()
	ps := pubsub.NewPubSub()
	go func() {
		defer ps.UnSub()
		defer func() {
			server.ClientNum--
			time.Sleep(time.Second * 3)
			server.shutdownIfNeed()
		}()
		server.ClientNum++

		// read
		for {
			_, b, e := c.ReadMessage()
			if e != nil {
				return
			}
			server.handleMsg(b)
		}
	}()

	// write
	ps.Sub(server.ChanID, func(v interface{}) {
		if s, ok := v.(string); ok {
			c.WriteMessage(websocket.TextMessage, []byte(s))
			return
		}

		b, e := json.Marshal(v)
		if e != nil {
			fmt.Println(e)
			return
		}
		c.WriteMessage(websocket.TextMessage, b)
	})
}

// shutdownIfNeed shutdown server if no client connected
func (w *WsServer) shutdownIfNeed() {
	if w.ClientNum < 1 && !w.DisableExitWhenAllClientClosed {
		w.shutdown()
	}
}

func (w *WsServer) handleMsg(b []byte) {

}

func (w *WsServer) pub(s string) {
	pubsub.Pub(w.ChanID, s)
}
