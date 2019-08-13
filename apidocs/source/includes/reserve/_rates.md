# Rates and Prices

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

`GET https://gateway.local/prices/:base/:quote`

## Get token rates from blockchain

```shell
curl -X GET "http://127.0.0.1:8000/getrates"
```

> sample response

```json
{
    "data": {
        "MCO": {
            "Valid": true,
            "Error": "",
            "Timestamp": "1515412582435",
            "ReturnTime": "1515412582710",
            "BaseBuy": 63.99319226272073,
            "CompactBuy": 21,
            "BaseSell": 0.014716371218820246,
            "CompactSell": -20,
            "Rate": 0,
            "Block": 2420849
        },
        "OMG": {
            "Valid": true,
            "Error": "",
            "Timestamp": "1515412582435",
            "ReturnTime": "1515412582710",
            "BaseBuy": 44.45707162223901,
            "CompactBuy": 30,
            "BaseSell": 0.021183301968644246,
            "CompactSell": -29,
            "Rate": 0,
            "Block": 2420849
        },
        "PAY": {
            "Valid": true,
            "Error": "",
            "Timestamp": "1515412582435",
            "ReturnTime": "1515412582710",
            "BaseBuy": 295.08854913901575,
            "CompactBuy": -13,
            "BaseSell": 0.003191406699999999,
            "CompactSell": 13,
            "Rate": 0,
            "Block": 2420849
        }
    },
    "success": true,
    "timestamp": "1515412583215",
    "version": 1515412582435
}
```

### HTTP Request

`GET http://gateway.local/v3/getrates`


## Get all token rates from blockchain

```shell
curl -X GET "https://gateway.local/get-all-rates"
```

> sample response

```json
{
    "data": [
        {
            "Version": 0,
            "Valid": true,
            "Error": "",
            "Timestamp": "1517280618739",
            "ReturnTime": "1517280619071",
            "Data": {
                "SNT": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280618739",
                    "ReturnTime": "1517280619071",
                    "BaseBuy": 4053.2170631085987,
                    "CompactBuy": 43,
                    "BaseSell": 0.000233599514875301,
                    "CompactSell": -3,
                    "Rate": 0,
                    "Block": 5635245
                }
            }
        },
        {
            "Version": 0,
            "Valid": true,
            "Error": "",
            "Timestamp": "1517280621738",
            "ReturnTime": "1517280622251",
            "Data": {
                "EOS": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 87.21360760013062,
                    "CompactBuy": 0,
                    "BaseSell": 0.0128686459657361,
                    "CompactSell": 0,
                    "Rate": 0,
                    "Block": 5635245
                },
                "ETH": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 0,
                    "CompactBuy": 32,
                    "BaseSell": 0,
                    "CompactSell": -14,
                    "Rate": 0,
                    "Block": 5635245
                },
                "KNC": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 307.05930436561505,
                    "CompactBuy": -34,
                    "BaseSell": 0.003084981280661941,
                    "CompactSell": 81,
                    "Rate": 0,
                    "Block": 5635245
                },
                "OMG": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 65.0580993582104,
                    "CompactBuy": 32,
                    "BaseSell": 0.014925950060437398,
                    "CompactSell": -14,
                    "Rate": 0,
                    "Block": 5635245
                },
                "SALT": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 152.3016783627643,
                    "CompactBuy": 9,
                    "BaseSell": 0.006196212698403499,
                    "CompactSell": 23,
                    "Rate": 0,
                    "Block": 5635245
                },
                "SNT": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 4053.2170631085987,
                    "CompactBuy": 43,
                    "BaseSell": 0.000233599514875301,
                    "CompactSell": -3,
                    "Rate": 0,
                    "Block": 5635245
                }
            }
        }
    ],
    "success": true
}
```

### HTTP Request

`GET https://gateway.local/v3/get-all-rates`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
fromTime | uint64 | false | nil | fromTime to get all rates
toTime | uint64 | false | nil | toTime to get all rates

