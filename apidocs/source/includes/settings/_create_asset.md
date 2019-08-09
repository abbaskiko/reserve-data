# Create Asset


## Create asset

```shell
curl -X POST "https://gateway.local/v3/create-asset"
-H 'Content-Type: application/json'
-d '{
    assets: [
        {
            "symbol": "OMG",
            "name": "Omisego",
            "address": "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07",
            "old_addresses": [],
            "decimals": 18,
            "tranferable": true,
            "set_rate": "exchange_feed",
            "rebalance": false,
            "is_quote": false,
            "pwi": {
                "ask": 121,
                "bid": 12
            },
            "rebalance_quadratic": {
                "a": 1,
                "b": 2,
                "c": 3
            },
            "exchange": [],
            target: 
        }
    ]
}'
```

> sample response

```json
{

}
```

### HTTP Request

`POST https://gateway.local/create-asset`

Param | Type | Required | Default | Description
----- | ---- | -------- | ------- | -----------

## Get pending create asset

```shell
curl -X -GET "https://gateway.local/create-asset"
```

> sample response

```json
{}
```

### HTTP Request

`GET https://gateway.local/`

Param | Type | Required | Default | Description
----- | ---- | -------- | ------- | -----------

## Get pending create asset by id

```shell
```

> sample response

```json
```

### HTTP Request

``

Param | Type | Required | Default | Description
----- | ---- | -------- | ------- | -----------

## Confirm create asset

```shell
```

> sample response

```json
```

### HTTP Request

``

Param | Type | Required | Default | Description
----- | ---- | -------- | ------- | -----------

## Reject create asset

```shell
```

> sample response

```json
```

### HTTP Request

``

Param | Type | Required | Default | Description
----- | ---- | -------- | ------- | -----------