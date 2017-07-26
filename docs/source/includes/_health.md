# Health

## Check Libra health

```shell
curl "http://libra.consul/ping"
```

> The above command returns JSON structured like this:

```json
"pong"
```

This endpoint gets the health of a Libra server.

### HTTP Request

`GET http://libra.consul/ping`