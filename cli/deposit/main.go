package main

import (
	"context"

	"os"

	"fmt"

	"github.com/novaprotocolio/backend/cli"
	"github.com/novaprotocolio/backend/connection"
	"github.com/novaprotocolio/backend/models"
	"github.com/novaprotocolio/sdk-backend/common"
	"github.com/novaprotocolio/sdk-backend/sdk"
	"github.com/novaprotocolio/sdk-backend/sdk/ethereum"
	"github.com/novaprotocolio/sdk-backend/utils"
	"github.com/novaprotocolio/sdk-backend/watcher"
	"github.com/ethereum/go-ethereum/crypto"

	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/onrik/ethrpc"
)

type DBTransactionHandler struct {
	w watcher.Watcher
}

func (handler DBTransactionHandler) Update(tx sdk.Transaction, timestamp uint64) {
	// hash := tx.GetHash()
	// txReceipt, _ := handler.w.Nova.GetTransactionReceipt(hash)
	// result := txReceipt.GetResult()

	recipientAddr := tx.GetTo()
	comparedAddr := strings.ToLower("0xFcc9d477AF8A7FE823Ecb24bbd541e779aa72F31")
	if recipientAddr == comparedAddr {
		// deposit
		utils.Infof("got Nova from %s", tx.GetFrom())
		deposit("0x94569C5a6B59003129d598dDc4060cF50908C980", handler.w.Nova)
	}

	// err := handler.w.QueueClient.Push(bts)

	// if err != nil {
	// 	utils.Errorf("Push event into Queue Errorf %v", err)
	// }

}

func deposit(toAddr string, nova sdk.Nova) {

	// send 1 ether
	value := ethrpc.Eth1()

	senderPrivKey, _ := crypto.HexToECDSA(os.Getenv("NSK_RELAYER_PK"))

	hash, err := nova.SendTransaction(toAddr, value, nil, senderPrivKey)

	if err != nil {
		fmt.Printf("error :%s \n", err.Error())
	} else {
		fmt.Println("hash", hash)
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
	nova := ethereum.NewEthereumNova(os.Getenv("ETHEREUM_BLOCKCHAIN_RPC_URL"), os.Getenv("NSK_HYBRID_EXCHANGE_ADDRESS"))
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
		Name:   common.NOVA_ENGINE_DEPOSIT_QUEUE_KEY,
		Client: client,
		Ctx:    ctx,
	})

	if err != nil {
		panic(err)
	}

	w := watcher.Watcher{
		Ctx:         ctx,
		Nova:        nova,
		KVClient:    kvStore,
		QueueClient: queue,
	}

	w.RegisterHandler(DBTransactionHandler{w})

	go utils.StartMetrics()

	// just for debug from start
	// w.KVClient.Set(common.NOVA_WATCHER_BLOCK_NUMBER_CACHE_KEY, "1", 0)

	w.Run()

	utils.Infof("Watcher Exit")
}
