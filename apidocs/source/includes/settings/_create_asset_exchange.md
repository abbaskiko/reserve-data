# Create asset exchange

## Create asset exchange

```shell
curl -X POST "https://gateway.local/v3/create-asset-exchange" \
-H 'Content-Type: application/json' \
-d '{
    "asset_exchanges": [
        {
            "asset_id": 3,
            "exchange_id": 1,
            "symbol": "KNC",
            "deposit_address": "",
            "min_deposit": 0,
            "withdraw_fee": 0,
            "target_recommended": 0,
            "target_ratio": 0
        }
    ]
}'
```

> sample response

```json
{
  "id": 5,
  "success": true
}
```

### HTTP Request

`POST https://gateway.local/v3/create-asset-exchange`


## Get pending asset exchange


```shell
curl -X GET "https://gateway.local/v3/create-asset-exchange"
```

> sample response

```json
{
  "data": [
    {
      "id": 5,
      "created": "2019-08-13T06:52:42.069694Z",
      "data": {
        "asset_exchanges": [
          {
            "asset_id": 3,
            "exchange_id": 1,
            "symbol": "KNC",
            "deposit_address": "0xc5094e852f71346df9ed6795be6f014994c43e09",
            "min_deposit": 0,
            "withdraw_fee": 0,
            "target_recommended": 0,
            "target_ratio": 0
          }
        ]
      }
    }
  ],
  "success": true
}
```

### HTTP Request

`GET https://gateway.local/v3/create-asset-exchange`


## Confirm pending asset exchange

```shell
curl -X PUT "https://gateway.local/v3/create-asset-exchange/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`PUT https://gateway.local/v3/create-asset-exchange/:asset_exchange_id`


## Reject pending asset exchange

```shell
curl -X DELETE "https://gateway.local/v3/create-asset-exchange/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`DELETE https://gateway.local/v3/create-asset-exchange/:asset_exchange_id`