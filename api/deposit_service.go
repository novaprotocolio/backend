package api

import (
	"errors"
	"github.com/novaprotocolio/backend/models"
	"os"
	"strconv"
)

const (
	PENDING = "PENDING"
	SUCCESS = "SUCCESS"
	FAILED  = "FAILED"
	DEPOSIT = "deposit"
)

func GetSchemaVersion(p Param) (interface{}, error) {
	//req := p.(*DepositReq)

	schema, err := strconv.Atoi(os.Getenv("SWAP_SCHEMA_VERSION"))
	if err != nil {
		schema = -1
	}

	return &DepositGetSchemaResp{
		Schema: uint64(schema),
	}, nil
}

func GetGenerateAddress(p Param) (interface{}, error) {
	baseAddr := p.(*DepositGenAddrResq).UserAddr
	if baseAddr == "" {
		return &DepositGenAddrResp{Status: "User addr is empty"}, nil
	}

	var rowFromDb models.AssociationRow
	models.DB.Table("deposit").Where("base_address = ?", baseAddr).First(&rowFromDb)

	if rowFromDb.BaseAddress != "" && rowFromDb.Status == PENDING {
		addr := rowFromDb.QuoteAddress
		return &DepositGenAddrResp{
			Status:        "Transaction is pending",
			GeneratedAddr: addr}, nil
	} else if rowFromDb.BaseAddress != "" {
		models.DB.Table(DEPOSIT).Delete(rowFromDb)
	}

	row := models.AssociationRow{BaseAddress: baseAddr, Status: PENDING}
	models.DB.Table(DEPOSIT).Create(&row)

	id := row.ID
	address, err := models.KeyGenerator.Generate(id)
	if err != nil {
		return nil, errors.New("Error when generate key: " + err.Error())
	}
	models.DB.Table(DEPOSIT).Model(&row).Update("quote_address", address.String())

	return &DepositGenAddrResp{
		GeneratedAddr: row.QuoteAddress,
		Status:        "Key is generated successful",
	}, nil
}

func GetDepositHistory(p Param) (interface{}, error) {
	baseAddr := p.(*DepositHistoryResq).UserAddr
	var rowFromDb models.AssociationRow
	models.DB.Table(DEPOSIT).Where("base_address = ?", baseAddr).First(&rowFromDb)

	var history []string
	if rowFromDb.BaseAddress != "" {
		history = []string(rowFromDb.TxtReciver)
	}

	return &DepositHistoryResp{
		History: history,
	}, nil
}
