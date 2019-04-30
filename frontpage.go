package frontpage

import "github.com/StevenZack/fasthttp"

type FrontPage struct {
	Router                       *fasthttp.Router
	Port                         string
	DisableExitWhenAllPageClosed bool
	wsCounter                    int
	fnMap                        map[string]Fn
}
