## Build docker

```bash
docker build -t novaprotocolio/nova-scaffold-dex-backend:latest .
```

## Hot reload

```bash
cd novalex-dex-starterkit
yarn watch api
```

## Try creating a new market

        go run cli/admincli/main.go market new HOT-Tri \
          --baseTokenAddress=0x4c4fa7e8ea4cfcfc93deae2c0cff142a1dd3a218 \
          --quoteTokenAddress=0x1D52a52f5996FDff37317a34EBFbeC7345Be3b55

- The base token is the first symbol (HOT above), the quote token is the second symbol (Tri above).
- You could try this with different token symbols and contract addresses.
- This creates a market with "default parameters" for fees, decimals, etc.

        go run cli/admincli/main.go market publish HOT-Tri

- This makes the market viewable on the frontend

## Setup environment with testing

Can run docker-compose anywhere within the project

```bash
# create schema
docker-compose exec db bash -c 'psql -U postgres -d postgres < /docker-entrypoint-initdb.d/0001-init.up.sql'

# seed data
docker-compose exec db bash -c 'psql -U postgres -d postgres < /docker-entrypoint-initdb.d/0002-seed.sql'

# run query
docker-compose exec db psql -U postgres -d postgres -c 'SELECT * FROM markets'

# run a more complex query
docker-compose exec db psql -U postgres -d postgres -c "select sum(amount) as locked_balance from orders where status='pending' and trader_address='0xe36ea790bc9d7ab70c55260c66d52b1eca985f84' and market_id like 'DAI-%' and side = 'sell'"

# can also update once using init sql file
docker-compose exec db psql -U postgres -d postgres -c "update markets set base_token_address='0x6F7ccbaCf134d826500ebCC574278cfC8aC5998d', quote_token_address='0x48690560139fCc885AD2B291f196c1908bc54281' where id='HOT-WETH'"

docker-compose exec db psql -U postgres -d postgres -c "SELECT * FROM markets where id='HOT-WETH'"

docker-compose exec db psql -U postgres -d postgres -c "update markets set base_token_address='0x6F7ccbaCf134d826500ebCC574278cfC8aC5998d', quote_token_address='0x31D7A88aF82D915eA4E74bbe1D95099546f596Cc' where id='HOT-DAI'"

docker-compose exec db psql -U postgres -d postgres -c "SELECT * FROM markets where id='HOT-DAI'"

docker-compose exec db psql -U postgres -d postgres -c "update markets set base_token_address='0x48690560139fCc885AD2B291f196c1908bc54281', quote_token_address='0x31D7A88aF82D915eA4E74bbe1D95099546f596Cc' where id='WETH-DAI'"

docker-compose exec db psql -U postgres -d postgres -c "SELECT * FROM markets where id='WETH-DAI'"

docker-compose exec db psql -U postgres -d postgres -c "update tokens set address='0x6F7ccbaCf134d826500ebCC574278cfC8aC5998d' where symbol='HOT'"

docker-compose exec db psql -U postgres -d postgres -c "update tokens set address='0x31D7A88aF82D915eA4E74bbe1D95099546f596Cc' where symbol='DAI'"

docker-compose exec db psql -U postgres -d postgres -c "update tokens set address='0x48690560139fCc885AD2B291f196c1908bc54281' where symbol='WETH'"
```
