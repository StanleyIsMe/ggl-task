package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"ggltask/internal/api"
	apiCfg "ggltask/internal/api/config"
	"ggltask/pkg/config"
	"ggltask/pkg/shutdown"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const shutdownGracePeriod = 30 * time.Second

// @title           Gogolook Task API
// @version         1.0
// @description     API Server for Gogolook interview task

// @contact.name   Stanley Hsieh
// @contact.email  grimmh6838@gmail.com

// @license.name  GNU General Public License v3.0
// @license.url   https://www.gnu.org/licenses/gpl-3.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	mainCtx, mainStopCtx := context.WithCancel(context.Background())

	cfg, err := config.LoadWithEnv[apiCfg.Config](mainCtx, "./config/api")
	if err != nil {
		log.Fatal().Err(err).Msg("load config failed")
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if cfg.PrettyLog {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	shutdownHandler := shutdown.New(
		&logger,
		shutdown.WithGracePeriodDuration(shutdownGracePeriod),
	)

	app := api.NewAPI(cfg, shutdownHandler, &logger)
	if err = app.Start(mainCtx); err != nil {
		logger.Fatal().Err(err).Msg("api serve failed with an error")
	}

	if err := shutdownHandler.Listen(
		mainCtx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	); err != nil {
		logger.Fatal().Err(err).Msg("graceful shutdown failed.. forcing exit.")
	}

	mainStopCtx()
}
