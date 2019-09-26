# Stable tokens and BTC

## Set stable token params

```shell
curl -X POST "https://gateway.local/v3/setting-change-stable"
-H 'Content-Type: application/x-www-form-urlencoded'
-d '{
    "change_list": [{
      "type": "update_stable_token_params",
      "data": {
        "params": {
          "DGX": {
            "AskSpread": 50,
            "BidSpread": 50,
            "PriceUpdateThreshold": 0.1
          }
        }
      }
    }]
}'
```

> sample response

```json
{
    "id": 1,
    "success": true
}
```

### HTTP request

`POST https://gateway.local/v3/setting-change-stable`

## Confirm stable token params

```shell
curl -X PUT "https://gateway.local/v3/setting-change-stable/16"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`POST https://gateway.local/v3/setting-change-stable/:change_id`

## Reject stable token params

```shell
curl -X DELETE "https://gateway.local/v3/setting-change-stable/16"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`DELETE https://gateway.local/v3/setting-change-stable/:change_id`

## Get pending stable token params

```shell
curl -X GET "https://gateway.local/v3/setting-change-stable"
```

> sample response

```json
{
    "data": [{
        "id": 24,
        "created": "2019-09-09T10:16:17.768309Z",
        "change_list": [{
            "type": "update_stable_token_params",
            "data": {
                "params": {
                    "DGX": {
                        "AskSpread": 50,
                        "BidSpread": 50,
                        "PriceUpdateThreshold": 0.1
                    }
                }
            }
        }]
    }],
    "success": true
}
```

### HTTP Request

`GET https://gateaway.local/v3/setting-change-stable`

## Get stable token params

```shell
curl -X GET "https://gateway.local/v3/stable-token-params"
```

> sample response

```json
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

### HTTP Request

`GET https://gateway.local/stable-token-params`

## Get gold data

```shell
curl -X GET "https://gateway.local/v3/gold-feed"
```

> sample response

```json
{
    "data": {
        "Timestamp": 1526923808631,
        "DGX": {
            "Valid": true,
            "Timestamp": 0,
            "success": "",
            "data": [
                {
                    "symbol": "DGXETH",
                    "price": 0.06676463,
                    "time": 1526923801
                },
                {
                    "symbol": "ETHUSD",
                    "price": 694.4,
                    "time": 1526923801
                },
                {
                    "symbol": "ETHSGD",
                    "price": 931.89,
                    "time": 1526923801
                },
                {
                    "symbol": "DGXUSD",
                    "price": 46.36,
                    "time": 1526923801
                },
                {
                    "symbol": "EURUSD",
                    "price": 1.17732,
                    "time": 1526923801
                },
                {
                    "symbol": "USDSGD",
                    "price": 1.34201,
                    "time": 1526923801
                },
                {
                    "symbol": "XAUUSD",
                    "price": 1291.468,
                    "time": 1526923801
                },
                {
                    "symbol": "USDJPY",
                    "price": 111.061,
                    "time": 1526923801
                }
            ],
            "Error": ""
        },
        "OneForgeETH": {
            "Value": 1.85646,
            "Text": "1 XAU is worth 1.85646 ETH",
            "Timestamp": 1526923803,
            "Error": false,
            "Message": ""
        },
        "OneForgeUSD": {
            "Value": 1291.57,
            "Text": "1 XAU is worth 1291.57 USD",
            "Timestamp": 1526923803,
            "Error": false,
            "Message": ""
        },
        "GDAX": {
            "Valid": true,
            "Error": "",
            "trade_id": 34527604,
            "price": "695.56000000",
            "size": "0.00894700",
            "bid": "695.55",
            "ask": "695.56",
            "volume": "50497.82498957",
            "time": "2018-05-21T17:30:04.729000Z"
        },
        "Kraken": {
            "Valid": true,
            "network_error": "",
            "error": [],
            "result": {
                "XETHZUSD": {
                    "a": [
                        "696.66000",
                        "1",
                        "1.000"
                    ],
                    "b": [
                        "696.33000",
                        "4",
                        "4.000"
                    ],
                    "c": [
                        "696.33000",
                        "0.10776064"
                    ],
                    "v": [
                        "13536.83019524",
                        "16999.30348103"
                    ],
                    "p": [
                        "707.93621",
                        "710.18316"
                    ],
                    "t": [
                        5361,
                        8276
                    ],
                    "l": [
                        "693.97000",
                        "693.97000"
                    ],
                    "h": [
                        "721.38000",
                        "724.80000"
                    ],
                    "o": "715.65000"
                }
            }
        },
        "Gemini": {
            "Valid": true,
            "Error": "",
            "bid": "694.50",
            "ask": "695.55",
            "volume": {
                "ETH": "11418.5646926",
                "USD": "8064891.13775284649999999999999999999704534",
                "timestamp": 1526923800000
            },
            "last": "695.36"
        }
    },
    "success": true
}
```

### HTTP Request

`GET https://gateway.local/v3/gold-feed`


## Get BTC data

```shell
curl -X GET "https://gateway.local/v3/btc-feed"
```

> sample response

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

### HTTP Request

`GET https://gateway.local/v3/btc-feed`