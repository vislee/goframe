package util

import (
	"sync"
)


type GoSrv struct {
	sync.WaitGroup
}


//Wrap go the func
func (s *GoSrv) Wrap(cb func()) {
	s.Add(1)
	go func() {
		cb()
		s.Done()
	}()
}
