package dydx_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	dydx "github.com/magiclars-off/go-dydx"
	"github.com/magiclars-off/go-dydx/helpers"
	"github.com/magiclars-off/go-dydx/private"
	"github.com/magiclars-off/go-dydx/realtime"
	"github.com/magiclars-off/go-dydx/types"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultHost     = "http://localhost:8080"
	EthereumAddress = ""
	StarkKey        = ""
)

var userID int64 = 0

var options = types.Options{
	NetworkId:                 types.NetworkIdGoerli,
	Host:                      types.ApiHostGoerli,
	DefaultEthereumAddress:    EthereumAddress,
	StarkPublicKey:            "",
	StarkPrivateKey:           "",
	StarkPublicKeyYCoordinate: "",
	ApiKeyCredentials: &types.ApiKeyCredentials{
		Key:        "",
		Secret:     "",
		Passphrase: "",
	},
}

func TestConnect(t *testing.T) {
	client := dydx.New(options)

	_, err := client.Private.GetAccount(client.Private.DefaultAddress)
	assert.NoError(t, err)

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	r := make(chan realtime.Response)

	go realtime.Connect(ctx, r, []string{realtime.ACCOUNT, realtime.TRADES}, []string{"BTC-USD"}, client.Private, nil)

	for {
		select {
		case v := <-r:
			switch v.Channel {
			case realtime.ACCOUNT:
				fmt.Println(v.Account)
			case realtime.TRADES:
				fmt.Println(v.Trades)
			case realtime.ERROR:
				log.Println(v.Results)
				goto EXIT
			}

		}
	}

EXIT:
	cancel()
}

func TestUsers(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Private.GetUsers()
	assert.NoError(t, err)

	fmt.Printf("%v", res)

	fmt.Printf("makerFee: %s, takeFee: %s\n", res.User.MakerFeeRate, res.User.TakerFeeRate)
}

func TestCreateOrder(t *testing.T) {
	client := dydx.New(options)

	acc, err := client.Private.GetAccount(client.Private.DefaultAddress)
	assert.NoError(t, err)
	positionId := acc.Account.PositionId

	o := &private.ApiOrder{
		ApiBaseOrder: private.ApiBaseOrder{Expiration: helpers.ExpireAfter(5 * time.Minute)},
		Market:       "ETH-USD",
		Side:         "BUY",
		Type:         "LIMIT",
		Size:         "1",
		Price:        "1000",
		ClientId:     helpers.RandomClientId(),
		TimeInForce:  "GTT",
		PostOnly:     true,
		LimitFee:     "0.01",
	}
	res, err := client.Private.CreateOrder(o, positionId)
	assert.NoError(t, err)

	fmt.Printf("%v", res)
}

// important!! Withdraw has not done any actual testing
// func TestWithdrawFast(t *testing.T) {
// 	client := dydx.New(options)
// 	res, err := client.Private.WithdrawFast(&private.WithdrawalParam{})
// 	assert.NoError(t, err)

// 	fmt.Printf("%v", res)
// }

func TestGetHistoricalPnL(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Private.GetHistoricalPnL(&private.HistoricalPnLParam{})
	assert.NoError(t, err)

	fmt.Printf("%v", res)
}

func TestGetTradingRewards(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Private.GetTradingRewards(&private.TradingRewardsParam{
		Epoch: 8,
	})
	assert.NoError(t, err)

	fmt.Printf("%v", res)
}
