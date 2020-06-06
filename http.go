package jpl

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
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
		Constructor: newHTTPRequest,
		Methods: map[string]Method{
			"toDict": httpRequestToDict,
		},
		Properties: map[string]Property{
			"method": {
				Name:    "method",
				Adapter: textToString,
			},
			"url": {
				Name:    "url",
				Adapter: textToString,
			},
			"body": {
				Name:    "body",
				Adapter: textToString,
				Default: "",
			},
			"header": {
				Name:    "header",
				Adapter: dictToMap,
				Default: map[string]string{},
			},
		},
	}
}

func makeHTTPClient() AnyClass {
	return AnyClass{
		Name:        "HttpClient",
		Constructor: newHTTPClient,
		Methods: map[string]Method{
			"do": doHTTPRequest,
		},
		Properties: map[string]Property{},
	}
}

func makeHTTPResponse() AnyClass {
	return AnyClass{
		Name:        "HttpResponse",
		Constructor: nil,
		Methods: map[string]Method{
			"toDict": httpResponseToDict,
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

func newHTTPClient(cls AnyClass, args []Expr, ctxName string) Expr {
	//engine.debug("HttpClient", "app", "new", args)
	if len(args) < 2 {
		return errID
	}
	_, err := cls.adaptDict(args[1].Eval())
	if err != nil {
		return errID
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	//engine.debug("HttpServer", "app", "new", app)
	return &Any{Name: "HttpClient", CtxName: ctxName, Value: client}
}

func newHTTPRequest(cls AnyClass, args []Expr, ctxName string) Expr {
	engine.debug("HttpRequest", "new", args)
	if len(args) < 2 {
		return errID
	}
	props, err := cls.adaptDict(args[1].Eval())
	if err != nil {
		engine.debug("HttpRequest", "new", "error", err)
		return errID
	}
	body := strings.NewReader(props["body"].(string))
	req, err := http.NewRequest(props["method"].(string), props["url"].(string), body)
	if err != nil {
		engine.debug("HttpRequest", "new", "error", err)
		return errID
	}
	for key, item := range props["header"].(map[string]string) {
		req.Header.Set(key, item)
	}
	engine.debug("HttpRequest", "new", req)
	return &Any{Name: "HttpRequest", CtxName: ctxName, Value: req}
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

func httpResponseToDict(any *Any, args []Expr, ctxName string) Expr {
	resp, ok := any.Value.(*http.Response)
	if !ok {
		return errID
	}
	header := map[string]Expr{}
	for key, item := range resp.Header {
		a := make([]Expr, len(item))
		for i, h := range item {
			a[i] = &Text{Name: "Text", Value: h, CtxName: ctxName}
		}
		header[key] = &Alist{Name: "Alist", Value: a, CtxName: ctxName}
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var closeExpr = falseID
	if resp.Close {
		closeExpr = trueID
	}
	var uncomprExpr = falseID
	if resp.Uncompressed {
		closeExpr = trueID
	}
	//resp.Body.Close()
	return &Dict{Name: "Dict", CtxName: ctxName,
		Value: map[string]Expr{
			"request":      &Any{Name: "HttpRequest", Value: resp.Request, CtxName: ctxName},
			"status":       &Text{Name: "Text", Value: resp.Status, CtxName: ctxName},
			"code":         &Int{Name: "Int", Value: resp.StatusCode, CtxName: ctxName},
			"close":        closeExpr,
			"proto":        &Text{Name: "Text", Value: resp.Proto, CtxName: ctxName},
			"protoMajor":   &Int{Name: "Int", Value: resp.ProtoMajor, CtxName: ctxName},
			"protoMinor":   &Int{Name: "Int", Value: resp.ProtoMinor, CtxName: ctxName},
			"header":       &Dict{Name: "Dict", Value: header, CtxName: ctxName},
			"length":       &Int{Name: "Int", Value: int(resp.ContentLength), CtxName: ctxName},
			"uncompressed": uncomprExpr,
			"body":         &Text{Name: "Text", Value: string(body), CtxName: ctxName},
		},
	}
}

func doHTTPRequest(any *Any, args []Expr, ctxName string) Expr {
	engine.debug("HttpRequest", "do", args)
	if len(args) < 3 {
		return errID
	}
	client, ok := any.Value.(*http.Client)
	if !ok {
		engine.debug("HttpClient", "do", "error any")
		return errID
	}
	engine.debug("HttpClient", "do", client)
	reqAny, ok := args[2].Eval().(*Any)
	if !ok {
		engine.debug("HttpClient", "do", "error args[2]")
		return errID
	}
	req, ok := reqAny.Value.(*http.Request)
	engine.debug("HttpClient", "do", client)
	resp, err := client.Do(req)
	if err != nil {
		engine.debug("HttpClient", "do", "error", err)
		return errID
	}
	return &Any{Name: "HttpResponse", CtxName: ctxName, Value: resp}
}
