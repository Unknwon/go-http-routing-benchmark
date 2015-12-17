// Copyright 2014 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/emicklei/go-restful"
	"github.com/go-martini/martini"
	"github.com/gorilla/mux"
	"github.com/revel/revel"
	"github.com/robfig/pathtree"
	goji "github.com/zenazn/goji/web"
	"gopkg.in/macaron.v1"
)

type route struct {
	method string
	path   string
}

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

var nullLogger *log.Logger

func init() {
	// beego sets it to runtime.NumCPU()
	// Currently none of the contestors does concurrent routing
	runtime.GOMAXPROCS(1)

	// makes logging 'webscale' (ignores them)
	log.SetOutput(new(mockResponseWriter))
	nullLogger = log.New(new(mockResponseWriter), "", 0)

	initBeego()
	initMartini()
	initRevel()
}

// Common
func httpHandlerFunc(w http.ResponseWriter, r *http.Request) {}

// beego
func beegoHandler(ctx *context.Context) {}

func beegoHandlerWrite(ctx *context.Context) {
	ctx.WriteString(ctx.Input.Param(":name"))
}

func initBeego() {
	beego.RunMode = "prod"
	beego.BeeLogger.Close()
}

func loadBeego(routes []route) http.Handler {
	re := regexp.MustCompile(":([^/]*)")
	app := beego.NewControllerRegister()
	for _, route := range routes {
		route.path = re.ReplaceAllString(route.path, ":$1")
		switch route.method {
		case "GET":
			app.Get(route.path, beegoHandler)
		case "POST":
			app.Post(route.path, beegoHandler)
		case "PUT":
			app.Put(route.path, beegoHandler)
		case "PATCH":
			app.Patch(route.path, beegoHandler)
		case "DELETE":
			app.Delete(route.path, beegoHandler)
		default:
			panic("Unknow HTTP method: " + route.method)
		}
	}
	return app
}

func loadBeegoSingle(method, path string, handler beego.FilterFunc) http.Handler {
	app := beego.NewControllerRegister()
	switch method {
	case "GET":
		app.Get(path, handler)
	case "POST":
		app.Post(path, handler)
	case "PUT":
		app.Put(path, handler)
	case "PATCH":
		app.Patch(path, handler)
	case "DELETE":
		app.Delete(path, handler)
	default:
		panic("Unknow HTTP method: " + method)
	}
	return app
}

// goji
func gojiFuncWrite(c goji.C, w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, c.URLParams["name"])
}

func loadGoji(routes []route) http.Handler {
	mux := goji.New()
	for _, route := range routes {
		switch route.method {
		case "GET":
			mux.Get(route.path, httpHandlerFunc)
		case "POST":
			mux.Post(route.path, httpHandlerFunc)
		case "PUT":
			mux.Put(route.path, httpHandlerFunc)
		case "PATCH":
			mux.Patch(route.path, httpHandlerFunc)
		case "DELETE":
			mux.Delete(route.path, httpHandlerFunc)
		default:
			panic("Unknown HTTP method: " + route.method)
		}
	}
	return mux
}

func loadGojiSingle(method, path string, handler interface{}) http.Handler {
	mux := goji.New()
	switch method {
	case "GET":
		mux.Get(path, handler)
	case "POST":
		mux.Post(path, handler)
	case "PUT":
		mux.Put(path, handler)
	case "PATCH":
		mux.Patch(path, handler)
	case "DELETE":
		mux.Delete(path, handler)
	default:
		panic("Unknow HTTP method: " + method)
	}
	return mux
}

// go-restful
func goRestfulHandlerWrite(r *restful.Request, w *restful.Response) {
	io.WriteString(w, r.Request.URL.Query().Get("name"))
}

func goRestfulHandler(r *restful.Request, w *restful.Response) {}

func loadGoRestful(routes []route) http.Handler {
	wsContainer := restful.NewContainer()
	ws := new(restful.WebService)

	for _, route := range routes {
		switch route.method {
		case "GET":
			ws.Route(ws.GET(route.path).To(goRestfulHandler))
		case "POST":
			ws.Route(ws.POST(route.path).To(goRestfulHandler))
		case "PUT":
			ws.Route(ws.PUT(route.path).To(goRestfulHandler))
		case "PATCH":
			ws.Route(ws.PATCH(route.path).To(goRestfulHandler))
		case "DELETE":
			ws.Route(ws.DELETE(route.path).To(goRestfulHandler))
		default:
			panic("Unknow HTTP method: " + route.method)
		}
	}
	wsContainer.Add(ws)
	return wsContainer
}

func loadGoRestfulSingle(method, path string, handler restful.RouteFunction) http.Handler {
	wsContainer := restful.NewContainer()
	ws := new(restful.WebService)
	switch method {
	case "GET":
		ws.Route(ws.GET(path).To(handler))
	case "POST":
		ws.Route(ws.POST(path).To(handler))
	case "PUT":
		ws.Route(ws.PUT(path).To(handler))
	case "PATCH":
		ws.Route(ws.PATCH(path).To(handler))
	case "DELETE":
		ws.Route(ws.DELETE(path).To(handler))
	default:
		panic("Unknow HTTP method: " + method)
	}
	wsContainer.Add(ws)
	return wsContainer
}

// gorilla/mux
func gorillaHandlerWrite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	io.WriteString(w, params["name"])
}

func loadGorillaMux(routes []route) http.Handler {
	re := regexp.MustCompile(":([^/]*)")
	m := mux.NewRouter()
	for _, route := range routes {
		m.HandleFunc(
			re.ReplaceAllString(route.path, "{$1}"),
			httpHandlerFunc,
		).Methods(route.method)
	}
	return m
}

func loadGorillaMuxSingle(method, path string, handler http.HandlerFunc) http.Handler {
	m := mux.NewRouter()
	m.HandleFunc(path, handler).Methods(method)
	return m
}

// Macaron
func macaronHandler(_ *macaron.Context) {}

func macaronHandlerWrite(c *macaron.Context) string {
	return c.Params("name")
}

func loadMacaron(routes []route) http.Handler {
	m := macaron.New()
	for _, route := range routes {
		switch route.method {
		case "GET":
			m.Get(route.path, martiniHandler)
		case "POST":
			m.Post(route.path, martiniHandler)
		case "PUT":
			m.Put(route.path, martiniHandler)
		case "PATCH":
			m.Patch(route.path, martiniHandler)
		case "DELETE":
			m.Delete(route.path, martiniHandler)
		default:
			panic("Unknow HTTP method: " + route.method)
		}
	}
	return m
}

func loadMacaronSingle(method, path string, handler interface{}) http.Handler {
	m := macaron.New()
	switch method {
	case "GET":
		m.Get(path, handler)
	case "POST":
		m.Post(path, handler)
	case "PUT":
		m.Put(path, handler)
	case "PATCH":
		m.Patch(path, handler)
	case "DELETE":
		m.Delete(path, handler)
	default:
		panic("Unknow HTTP method: " + method)
	}
	return m
}

// Martini
func martiniHandler() {}

func martiniHandlerWrite(params martini.Params) string {
	return params["name"]
}

func initMartini() {
	martini.Env = martini.Prod
}

func loadMartini(routes []route) http.Handler {
	router := martini.NewRouter()
	for _, route := range routes {
		switch route.method {
		case "GET":
			router.Get(route.path, martiniHandler)
		case "POST":
			router.Post(route.path, martiniHandler)
		case "PUT":
			router.Put(route.path, martiniHandler)
		case "PATCH":
			router.Patch(route.path, martiniHandler)
		case "DELETE":
			router.Delete(route.path, martiniHandler)
		default:
			panic("Unknow HTTP method: " + route.method)
		}
	}
	martini := martini.New()
	martini.Action(router.Handle)
	return martini
}

func loadMartiniSingle(method, path string, handler interface{}) http.Handler {
	router := martini.NewRouter()
	switch method {
	case "GET":
		router.Get(path, handler)
	case "POST":
		router.Post(path, handler)
	case "PUT":
		router.Put(path, handler)
	case "PATCH":
		router.Patch(path, handler)
	case "DELETE":
		router.Delete(path, handler)
	default:
		panic("Unknow HTTP method: " + method)
	}

	martini := martini.New()
	martini.Action(router.Handle)
	return martini
}

// Revel (Router only)
// In the following code some Revel internals are modelled.
// The original revel code is copyrighted by Rob Figueiredo.
// See https://github.com/revel/revel/blob/master/LICENSE
type RevelController struct {
	*revel.Controller
	router *revel.Router
}

func (rc *RevelController) Handle() revel.Result {
	return revelResult{}
}

func (rc *RevelController) HandleWrite() revel.Result {
	return rc.RenderText(rc.Params.Get("name"))
}

type revelResult struct{}

func (rr revelResult) Apply(req *revel.Request, resp *revel.Response) {}

func (rc *RevelController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Dirty hacks, do NOT copy!
	revel.MainRouter = rc.router

	upgrade := r.Header.Get("Upgrade")
	if upgrade == "websocket" || upgrade == "Websocket" {
		panic("Not implemented")
	} else {
		var (
			req  = revel.NewRequest(r)
			resp = revel.NewResponse(w)
			c    = revel.NewController(req, resp)
		)
		req.Websocket = nil
		revel.Filters[0](c, revel.Filters[1:])
		if c.Result != nil {
			c.Result.Apply(req, resp)
		} else if c.Response.Status != 0 {
			panic("Not implemented")
		}
		// Close the Writer if we can
		if w, ok := resp.Out.(io.Closer); ok {
			w.Close()
		}
	}
}

func initRevel() {
	// Only use the Revel filters required for this benchmark
	revel.Filters = []revel.Filter{
		revel.RouterFilter,
		revel.ParamsFilter,
		revel.ActionInvoker,
	}

	revel.RegisterController((*RevelController)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Handle",
			},
			&revel.MethodType{
				Name: "HandleWrite",
			},
		})
}

func loadRevel(routes []route) http.Handler {
	router := revel.NewRouter("")

	// parseRoutes
	var rs []*revel.Route
	for _, r := range routes {
		rs = append(rs, revel.NewRoute(r.method, r.path, "RevelController.Handle", "", "", 0))
	}
	router.Routes = rs

	// updateTree
	router.Tree = pathtree.New()
	for _, r := range router.Routes {
		err := router.Tree.Add(r.TreePath, r)
		// Allow GETs to respond to HEAD requests.
		if err == nil && r.Method == "GET" {
			err = router.Tree.Add("/HEAD"+r.Path, r)
		}
		// Error adding a route to the pathtree.
		if err != nil {
			panic(err)
		}
	}

	rc := new(RevelController)
	rc.router = router
	return rc
}

func loadRevelSingle(method, path, action string) http.Handler {
	router := revel.NewRouter("")

	route := revel.NewRoute(method, path, action, "", "", 0)
	if err := router.Tree.Add(route.TreePath, route); err != nil {
		panic(err)
	}

	rc := new(RevelController)
	rc.router = router
	return rc
}

// Usage notice
func main() {
	fmt.Println("Usage: go test -bench=. -timeout=20m")
	os.Exit(1)
}
