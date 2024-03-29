package adminapi

import (
	"context"
	"github.com/novaprotocolio/backend/connection"
	"github.com/novaprotocolio/backend/models"
	"github.com/novaprotocolio/sdk-backend/common"
	"github.com/novaprotocolio/sdk-backend/sdk/ethereum"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"os"
	"time"	
)

var queueService common.IQueue
var healthCheckService IHealthCheckMonitor
var erc20Service ethereum.IErc20

func loadRoutes(e *echo.Echo) {
	e.Add("GET", "/markets", ListMarketsHandler)
	e.Add("POST", "/markets", CreateMarketHandler)
	e.Add("POST", "/markets/approve", ApproveMarketHandler)
	e.Add("PUT", "/markets", EditMarketHandler)
	e.Add("DELETE", "/orders/:order_id", DeleteOrderHandler)
	e.Add("GET", "/orders", GetOrdersHandler)
	e.Add("GET", "/trades", GetTradesHandler)
	e.Add("GET", "/balances", GetBalancesHandler)
	e.Add("GET", "/status", GetStatusHandler)
	e.Add("POST", "/restart_engine", RestartEngineHandler)
}

func newEchoServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	loadRoutes(e)

	return e
}

func StartServer(ctx context.Context) {
	//init database
	models.Connect(os.Getenv("NSK_DATABASE_URL"))

	//init health check service
	healthCheckService = NewHealthCheckService(nil)

	//init erc20 service
	erc20Service = ethereum.NewErc20Service(nil)

	//init event queue
	queueService, _ = common.InitQueue(
		&common.RedisQueueConfig{
			Name:   common.NOVA_ENGINE_EVENTS_QUEUE_KEY,
			Ctx:    ctx,
			Client: connection.NewRedisClient(os.Getenv("NSK_REDIS_URL")),
		},
	)

	addr := ":3003"
	if port, ok := os.LookupEnv("NSK_PORT"); ok {
		addr = ":" + port
	}

	e := newEchoServer()
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("shutting down the server: %v", err)
			panic(err)
		}
	}()

	<-ctx.Done()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
