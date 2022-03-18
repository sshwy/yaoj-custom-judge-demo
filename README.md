# Yaoj Custom Judge Demo

api list：

| Route                        | Description     |
| ---------------------------- | --------------- |
| /                            | welcome message |
| /jsonrpc                     | json RPC route  |
| /files?name={name}&ext={ext} | upload files    |

e. g.

```bash
curl -i -X GET \
 'http://localhost:3000'

curl -i -X POST \
   -H "Content-Type:text/x-csrc" \
   -T "./a.in" \
 'http://localhost:3000/files?name=test.in'

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "jsonrpc":"2.0",
  "method":"JudgeService.CustomTest",
  "params":[
     { "src": "test.c", "input": "test.in" }
  ],
  "id": 1
}' \
 'http://localhost:3000/jsonrpc'
```



