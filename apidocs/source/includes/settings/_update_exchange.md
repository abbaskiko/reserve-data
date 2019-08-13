# Update exchange

## Create update exchange

```shell
curl -X POST "https://gateway.local/v3/update-exchange"
-H 'Content-Type: application/json'
-d '{

}'
```

> sample response

```json
```

### HTTP Request

`POST https://gateway.local/v3/update-exchange`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------

## Get pending update asset exchange


```shell
curl -X GET "https://gateway.local/v3/update-exchange"
```

> sample response

```json
```

### HTTP Request

`GET https://gateway.local/v3/update-exchange`


## Confirm pending update asset exchange

```shell
curl -X PUT "https://gateway.local/v3/update-exchange/1"
```

> sample response

```json
```

### HTTP Request

`PUT https://gateway.local/v3/update-exchange/:exchange_id`


## Reject pending update asset exchange 

```shell
curl -X DELETE "https://gateway.local/v3/update-exchange/1"
```

> sample response

```json
```

### HTTP Request

`DELETE https://gateway.local/v3/update-exchange/:exchange_id`