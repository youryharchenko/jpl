package jpl

import (
	"sync"
	"time"
)

// Constants -
const (
	CLOSE = "close"
)

var (
	closeID = &ID{Value: CLOSE, Name: "ID"}
)

var actors = map[string]*Actor{}
var actorsLock = sync.RWMutex{}
var waitGroup sync.WaitGroup
var stopCh = make(chan struct{})

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
	engine.current[id.Value] = &Context{parent: engine.current["main"].clone(), vars: map[string]Expr{}}
	handler.ChangeContext(id.Value)
	actor := &Actor{id: id.Value, chBox: make(chan Expr, 0), handler: handler}
	actorsLock.Lock()
	actors[id.Value] = actor
	actorsLock.Unlock()
	waitGroup.Add(1)
	go func(actor *Actor) {
		defer func() {
			engine.debug("actor", actor.id, "defer")
			waitGroup.Done()
		}()
		for {
			select {
			case <-stopCh:
				return
			default:
			}
			engine.debug("actor", actor.id, "waiting...")
			var e Expr
			select {
			case <-stopCh:
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
	waitGroup.Wait()
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
	actorsLock.RLock()
	actor, ok := actors[idTo.Value]
	actorsLock.RUnlock()
	if !ok {
		return undefID
	}
	engine.debug("send", idFrom.Value, "to", actor.id, "sending ...", e)
	select {
	case <-stopCh:
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
	close(stopCh)
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
