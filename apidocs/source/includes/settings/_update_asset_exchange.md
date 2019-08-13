# Update asset exchange

## Create update asset exchange 

```shell
curl -X POST "https://gateway.local/v3/update-asset-exchange"
-H 'Content-Type: application/json'
-d '{

}'
```

> sample response

```json
```

### HTTP Request

`POST https://gateway.local/v3/update-asset-exchange`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------

## Get pending update asset exchange


```shell
curl -X GET "https://gateway.local/v3/update-asset-exchange"
```

> sample response

```json
```

### HTTP Request

`GET https://gateway.local/v3/update-asset-exchange`


## Confirm pending update asset exchange

```shell
curl -X PUT "https://gateway.local/v3/update-asset-exchange/1"
```

> sample response

```json
```

### HTTP Request

`PUT https://gateway.local/v3/update-asset-exchange/:update_asset_exchange_id`


## Reject pending update asset exchange 

```shell
curl -X DELETE "https://gateway.local/v3/update-asset-exchange/1"
```

> sample response

```json
```

### HTTP Request

`DELETE https://gateway.local/v3/update-asset-exchange/:update_asset_exchange_id`