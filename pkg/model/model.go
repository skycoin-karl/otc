package model

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/skycoin-karl/otc/pkg/actor"
	"github.com/skycoin-karl/otc/pkg/currencies"
	"github.com/skycoin-karl/otc/pkg/monitor"
	"github.com/skycoin-karl/otc/pkg/otc"
	"github.com/skycoin-karl/otc/pkg/scanner"
	"github.com/skycoin-karl/otc/pkg/sender"
)

var ErrReqMissing error = errors.New("request missing")

type Model struct {
	sync.RWMutex

	Workers *Workers
	Logger  *log.Logger
	Router  *actor.Actor
	Lookup  map[string]*otc.Request

	paused bool
}

func New(curs *currencies.Currencies) *Model {
	workers := &Workers{
		Scanner: actor.New(
			log.New(os.Stdout, "[SCANNER] ", log.LstdFlags),
			scanner.Task(curs),
		),
		Sender: actor.New(
			log.New(os.Stdout, " [SENDER] ", log.LstdFlags),
			sender.Task(curs),
		),
		Monitor: actor.New(
			log.New(os.Stdout, "[MONITOR] ", log.LstdFlags),
			monitor.Task(curs),
		),
	}

	model := &Model{
		Workers: workers,
		Logger:  log.New(os.Stdout, "    [OTC] ", log.LstdFlags),
		Router: actor.New(
			log.New(os.Stdout, "  [MODEL] ", log.LstdFlags),
			Task(workers),
		),
		Lookup: make(map[string]*otc.Request),
	}
	model.Start()

	return model
}

func (m *Model) Start() {
	go m.Router.Run(time.Second * 5)
	go m.Workers.Scanner.Run(time.Second * 5)
	go m.Workers.Sender.Run(time.Second * 5)
	go m.Workers.Monitor.Run(time.Second * 5)

	go func() {
		for {
			<-time.After(time.Second * 5)

			m.Logger.Printf(
				"(%d) = [%d] + [%d] + [%d]\n",
				m.Router.Count(),
				m.Workers.Scanner.Count(),
				m.Workers.Sender.Count(),
				m.Workers.Monitor.Count(),
			)
		}
	}()
}

func (m *Model) Status(iden string) (otc.Status, int64, error) {
	m.RLock()
	defer m.RUnlock()
	if m.Lookup[iden] == nil {
		return "", 0, ErrReqMissing
	}

	m.Lookup[iden].Lock()
	defer m.Lookup[iden].Unlock()
	return m.Lookup[iden].Status, m.Lookup[iden].Times.UpdatedAt, nil
}

func (m *Model) Load(req *otc.Request) error {
	m.Lock()
	defer m.Unlock()
	m.Lookup[req.Iden()] = req

	req.Lock()
	defer req.Unlock()

	work := &otc.Work{
		Request: req,
		Done:    make(chan *otc.Result, 1),
	}
	work.Done <- &otc.Result{
		Finished: time.Now().UTC().Unix(),
		Err:      nil,
	}
	m.Router.Add(work)

	return nil
}

func (m *Model) Add(req *otc.Request) error {
	m.Lock()
	defer m.Unlock()
	m.Lookup[req.Iden()] = req

	req.Lock()
	defer req.Unlock()

	res := &otc.Result{
		Finished: time.Now().UTC().Unix(),
		Err:      nil,
	}
	if err := Save(req, res); err != nil {
		return err
	}
	if req.Status == otc.NEW {
		req.Status = otc.DEPOSIT
	}

	work := &otc.Work{
		Request: req,
		Done:    make(chan *otc.Result, 1),
	}
	work.Done <- res

	m.Router.Add(work)

	return nil
}

func (m *Model) Paused() bool {
	m.RLock()
	defer m.RUnlock()
	return m.paused
}

func (m *Model) Pause() {
	m.Lock()
	defer m.Unlock()
	m.paused = true
}

func (m *Model) Unpause() {
	m.Lock()
	defer m.Unlock()
	m.paused = false
}
