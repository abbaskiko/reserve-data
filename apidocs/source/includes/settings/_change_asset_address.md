# Change asset address


## Create change asset address

```shell
curl -X POST "https://gateway.local/v3/change-asset-address"
-H 'Content-Type: application/json'
-d '{
    "assets": [
        {
            "id": 1,
            "address": "0xC7DC5C95728d9ca387239Af0A49b7BCe8927d309"
        }
    ]
}'
```

> sample json

```json
{
    "success": true
}
```

### HTTP Request

`POST https://gateway.local/v3/change-asset-address`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------


## Get pending change asset address

```shell
curl -X GET "https://gateway.local/v3/change-asset-address"
```

> sample response

```json
```

### HTTP Request

`GET https://gateway.local/v3/change-asset-address`


## Confirm pending change asset address

```shell
curl -X PUT "https://gateway.local/v3/change-asset-address/1"
```

> sample response

```json
```

### HTTP Request

`PUT https://gateway.local/v3/change-asset-address/:asset_id`

## Reject pendign change asset address

```shell
curl -X DELETE "https://gateway.local/v3/change-asset-address/1"
```

> sample response

### HTTP Request

`DELETE https://gateway.local/v3/change-asset-address/:asset_id`