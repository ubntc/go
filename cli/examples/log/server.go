package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/ubntc/go/cli/cli"
	"github.com/ubntc/go/cli/loggers/stdlogger"
)

// Server is an dummy server.
type Server struct {
	sync.WaitGroup
}

// Serve starts the server to test standard log usage.
func (s *Server) Serve(ctx context.Context) {
	s.Add(1)
	defer s.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("server stopped")
			return
		case <-time.After(1000 * time.Millisecond):
			log.Println("server alive")
		}
	}
}

// Shutdown waits for the server to stop.
func (s *Server) Shutdown() {
	s.Wait()
}

func main() {
	var (
		interactive = flag.Bool("i", false, "interactive mode")
		noClock     = flag.Bool("n", false, "don't display the clock")
	)
	flag.Parse()

	var opt []cli.Option
	if *interactive {
		opt = append(opt, cli.WithQuit())
	}
	if *noClock {
		opt = append(opt, cli.WithoutClock())
	}

	if *interactive {
		cli.SetupLogging(stdlogger.Setup)
	} else {
		log.SetFlags(log.LstdFlags | log.LUTC) // use unix time if running on non-interactive server
		log.Println("setting logger to UTC")
	}

	ctx, cancel := cli.WithSigWait(context.Background(), opt...)
	defer cancel()

	var srv Server
	go srv.Serve(ctx)

	<-ctx.Done()
	srv.Shutdown()
}
