package api

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestCreateSchema(t *testing.T) {
	resp, _ := GetSchemaVersion(nil)

	schema := resp.(*DepositGetSchemaResp).Schema
	//t.Log("shema: " + string(schema))
	//assert.EqualValues(t, 1, 1)
	assert.EqualValues(t, os.Getenv("SWAP_SCHEMA_VERSION"), string(strconv.Itoa(int(schema))))
}

func TestGetGenerateAddress(t *testing.T) {

}
