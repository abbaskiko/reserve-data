# Create trading pair

## Create trading pair 

```shell
curl -X POST "https://gateway.local/v3/create-trading-pair"
-H 'Content-Type: application/json'
-d '{

}'
```

> sample response

```json
```

### HTTP Request

`POST https://gateway.local/v3/create-trading-pair`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------

## Get pending trading pair 


```shell
curl -X GET "https://gateway.local/v3/create-trading-pair"
```

> sample response

```json
```

### HTTP Request

`GET https://gateway.local/v3/create-trading-pair`


## Confirm pending trading pair 

```shell
curl -X PUT "https://gateway.local/v3/create-trading-pair/1"
```

> sample response

```json
```

### HTTP Request

`PUT https://gateway.local/v3/create-trading-pair/:trading_pair_id`


## Reject pending trading pair 

```shell
curl -X DELETE "https://gateway.local/v3/create-trading-pair/1"
```

> sample response

```json
```

### HTTP Request

`DELETE https://gateway.local/v3/create-trading-pair/:trading_pair_id`