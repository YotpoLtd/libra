# Backends

## Get All Backends

```shell
curl "http://libra.consul/backends"
```

> The above command returns JSON structured like this:

```json
[
  {
    "name": "other-backend",
    "kind": "graphite"
  },
  {
    "name": "test-backend",
    "kind": "cloudwatch"
  }
]
```

This endpoint retrieves all configured backends.

### HTTP Request

`GET http://libra.consul/backends`