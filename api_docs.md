## API for reserve data

### Get open orders in exchange

```
localhost://open-orders?exchange=binance
GET request
Query params:
    - exchange: string (ex: binance) - optional
```

response:

```json
{
    "success": true,
    "data": {
        "binance": [
            {
                "OrderID": "123132",
                "Price": 0.43432
            }
        ],
        "huobi": [
            {
                "ID": "",
                "Base": "KNC",
                "Quote": "ETH",
                "OrderID": "71167643340",
                "Price": 0.0018,
                "OrigQty": 108.33,
                "ExecutedQty": 0,
                "TimeInForce": "",
                "Type": "",
                "Side": "",
                "StopPrice": "",
                "IcebergQty": "",
                "Time": 0
            }
        ]
    }
}
```

### Cancel order by order id (signing require)

```
<host>:8000/cancel-order-by-order-id
POST request
Form params:
  - exchange_id: string (ex: binance)
  - order_id: string (this is order id get from open orders)
```

sample response:
```json
{
    "reason": "UNKNOWN_ORDER",
    "success": false
}
```