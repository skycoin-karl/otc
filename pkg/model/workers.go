package model

import (
	"github.com/skycoin/services/otc/pkg/actor"
	"github.com/skycoin/services/otc/pkg/otc"
)

type Workers struct {
	Scanner *actor.Actor
	Sender  *actor.Actor
	Monitor *actor.Actor
}

func (w *Workers) Route(work *otc.Work) {
	switch work.Request.Status {
	case otc.DEPOSIT:
		w.Scanner.Add(work)
	case otc.SEND:
		w.Sender.Add(work)
	case otc.CONFIRM:
		w.Monitor.Add(work)
	}
}
