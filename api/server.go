package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/novaprotocolio/backend/connection"
	"github.com/novaprotocolio/backend/models"
	"github.com/novaprotocolio/sdk-backend/common"
	"github.com/novaprotocolio/sdk-backend/sdk"
	"github.com/novaprotocolio/sdk-backend/sdk/ethereum"
	"github.com/novaprotocolio/sdk-backend/utils"
)

var CacheService common.IKVStore
var QueueService common.IQueue

func loadRoutes(e *echo.Echo) {
	e.Use(initNovaApiContext)

	webAsset := os.Getenv("NOVA_WEB_ASSETS")
	if webAsset == "" {
		e.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Novalex Decentralized Exchange (v1.0)!")
		})
	} else {
		e.File("/", webAsset+"/index.html")
		e.File("/favicon.ico", webAsset+"/favicon.ico")
		e.File("/config.js", webAsset+"/config.js")
		e.Static("/static", webAsset+"/static")
	}

	addRoute(e, "GET", "/v2/orderbook/:marketID/:orderID", &OrderbookV2Req{}, GetOrderbookV2)
	addRoute(e, "POST", "/v2/orderbook/processOrder", &OrderbookMsgV2Req{}, ProcessOrderbookV2)
	addRoute(e, "POST", "/v2/orderbook/cancelOrder", &OrderbookCancelMsgV2Req{}, CancelOrderbookV2)
	addRoute(e, "GET", "/v2/orderbook/bestAskList/:marketID", &OrderbookReq{}, BestAskListV2)
	addRoute(e, "GET", "/v2/orderbook/bestBidList/:marketID", &OrderbookReq{}, BestBidListV2)

	addRoute(e, "GET", "/markets", nil, GetMarkets)
	addRoute(e, "GET", "/markets/:marketID/orderbook", &OrderbookReq{}, GetOrderbook)
	addRoute(e, "GET", "/markets/:marketID/trades", &QueryTradeReq{}, GetAllTrades)

	addRoute(e, "GET", "/markets/:marketID/trades/mine", &QueryTradeReq{}, GetAccountTrades, authMiddleware)
	addRoute(e, "GET", "/markets/:marketID/candles", &CandlesReq{}, GetTradingView)
	addRoute(e, "GET", "/fees", &FeesReq{}, GetFees)

	addRoute(e, "GET", "/orders", &QueryOrderReq{}, GetOrders, authMiddleware)
	addRoute(e, "GET", "/orders/:orderID", &QuerySingleOrderReq{}, GetSingleOrder, authMiddleware)
	addRoute(e, "POST", "/orders/build", &BuildOrderReq{}, BuildOrder, authMiddleware)
	addRoute(e, "POST", "/orders", &PlaceOrderReq{}, PlaceOrder, authMiddleware)
	addRoute(e, "DELETE", "/orders/:orderID", &CancelOrderReq{}, CancelOrder, authMiddleware)
	addRoute(e, "GET", "/account/lockedBalances", &LockedBalanceReq{}, GetLockedBalance, authMiddleware)

	addRoute(e, "GET", "/deposit/schema", &DepositGetSchemaResq{}, GetSchemaVersion)
	addRoute(e, "GET", "/deposit/generate-address", &DepositGenAddrResq{}, GetGenerateAddress)
	addRoute(e, "GET", "/deposit/history", &DepositHistoryResq{}, GetDepositHistory)
}

func addRoute(e *echo.Echo, method, url string, param Param, handler func(p Param) (interface{}, error), middlewares ...echo.MiddlewareFunc) {
	e.Add(method, url, commonHandler(param, handler), middlewares...)
}

type Response struct {
	Status int         `json:"status"`
	Desc   string      `json:"desc"`
	Data   interface{} `json:"data,omitempty"`
}

func commonResponse(c echo.Context, data interface{}) error {
	return c.String(http.StatusOK, utils.ToJsonString(Response{
		Status: 0,
		Desc:   "success",
		Data:   data,
	}))
}

func errorHandler(err error, c echo.Context) {
	e := c.Echo()

	var status int
	var desc string

	if apiError, ok := err.(*ApiError); ok {
		status = apiError.Code
		desc = apiError.Desc
	} else if errors, ok := err.(validator.ValidationErrors); ok {
		status = -1
		desc = buildErrorMessage(errors)
	} else if e.Debug {
		status = -1
		desc = err.Error()
	} else {
		status = -1
		fmt.Println("err:", err)
		desc = "something wrong"
	}

	// Send response
	if !c.Response().Committed {
		err = c.JSON(http.StatusOK, Response{
			Status: status,
			Desc:   desc,
		})

		if err != nil {
			e.Logger.Error(err)
		}
	}
}

func getEchoServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	// open Debugf will return server errors details in json body
	// e.Debugf = true

	e.HTTPErrorHandler = errorHandler
	e.Use(recoverHandler)

	if os.Getenv("NSK_LOG_LEVEL") == "INFO" {
		e.Use(middleware.Logger())
	} else if os.Getenv("NSK_LOG_LEVEL") == "DEBUG" {
		// The BodyDump plugin is used for debug, you can uncomment these lines to see request and response body
		// More details goes https://echo.labstack.com/middleware/body-dump
		//
		e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			utils.Debugf("Header: %s", c.Request().Header)
			utils.Debugf("Url: %s, Request Body: %s; Response Body: %s", c.Request().RequestURI, string(reqBody), string(resBody))
		}))
	}	

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Jwt-Authentication", "Nova-Authentication"},
	}))

	loadRoutes(e)

	return e
}

var nova sdk.Nova

func StartServer(ctx context.Context, startMetric func()) {

	// init redis
	redisClient := connection.NewRedisClient(os.Getenv("NSK_REDIS_URL"))

	// init blockchain
	nova = ethereum.NewEthereumNova(os.Getenv("NSK_BLOCKCHAIN_RPC_URL"), os.Getenv("NSK_HYBRID_EXCHANGE_ADDRESS"))

	//init database
	models.Connect(os.Getenv("NSK_DATABASE_URL"))

	//Config key generator
	models.ConfigKeyGenerator(os.Getenv("NSK_MASTER_KEY"))

	CacheService, _ = common.InitKVStore(
		&common.RedisKVStoreConfig{
			Ctx:    ctx,
			Client: redisClient,
		},
	)

	QueueService, _ = common.InitQueue(
		&common.RedisQueueConfig{
			Name:   common.NOVA_ENGINE_EVENTS_QUEUE_KEY,
			Ctx:    ctx,
			Client: redisClient,
		},
	)

	addr := ":3001"
	if port, ok := os.LookupEnv("NSK_PORT"); ok {
		addr = ":" + port
	}

	e := getEchoServer()

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

	go startMetric()
	<-ctx.Done()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func recoverHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				stack := make([]byte, 2048)
				length := runtime.Stack(stack, false)
				utils.Errorf("unhandled error: %v %s", err, stack[:length])
				c.Error(err)
			}
		}()
		return next(c)
	}
}
