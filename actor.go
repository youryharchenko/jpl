package jpl

import (
	"time"
)

/*
// Constants -
const (
	CLOSE = "close"
)

var (
	closeID = &ID{Value: CLOSE, Name: "ID"}
)
*/
/*
var actors = map[string]*Actor{}
var actorsLock = sync.RWMutex{}
var waitGroup sync.WaitGroup
var stopCh = make(chan struct{})
*/
func actorFuncs() map[string]Func {
	return map[string]Func{
		"actor": actor,
		"wait":  wait,
		"send":  send,
		"stop":  stopAll,
		"sleep": sleep,
	}
}

func actor(args []Expr, ctxName string) Expr {
	engine.debug("actor", args)
	if len(args) != 2 {
		return errID
	}
	id, ok := args[0].Eval().(*ID)
	if !ok {
		return errID
	}
	handler, ok := args[1].Eval().(*Lamb)
	if !ok {
		return errID
	}
	//engine.current[id.Value] = &Context{parent: engine.current["main"].clone(), vars: map[string]Expr{}}
	clon, _ := engine.current.Load("main") // engine.current["main"].clone()
	engine.current.Store(id.Value, &Context{parent: clon.(*Context), vars: map[string]Expr{}})
	handler.ChangeContext(id.Value)
	actor := &Actor{id: id.Value, chBox: make(chan Expr, 0), handler: handler}
	engine.actorsLock.Lock()
	engine.actors[id.Value] = actor
	engine.actorsLock.Unlock()
	engine.waitGroup.Add(1)
	go func(actor *Actor) {
		defer func() {
			engine.debug("actor", actor.id, "defer")
			engine.waitGroup.Done()
		}()
		for {
			select {
			case <-engine.stopCh:
				return
			default:
			}
			engine.debug("actor", actor.id, "waiting...")
			var e Expr
			select {
			case <-engine.stopCh:
				return
			case e = <-actor.chBox:
				engine.debug("actor", actor.id, "apply handler", e)
				res := applyFunc(actor.id, actor.handler, []Expr{e})
				//time.Sleep(time.Millisecond * 100)
				engine.debug("actor", actor.id, "result", res)
			}
		}
	}(actor)
	return trueID
}

func wait(args []Expr, ctxName string) Expr {
	if len(args) != 0 {
		return errID
	}
	engine.waitGroup.Wait()
	return nullID
}

func send(args []Expr, ctxName string) Expr {
	if len(args) != 3 {
		return errID
	}
	idFrom, ok := args[0].Eval().(*ID)
	if !ok {
		return errID
	}
	idTo, ok := args[1].Eval().(*ID)
	if !ok {
		return errID
	}
	e := args[2].Eval()
	engine.actorsLock.RLock()
	actor, ok := engine.actors[idTo.Value]
	engine.actorsLock.RUnlock()
	if !ok {
		return undefID
	}
	engine.debug("send", idFrom.Value, "to", actor.id, "sending ...", e)
	select {
	case <-engine.stopCh:
		engine.debug("send", idFrom.Value, "to", actor.id, "aborted")
	case actor.chBox <- e:
		engine.debug("sent", idFrom.Value, "to", actor.id, e)
	}
	return nullID
}

func stopAll(args []Expr, ctxName string) Expr {
	if len(args) != 0 {
		return errID
	}
	close(engine.stopCh)
	return nullID
}

func sleep(args []Expr, ctxName string) Expr {
	if len(args) != 1 {
		return errID
	}
	i, ok := args[0].Eval().(*Int)
	if !ok {
		return errID
	}
	time.Sleep(time.Millisecond * time.Duration(i.Value))
	return nullID
}

// Actor -
type Actor struct {
	id      string
	chBox   chan Expr
	handler *Lamb
}
