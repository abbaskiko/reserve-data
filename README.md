# Data fetcher for KyberNetwork reserve
[![Go Report Card](https://goreportcard.com/badge/github.com/KyberNetwork/reserve-data)](https://goreportcard.com/report/github.com/KyberNetwork/reserve-data)
[![Build Status](https://travis-ci.org/KyberNetwork/reserve-data.svg?branch=develop)](https://travis-ci.org/KyberNetwork/reserve-data)

This repo is contains two components:

- core: 
	- interacts with blockchain to get/set rates for tokens pair
	- buy/sell with centralized exchanges (binance, huobi, etc)
(For more detail, take a look to interface ReserveCore in intefaces.go)

- stat:  
	- fetch tradelogs from blockchain and do aggregation and save its data to database and allow client to query

(For more detail, find ReserveStat interface in interfaces.go)

## Compile it

```shell
cd cmd && go build -v
```

a `cmd` executable file will be created in `cmd` module.

## Run the reserve data

1. You need to prepare a `config.json` file inside `cmd` module. The file is described in later section.
2. You need to prepare a JSON keystore file inside `cmd` module. It is the keystore for the reserve owner.
3. Make sure your working directory is `cmd`. Run `KYBER_EXCHANGES=binance,huobi ./cmd` in dev mode.

### Manual

```shell
cd cmd
```

- Run core only

```shell
KYBER_EXCHANGES="binance,huobi" KYBER_ENV=production ./cmd server --log-to-stdout
```

- Run stat only

```shell
KYBER_ENV=production ./cmd server --log-to-stdout --enable-stat --no-core
```

### Docker (recommended)

This repository will build docker images and public on [docker hub](https://hub.docker.com/r/kybernetwork/reserve-data/tags/), you can pull image from docker hub and run:

- Run core only

```shell
docker run -p 8000:8000 -v /location/of/config.json:/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json -e KYBER_EXCHANGES="binance,huobi" KYBER_ENV="production" kybernetwork/reserve-data:develop server --log-to-stdout
```

- Run stat only 

```shell
docker run -p 8000:8000 -v /location/of/config.json:/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json -e KYBER_ENV="production" kybernetwork/reserve-data:develop server --enable-stat --no-core --log-to-stdout
```

**Note** : 

- KYBER_ENV includes "dev, simulation and production", different environment mode uses different settings (check cmd folder for settings file).  

- reserve-data requires config.json file to run, so you need to -v (mount config.json file to docker) so it can run.

## Config file

sample:

```json
{
  "binance_key": "your binance key",
  "binance_secret": "your binance secret",
  "huobi_key": "your huobi key",
  "huobi_secret_key": "your huobi secret",
  "kn_secret": "secret key for people to sign their requests to our apis. It is ignored in dev mode.",
  "kn_readonly": "read only key for people to sign their requests, this key can read everything but cannot execute anything",
  "kn_configuration": "key for people to sign their requests, this key can read everything and set configuration such as target quantity",
  "kn_confirm_configuration": "key for people to sign ther requests, this key can read everything and confirm target quantity, enable/disable setrate or rebalance",
  "keystore_path": "path to the JSON keystore file, recommended to be absolute path",
  "passphrase": "passphrase to unlock the JSON keystore",
  "keystore_deposit_path": "path to the JSON keystore file that will be used to deposit",
  "passphrase_deposit": "passphrase to unlock the JSON keystore",
  "keystore_intermediator_path": "path to JSON keystore file that will be used to deposit to Huobi",
  "passphrase_intermediate_account": "passphrase to unlock JSON keystore",
  "aws_access_key_id": "your aws key ID",
  "aws_secret_access_key": "your aws scret key",
  "aws_expired_stat_data_bucket_name" : "AWS bucket for expired stat data (already created)",
  "aws_expired_reserve_data_bucket_name" : "AWS bucket for expired reserve data (already created)",
  "aws_log_bucket_name" :"AWS bucket for log backup(already created)",
  "aws_region":"AWS region"
}
```

## APIs

### Get token rates from blockchain

```shell
<host>:8000/getrates
```

eg:

```shell
curl -X GET "http://127.0.0.1:8000/getrates"
```
response:

```json
{
    "data": {
        "MCO": {
            "Valid": true,
            "Error": "",
            "Timestamp": "1515412582435",
            "ReturnTime": "1515412582710",
            "BaseBuy": 63.99319226272073,
            "CompactBuy": 21,
            "BaseSell": 0.014716371218820246,
            "CompactSell": -20,
            "Rate": 0,
            "Block": 2420849
        },
        "OMG": {
            "Valid": true,
            "Error": "",
            "Timestamp": "1515412582435",
            "ReturnTime": "1515412582710",
            "BaseBuy": 44.45707162223901,
            "CompactBuy": 30,
            "BaseSell": 0.021183301968644246,
            "CompactSell": -29,
            "Rate": 0,
            "Block": 2420849
        },
        "PAY": {
            "Valid": true,
            "Error": "",
            "Timestamp": "1515412582435",
            "ReturnTime": "1515412582710",
            "BaseBuy": 295.08854913901575,
            "CompactBuy": -13,
            "BaseSell": 0.003191406699999999,
            "CompactSell": 13,
            "Rate": 0,
            "Block": 2420849
        }
    },
    "success": true,
    "timestamp": "1515412583215",
    "version": 1515412582435
}
```


### Get all token rates from blockchain

```shell
<host>:8000/get-all-rates
```

url params:
*fromTime*: optional, get all rates from this timepoint (millisecond)
*toTime*: optional, get all rates to this timepoint (millisecond)

eg:

```shell
curl -X GET "http://127.0.0.1:8000/get-all-rates"
```

response

```json
{
    "data": [
        {
            "Version": 0,
            "Valid": true,
            "Error": "",
            "Timestamp": "1517280618739",
            "ReturnTime": "1517280619071",
            "Data": {
                "SNT": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280618739",
                    "ReturnTime": "1517280619071",
                    "BaseBuy": 4053.2170631085987,
                    "CompactBuy": 43,
                    "BaseSell": 0.000233599514875301,
                    "CompactSell": -3,
                    "Rate": 0,
                    "Block": 5635245
                }
            }
        },
        {
            "Version": 0,
            "Valid": true,
            "Error": "",
            "Timestamp": "1517280621738",
            "ReturnTime": "1517280622251",
            "Data": {
                "EOS": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 87.21360760013062,
                    "CompactBuy": 0,
                    "BaseSell": 0.0128686459657361,
                    "CompactSell": 0,
                    "Rate": 0,
                    "Block": 5635245
                },
                "ETH": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 0,
                    "CompactBuy": 32,
                    "BaseSell": 0,
                    "CompactSell": -14,
                    "Rate": 0,
                    "Block": 5635245
                },
                "KNC": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 307.05930436561505,
                    "CompactBuy": -34,
                    "BaseSell": 0.003084981280661941,
                    "CompactSell": 81,
                    "Rate": 0,
                    "Block": 5635245
                },
                "OMG": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 65.0580993582104,
                    "CompactBuy": 32,
                    "BaseSell": 0.014925950060437398,
                    "CompactSell": -14,
                    "Rate": 0,
                    "Block": 5635245
                },
                "SALT": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 152.3016783627643,
                    "CompactBuy": 9,
                    "BaseSell": 0.006196212698403499,
                    "CompactSell": 23,
                    "Rate": 0,
                    "Block": 5635245
                },
                "SNT": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 4053.2170631085987,
                    "CompactBuy": 43,
                    "BaseSell": 0.000233599514875301,
                    "CompactSell": -3,
                    "Rate": 0,
                    "Block": 5635245
                }
            }
        }
    ],
    "success": true
}
```


### Get trade history for an account (signing required)

```shell
  <host>:8000/tradehistory  
  params:
  - fromTime: millisecond (required)
  - toTime: millisecond (required)
  Restriction: toTime - fromTime <= 3 days (in millisecond)
```

eg:

```shell
curl -X GET "http://localhost:8000/tradehistoryfromTime=1516116380102&toTime=18446737278344972745"
```

response:

```json
{"data":{"Version":1517298257114,"Valid":true,"Timestamp":"1517298257115","Data":{"binance":{"EOS-ETH":[],"KNC-ETH":[{"ID":"548002","Price":0.003038,"Qty":50,"Type":"buy","Timestamp":1516116380102},{"ID":"548003","Price":0.0030384,"Qty":7,"Type":"buy","Timestamp":1516116380102},{"ID":"548004","Price":0.003043,"Qty":16,"Type":"buy","Timestamp":1516116380102},{"ID":"548005","Price":0.0030604,"Qty":29,"Type":"buy","Timestamp":1516116380102},{"ID":"548006","Price":0.003065,"Qty":29,"Type":"buy","Timestamp":1516116380102},{"ID":"548007","Price":0.003065,"Qty":130,"Type":"buy","Timestamp":1516116380102}],"OMG-ETH":[{"ID":"123980","Price":0.020473,"Qty":48,"Type":"buy","Timestamp":1512395498231},{"ID":"130518","Price":0.021022,"Qty":13.49,"Type":"buy","Timestamp":1512564108827},{"ID":"130706","Price":0.020202,"Qty":9.93,"Type":"sell","Timestamp":1512569059460},{"ID":"140078","Price":0.019098,"Qty":11.07,"Type":"buy","Timestamp":1512714826339},{"ID":"140157","Price":0.019053,"Qty":7.68,"Type":"sell","Timestamp":1512716338997},{"ID":"295923","Price":0.020446,"Qty":4,"Type":"buy","Timestamp":1514360742162}],"SALT-ETH":[],"SNT-ETH":[]},"bittrex":{"OMG-ETH":[{"ID":"eb948865-6261-4991-8615-b36c8ccd1256","Price":0.01822057,"Qty":1,"Type":"buy","Timestamp":18446737278344972745}],"SALT-ETH":[],"SNT-ETH":[]}}},"success":true}
```



### Get exchange balances, reserve balances, pending activities at once (signing required)

```shell
<host>:8000/authdata
```

eg:

```shell
curl -X GET "http://localhost:8000/authdata"
```

response:

```json
{"data":{"Valid":true,"Error":"","Timestamp":"1514114408227","ReturnTime":"1514114408810","ExchangeBalances":{"bittrex":{"Valid":true,"Error":"","Timestamp":"1514114408226","ReturnTime":"1514114408461","AvailableBalance":{"ETH":0.10704306,"OMG":2.97381136},"LockedBalance":{"ETH":0,"OMG":0},"DepositBalance":{"ETH":0,"OMG":0}}},"ReserveBalances":{"ADX":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"BAT":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"CVC":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"DGD":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"EOS":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"ETH":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":360169992138038352},"FUN":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"GNT":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"KNC":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"LINK":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"MCO":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"OMG":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":23818094310417195708},"PAY":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0}},"PendingActivities":[]},"block": 2345678, "success":true,"timestamp":"1514114409088","version":39}
```

### Get all activityes (signing required)
```
<host>:8000/activities
GET request
url params: 
  fromTime: from timepoint - uint64, unix millisecond (optional if empty then get from first activity)
  toTime: to timepoint - uint64, unix millisecond (optional if empty then get to last activity)
```
Note: `fromTime` and `toTime` shouldn't be included into signing message.
### Get immediate pending activities (signing required)
```
<host>:8000/immediate-pending-activities
GET request
```

### Store processed data (signing required)
```
<host>:8000/metrics
POST request
form params:
  - timestamp: uint64, unix millisecond
  - data: string, in format of <token>_afpmid_spread|<token>_afpmid_spread|..., eg. OMG_0.4_5|KNC_1_2
```

### Get processed data (signing required)
```
<host>:8000/metrics
GET request
url params:
  - tokens: string, list of tokens to get data about, in format of <token_id>-<token_id>..., eg. OMG_DGD_KNC
  - from: uint64, unix millisecond
  - to: uint64, unix millisecond
```

response:
```json
{
    "data": {
        "DGD": [
            {
                "Timestamp": 19,
                "AfpMid": 4,
                "Spread": 5
            }
        ],
        "OMG": [
            {
                "Timestamp": 19,
                "AfpMid": 0.9,
                "Spread": 1
            }
        ]
    },
    "returnTime": 1514966512560,
    "success": true,
    "timestamp": 1514966512549
}
```
Returned data will only include datas that have timestamp in range of `[from, to]`


### Get pending pwis equation (signing required)
```
<host>:8000/pending-pwis-equation
GET request
```

response:
```
  {
    "success": true,
    "data":{"id":1517396850670,"data":"EOS_750_500_0.25|ETH_750_500_0.25|KNC_750_500_0.25|OMG_750_500_0.25|SALT_750_500_0.25"}
  }
```

### Get pwis equation (signing required)
```
<host>:8000/pwis-equation
GET request
```

response:
```
  {
    "success": true,
    "data":{"id":1517396850670,"data":"EOS_750_500_0.25|ETH_750_500_0.25|KNC_750_500_0.25|OMG_750_500_0.25|SALT_750_500_0.25"}
  }
```

### Set pwis equation (signing required)
```
<host>:8000/set-pwis-equation
POST request
form params:
  - data: required, string, must sort by token id by ascending order
  - id: optional, required to confirm target quantity
```
eg:
```
curl -X POST \
  http://localhost:8000/set-pwis-equation \
  -H 'content-type: multipart/form-data' \
  -F data= EOS_750_500_0.25|ETH_750_500_0.25|KNC_750_500_0.25|OMG_750_500_0.25|SALT_750_500_0.25 \
  -F id=1517396850670
```
response
```
  {
    "success": true,
  }
```

### Confirm pwis equation (signing required)
```
<host>:8000/confirm-pwis-equation
POST request
form params:
  - data: required, string, must sort by token id by ascending order
```
eg:
```
curl -X POST \
  http://localhost:8000/confirm-pwis-equation \
  -H 'content-type: multipart/form-data' \
  -F data=EOS_750_500_0.25|ETH_750_500_0.25|KNC_750_500_0.25|OMG_750_500_0.25|SALT_750_500_0.25
```
response
```
  {
    "success": true,
  }
```

### Get exchanges status
```
<host>:8000/get-exchange-status
GET request
```

eg:
```
curl -x GET http://localhost:8000/get-exchange-status
```

response:
```
{"data":{"binance":{"timestamp":1521532176702,"status":true},"bittrex":{"timestamp":1521532176704,"status":true},"huobi":{"timestamp":1521532176703,"status":true}},"success":true}
```

### Update exchanges status
```
<host>:8000/update-exchange-status
POST request

params: 
exchange (string): exchange name (eg: 'binance')
status (bool): true (up), false (down)
timestamp (integer): timestamp of the exchange status
```

eg:
```
curl -X POST \
  http://localhost:8000/update-exchange-status \
  -H 'content-type: multipart/form-data' \
  -F exchange=binance \
  -F status=false
```
### Update Price Analytic Data - (signing required) set a record marking the condition because of which the set price is called. 
```
<host>:8000/update-price-analytic-data
POST request
params:
 - timestamp - the timestamp of the action (real time ) in millisecond
 - value - the json enconded object to save. 

Note: the data sent over must be encoded in Json in order to make it valid for output operation
  In Python, the data would be encoded as:
   data = {"timestamp": timestamp, "value": json.dumps(analytic_data)} 
 ```

response:
```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```

### Get Price Analytic Data - (signing required) list of price analytic data, sorted by timestamp 
```
<host>:8000/get-get-price-analytic-data
GET request
params:
 - fromTime (integer) - from timestamp (millisecond)
 - toTime (integer) - to timestamp (millisecond)
```
example:
```
curl -x GET \
  http://localhost:8000/get-price-analytic-data?fromTime=1522753160000&toTime=1522755792000
```
 
response:
```
{
  "data": [
    {
      "Timestamp": 1522755271000,
      "Data": {
        "block_expiration": false,
        "trigger_price_update": true,
        "triggering_tokens_list": [
          {
            "ask_price": 0.002,
            "bid_price": 0.003,
            "mid afp_old_price": 0.34555,
            "mid_afp_price": 0.6555,
            "min_spread": 0.233,
            "token": "OMG"
          },
          {
            "ask_price": 0.004,
            "bid_price": 0.005,
            "mid afp_old_price": 0.21555,
            "mid_afp_price": 0.4355,
            "min_spread": 0.133,
            "token": "KNC"
          }
        ]
      }
    }
  ],
  "success": true
}
```

### Get exchange notifications
```
<host>:8000/exchange-notifications
GET request
```

response:
```
{"data":{"binance":{"trade":{"OMG":{"fromTime":123,"toTime":125,"isWarning":true,"msg":"3 times"}}}},"success":true}
```

### set stable token params - (signing required)
```
<host>:8000/set-stable-token-params
POST request
URL Params:
  - value (string) : the json enconded string, represent a map (string : interface)
```


response:
```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```
### confirm stable token params - (signing required)
```
<host>:8000/confirm-stable-token-params
POST request
URL Params:
  - value (string) : the json enconded string, represent a map (string : interface), must be equal to current pending.
```


response:
```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```

### reject stable token params - (signing required)
```
<host>:8000/reject-stable-token-params
POST request
URL Params:
  nil
```


response:
```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```
### Get pending stable token params- (signing required) return the current pending stable token params
```
<host>:8000/pending-stable-token-params
GET request
params:
  - nonce (uint64) : the nonce to conform to signing requirement
```
example:
```
curl -x GET \
  http://localhost:8000/pending-token-params?nonce=111111
```
 
response:
```
{
  "data": {
    "DGX": {
      "AskSpread": 50,
      "BidSpread": 50,
      "PriceUpdateThreshold": 0.1
    }
  },
  "success": true
}
```

### Get stable token params- (signing required) return the current confirmed stable token params
```
<host>:8000/stable-token-params
GET request
params:
  - nonce (uint64) : the nonce to conform to signing requirement
```
example:
```
curl -x GET \
  http://localhost:8000/stable-token-params?nonce=111111
```
 
response:
```
{
  "data": {
    "DGX": {
      "AskSpread": 50,
      "BidSpread": 50,
      "PriceUpdateThreshold": 0.1
    }
  },
  "success": true
}
```
### Get gold data
```
<host>:8000/gold-feed
```
response:
```
{"data":{"Timestamp":1526923808631,"DGX":{"Valid":true,"Timestamp":0,"success":"","data":[{"symbol":"DGXETH","price":0.06676463,"time":1526923801},{"symbol":"ETHUSD","price":694.4,"time":1526923801},{"symbol":"ETHSGD","price":931.89,"time":1526923801},{"symbol":"DGXUSD","price":46.36,"time":1526923801},{"symbol":"EURUSD","price":1.17732,"time":1526923801},{"symbol":"USDSGD","price":1.34201,"time":1526923801},{"symbol":"XAUUSD","price":1291.468,"time":1526923801},{"symbol":"USDJPY","price":111.061,"time":1526923801}],"Error":""},"OneForgeETH":{"Value":1.85646,"Text":"1 XAU is worth 1.85646 ETH","Timestamp":1526923803,"Error":false,"Message":""},"OneForgeUSD":{"Value":1291.57,"Text":"1 XAU is worth 1291.57 USD","Timestamp":1526923803,"Error":false,"Message":""},"GDAX":{"Valid":true,"Error":"","trade_id":34527604,"price":"695.56000000","size":"0.00894700","bid":"695.55","ask":"695.56","volume":"50497.82498957","time":"2018-05-21T17:30:04.729000Z"},"Kraken":{"Valid":true,"network_error":"","error":[],"result":{"XETHZUSD":{"a":["696.66000","1","1.000"],"b":["696.33000","4","4.000"],"c":["696.33000","0.10776064"],"v":["13536.83019524","16999.30348103"],"p":["707.93621","710.18316"],"t":[5361,8276],"l":["693.97000","693.97000"],"h":["721.38000","724.80000"],"o":"715.65000"}}},"Gemini":{"Valid":true,"Error":"","bid":"694.50","ask":"695.55","volume":{"ETH":"11418.5646926","USD":"8064891.13775284649999999999999999999704534","timestamp":1526923800000},"last":"695.36"}},"success":true}
```

### Get BTC data

```
<host>:8000/btc-feed
response:
```json
{
  "data": {
    "Timestamp": 1541571292437,
    "bitfinex": {
      "Valid": true,
      "Error": "",
      "mid": "0.0332995",
      "bid": "0.03329",
      "ask": "0.033309",
      "last_price": "0.033299",
      "low": "0.032559",
      "high": "0.034036",
      "volume": "29077.44919025",
      "timestamp": "1541571292.1917806"
    },
    "binance": {
      "Valid": true,
      "Error": "",
      "symbol": "ETHBTC",
      "bidPrice": "0.03328500",
      "bidQty": "2.04600000",
      "askPrice": "0.03328700",
      "askQty": "0.17800000"
    }
  },
  "success": true
}
```