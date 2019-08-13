# Update asset exchange

## Create update asset exchange 

```shell
curl -X POST "https://gateway.local/v3/update-asset-exchange" \
-H 'Content-Type: application/json' \
-d '{
    "asset_exchanges": [
        {
            "id": 1,
            "symbol": "KNC",
            "deposit_address": "0xC5094e852F71346dF9eD6795Be6F014994C43e09",
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
  "id": 2,
  "success": true
}
```

### HTTP Request

`POST https://gateway.local/v3/update-asset-exchange`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
asset_exchanges | arrays of asset exchange | true | nil | array of asset exchanges

## Get pending update asset exchange


```shell
curl -X GET "https://gateway.local/v3/update-asset-exchange"
```

> sample response

```json
{
  "data": [
    {
      "id": 3,
      "created": "2019-08-13T06:40:08.17648Z",
      "data": {
        "asset_exchanges": [
          {
            "id": 1,
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

`GET https://gateway.local/v3/update-asset-exchange`


## Confirm pending update asset exchange

```shell
curl -X PUT "https://gateway.local/v3/update-asset-exchange/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`PUT https://gateway.local/v3/update-asset-exchange/:update_asset_exchange_id`


## Reject pending update asset exchange 

```shell
curl -X DELETE "https://gateway.local/v3/update-asset-exchange/1"
```

> sample response

```json
{
    "success":true
}
```

### HTTP Request

`DELETE https://gateway.local/v3/update-asset-exchange/:update_asset_exchange_id`