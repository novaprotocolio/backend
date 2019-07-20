package main

import (
	"context"
	"encoding/json"

	"os"

	"github.com/novaprotocolio/backend/cli"
	"github.com/novaprotocolio/backend/connection"
	"github.com/novaprotocolio/backend/models"
	"github.com/novaprotocolio/sdk-backend/common"
	"github.com/novaprotocolio/sdk-backend/sdk"
	"github.com/novaprotocolio/sdk-backend/sdk/ethereum"
	"github.com/novaprotocolio/sdk-backend/utils"
	"github.com/novaprotocolio/sdk-backend/watcher"

	_ "github.com/joho/godotenv/autoload"
)

type DBTransactionHandler struct {
	w watcher.Watcher
}

func (handler DBTransactionHandler) Update(tx sdk.Transaction, timestamp uint64) {
	launchLog := models.LaunchLogDao.FindByHash(tx.GetHash())

	if launchLog == nil {
		utils.Debugf("Skip useless transaction %s", tx.GetHash())
		return
	}

	if launchLog.Status != common.STATUS_PENDING {
		utils.Infof("LaunchLog is not pending %s, skip", launchLog.Hash.String)
		return
	}

	if launchLog != nil {
		txReceipt, _ := handler.w.Nova.GetTransactionReceipt(tx.GetHash())
		result := txReceipt.GetResult()
		hash := tx.GetHash()
		transaction := models.TransactionDao.FindTransactionByID(launchLog.ItemID)
		utils.Infof("Transaction %s result is %+v", tx.GetHash(), result)

		var status string

		if result {
			status = common.STATUS_SUCCESSFUL
		} else {
			status = common.STATUS_FAILED
		}

		//approve event should not process with engine, so update and return
		if launchLog.ItemType == "novaApprove" {
			launchLog.Status = status
			err := models.LaunchLogDao.UpdateLaunchLog(launchLog)
			if err != nil {
				panic(err)
			}
			return
		}
		event := &common.ConfirmTransactionEvent{
			Event: common.Event{
				Type:     common.EventConfirmTransaction,
				MarketID: transaction.MarketID,
			},
			Hash:      hash,
			Status:    status,
			Timestamp: timestamp,
		}

		bts, _ := json.Marshal(event)

		err := handler.w.QueueClient.Push(bts)

		if err != nil {
			utils.Errorf("Push event into Queue Errorf %v", err)
		}
	}
}

func main() {
	ctx, stop := context.WithCancel(context.Background())

	go cli.WaitExitSignal(stop)

	// Init Database Client
	models.Connect(os.Getenv("NSK_DATABASE_URL"))

	// Init Redis client
	client := connection.NewRedisClient(os.Getenv("NSK_REDIS_URL"))

	// Init Blockchain Client
	nova := ethereum.NewEthereumNova(os.Getenv("NSK_BLOCKCHAIN_RPC_URL"), os.Getenv("NSK_HYBRID_EXCHANGE_ADDRESS"))
	if os.Getenv("NSK_LOG_LEVEL") == "DEBUG" {
		nova.EnableDebug(true)
	}

	// init Key/Value Store
	kvStore, err := common.InitKVStore(&common.RedisKVStoreConfig{
		Ctx:    ctx,
		Client: client,
	})

	if err != nil {
		panic(err)
	}

	// Init Queue
	// There is no block call of redis, so we share the client here.
	queue, err := common.InitQueue(&common.RedisQueueConfig{
		Name:   common.NOVA_ENGINE_EVENTS_QUEUE_KEY,
		Client: client,
		Ctx:    ctx,
	})

	if err != nil {
		panic(err)
	}

	w := watcher.Watcher{
		Ctx:         ctx,
		Nova:       nova,
		KVClient:    kvStore,
		QueueClient: queue,
	}

	w.RegisterHandler(DBTransactionHandler{w})

	go utils.StartMetrics()

	w.Run()

	utils.Infof("Watcher Exit")
}
