package main

import (
	"context"
	"github.com/novaprotocolio/backend/cli"
	"github.com/novaprotocolio/backend/connection"
	"github.com/novaprotocolio/sdk-backend/common"
	"github.com/novaprotocolio/sdk-backend/utils"
	"github.com/novaprotocolio/sdk-backend/websocket"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := context.WithCancel(context.Background())

	redisClient := connection.NewRedisClient(os.Getenv("NSK_REDIS_URL"))
	redisClient = redisClient.WithContext(ctx)

	go cli.WaitExitSignal(stop)

	// new a source queue
	queue, err := common.InitQueue(&common.RedisQueueConfig{
		Name:   common.NOVA_WEBSOCKET_MESSAGES_QUEUE_KEY,
		Ctx:    ctx,
		Client: redisClient,
	})

	if err != nil {
		panic(err)
	}

	addr := ":3002"
	if port, ok := os.LookupEnv("NSK_PORT"); ok {
		addr = ":" + port
	}

	// new a websocket server
	wsServer := websocket.NewWSServer(addr, queue)

	websocket.RegisterChannelCreator(
		common.MarketChannelPrefix,
		websocket.NewMarketChannelCreator(&websocket.DefaultHttpSnapshotFetcher{
			ApiUrl: os.Getenv("NSK_API_URL"),
		}),
	)

	// Start the server
	// It will block the current process to listen on the `addr` your provided.
	go utils.StartMetrics()
	wsServer.Start(ctx)

	return 0
}
