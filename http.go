package jpl

import (
	"io/ioutil"
	"net/http"
)

func makeHTTPServer() AnyClass {
	return AnyClass{
		Name:        "HttpServer",
		Constructor: newHTTPServer,
		Methods: map[string]Method{
			"run": runHTTPServer,
		},
		Properties: map[string]Property{
			"addr": {
				Name:    "addr",
				Adapter: textToString,
			},
			"handler": {
				Name:    "handler",
				Adapter: exprToExpr,
			},
		},
	}
}

func makeHTTPRequest() AnyClass {
	return AnyClass{
		Name:        "HttpRequest",
		Constructor: nil,
		Methods: map[string]Method{
			"toDict": httpRequestToDict,
		},
		Properties: map[string]Property{},
	}
}

// HTTPHandler -
type HTTPHandler struct {
	ctxName string
	handler Expr
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e := applyFunc(h.ctxName, h.handler, []Expr{
		&Any{Name: "HttpRequest", Value: r, CtxName: h.ctxName},
	})
	resp, ok := e.(*Text)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(resp.Value))
}

func newHTTPServer(cls AnyClass, args []Expr, ctxName string) Expr {
	//engine.debug("HttpServer", "app", "new", args)
	if len(args) < 2 {
		return errID
	}
	props, err := cls.adaptDict(args[1].Eval())
	if err != nil {
		return errID
	}

	server := &http.Server{
		Addr: props["addr"].(string),
		Handler: &HTTPHandler{
			handler: props["handler"].(Expr),
			ctxName: ctxName,
		},
	}
	//engine.debug("HttpServer", "app", "new", app)
	return &Any{Name: "HttpServer", CtxName: ctxName, Value: server}
}

func runHTTPServer(any *Any, args []Expr, ctxName string) Expr {
	engine.debug("HttpServer", "run", args)
	if len(args) < 2 {
		return errID
	}
	server, ok := any.Value.(*http.Server)
	if !ok {
		return errID
	}
	engine.debug("HttpServer", "run", server)
	err := server.ListenAndServe()
	if err != nil {
		engine.debug("HttpServer", "error", err)
		return errID
	}
	return nullID
}

func httpRequestToDict(any *Any, args []Expr, ctxName string) Expr {
	req, ok := any.Value.(*http.Request)
	if !ok {
		return errID
	}
	header := map[string]Expr{}
	for key, item := range req.Header {
		a := make([]Expr, len(item))
		for i, h := range item {
			a[i] = &Text{Name: "Text", Value: h, CtxName: ctxName}
		}
		header[key] = &Alist{Name: "Alist", Value: a, CtxName: ctxName}
	}
	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	return &Dict{Name: "Dict", CtxName: ctxName,
		Value: map[string]Expr{
			"host":       &Text{Name: "Text", Value: req.Host, CtxName: ctxName},
			"method":     &Text{Name: "Text", Value: req.Method, CtxName: ctxName},
			"url":        &Text{Name: "Text", Value: req.URL.String(), CtxName: ctxName},
			"proto":      &Text{Name: "Text", Value: req.Proto, CtxName: ctxName},
			"protoMajor": &Int{Name: "Int", Value: req.ProtoMajor, CtxName: ctxName},
			"protoMinor": &Int{Name: "Int", Value: req.ProtoMinor, CtxName: ctxName},
			"header":     &Dict{Name: "Dict", Value: header, CtxName: ctxName},
			"length":     &Int{Name: "Int", Value: int(req.ContentLength), CtxName: ctxName},
			"remote":     &Text{Name: "Text", Value: req.RemoteAddr, CtxName: ctxName},
			"body":       &Text{Name: "Text", Value: string(body), CtxName: ctxName},
		},
	}
}
