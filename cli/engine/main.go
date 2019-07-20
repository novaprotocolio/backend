package main

import (
	_ "github.com/joho/godotenv/autoload"
)

import (
	"context"
	"github.com/novaprotocolio/backend/cli"
	"github.com/novaprotocolio/backend/dex_engine"
	"github.com/novaprotocolio/sdk-backend/utils"
	"os"
)

func run() int {
	ctx, stop := context.WithCancel(context.Background())
	go cli.WaitExitSignal(stop)

	dex_engine.Run(ctx, utils.StartMetrics)
	return 0
}

func main() {
	os.Exit(run())
}
