package main

import (
	"context"
	"github.com/novaprotocolio/backend/api"
	"github.com/novaprotocolio/backend/cli"
	"github.com/novaprotocolio/sdk-backend/utils"
	_ "github.com/joho/godotenv/autoload"
	"os"
)


func run() int {
	ctx, stop := context.WithCancel(context.Background())	
	go cli.WaitExitSignal(stop)
	api.StartServer(ctx, utils.StartMetrics)

	return 0
}

func main() {
	os.Exit(run())
}
