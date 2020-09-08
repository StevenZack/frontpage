package frontpage

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/StevenZack/openurl"

	"github.com/StevenZack/frontpage/util"
	"github.com/StevenZack/frontpage/views"
	"github.com/StevenZack/mux"
)

type FrontPage struct {
	Server   *mux.Server
	Vars     *Vars
	WsServer *WsServer
	binder   *binder
}

func New(html []byte, port int) *FrontPage {
	fp := &FrontPage{
		Vars: NewVars(),
	}
	if port > 0 {
		fp.Vars.Addr = ":" + strconv.Itoa(port)
	}

	fp.Server = mux.NewServer(fp.Vars.Addr)
	fp.binder = newBinder(fp.Vars)
	fp.WsServer = newWsServer(fp.Server.Stop)

	// handlers
	fp.HandleFunc("/fp/ws", fp.WsServer.ws)
	fp.HandleFunc("/fp/var.js", func(w http.ResponseWriter, Server *http.Request) {
		mux.SetJsHeader(w)
		t := template.New("var.js")
		t.Parse(views.Str_VarJs)
		t.Execute(w, fp.Vars)
	})
	fp.Server.HandleMultiReqs("/fp/call/", fp.binder.handleCall)
	fp.HandleHtml("/", string(html))

	return fp
}

func (f *FrontPage) HandleHtml(path string, html string) {
	f.Server.HandleFunc(path, func(w http.ResponseWriter, Server *http.Request) {
		s, e := util.AddHead(html, `<script src="/fp/var.js" type="text/javascript"></script>`)
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		mux.WriteHtml(w, s)
	})
}

func AddVarsJs(html string) (string, error) {
	s, e := util.AddHead(html, `<script src="/fp/var.js" type="text/javascript"></script>`)
	return s, e
}

func (f *FrontPage) HandleHtmlFunc(path string, fn func(w http.ResponseWriter, Server *http.Request) string) {
	f.Server.HandleFunc(path, func(w http.ResponseWriter, Server *http.Request) {
		s, e := util.AddHead(fn(w, Server), `<script src="/fp/var.js" type="text/javascript"></script>`)
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		mux.WriteHtml(w, s)
	})
}

func (f *FrontPage) HandleJs(path, js string) {
	f.Server.HandleFunc(path, func(w http.ResponseWriter, Server *http.Request) {
		mux.SetJsHeader(w)
		w.Write([]byte(js))
	})
}

func (f *FrontPage) HandleCss(path, css string) {
	f.Server.HandleFunc(path, func(w http.ResponseWriter, Server *http.Request) {
		mux.SetCssHeader(w)
		w.Write([]byte(css))
	})
}

func (f *FrontPage) HandleFunc(path string, handler func(w http.ResponseWriter, Server *http.Request)) {
	f.Server.HandleFunc(path, handler)
}

func (f *FrontPage) Run() error {
	fmt.Println("Listened on http://" + f.Vars.Addr)
	return f.Server.ListenAndServe()
}

func (f *FrontPage) Bind(v interface{}) {
	f.binder.bind(v)
}

func (f *FrontPage) Eval(s string) {
	f.WsServer.pub(s)
}

func (f *FrontPage) Open() {
	openurl.Open("http://" + f.Vars.Addr)
}

func (f *FrontPage) OpenApp() {
	openurl.OpenApp("http://" + f.Vars.Addr)
}
