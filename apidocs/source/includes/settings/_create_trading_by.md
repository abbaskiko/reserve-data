# Create trading by

## Create trading by 

```shell
curl -X POST "https://gateway.local/v3/create-trading-by" \
-H 'Content-Type: application/json' \
-d '{
    "TradingBys": [
        {
            "asset_id": 3,
            "trading_pair_id": 1
        }
    ]
}'
```

> sample response

```json
{
  "id": 7,
  "success": true
}
```

### HTTP Request

`POST https://gateway.local/v3/create-trading-by`


## Get pending trading by


```shell
curl -X GET "https://gateway.local/v3/create-trading-by"
```

> sample response

```json
{
  "data": [
    {
      "id": 7,
      "created": "2019-08-13T07:31:28.197771Z",
      "data": {
        "TradingBys": [
          {
            "asset_id": 1,
            "trading_pair_id": 1
          }
        ]
      }
    }
  ],
  "success": true
}
```

### HTTP Request

`GET https://gateway.local/v3/create-trading-by`


## Confirm pending trading by 

```shell
curl -X PUT "https://gateway.local/v3/create-trading-by/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`PUT https://gateway.local/v3/create-trading-by/:trading_by_id`


## Reject pending trading by 

```shell
curl -X DELETE "https://gateway.local/v3/create-trading-by/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`DELETE https://gateway.local/v3/create-trading-by/:trading_by_id`