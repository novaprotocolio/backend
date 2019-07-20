package api

import (
	"os"
	"strconv"
)

func GetSchemaVersion(p Param) (interface{}, error) {
	//req := p.(*DepositReq)

	schema, err := strconv.Atoi(os.Getenv("SWAP_SCHEMA_VERSION"))
	if err != nil {
		schema = -1
	}

	return &DepositResp{
		Schema: uint64(schema),
	}, nil
}
