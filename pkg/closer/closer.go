package closer

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	"github.com/sanchey92/jwt-example/internal/logger"
)

var globalCloser = newCloser(syscall.SIGINT, syscall.SIGTERM)

func Add(fn ...func() error) {
	globalCloser.add(fn...)
}

func Wait() {
	globalCloser.wait()
}

func CloseAll() {
	globalCloser.closeAll()
}

type Closer struct {
	once  sync.Once
	mu    sync.Mutex
	done  chan struct{}
	funcs []func() error
}

func newCloser(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.closeAll()

		}()
	}
	return c
}

func (c *Closer) add(fn ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, fn...)
	c.mu.Unlock()
}

func (c *Closer) wait() {
	<-c.done
}

func (c *Closer) closeAll() {
	log := logger.GetLogger()

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))

		for _, fn := range funcs {
			go func(f func() error) {
				errs <- f()
			}(fn)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Error("Error returned from closer", zap.Error(err))
			}
		}

		log.Info("Graceful shutdown")
	})
}
