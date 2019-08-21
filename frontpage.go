package frontpage

import (
	"fmt"

	"github.com/StevenZack/frontpage/util"

	"github.com/StevenZack/fasthttp"
	"github.com/StevenZack/openurl"
	"github.com/StevenZack/tools/strToolkit"
	"github.com/StevenZack/websocket"
)

type FrontPage struct {
	Router                       *fasthttp.Router
	Port                         string
	chanID                       string
	DisableExitWhenAllPageClosed bool
	wsCounter                    int
	isRunning                    bool
	fnMap                        map[string]Fn
	fnOpen                       func(string)
	verbose                      bool
}

var (
	upgrader = websocket.FastHTTPUpgrader{}
)

func New(str string) *FrontPage {
	port := strToolkit.RandomPort()
	r := fasthttp.NewRouter()
	r.HandleFunc("/", func(cx *fasthttp.RequestCtx) {
		s, e := util.AddHead(str, `<script src="/var.js" type="text/javascript"></script>`)
		if e != nil {
			cx.Error(e.Error(), fasthttp.StatusBadRequest)
			return
		}
		cx.SetHtmlHeader()
		cx.WriteString(s)
	})
	fp := &FrontPage{
		Router:  r,
		Port:    port,
		chanID:  "frontpage/" + strToolkit.NewToken(),
		fnMap:   make(map[string]Fn),
		verbose: true,
	}
	r.HandleFunc("/ws", fp.ws)
	r.HandleFunc("/var.js", fp.varjs)
	return fp
}

func (f *FrontPage) SetVerbose(b bool) {
	f.verbose = b
}

func (f *FrontPage) HandleHTML(url, s string) {
	f.Router.HandleFunc(url, func(cx *fasthttp.RequestCtx) {
		s, e := util.AddHead(s, `<script src="/var.js" type="text/javascript"></script>`)
		if e != nil {
			cx.Error(e.Error(), fasthttp.StatusBadRequest)
			return
		}
		cx.SetHtmlHeader()
		cx.WriteString(s)
	})
}

func (f *FrontPage) HandleFunc(url string, fn func(cx *fasthttp.RequestCtx)) {
	f.Router.HandleFunc(url, fn)
}

func (f *FrontPage) HandleCSS(url string, css string) {
	f.Router.HandleFunc(url, func(cx *fasthttp.RequestCtx) {
		cx.SetCssHeader()
		cx.WriteString(css)
	})
}

func (f *FrontPage) HandleJS(url string, js string) {
	f.Router.HandleFunc(url, func(cx *fasthttp.RequestCtx) {
		cx.SetJsHeader()
		cx.WriteString(js)
	})
}

func (f *FrontPage) run() error {
	defer func() {
		f.isRunning = false
	}()
	f.isRunning = true
	return f.Router.ListenAndServe(":" + f.Port)
}
func (f *FrontPage) Run() error {
	if f.verbose {
		fmt.Println("listened on http://localhost:" + f.Port)
	}
	f.Open()
	return f.run()
}

func (f *FrontPage) RunBrowser() error {
	if f.verbose {
		fmt.Println("listened on http://localhost:" + f.Port)
	}
	f.OpenBrowser()
	return f.run()
}

func (f *FrontPage) SetOpenFn(fn func(string)) {
	f.fnOpen = fn
}

func (f *FrontPage) Start() {
	go f.Run()
}

func (f *FrontPage) StartBrowser() {
	go f.RunBrowser()
}

func (f *FrontPage) Open() {
	if f.fnOpen != nil {
		f.fnOpen("http://localhost:" + f.Port)
		return
	}
	openurl.OpenApp("http://localhost:" + f.Port)
}

func (f *FrontPage) OpenBrowser() {
	if f.fnOpen != nil {
		f.fnOpen("http://localhost:" + f.Port)
		return
	}
	openurl.Open("http://localhost:" + f.Port)
}

func (f *FrontPage) Shutdown() {
	f.Eval("window.close()")
	go f.Router.GetServer().Shutdown()
}

func (f *FrontPage) IsRunning() bool {
	return f.isRunning
}
