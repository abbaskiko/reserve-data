# Create asset exchange

## Create asset exchange

```shell
curl -X POST "https://gateway.local/v3/create-asset-exchange"
-H 'Content-Type: application/json'
-d '{

}'
```

> sample response

```json
```

### HTTP Request

`POST https://gateway.local/v3/create-asset-exchange`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------

## Get pending asset exchange


```shell
curl -X GET "https://gateway.local/v3/create-asset-exchange"
```

> sample response

```json
```

### HTTP Request

`GET https://gateway.local/v3/create-asset-exchange`


## Confirm pending asset exchange

```shell
curl -X PUT "https://gateway.local/v3/create-asset-exchange/1"
```

> sample response

```json
```

### HTTP Request

`PUT https://gateway.local/v3/create-asset-exchange/:asset_exchange_id`


## Reject pending asset exchange

```shell
curl -X DELETE "https://gateway.local/v3/create-asset-exchange/1"
```

> sample response

```json
```

### HTTP Request

`DELETE https://gateway.local/v3/create-asset-exchange/:asset_exchange_id`