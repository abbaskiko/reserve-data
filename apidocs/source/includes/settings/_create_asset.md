# Create Asset


## Create asset

```shell
curl -X POST "https://gateway.local/v3/create-asset" \
-H 'Content-Type: application/json' \
-d '{
    "assets": [
        {
            "symbol": "OMG",
            "name": "Omisego",
            "address": "0xdd974d5c2e2928dea5f71b9825b8b646686bd200",
            "decimals": 18,
            "transferable": true,
            "set_rate": "exchange_feed",
            "rebalance": true,
            "is_quote": true,
            "pwi": {
                "ask": {
                    "a": 13.88888,
                    "b": -0.11111,
                    "c": 0,
                    "min_min_spread": 0.005,
                    "price_multiply_factor": 0.45
                },
                "bid": {
                    "a": 13.88888,
                    "b": -0.11111,
                    "c": 0,
                    "min_min_spread": 0.005,
                    "price_multiply_factor": 0.45
                }
            },
            "rebalance_quadratic": {
                "a": 0.000001754386,
                "b": 0.0004894737,
                "c": 0.9995088
            },
            "exchanges": [
                {
                    "id": 5,
                    "asset_id": 3,
                    "exchange_id": 0,
                    "symbol": "OMG",
                    "deposit_address": "0x023ab1f7acaad1f7a01d3bfa4afd2ab575780090",
                    "min_deposit": 0,
                    "withdraw_fee": 18,
                    "target_recommended": 0,
                    "target_ratio": 0,
                    "trading_pairs": []
                }
            ],
            "target": {
                "total": 70000,
                "reserve": 48000,
                "rebalance_threshold": 0.33,
                "transfer_threshold": 0.25
            }
        }
    ]
}'
```

> sample response

```json
{
  "id": 8,
  "success": true
}
```

### HTTP Request

`POST https://gateway.local/v3/create-asset`

## Get pending create asset

```shell
curl -X -GET "https://gateway.local/v3/create-asset"
```

> sample response

```json
{
  "data": [
    {
      "id": 8,
      "created": "2019-08-13T07:47:22.680541Z",
      "data": {
        "assets": [
          {
            "symbol": "OMG",
            "name": "Omisego",
            "address": "0xdd974d5c2e2928dea5f71b9825b8b646686bd200",
            "old_addresses": null,
            "decimals": 18,
            "transferable": true,
            "set_rate": "exchange_feed",
            "rebalance": true,
            "is_quote": true,
            "pwi": {
              "ask": {
                "a": 13.88888,
                "b": -0.11111,
                "c": 0,
                "min_min_spread": 0.005,
                "price_multiply_factor": 0.45
              },
              "bid": {
                "a": 13.88888,
                "b": -0.11111,
                "c": 0,
                "min_min_spread": 0.005,
                "price_multiply_factor": 0.45
              }
            },
            "rebalance_quadratic": {
              "a": 0.000001754386,
              "b": 0.0004894737,
              "c": 0.9995088
            },
            "exchanges": [
              {
                "id": 5,
                "asset_id": 3,
                "exchange_id": 0,
                "symbol": "OMG",
                "deposit_address": "0x023ab1f7acaad1f7a01d3bfa4afd2ab575780090",
                "min_deposit": 0,
                "withdraw_fee": 18,
                "target_recommended": 0,
                "target_ratio": 0
              }
            ],
            "target": {
              "total": 70000,
              "reserve": 48000,
              "rebalance_threshold": 0.33,
              "transfer_threshold": 0.25
            }
          }
        ]
      }
    }
  ],
  "success": true
}
```

### HTTP Request

`GET https://gateway.local/v3/create-asset`


## Get pending create asset by id

```shell
curl -X GET "https://gateway.local/v3/create-asset/8"
```

> sample response

```json
{
  "data": {
    "id": 8,
    "created": "2019-08-13T07:47:22.680541Z",
    "data": {
      "assets": [
        {
          "symbol": "OMG",
          "name": "Omisego",
          "address": "0xdd974d5c2e2928dea5f71b9825b8b646686bd200",
          "old_addresses": null,
          "decimals": 18,
          "transferable": true,
          "set_rate": "exchange_feed",
          "rebalance": true,
          "is_quote": true,
          "pwi": {
            "ask": {
              "a": 13.88888,
              "b": -0.11111,
              "c": 0,
              "min_min_spread": 0.005,
              "price_multiply_factor": 0.45
            },
            "bid": {
              "a": 13.88888,
              "b": -0.11111,
              "c": 0,
              "min_min_spread": 0.005,
              "price_multiply_factor": 0.45
            }
          },
          "rebalance_quadratic": {
            "a": 0.000001754386,
            "b": 0.0004894737,
            "c": 0.9995088
          },
          "exchanges": [
            {
              "id": 5,
              "asset_id": 3,
              "exchange_id": 0,
              "symbol": "OMG",
              "deposit_address": "0x023ab1f7acaad1f7a01d3bfa4afd2ab575780090",
              "min_deposit": 0,
              "withdraw_fee": 18,
              "target_recommended": 0,
              "target_ratio": 0
            }
          ],
          "target": {
            "total": 70000,
            "reserve": 48000,
            "rebalance_threshold": 0.33,
            "transfer_threshold": 0.25
          }
        }
      ]
    }
  },
  "success": true
}
```

### HTTP Request

`GET https://gateway.local/v3/create-asset/:asset_id`


## Confirm create asset

```shell
curl -X PUT "https://gateway.local/v3/create-asset/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`PUT https://gateway.local/v3/create-asset/:asset_id`

## Reject create asset

```shell
curl -X DELETE "https://gateway.local/v3/create-asset/1"
```

> sample response

```json
{
    "success": true
}
```

### HTTP Request

`DELETE https://gateway.local/v3/create-asset/1`