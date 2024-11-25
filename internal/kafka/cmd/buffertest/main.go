package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/ubntc/go/cli/loggers/zerologger"
	bt "github.com/ubntc/go/internal/kafka/buffertest"
)

func main() {
	zerologger.Setup(os.Stderr, "")
	cfg := bt.LoadConfig()

	log.Info().Strs("brokers", cfg.Writer.Brokers).Msg("starting buffertest")
	numHandled, err := bt.Run(cfg)

	switch {
	case err != nil:
		log.Fatal().Err(err).Msg("[FAIL] buffertest failed with errors")
	case numHandled < cfg.NumEvents:
		log.Fatal().Msg("[FAIL] buffertest failed with missing events")
	default:
		log.Info().Msg("[OK] buffertest successsful")
	}
}
