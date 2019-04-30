package frontpage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/StevenZack/fasthttp"
	"github.com/StevenZack/pubsub"
	"github.com/StevenZack/websocket"
)

func (fp *FrontPage) ws(cx *fasthttp.RequestCtx) {
	e := upgrader.Upgrade(cx, func(c *websocket.Conn) {
		defer c.Close()
		ps := pubsub.NewPubSub()
		go func() {
			defer ps.UnSub()
			defer func() {
				fp.wsCounter--
				if fp.wsCounter < 1 && !fp.DisableExitWhenAllPageClosed {
					time.Sleep(time.Second)
					fp.Router.GetServer().Shutdown()
				}
			}()
			fp.wsCounter++
			for {
				_, b, e := c.ReadMessage()
				if e != nil {
					fmt.Println("read err:", e)
					return
				}
				fmt.Println("got ws msg:", string(b))
				fp.handleMsg(b)
			}
		}()
		ps.Sub(chanID, func(v interface{}) {
			if s, ok := v.(string); ok {
				c.WriteMessage(websocket.TextMessage, []byte(s))
			} else {
				b, e := json.Marshal(v)
				if e != nil {
					fmt.Println("marshal error :", e)
					return
				}
				c.WriteMessage(websocket.TextMessage, b)
			}
		})
	})
	if e != nil {
		cx.Error(`upgrade : `+e.Error(), fasthttp.StatusInternalServerError)
		return
	}
}
