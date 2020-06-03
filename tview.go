package jpl

import (
	"github.com/rivo/tview"
)

func makeTviewApp() AnyClass {
	return AnyClass{
		Name:        "TviewApp",
		Constructor: newApp,
		Methods: map[string]Method{
			"run":  runApp,
			"stop": stopApp,
		},
		Properties: map[string]Property{
			"root": {
				Name:    "root",
				Adapter: anyToInterface,
			},
			"mouse": {
				Name:    "mouse",
				Adapter: idToBool,
			},
		},
	}
}

func makeTviewBox() AnyClass {
	return AnyClass{
		Name:        "TviewBox",
		Constructor: newBox,
		Methods:     map[string]Method{},
		Properties: map[string]Property{
			"title": {
				Name:    "title",
				Adapter: textToString,
			},
			"border": {
				Name:    "border",
				Adapter: idToBool,
			},
		},
	}
}

func makeTviewPages() AnyClass {
	return AnyClass{
		Name:        "TviewPages",
		Constructor: newPages,
		Methods: map[string]Method{
			"switch": switchToPage,
		},
		Properties: map[string]Property{
			"resize": {
				Name:    "resize",
				Adapter: idToBool,
			},
			"visible": {
				Name:    "visible",
				Adapter: idToBool,
			},
			"title": {
				Name:    "title",
				Adapter: textToString,
			},
			"page": {
				Name:    "page",
				Adapter: anyToInterface,
			},
		},
	}
}

func makeTviewModal() AnyClass {
	return AnyClass{
		Name:        "TviewModal",
		Constructor: newModal,
		Methods:     map[string]Method{},
		Properties: map[string]Property{
			"buttons": {
				Name:    "buttons",
				Adapter: alistToStrings,
			},
			"text": {
				Name:    "text",
				Adapter: textToString,
			},
			"done": {
				Name:    "done",
				Adapter: exprToExpr,
			},
		},
	}
}

func newApp(cls AnyClass, args []Expr, ctxName string) Expr {
	//engine.debug("tview", "app", "new", args)
	if len(args) < 2 {
		return errID
	}
	props, err := cls.adaptDict(args[1].Eval())
	if err != nil {
		return errID
	}

	app := tview.NewApplication()
	app.SetRoot(props["root"].(tview.Primitive), true)
	app.EnableMouse(props["mouse"].(bool))
	//engine.debug("tview", "app", "new", app)
	return &Any{Name: "TviewApp", CtxName: ctxName, Value: app}
}

func newBox(cls AnyClass, args []Expr, ctxName string) Expr {
	if len(args) < 2 {
		return errID
	}
	props, err := cls.adaptDict(args[1].Eval())
	if err != nil {
		return errID
	}
	box := tview.NewBox()
	box.SetBorder(props["border"].(bool))
	box.SetTitle(props["title"].(string))
	return &Any{Name: "TviewBox", Value: box, CtxName: ctxName}
}

func newModal(cls AnyClass, args []Expr, ctxName string) Expr {
	//engine.debug("tview", "modal", "new", args)
	if len(args) < 2 {
		return errID
	}
	props, err := cls.adaptDict(args[1].Eval())
	if err != nil {
		return errID
	}
	modal := tview.NewModal()
	modal.AddButtons(props["buttons"].([]string))
	modal.SetText(props["text"].(string))
	modal.SetDoneFunc(func(index int, label string) {
		applyFunc(ctxName, props["done"].(Expr), []Expr{
			&Int{Name: "Num", Value: index, CtxName: ctxName},
			&Text{Name: "Text", Value: label, CtxName: ctxName},
		})
	})
	//engine.debug("tview", "modal", "new", modal)
	return &Any{Name: "TviewModal", Value: modal, CtxName: ctxName}
}

func newPages(cls AnyClass, args []Expr, ctxName string) Expr {
	if len(args) < 2 {
		return errID
	}
	list, ok := args[1].Eval().(*Alist)
	if !ok {
		return errID
	}
	pages := tview.NewPages()
	for _, item := range list.Value {

		props, err := cls.adaptDict(item)
		if err != nil {
			return errID
		}
		pages.AddPage(props["title"].(string), props["page"].(tview.Primitive),
			props["resize"].(bool), props["visible"].(bool))
	}
	return &Any{Name: "TviewPages", Value: pages, CtxName: ctxName}
}

func runApp(any *Any, args []Expr, ctxName string) Expr {
	//engine.debug("tview", "app", "run", args)
	if len(args) < 2 {
		return errID
	}
	app, ok := any.Value.(*tview.Application)
	if !ok {
		return errID
	}
	//engine.debug("tview", "app", "run", app)
	err := app.Run()
	if err != nil {
		return errID
	}
	return nullID
}

func stopApp(any *Any, args []Expr, ctxName string) Expr {
	//engine.debug("tview", "app", "run", args)
	if len(args) < 2 {
		return errID
	}
	app, ok := any.Value.(*tview.Application)
	if !ok {
		return errID
	}
	//engine.debug("tview", "app", "run", app)
	app.Stop()
	return nullID
}

func switchToPage(any *Any, args []Expr, ctxName string) Expr {
	engine.debug("tview", "pages", "switch", any, args)
	if len(args) < 3 {
		return errID
	}
	pages, ok := any.Value.(*tview.Pages)
	if !ok {
		return errID
	}
	name, ok := args[2].Eval().(*Text)
	if !ok {
		return errID
	}
	engine.debug("tview", "pages", "switch", name, pages)
	pages.SwitchToPage(name.Value)
	return nullID
}
