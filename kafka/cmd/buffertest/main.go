package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/ubntc/go/cli/loggers/zerologger"
	bt "github.com/ubntc/go/kafka/internal/buffertest"
)

func main() {
	zerologger.Setup(os.Stderr, "")
	cfg := bt.LoadConfig()

	log.Info().Msg("starting buffertest")
	numHandled, err := bt.Run(cfg)

	if numHandled != cfg.NumEvents || err != nil {
		log.Fatal().Msg("[FAIL] buffertest failed")
	} else {
		log.Info().Msg("[OK] buffertest successsful")
	}
}
