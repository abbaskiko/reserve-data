# GasThreshold

## Set GasThreshold

``` shell
curl -X POST "https://gateway.local/v3/gas-threshold" \
-H 'Content-Type: application/json' \
-d '{
          "high": 10.0
		  "low": 5.345
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

`POST https://gateway.local/v3/gas-threshold`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
high | float64 | yes |  | the high value
low | float64 | yes |  | the low value
<aside class="notice">Rebalance key is required</aside>

## Get rate trigger period


```shell
curl -X GET "https://gateway.local/v3/gas-threshold"
```

> sample response

```json
{
  "success": true,
  "data": {
    "eth-gas-station": {
      "fast": 100,
      "standard": 80,
      "low": 50
    },
    "high": 150,
    "low": 120
  }
}
```

### HTTP Request

`GET https://gateway.local/v3/rate-trigger-period`
<aside class="notice">All keys are accepted</aside>