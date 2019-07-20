package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"	
)

func GetOrderbookV2(p Param) (interface{}, error) {

	params := p.(*OrderbookV2Req)	
	fmt.Printf("param : marketID:%v, orderID:%v\n", params.MarketID, params.OrderID)
	orderbook := nova.GetOrderbook(params.MarketID, params.OrderID)

	var snapshot map[string]interface{}
	err := json.Unmarshal(orderbook, &snapshot)

	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func ProcessOrderbookV2(p Param) (interface{}, error) {

	params := p.(*OrderbookMsgV2Req)		
	
	payload := params.ToQuote()
	payload["timestamp"] = strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)	
	fmt.Printf("payload : %v\n", payload)
	var snapshot map[string]interface{}
	
	result := nova.ProcessOrderbook(payload)	

	err := json.Unmarshal(result, &snapshot)

	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func CancelOrderbookV2(p Param) (interface{}, error) {

	params := p.(*OrderbookCancelMsgV2Req)		
	var ret = false
	payload := params.ToQuote()
	fmt.Printf("payload : %v\n", payload)
	ret = nova.CancelOrderbook(payload)
	return ret, nil
}

func BestAskListV2(p Param)(interface{}, error) {
	params := p.(*OrderbookReq)		
	
	var snapshot []map[string]interface{}
	
	result := nova.BestAskList(params.MarketID)	

	err := json.Unmarshal(result, &snapshot)

	

	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func BestBidListV2(p Param)(interface{}, error) {
	params := p.(*OrderbookReq)			
	
	var snapshot []map[string]interface{}
	
	result := nova.BestBidList(params.MarketID)	

	err := json.Unmarshal(result, &snapshot)

	if err != nil {
		return nil, err
	}

	return snapshot, nil
}
