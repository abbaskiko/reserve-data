---
title: API Reference

language_tabs: # must be one of https://git.io/vQNgJ
  - shell

toc_footers:
  - <a href='https://github.com/lord/slate'>Documentation Powered by Slate</a>

includes:
  - exchanges/exchanges
  - exchanges/rebalance
  - errors

search: true
---

# Introduction

# Authentication
Authentication follow: https://tools.ietf.org/html/draft-cavage-http-signatures-10

Required headers:

- **Digest**
- **Authorization**
- **Signature**
- **Nonce**

## Get time server

```shell
curl -X GET "http://gateway.local/timeserver"
```

> sample response

```json
{
  "data": "1517479497447",
  "success": true
}
```

### HTTP request

`GET https://gateway.local/timeserver`

