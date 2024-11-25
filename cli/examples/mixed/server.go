package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/ubntc/go/cli/cli"
	"github.com/ubntc/go/cli/cli/config"
	"github.com/ubntc/go/cli/loggers/stdlogger"
	"github.com/ubntc/go/cli/loggers/zerologger"
)

// Server is an dummy server.
type Server struct {
	sync.WaitGroup
	status      string
	logInterval time.Duration
}

// Serve starts the server to test standard log usage.
func (s *Server) Serve(ctx context.Context) {
	s.status = "started"
	s.Add(1)
	defer s.Done()
	defer func() {
		s.status = "stopped"
	}()
	for {
		select {
		case <-ctx.Done():
			zlog.Info().Msg("server stopped")
			return
		case <-time.After(s.logInterval):
			s.Status(ctx)
		}
	}
}

func (s *Server) Status(context.Context) {
	zlog.Info().Str("status", s.status).Msg("zlog: server status")
	log.Println("log:  server status", s.status)
	fmt.Println("fmt:  server status", s.status)
}

func (s *Server) PrintStatus(context.Context) {
	fmt.Println("server", s.status)
}

// Shutdown waits for the server to stop.
func (s *Server) Shutdown(context.Context) {
	s.Wait()
}

func fmtPrint(context.Context) {
	fmt.Println("fmt.Println single-line")
	fmt.Println("fmt.Println\nmulti-\nline")
}

func logPrint(context.Context) {
	log.Println("log.Println single-line")
	log.Println("log.Println\nmulti-\nline")
}

func zeroPrint(context.Context) {
	zlog.Print("log.Println single-line")
	zlog.Print("log.Println\nmulti-\nline")
}

func help(context.Context) {
	fmt.Println(cli.GetCommands().Help())
}

func main() {
	var (
		interactive = flag.Bool("i", false, "interactive mode (also required for clock, raw, quit, and CR flags)")
		useZeroLog  = flag.Bool("z", false, "setup zerolog in interactive mode")
		stdLog      = flag.Bool("s", false, "setup stdlog in interactive mode")

		showClock = flag.Bool("c", false, "don't display the clock")
		rawTerm   = flag.Bool("raw", false, "set term to raw mode")
		crFix     = flag.Bool("cr", false, "prepend CR to NL")
		useQuit   = flag.Bool("q", false, "use Quit keys")

		verbose = flag.Bool("v", false, "more logs")
		debug   = flag.Bool("x", false, "debug mode")
	)
	flag.Parse()

	cli.GetTerm().SetVerbose(*verbose)
	cli.GetTerm().SetDebug(*debug)

	var srv Server
	srv.logInterval = time.Second

	var cmds cli.Commands
	cfg := config.Server()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if *interactive {
		cmds = []cli.Command{
			{Name: "fmt.Print", Key: 'f', Fn: fmtPrint},
			{Name: "log.Print", Key: 'l', Fn: logPrint},
			{Name: "zerolog.Print", Key: 'z', Fn: zeroPrint},
			{Name: "help", Key: 'h', Fn: help},
			{Name: "status", Key: 's', Fn: srv.Status},
			{Name: "print status", Key: 'p', Fn: srv.PrintStatus},
		}
		srv.logInterval = time.Second * 10

		// override all defaults
		cfg.ShowClock = *showClock
		cfg.PrependCR = *crFix
		cfg.MakeTermRaw = *rawTerm

		cfg.WithQuit = *useQuit

		if !cfg.WithQuit {
			cmds = append(cmds, cli.Command{
				Name: "custom quit", Key: 'q',
				Fn: func(context.Context) { cancel() },
			})
		}

		if *useZeroLog {
			cli.SetupLogging(zerologger.Setup)
			zlog.Info().Msg("setup zlog")
		} else if *stdLog {
			cli.SetupLogging(stdlogger.Setup)
			log.Println("setup stdlog")
		}

	} else {
		// use unix time if running on non-interactive server
		log.SetFlags(log.LstdFlags | log.LUTC)
		log.Println("setting standard logger to UTC")
	}

	tctx, tcancel := cli.StartTerm(ctx, cfg, cmds...)
	defer tcancel()

	go srv.Serve(tctx)

	<-ctx.Done()
	srv.Shutdown(tctx)
}
