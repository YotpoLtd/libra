# Restarting

## Restart a Nomad job

```shell
curl -X POST \
  http://libra.consul/restart \
  -H 'content-type: application/json' \
  -d '{
	      "job": "nginx",
        "group": "nginx",
        "task": "nginx",
        "image": "nginx:latest"
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
group | string | The name of the Nomad group to restart
task | string | The name of the Nomad task to restart
image | string | The Docker image that the Nomad job will pull down on restart