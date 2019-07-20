package main

import (
	_ "github.com/joho/godotenv/autoload"
)

import (
	"context"
	"github.com/NovaProtocol/backend/cli"
	"github.com/NovaProtocol/backend/dex_engine"
	"github.com/NovaProtocol/sdk-backend/utils"
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
