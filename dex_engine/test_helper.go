package dex_engine

import "os"

var User1PrivateKey string

func setEnvs() {
	_ = os.Setenv("NSK_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&TimeZone=Europe/Paris")
	_ = os.Setenv("NSK_REDIS_URL", "redis://redis:6379/0")
	_ = os.Setenv("NSK_BLOCKCHAIN_RPC_URL", "http://127.0.0.1:8545")
	_ = os.Setenv("NSK_WETH_TOKEN_ADDRESS", "0x7Cb242e4f8EE3FE4F1f244209c05B794F192353E")
	_ = os.Setenv("NSK_USD_TOKEN_ADDRESS", "0xbc3524faa62d0763818636d5e400f112279d6cc0")
	_ = os.Setenv("NSK_NOVA_TOKEN_ADDRESS", "0x224E34A640FC4108FABDb201eD85D909059105fA")
	_ = os.Setenv("NSK_PROXY_ADDRESS", "0x1D52a52f5996FDff37317a34EBFbeC7345Be3b55")
	_ = os.Setenv("NSK_HYBRID_EXCHANGE_ADDRESS", "0x179fd00c328d4ecdb5043c8686d377a24ede9d11")
	_ = os.Setenv("NSK_PROXY_MODE", "deposit")
	_ = os.Setenv("NSK_LOG_LEVEL", "DEBUG")
	_ = os.Setenv("NSK_RELAYER_ADDRESS", "0x93388b4efe13b9b18ed480783c05462409851547")
	_ = os.Setenv("NSK_RELAYER_PK", "95b0a982c0dfc5ab70bf915dcf9f4b790544d25bc5e6cff0f38a59d0bba58651")
	_ = os.Setenv("NSK_CHAIN_ID", "50")
	_ = os.Setenv("NSK_WEB3_URL", "http://127.0.0.1:8545")

	User1PrivateKey = "b7a0c9d2786fc4dd080ea5d619d36771aeb0c8c26c290afd3451b92ba2b7bc2c"
}
