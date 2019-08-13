# Update asset

## Create update asset 

```shell
curl -X POST "https://gateway.local/v3/update-asset" \
-H 'Content-Type: application/json' \
-d '{
    "assets": [
        {
            "asset_id": 3,
            "symbol": "KNC",
            "name": "Kyber Network Crystal",
            "address": "",
            "decimals": 18,
            "transferable": true,
            "set_rate": "exchange_fee",
            "rebalance": true,
            "is_quote": true,
            "pwi": {
                "ask": {
                    "a": 0.3241,
                    "b": 0.342472,
                    "c": 0.3414,
                    "min_min_spread": 0.34241,
                    "price_multiply_factor": 0.43141
                },
                "bid": {
                    "a": 0.3241,
                    "b": 0.342472,
                    "c": 0.3414,
                    "min_min_spread": 0.34241,
                    "price_multiply_factor": 0.43141
                }
            },
            "rebalance_quadratic": {
                    "a": 0.3241,
                    "b": 0.342472,
                    "c": 0.3414
            },
            "target": {
                "total": 4134242.432,
                "reserve": 34535,
                "rebalance_threshold": 343.43,
                "transfer_threshold": 4353.4353
            }
        }
    ]
}'
```

> sample response

```json
{
    "id": 1,
    "success": true
}
```

### HTTP Request

`POST https://gateway.local/v3/update-asset`


## Get pending update asset exchange


```shell
curl -X GET "https://gateway.local/v3/update-asset"
```

> sample response

```json
{
    "assets": [
        {
            "asset_id": 3,
            "symbol": "KNC",
            "name": "Kyber Network Crystal",
            "address": "",
            "decimals": 18,
            "transferable": true,
            "set_rate": "exchange_fee",
            "rebalance": true,
            "is_quote": true,
            "pwi": {
                "ask": {
                    "a": 0.3241,
                    "b": 0.342472,
                    "c": 0.3414,
                    "min_min_spread": 0.34241,
                    "price_multiply_factor": 0.43141
                },
                "bid": {
                    "a": 0.3241,
                    "b": 0.342472,
                    "c": 0.3414,
                    "min_min_spread": 0.34241,
                    "price_multiply_factor": 0.43141
                }
            },
            "rebalance_quadratic": {
                    "a": 0.3241,
                    "b": 0.342472,
                    "c": 0.3414
            },
            "target": {
                "total": 4134242.432,
                "reserve": 34535,
                "rebalance_threshold": 343.43,
                "transfer_threshold": 4353.4353
            }
        }
    ]
}
```

### HTTP Request

`GET https://gateway.local/v3/update-asset`


## Confirm pending update asset 

```shell
curl -X PUT "https://gateway.local/v3/update-asset/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`PUT https://gateway.local/v3/update-asset/:asset_id`


## Reject pending update asset

```shell
curl -X DELETE "https://gateway.local/v3/update-asset/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`DELETE https://gateway.local/v3/update-asset/:asset_id`