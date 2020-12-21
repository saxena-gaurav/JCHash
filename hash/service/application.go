package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
)

// Application stores dependencies and configuration
type Application struct {
	config     *Config
	router     *http.ServeMux
	inShutdown uint64
	waitGroup  *sync.WaitGroup
	errChan    chan error
}

// NewApplication returns an *Application.
func NewApplication(config *Config) *Application {
	return &Application{
		config:    config,
		router:    http.NewServeMux(),
		waitGroup: &sync.WaitGroup{},
		errChan:   make(chan error, 1),
	}
}

// Start begins the http.Server and waits for the server to stop.
func (a *Application) Start() error {
	server := &http.Server{
		Addr: a.config.listenAddr(),
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		a.Quit()
	}()

	a.waitGroup.Add(1)
	go func() {
		log.Printf("Server is listening at %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.QuitWithError(err)
		}
	}()

	a.waitGroup.Wait()

	if err := <-a.errChan; err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Server shutting down...\n")
	server.SetKeepAlivesEnabled(true)
	if err := server.Shutdown(context.Background()); err != nil {
		log.Println(err)
		return err
	}

	log.Println("Server stopped gracefully")
	return nil
}

// InShutdown returns whether or not Application is in the process of shutting down the http.Server
func (a *Application) InShutdown() bool {
	return atomic.LoadUint64(&a.inShutdown) == 1
}

// Quit calls Application.QuitWithError with a nil value indicating that the Application
// is quitting under expected circumstances.
func (a *Application) Quit() {
	a.QuitWithError(nil)
}

// QuitWithError starts the shutdown process.
func (a *Application) QuitWithError(err error) {
	if !a.InShutdown() {
		atomic.StoreUint64(&a.inShutdown, 1)
		a.waitGroup.Done()
		a.errChan <- err
	}
}
