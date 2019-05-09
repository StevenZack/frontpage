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

func Run(str string) error {
	return New(str).Run()
}

func RunBrowser(str string) error {
	return New(str).RunBrowser()
}

func New(str string) *FrontPage {
	port := strToolkit.RandomPort()
	r := fasthttp.NewRouter()
	r.HandleFunc("/", func(cx *fasthttp.RequestCtx) {
		cx.SetHtmlHeader()
		cx.WriteString(str)
	})
	fp := &FrontPage{
		Router: r,
		Port:   port,
		chanID: "frontpage/" + strToolkit.NewToken(),
		fnMap:  make(map[string]Fn),
	}
	r.HandleFunc("/ws", fp.ws)
	r.HandleFunc("/var.js", fp.varjs)
	return fp
}
func (f *FrontPage) run() error {
	defer func() {
		f.isRunning = false
	}()
	f.isRunning = true
	return f.Router.ListenAndServe(":" + f.Port)
}
func (f *FrontPage) Run() error {
	fmt.Println("listened on http://localhost:" + f.Port)
	f.Open()
	return f.run()
}

func (f *FrontPage) RunBrowser() error {
	fmt.Println("listened on http://localhost:" + f.Port)
	f.OpenBrowser()
	return f.run()
}

func (f *FrontPage) Start() {
	go f.Run()
}

func (f *FrontPage) StartBrowser() {
	go f.RunBrowser()
}

func (f *FrontPage) Open() {
	openurl.OpenApp("http://localhost:" + f.Port)
}

func (f *FrontPage) OpenBrowser() {
	openurl.Open("http://localhost:" + f.Port)
}

func (f *FrontPage) Shutdown() {
	f.Eval("window.close()")
	go f.Router.GetServer().Shutdown()
}

func (f *FrontPage) IsRunning() bool {
	return f.isRunning
}
