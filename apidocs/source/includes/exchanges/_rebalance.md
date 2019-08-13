# Rebalance

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