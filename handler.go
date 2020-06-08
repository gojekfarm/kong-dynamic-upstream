package main

import (
	"fmt"
	"sync"

	"github.com/Kong/go-pdk"
	"github.com/gojekfarm/kong-dyamic-upstream/dynamicupstream"
)

func init() {

}

//Handler is executed per request
type Handler struct {
	mu        sync.Mutex // guards balance
	Upstreams string
}

//New Creates an instance of a Handler
func New() interface{} {
	return &Handler{}
}

//Access is the entry point per request
func (conf *Handler) Access(kong *pdk.PDK) {
	r, _ := kong.Router.GetRoute()
	routeName := r.Name

	conf.mu.Lock()
	expr, port := dynamicupstream.RouteConfig(routeName, conf.Upstreams)
	path, err := kong.Request.GetPath()
	if err != nil {
		kong.Log.Err(err)
		return
	}

	kong.Log.Debug(fmt.Sprintf("the path is %s, handled by Expression %s, Port %s", path, expr, port))

	du, err := dynamicupstream.New(expr, port)

	if err != nil {
		kong.Log.Info("Error in initializing dynamic upstream")
		kong.Log.Err(err)
		kong.Nginx.Ask("kong.response.exit", 500)
		return
	}
	err = du.Target(&dynamicupstream.KongPDKWrapper{Kong: kong})
	if err != nil {
		kong.Log.Err("Error in targetting dynamic upstream")
		kong.Log.Err(err)
		kong.Nginx.Ask("kong.response.exit", 500)
		return
	}
	conf.mu.Unlock()
}
