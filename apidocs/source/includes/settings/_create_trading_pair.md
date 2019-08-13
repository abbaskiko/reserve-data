# Create trading pair

## Create trading pair 

```shell
curl -X POST "https://gateway.local/v3/create-trading-pair" \
-H 'Content-Type: application/json' \
-d '{
    "trading_pairs": [
        {
            "base": 1,
            "quote": 3,
            "price_precision": 10,
            "amount_precision": 10,
            "amount_limit_min": 1,
            "amount_limit_max": 20,
            "price_limit_min": 1,
            "price_limit_max": 1000,
            "min_notional": 0.23132,
            "exchange_id": 1
        }
    ]
}'
```

> sample response

```json
{
    "id": 1,
    "success": true,
}
```

### HTTP Request

`POST https://gateway.local/v3/create-trading-pair`


## Get pending trading pair 


```shell
curl -X GET "https://gateway.local/v3/create-trading-pair"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`GET https://gateway.local/v3/create-trading-pair`


## Confirm pending trading pair 

```shell
curl -X PUT "https://gateway.local/v3/create-trading-pair/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`PUT https://gateway.local/v3/create-trading-pair/:trading_pair_id`


## Reject pending trading pair 

```shell
curl -X DELETE "https://gateway.local/v3/create-trading-pair/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`DELETE https://gateway.local/v3/create-trading-pair/:trading_pair_id`