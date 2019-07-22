**Hot reload**

```bash
cd novalex-dex-starterkit
yarn watch api
```

**Try creating a new market**

        go run cli/admincli/main.go market new HOT-Tri \
          --baseTokenAddress=0x224E34A640FC4108FABDb201eD85D909059105fA \
          --quoteTokenAddress=0x1D52a52f5996FDff37317a34EBFbeC7345Be3b55

- The base token is the first symbol (HOT above), the quote token is the second symbol (Tri above).
- You could try this with different token symbols and contract addresses.
- This creates a market with "default parameters" for fees, decimals, etc.

        go run cli/admincli/main.go market publish HOT-Tri

- This makes the market viewable on the frontend
