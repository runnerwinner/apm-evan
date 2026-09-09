package dogapm

import (
	"os"
	"os/signal"
	"syscall"
)
type starter interface {
	Start() error
}

type closer interface {
	Close() error
}

var (
	globalStarters = make([]starter, 0)
	globalClosers = make([]closer, 0)
)

type endPoint struct {
	stop chan int 
}

var EndPoint = &endPoint{
	stop: make(chan int, 1),
}

func (e *endPoint) Start() {
	for _, s := range globalStarters {
		_ = s.Start()
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		e.Shutdown()
	}()
	<-e.stop
}

func (e *endPoint) Shutdown() {
	for _, c := range globalClosers {
		_ = c.Close()
	}
	e.stop <- 1
}

func (e *endPoint) Close() {
	for _, c := range globalClosers {
		_ = c.Close()
	}

}