# Restarting

## Restart a Nomad job

```shell
curl -X POST \
  http://libra.consul/restart \
  -H 'content-type: application/json' \
  -d '{
	"job": "nginx"
}'
```

> The above command returns JSON structured like this:

```json
{
  "eval": "76e58486-0fd3-c2d9-f442-2996025ea814"
}
```

This endpoint will restart a Nomad job.

### HTTP Request

`POST http://libra.consul/restart`

### JSON Parameters

Parameter | Type | Description
--------- | ---- | -----------
job | string | The name of the Nomad job to restart