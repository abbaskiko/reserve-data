## Get prices for all base-quote pairs

```shell
curl -X -GET "http://gateway.local/prices"
```

> sample response

```json
{
    "data": {
        "MCO-ETH": {
            "bittrex": {
                "Valid": true,
                "Error": "",
                "Timestamp": "1514114579228",
                "Bids": [
                    {
                        "Quantity": 142.93534777,
                        "Rate": 0.02437378
                    },
                    {
                        "Quantity": 1.21116959,
                        "Rate": 0.02437377
                    },
                    {
                        "Quantity": 1.63701658,
                        "Rate": 0.02437376
                    }
                ],
                "Asks": [
                    {
                        "Quantity": 15.39680469,
                        "Rate": 0.02503471
                    },
                    {
                        "Quantity": 18.71484714,
                        "Rate": 0.02503534
                    },
                    {
                        "Quantity": 93.57423573,
                        "Rate": 0.02503537
                    }
                ],
                "ReturnTime": "1514114579481"
            }
        },
        "OMG-ETH": {
            "bittrex": {
                "Valid": true,
                "Error": "",
                "Timestamp": "1514114579228",
                "Bids": [
                    {
                        "Quantity": 5.49,
                        "Rate": 0.019857
                    },
                    {
                        "Quantity": 13.62550123,
                        "Rate": 0.0197758
                    },
                    {
                        "Quantity": 10,
                        "Rate": 0.01976677
                    },
                    {
                        "Quantity": 6.92629385,
                        "Rate": 0.01970274
                    }
                ],
                "Asks": [
                    {
                        "Quantity": 6.73770653,
                        "Rate": 0.02025768
                    },
                    {
                        "Quantity": 7.49193537,
                        "Rate": 0.02025774
                    },
                    {
                        "Quantity": 1.48831433,
                        "Rate": 0.02025781
                    }
                ],
                "ReturnTime": "1514114579575"
            }
        }
    },
    "success": true,
    "timestamp": "1514114582015",
    "version": 64
}
```

### HTTP Request

`GET https://gateway.local/prices`

## Get prices for a specific base-quote pair

```shell
curl -X GET "https://gateway.local/prices/omg/eth"
```

> sample response


### HTTP Request

`GET https://gateway.local/:base/:quote`