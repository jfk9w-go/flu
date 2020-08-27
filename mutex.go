package flu

import "sync"

type Unlocker interface {
	Unlock()
}

type UnlockerFunc func()

func (fun UnlockerFunc) Unlock() {
	fun()
}

type RWMutex struct {
	sync.RWMutex
}

func (mu *RWMutex) RLock() Unlocker {
	mu.RWMutex.RLock()
	return UnlockerFunc(mu.RUnlock)
}

func (mu *RWMutex) Lock() Unlocker {
	mu.RWMutex.Lock()
	return UnlockerFunc(mu.Unlock)
}

type Mutex struct {
	sync.Mutex
}

func (mu *Mutex) Lock() Unlocker {
	mu.Mutex.Lock()
	return UnlockerFunc(mu.Unlock)
}
