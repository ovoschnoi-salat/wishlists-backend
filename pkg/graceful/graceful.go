package graceful

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
)

// Config - конфигурация для Runner.
type Config struct {
	// SignalsToExit - количество сигналов, которые нужно получить для экстренного завершения программы.
	// Ноль и отрицательное значение означают, что программа будет дожидаться завершения всех запущенных сервисов.
	SignalsToExit int
}

func DefaultConfig() *Config {
	return &Config{SignalsToExit: 3}
}

// GetRunner - возвращает Runner с указанным конфигурационным объектом.
func (c *Config) GetRunner() *Runner {
	return &Runner{cfg: c, sigChan: make(chan os.Signal, c.SignalsToExit)}
}

// RunnableService - сервис, функция Run которого завершается только после завершения работы сервиса.
type RunnableService interface {
	Run() error // Run возвращает ошибку, если она возникла в процессе выполнения.
	Stop()      // Stop вызывается для завершения работы сервиса.
}

// Stopper - сервис, который может быть остановлен.
type Stopper interface {
	Stop() // Stop вызывается для завершения работы сервиса.
}

// Runner - управляет запуском и остановкой сервиса.
type Runner struct {
	cfg      *Config
	services []RunnableService // Сервис для запуска.
	errChan  chan error        // Канал для передачи ошибок от сервисов. nil означает, что сервис завершил работу без ошибок.
	sigChan  chan os.Signal    // Канал для получения сигналов.
}

func NewDefaultRunner() *Runner {
	return DefaultConfig().GetRunner()
}

func (g *Runner) Run(s ...RunnableService) error {
	g.services = s
	for _, service := range s {
		go func() {
			g.errChan <- service.Run()
		}()
	}
	return g.waitForResult()
}

// Run запускает все сервисы и возвращает ошибки, если они возникли при работе сервисов.
func (g *Runner) waitForResult() (resErr error) {
	signal.Ignore(syscall.SIGHUP)
	signal.Notify(g.sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	g.errChan = make(chan error, len(g.services))

	sigCount := 0
	stopFunc := sync.OnceFunc(g.stopServices)
	running := len(g.services)
	for {
		select {
		case sig := <-g.sigChan:
			log.Info().Str("signal", sig.String()).Msg("received signal")
			go stopFunc()
			sigCount++
			if sigCount == g.cfg.SignalsToExit {
				return errors.Join(resErr, errors.New("force exit"))
			}
		case err := <-g.errChan:
			go stopFunc()
			if err != nil {
				resErr = errors.Join(resErr, err)
			}
			running--
			if running == 0 {
				return
			}
		}
	}
}

func (g *Runner) stopServices() {
	for i := len(g.services) - 1; i >= 0; i-- {
		g.services[i].Stop()
	}
}
