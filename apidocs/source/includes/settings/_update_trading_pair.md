# Update trading pair


## Create update trading pair 

```shell
curl -X POST "https://gateway.local/v3/update-trading-pair" \
-H 'Content-Type: application/json' \
-d '{

}'
```

> sample response

```json
```

### HTTP Request

`POST https://gateway.local/v3/update-trading-pair`


## Get pending update trading pair 


```shell
curl -X GET "https://gateway.local/v3/update-trading-pair"
```

> sample response

```json
```

### HTTP Request

`GET https://gateway.local/v3/update-trading-pair`


## Confirm pending update trading pair 

```shell
curl -X PUT "https://gateway.local/v3/update-pair/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`PUT https://gateway.local/v3/update-trading-pair/:trading_pair_id`


## Reject pending update trading pair 

```shell
curl -X DELETE "https://gateway.local/v3/update-trading-pair/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`DELETE https://gateway.local/v3/update-exchange/:trading_pair_id`