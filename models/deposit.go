package models

import (
	"github.com/lib/pq"
	"github.com/novaprotocolio/sdk-backend/utils"
)

var KeyGenerator *utils.AddressGenerator

type AssociationRow struct {
	ID                 uint64
	BaseAddress        string
	QuoteAddress       string
	Status             string
	BlockNumber        int
	CurrentBlockNumber int
	TxtReciver         pq.StringArray
}

func ConfigKeyGenerator(masterKey string) {
	var err error
	KeyGenerator, err = utils.NewAddressGenerator(masterKey)
	if err != nil {
		panic("Error when create KeyGenerator from master key: " + masterKey + " " + err.Error())
	}
}
