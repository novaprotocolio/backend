package dex_engine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/novaprotocolio/backend/connection"
	"github.com/novaprotocolio/backend/models"
	"github.com/novaprotocolio/sdk-backend/common"
	"github.com/novaprotocolio/sdk-backend/engine"
	"github.com/novaprotocolio/sdk-backend/sdk/ethereum"
	"github.com/novaprotocolio/sdk-backend/utils"
	"os"
	"strings"
	"sync"
)

type RedisOrderbookSnapshotHandler struct {
	kvStore common.IKVStore
}

func (handler RedisOrderbookSnapshotHandler) Update(key string, bookSnapshot *common.SnapshotV2) sync.WaitGroup {
	bts, err := json.Marshal(bookSnapshot)
	if err != nil {
		panic(err)
	}

	_ = handler.kvStore.Set(key, string(bts), 0)

	return sync.WaitGroup{}
}

type RedisOrderbookActivitiesHandler struct {
}

func (handler RedisOrderbookActivitiesHandler) Update(webSocketMessages []common.WebSocketMessage) sync.WaitGroup {
	for _, msg := range webSocketMessages {
		if strings.HasPrefix(msg.ChannelID, "Market#") {
			pushMessage(msg)
		}
	}

	return sync.WaitGroup{}
}

type DexEngine struct {
	// global ctx, if this ctx is canceled, queue handlers should exit in a short time.
	ctx context.Context

	// all redis queues handlers
	marketHandlerMap map[string]*MarketHandler
	eventQueue       common.IQueue

	// Wait for all queue handler exit gracefully
	Wg sync.WaitGroup

	NovaEngine *engine.Engine
}

func NewDexEngine(ctx context.Context) *DexEngine {
	// init redis
	redis := connection.NewRedisClient(os.Getenv("NSK_REDIS_URL"))

	// init websocket queue
	wsQueue, _ := common.InitQueue(
		&common.RedisQueueConfig{
			Name:   common.NOVA_WEBSOCKET_MESSAGES_QUEUE_KEY,
			Ctx:    ctx,
			Client: redis,
		},
	)
	InitWsQueue(wsQueue)

	// init event queue
	eventQueue, _ := common.InitQueue(
		&common.RedisQueueConfig{
			Name:   common.NOVA_ENGINE_EVENTS_QUEUE_KEY,
			Client: redis,
			Ctx:    ctx,
		})

	e := engine.NewEngine(context.Background())

	// setup handler for nova engine
	kvStore, _ := common.InitKVStore(&common.RedisKVStoreConfig{Ctx: ctx, Client: redis})
	snapshotHandler := RedisOrderbookSnapshotHandler{kvStore: kvStore}
	e.RegisterOrderbookSnapshotHandler(snapshotHandler)

	activityHandler := RedisOrderbookActivitiesHandler{}
	e.RegisterOrderbookActivitiesHandler(activityHandler)

	engine := &DexEngine{
		ctx:              ctx,
		eventQueue:       eventQueue,
		marketHandlerMap: make(map[string]*MarketHandler),
		Wg:               sync.WaitGroup{},

		NovaEngine: e,
	}

	markets := models.MarketDao.FindPublishedMarkets()
	for _, market := range markets {
		_, err := engine.newMarket(market.ID)
		if err != nil {
			panic(err)
		}
	}

	return engine
}

func (e *DexEngine) newMarket(marketId string) (marketHandler *MarketHandler, err error) {
	_, ok := e.marketHandlerMap[marketId]

	if ok {
		err = fmt.Errorf("open market fail, market [%s] already exist", marketId)
		return
	}

	market := models.MarketDao.FindMarketByID(marketId)
	if market == nil {
		err = fmt.Errorf("open market fail, market [%s] not found", marketId)
		return
	}

	if !market.IsPublished {
		err = fmt.Errorf("open market fail, market [%s] not published", marketId)
		return
	}

	marketHandler, err = NewMarketHandler(e.ctx, market, e.NovaEngine)
	if err != nil {
		return
	}

	e.marketHandlerMap[market.ID] = marketHandler
	utils.Infof("market %s init done", marketHandler.market.ID)
	return
}

func (e *DexEngine) closeMarket(marketId string) {
	_, ok := e.marketHandlerMap[marketId]
	if !ok {
		utils.Errorf("close market fail, market [%s] not found", marketId)
		return
	}

	marketHandler := e.marketHandlerMap[marketId]
	delete(e.marketHandlerMap, marketId)
	marketHandler.Stop()
	return
}

func runMarket(e *DexEngine, marketHandler *MarketHandler) {
	e.Wg.Add(1)

	go func() {
		defer e.Wg.Done()

		utils.Infof("%s market handler is running", marketHandler.market.ID)
		defer utils.Infof("%s market handler is stopped", marketHandler.market.ID)

		marketHandler.Run()
	}()
}

func (e *DexEngine) start() {
	for i := range e.marketHandlerMap {
		marketHandler := e.marketHandlerMap[i]
		runMarket(e, marketHandler)
	}

	go func() {
		for {
			select {
			case <-e.ctx.Done():
				for _, handler := range e.marketHandlerMap {
					close(handler.eventChan)
				}
				return
			default:
				data, err := e.eventQueue.Pop()
				if err != nil {
					panic(err)
				}
				var event common.Event
				err = json.Unmarshal(data, &event)
				if err != nil {
					utils.Errorf("wrong event format: %+v", err)
					continue
				}

				switch event.Type {
				case common.EventOpenMarket:
					marketHandler, err := e.newMarket(event.MarketID)
					if err == nil {
						runMarket(e, marketHandler)
					} else {
						utils.Errorf(err.Error())
					}
					break
				case common.EventCloseMarket:
					e.closeMarket(event.MarketID)
					break
				default:
					marketHandler, ok := e.marketHandlerMap[event.MarketID]
					if !ok {
						utils.Errorf("engine not support market [%s]", event.MarketID)
					}
					marketHandler.eventChan <- data
				}
			}
		}
	}()
}

var novaProtocol = ethereum.NewEthereumNova(os.Getenv("NSK_BLOCKCHAIN_RPC_URL"), os.Getenv("NSK_HYBRID_EXCHANGE_ADDRESS"))
// var novaProtocol = &ethereum.EthereumNovaProtocol{}


func Run(ctx context.Context, startMetrics func()) {
	utils.Infof("dex engine start...")

	//init database
	models.Connect(os.Getenv("NSK_DATABASE_URL"))

	//start dex engine
	dexEngine := NewDexEngine(ctx)
	dexEngine.start()
	go startMetrics()

	dexEngine.Wg.Wait()
	utils.Infof("dex engine stopped!")
}
