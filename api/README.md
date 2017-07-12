### `GET /backends`
```javascript
[
  {
    "name": "other-backend",
    "kind": "cloudwatch"
  },
  {
    "name": "test-backend",
    "kind": "cloudwatch"
  }
]
```

### `POST /scale`
```javascript
{
    "eval": "76e58486-0fd3-c2d9-f442-2996025ea814",
    "new_count": 3
}
```

### `POST /capacity`
```javascript
{
    "eval": "76e58486-0fd3-c2d9-f442-2996025ea814",
    "new_count": 3
}
```