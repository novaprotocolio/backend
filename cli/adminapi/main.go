package main

import (
	"context"
	"github.com/novaprotocolio/backend/admin/api"
	"github.com/novaprotocolio/backend/cli"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

func run() int {
	ctx, stop := context.WithCancel(context.Background())

	go cli.WaitExitSignal(stop)
	adminapi.StartServer(ctx)

	return 0
}

func main() {
	os.Exit(run())
}
