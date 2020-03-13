## API for reserve data

### Get open orders in exchange

```
localhost://open-orders?exchange=binance
GET request
Query params:
    - exchange: string (ex: binance)
```

response:

```json
    {
        "success": true,
        "data": [
            "OrderID": "123132",
            "Price": 0.43432,
        ]
    }
```

### Cancel order (signing required)
```
<host>:8000/cancelorder/
POST request
Form params:
  - exchange_id: string (ex: binance)
  - order_id: string (this is order id get from open orders)
```

response:
```json
{
    "reason": "UNKNOWN_ORDER",
    "success": false
}
```
