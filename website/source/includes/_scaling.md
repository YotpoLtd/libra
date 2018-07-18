# Scaling

## Scale a Nomad group

```shell
curl -X POST \
  http://libra.consul/scale \
  -H 'content-type: application/json' \
  -d '{
	"job": "nginx",
	"group": "nginx",
	"count": -1
}'
```

> The above command returns JSON structured like this:

```json
{
  "eval": "76e58486-0fd3-c2d9-f442-2996025ea814",
  "new_count": 3
}
```

This endpoint will increase or decrease the deesired count of a Nomad group.

### HTTP Request

`POST http://libra.consul/scale`

### JSON Parameters

Parameter | Type | Description
--------- | ---- | -----------
job | string | The name of the Nomad job to scale
group | string | The name of the Nomad group to scale
count | integer | The amount to scale the group up or down by, positive or negative

## Set the desired capacity of a Nomad group

```shell
curl -X POST \
  http://libra.consul/capacity \
  -H 'content-type: application/json' \
  -d '{
	"job": "nginx",
	"group": "nginx",
	"count": 3
}'
```

> The above command returns JSON structured like this:

```json
{
  "eval": "76e58486-0fd3-c2d9-f442-2996025ea814",
  "new_count": 3
}
```

This endpoint sets the desired count of a Nomad group.

### HTTP Request

`POST http://libra.consul/capacity`

### URL Parameters

Parameter | Type | Description
--------- | ---- | -----------
job | string | The name of the Nomad job to set the capacity of
group | string | The name of the Nomad group to set the capacity of
count | integer | The desired number of Nomad groups to run