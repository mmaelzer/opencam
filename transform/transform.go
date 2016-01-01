package transform

import "sync"

type Lockit struct {
	mutex    sync.Mutex
	Unlocked bool
}

func (l *Lockit) Lock() {
	l.set(false)
}

func (l *Lockit) Unlock() {
	l.set(true)
}

func (l *Lockit) set(lock bool) {
	l.mutex.Lock()
	l.Unlocked = lock
	l.mutex.Unlock()
}
