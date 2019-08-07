
## Trade

Create an order into centralized exchange

```shell
curl -X POST "http://gateway.local/trade/binance"
-H 'Content-Type: application/json'
-d '{
   "base": "KNC",
   "quote": "ETH",
   "rate": 0.443,
   "amount": 141,
   "type": "buy" 
}'
```

> sample response

```json
{
    "id": "19234634",
    "success": true,
    "done": 0,
    "remaining": 0.01,
    "finished": false
}
```

### HTTP Request

`POST https://gateway.local/trade/:exchange_id`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
base | string | true | nil | base asset
quote | string | true | nil | quote asset
rate | float64 | true | nil | rate to create order
amount | float64 | true | nil | order amount 
type | string | true | nil | order type (buy or sell)


## Cancel order 

```shell
curl -X POST "http://gateway.local/cancelorder/binance"
-H 'Content-Type: application/json'
-d '{
    "order_id": 43142
}'
```

> sample response

```json
{
    "success": true
}
```

### HTTP request

`POST https://gateway.local/cancelorder/:exchange`

Param | Type | Required | Default | Description
----- | ---- | -------- | ------- | -----------
order_id | string | true | nil | order id to be cancelled

## Withdraw

```shell
curl -X POST "https://gateway.local/withdraw/binance"
-H 'Content-type: application/json'
-d '{
    "amount": 41.42342,
    "asset": "ETH"
}'
```

> sample response

```json
    "success": true,
    "id": 1432423
```

### HTTP Request

`POST https://gateway.local/withdraw/:exchange`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
amount | string(big int) | true | nil | amount we want to withdraw
asset | uint64 (asset id) | true | nil | asset we want to withdraw


## Deposit

```shell
curl -X POST "https://gateway.local/deposit/binance"
-H 'Content-Type: application/json'
-d '{
    "asset": 1,
    "amount": "41342342"
}'
```

> sample response

```json
{
    "success": true,
    "id": 34142342
}
```

### HTTP Request

`POST https://gateway.local/deposit/:exchange`

Param | Type | Required | Default | Description
----- | ---- | -------- | ------- | -----------
asset | integer (asset id) | true | nil | asset id
amount | string (big integer) | true | nil | amount to deposit


## Get rebalance status
Get rebalance status, if reponse is *true* then rebalance is enable, the analytic can perform rebalance, else reponse is *false*, the analytic hold rebalance ability.

```shell
curl -X GET "http://gateway.local/rebalancestatus"
```

> sample response
```json
  {
    "success": true,
  }
```

### HTTP Request

`GET http://gateway.local/rebalancestatus`

## Hold rebalance

```shell
curl -X POST "https://gateway.local/holdrebalance"
```

> sample response
```json
  {
    "success": true
  }
```

## Enable rebalance

```shell
curl -X POST "https://gateway.local/enablerebalance"
```

> sample response
```json
  {
    "success": true
  }
```

### HTTP Request

`POST http://gateway.local/enablerebalance`

## Get setrate status
Get setrate status, if reponse is *true* then setrate is enable, the analytic can perform setrate, else reponse is *false*, the analytic hold setrate ability.

```shell
curl -X -GET "https://gateway.local/setratestatus"
```

> sammple response
```json
  {
    "success": true,
    "data": true
  }
```

### Hold setrate

```shell
curl -X POST "https://gateway.local/holdsetrate"
```

> sample response
```json
{
    "success": true
}
```

## Enable setrate

```shell
curl -X POST "https://gateway.local/enablesetrate"
```

> sample response
```json
{
    "success": true
}
```