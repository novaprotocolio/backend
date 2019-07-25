module github.com/novaprotocolio/backend

go 1.12

// replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190419153524-e8e3143a4f4a

// replace gopkg.in/go-playground/validator.v9 => github.com/go-playground/validator v9.28.0+incompatible

// replace gopkg.in/fsnotify.v1 => github.com/fsnotify/fsnotify v1.4.7

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422

// for local test only
replace github.com/ethereum/go-ethereum => ../novalex // v1.8.27

replace github.com/novaprotocolio/sdk-backend => ../sdk-backend

replace github.com/novaprotocolio/orderbook => ../orderbook

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/ethereum/go-ethereum v0.0.0-00010101000000-000000000000
	github.com/go-playground/validator v9.29.0+incompatible
	github.com/go-redis/redis v6.15.1+incompatible
	github.com/jinzhu/gorm v1.9.4
	github.com/jinzhu/now v1.0.0 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/lib/pq v1.0.0
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/novaprotocolio/orderbook v0.0.0-00010101000000-000000000000
	github.com/novaprotocolio/sdk-backend v0.0.39
	github.com/onrik/ethrpc v0.0.0-20190213081453-aa076c1849e6
	github.com/satori/go.uuid v1.2.0
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/stretchr/testify v1.3.0
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/urfave/cli v1.20.0
	go.uber.org/atomic v1.4.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.28.0
	gotest.tools v2.2.0+incompatible // indirect
)
