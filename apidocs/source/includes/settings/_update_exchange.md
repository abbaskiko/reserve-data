# Update exchange

## Create update exchange

```shell
curl -X POST "https://gateway.local/v3/update-exchange" \
-H 'Content-Type: application/json' \
-d '{
    "exchanges": [
        {
            "exchange_id": 1,
            "trading_fee_maker": 0,
            "trading_fee_taker": 0,
            "disable": true
        }
    ]
}'
```

> sample response

```json
{
  "id": 6,
  "success": true
}
```

### HTTP Request

`POST https://gateway.local/v3/update-exchange`


## Get pending update exchange 


```shell
curl -X GET "https://gateway.local/v3/update-exchange"
```

> sample response

```json
{
  "data": [
    {
      "id": 6,
      "created": "2019-08-13T07:25:49.869418Z",
      "data": {
        "exchanges": [
          {
            "exchange_id": 1,
            "trading_fee_maker": 0,
            "trading_fee_taker": 0,
            "disable": true
          }
        ]
      }
    }
  ],
  "success": true
}
```

### HTTP Request

`GET https://gateway.local/v3/update-exchange`


## Confirm pending update exchange

```shell
curl -X PUT "https://gateway.local/v3/update-exchange/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`PUT https://gateway.local/v3/update-exchange/:exchange_id`


## Reject pending update exchange 

```shell
curl -X DELETE "https://gateway.local/v3/update-exchange/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`DELETE https://gateway.local/v3/update-exchange/:exchange_id`