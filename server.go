package frontpage

import (
	"encoding/json"
	"fmt"

	"github.com/StevenZack/frontpage/views"
	"github.com/StevenZack/pubsub"

	"github.com/StevenZack/websocket"

	"github.com/StevenZack/openurl"

	"github.com/StevenZack/fasthttp"
	"github.com/StevenZack/tools/strToolkit"
)

var (
	upgrader   = websocket.FastHTTPUpgrader{}
)

const (
	chanID = "frontpage"
)

type FrontPage struct {
	Router *fasthttp.Router
	Port   string
}

func Run(str string) error {
	return New(str).Run()
}

func New(str string) *FrontPage {
	port := strToolkit.RandomPort()
	r := fasthttp.NewRouter()
	r.HandleFunc("/", func(cx *fasthttp.RequestCtx) {
		cx.SetHtmlHeader()
		cx.WriteString(str)
	})
	r.HandleFunc("/ws", ws)
	r.HandleFunc("/var.js", varjs)
	fp := &FrontPage{Router: r, Port: port}
	return fp
}

func (f *FrontPage) Run() error {
	fmt.Println("listened on http://localhost:" + f.Port)
	openurl.Open("http://localhost:" + f.Port)
	return f.Router.ListenAndServe(":" + f.Port)
}

func (f *FrontPage) RunInApp() error {
	fmt.Println("listened on http://localhost:" + f.Port)
	openurl.OpenApp("http://localhost:" + f.Port)
	return f.Router.ListenAndServe(":" + f.Port)
}

func ws(cx *fasthttp.RequestCtx) {
	e := upgrader.Upgrade(cx, func(c *websocket.Conn) {
		defer c.Close()
		ps := pubsub.NewPubSub()
		go func() {
			defer ps.UnSub()
			for {
				_, b, e := c.ReadMessage()
				if e != nil {
					return
				}
				handleMsg(b)
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


func varjs(cx *fasthttp.RequestCtx) {
	cx.SetJsHeader()
	cx.WriteString(views.Str_var)
}
