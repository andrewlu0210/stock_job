# 每日股價下載資訊

### 安裝
```Shell
> go get -u github.com/andrewlu0210/stock
```

### Create stock.yaml
stock.yaml放在執行目錄
```YAML
stock_db:
  host: xx.xx.xx.xx
  db_name: stockDB
  db_account: account
  db_passwd: passwd

csvRoot: /Users/user/Documents/stock_csv/stock_price
```

### Run Command
下載當日股價資訊
```
> go run main.go
```
下載特定日股價資訊
```
> go run main.go 20220615
```
從特定日下載股價資訊共10日
```
> go run main.go 20220615 10
```
