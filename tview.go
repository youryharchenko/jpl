package jpl

import "github.com/rivo/tview"

func tviewFuncs() map[string]Func {
	return map[string]Func{
		"tview": cmdTview,
	}
}

func cmdTview(args []Expr, ctxName string) Expr {
	if len(args) == 0 {
		return errID
	}
	cmd, ok := args[0].Eval().(*ID)
	if !ok {
		return errID
	}
	switch cmd.Value {
	case "app":
		return newApp(args, ctxName)
	case "run":
		return runApp(args, ctxName)
	case "box":
		return newBox(args, ctxName)
	case "modal":
		return newModal(args, ctxName)
	case "pages":
		return newPages(args, ctxName)
	}
	return nullID
}

func newApp(args []Expr, ctxName string) Expr {
	engine.debug("tview", "app", "new", args)
	if len(args) < 2 {
		return errID
	}
	dict, ok := args[1].Eval().(*Dict)
	if !ok {
		return errID
	}
	eroot, ok := dict.Value["root"]
	if !ok {
		return errID
	}
	root, ok := eroot.(*Any)
	if !ok {
		return errID
	}
	prim, ok := root.Value.(tview.Primitive)
	if !ok {
		return errID
	}
	emouse, ok := dict.Value["mouse"]
	if !ok {
		return errID
	}
	mouse, ok := emouse.(*ID)
	if !ok {
		return errID
	}
	app := tview.NewApplication()
	app.SetRoot(prim, true)
	app.EnableMouse(!mouse.Equals(falseID))
	engine.debug("tview", "app", "new", app)
	return &Any{Name: "TView.App", CtxName: ctxName, Value: app}
}

func newBox(args []Expr, ctxName string) Expr {
	if len(args) < 2 {
		return errID
	}
	dict, ok := args[1].Eval().(*Dict)
	if !ok {
		return errID
	}
	eborder, ok := dict.Value["border"]
	if !ok {
		return errID
	}
	border, ok := eborder.(*ID)
	if !ok {
		return errID
	}
	etitle, ok := dict.Value["title"]
	if !ok {
		return errID
	}
	title, ok := etitle.(*Text)
	if !ok {
		return errID
	}
	box := tview.NewBox()
	box.SetBorder(!border.Equals(falseID))
	box.SetTitle(title.Value)
	return &Any{Name: "TView.Box", Value: box, CtxName: ctxName}
}

func newModal(args []Expr, ctxName string) Expr {
	engine.debug("tview", "modal", "new", args)
	if len(args) < 2 {
		return errID
	}
	dict, ok := args[1].Eval().(*Dict)
	if !ok {
		return errID
	}
	ebuttons, ok := dict.Value["buttons"]
	if !ok {
		return errID
	}
	lbuttons, ok := ebuttons.(*Alist)
	if !ok {
		return errID
	}
	buttons := make([]string, len(lbuttons.Value))
	for i, button := range lbuttons.Value {
		buttons[i] = button.String()
	}
	etext, ok := dict.Value["text"]
	if !ok {
		return errID
	}
	text, ok := etext.(*Text)
	if !ok {
		return errID
	}
	edone, ok := dict.Value["done"]
	if !ok {
		return errID
	}
	modal := tview.NewModal()
	modal.AddButtons(buttons)
	modal.SetText(text.Value)
	modal.SetDoneFunc(func(index int, label string) {
		applyFunc(ctxName, edone, []Expr{
			&Int{Name: "Num", Value: index, CtxName: ctxName},
			&Text{Name: "Text", Value: label, CtxName: ctxName},
		})
	})
	engine.debug("tview", "modal", "new", modal)
	return &Any{Name: "TView.Modal", Value: modal, CtxName: ctxName}
}

func newPages(args []Expr, ctxName string) Expr {
	if len(args) < 2 {
		return errID
	}
	list, ok := args[1].Eval().(*Alist)
	if !ok {
		return errID
	}
	pages := tview.NewPages()
	for _, item := range list.Value {
		dict, ok := item.(*Dict)
		if !ok {
			return errID
		}
		epage, ok := dict.Value["page"]
		if !ok {
			return errID
		}
		page, ok := epage.(*Any)
		if !ok {
			return errID
		}
		prim, ok := page.Value.(tview.Primitive)
		if !ok {
			return errID
		}
		etitle, ok := dict.Value["title"]
		if !ok {
			return errID
		}
		title, ok := etitle.(*Text)
		if !ok {
			return errID
		}
		eresize, ok := dict.Value["resize"]
		if !ok {
			return errID
		}
		resize, ok := eresize.(*ID)
		if !ok {
			return errID
		}
		evisible, ok := dict.Value["visible"]
		if !ok {
			return errID
		}
		visible, ok := evisible.(*ID)
		if !ok {
			return errID
		}
		pages.AddPage(title.Value, prim, !resize.Equals(falseID), !visible.Equals(falseID))
	}

	return &Any{Name: "TView.Pages", Value: pages, CtxName: ctxName}
}

func runApp(args []Expr, ctxName string) Expr {
	engine.debug("tview", "app", "run", args)
	if len(args) < 2 {
		return errID
	}
	eapp, ok := args[1].Eval().(*Any)
	if !ok {
		return errID
	}
	app, ok := eapp.Value.(*tview.Application)
	if !ok {
		return errID
	}
	engine.debug("tview", "app", "run", app)
	err := app.Run()
	if err != nil {
		return errID
	}
	return nullID
}
