package event_watcher

import (
	"sync"

	"github.com/google/uuid"
)

var (
	Listeners = make(map[string]*Listener)
)

type Worker struct {
	Id      string
	Channel chan *Event
}

type Listener struct {
	*sync.Mutex
	chans map[string]*Worker
}

func NewWorker() *Worker {
	return &Worker{
		Id:      uuid.NewString(),
		Channel: make(chan *Event),
	}
}

func (l *Listener) Add(w *Worker) {
	if l.Mutex == nil {
		l.Mutex = new(sync.Mutex)
		l.chans = map[string]*Worker{}
	}
	l.Lock()
	defer l.Unlock()

	l.chans[w.Id] = w
}

func (l *Listener) Remove(w *Worker) {
	if l.Mutex == nil {
		l.Mutex = new(sync.Mutex)
		l.chans = map[string]*Worker{}
	}
	l.Lock()
	defer l.Unlock()

	delete(l.chans, w.Id)
}

func (l *Listener) Send(e *Event) {
	if l.Mutex == nil {
		l.Mutex = new(sync.Mutex)
		l.chans = map[string]*Worker{}
	}
	l.Lock()
	defer l.Unlock()

	for _, w := range l.chans {
		w.Channel <- e
	}
}
