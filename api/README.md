## GET
### `GET /ping`
```javascript
"welcome to libra, the Nomad auto-scaler"
```

### `GET /backends`
```javascript
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

### `GET /ping`
```javascript
"pong"
```

## POST
### `POST /capacity`
```javascript
{
    "eval": "76e58486-0fd3-c2d9-f442-2996025ea814",
    "new_count": 3
}
```

### `POST /scale`
```javascript
{
    "eval": "76e58486-0fd3-c2d9-f442-2996025ea814",
    "new_count": 3
}
```

