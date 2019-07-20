package dex_engine

import (
	"github.com/novaprotocolio/backend/models"
	"github.com/novaprotocolio/sdk-backend/sdk"
	"github.com/novaprotocolio/sdk-backend/utils"
)

func getNovaOrderFromModelOrder(orderJSON *models.OrderJSON) *sdk.Order {
	return sdk.NewOrderWithData(
		orderJSON.Trader,
		orderJSON.Relayer,
		orderJSON.BaseCurrency,
		orderJSON.QuoteCurrency,
		utils.DecimalToBigInt(orderJSON.BaseCurrencyHugeAmount),
		utils.DecimalToBigInt(orderJSON.QuoteCurrencyHugeAmount),
		utils.DecimalToBigInt(orderJSON.GasTokenHugeAmount),
		orderJSON.Data,
		orderJSON.Signature,
	)
}

func getNovaOrderHashHexFromOrderJson(orderJSON *models.OrderJSON) string {
	order := sdk.NewOrderWithData(
		orderJSON.Trader,
		orderJSON.Relayer,
		orderJSON.BaseCurrency,
		orderJSON.QuoteCurrency,
		utils.DecimalToBigInt(orderJSON.BaseCurrencyHugeAmount),
		utils.DecimalToBigInt(orderJSON.QuoteCurrencyHugeAmount),
		utils.DecimalToBigInt(orderJSON.GasTokenHugeAmount),
		orderJSON.Data,
		"",
	)

	return utils.Bytes2HexP(novaProtocol.GetOrderHash(order))
}
