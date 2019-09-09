# Rebalance

## Get rebalance status
Get rebalance status, if reponse is *true* then rebalance is enable, the analytic can perform rebalance, else reponse is *false*, the analytic hold rebalance ability.

```shell
curl -X GET "http://gateway.local/v3/rebalance-status"
```

> sample response

```json
  {
    "success": true
  }
```

### HTTP Request

`GET http://gateway.local/v3/rebalance-status`

## Hold rebalance

```shell
curl -X POST "https://gateway.local/v3/hold-rebalance"
```

> sample response

```json
  {
    "success": true
  }
```

### HTTP Request

`POST http://gateway.local/v3/hold-rebalance`
<aside class="notice">Confirm key is required</aside>

## Enable rebalance

```shell
curl -X POST "https://gateway.local/v3/enable-rebalance"
```

> sample response

```json
  {
    "success": true
  }
```

### HTTP Request

`POST http://gateway.local/v3/enable-rebalance`
<aside class="notice">Confirm key is required</aside>

## Get set-rate status
Get set-rate status, if response is *true* then set-rate is enable, the analytic can perform set-rate, else response is *false*, the analytic hold set-rate ability.

```shell
curl -X -GET "https://gateway.local/set-rate-status"
```

> sample response

```json
  {
    "success": true
  }
```

## Hold set-rate

```shell
curl -X POST "https://gateway.local/v3/hold-set-rate"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`POST http://gateway.local/v3/hold-set-rate`
<aside class="notice">Confirm key is required</aside>

## Enable setrate

```shell
curl -X POST "https://gateway.local/v3/enable-set-rate"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`POST http://gateway.local/v3/enable-set-rate`
<aside class="notice">Confirm key is required</aside>