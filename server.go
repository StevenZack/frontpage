package frontpage

import (
	"fmt"

	"github.com/StevenZack/websocket"

	"github.com/StevenZack/openurl"

	"github.com/StevenZack/fasthttp"
	"github.com/StevenZack/tools/strToolkit"
)

var (
	upgrader = websocket.FastHTTPUpgrader{}
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

func (f *FrontPage) RunApp() error {
	fmt.Println("listened on http://localhost:" + f.Port)
	openurl.OpenApp("http://localhost:" + f.Port)
	return f.Router.ListenAndServe(":" + f.Port)
}
