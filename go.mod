module github.com/novaprotocolio/backend

go 1.13

require (
	// github.com/ethereum/go-ethereum v1.9.6
	github.com/ethereum/go-ethereum v0.0.0-00010101000000-000000000000
	github.com/go-playground/validator v9.29.0+incompatible
	github.com/go-redis/redis v6.15.1+incompatible
	github.com/jinzhu/gorm v1.9.11
	github.com/joho/godotenv v1.3.0 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/lib/pq v1.2.0
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/novaprotocolio/orderbook v0.0.0-00010101000000-000000000000
	github.com/novaprotocolio/sdk-backend v0.0.0-00010101000000-000000000000
	github.com/satori/go.uuid v1.2.0
	github.com/shopspring/decimal v0.0.0-20191009025716-f1972eb1d1f5
	github.com/stretchr/testify v1.4.0
	gopkg.in/go-playground/validator.v9 v9.30.0
)

replace github.com/uber-go/atomic => go.uber.org/atomic v1.5.0

// for local test only
replace github.com/ethereum/go-ethereum => ../novalex // v1.8.27

replace github.com/novaprotocolio/sdk-backend => ../sdk-backend

replace github.com/novaprotocolio/orderbook => ../orderbook

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20190930215403-16217165b5de

replace github.com/urfave/cli => gopkg.in/urfave/cli.v1 v1.20.0
