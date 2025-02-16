# Perf test results for GET /api/info for one user 100k requests

hey -n 100000 -q 1000 -m GET -H "Authorization: Bearer `<jwt-token>`" http://localhost:8080/api/info

Summary:
  Total:        20.3726 secs
  Slowest:      0.1789 secs
  Fastest:      0.0007 secs
  Average:      0.0101 secs
  Requests/sec: 4908.5493

Response time histogram:
  0.001 [1]     |
  0.018 [94605] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.036 [5109]  |■■
  0.054 [220]   |
  0.072 [15]    |
  0.090 [0]     |
  0.108 [0]     |
  0.125 [3]     |
  0.143 [38]    |
  0.161 [7]     |
  0.179 [2]     |

Latency distribution:
  10% in 0.0053 secs
  25% in 0.0069 secs
  50% in 0.0092 secs
  75% in 0.0120 secs
  90% in 0.0157 secs
  95% in 0.0189 secs
  99% in 0.0277 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0007 secs, 0.1789 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0248 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0024 secs
  resp wait:    0.0101 secs, 0.0006 secs, 0.1506 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0007 secs

Status code distribution:
  [200] 100000 responses

# Info

{
  "coinHistory": {
    "received": [
      {
        "amount": 100,
        "fromUser": "R122111"
      },
      {
        "amount": 100,
        "fromUser": "R122111"
      },
      {
        "amount": 100,
        "fromUser": "R122111"
      },
      {
        "amount": 100,
        "fromUser": "R122111"
      },
      {
        "amount": 100,
        "fromUser": "R122111"
      },
      {
        "amount": 100,
        "fromUser": "R122111"
      },
      {
        "amount": 100,
        "fromUser": "R122111"
      },
      {
        "amount": 100,
        "fromUser": "R122111"
      },
      {
        "amount": 100,
        "fromUser": "R122111"
      },
      {
        "amount": 100,
        "fromUser": "R122111"
      }
    ],
    "sent": [
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 10,
        "toUser": "R122111"
      },
      {
        "amount": 1,
        "toUser": "R122111"
      },
      {
        "amount": 1,
        "toUser": "R122111"
      },
      {
        "amount": 1,
        "toUser": "R122111"
      },
      {
        "amount": 1,
        "toUser": "R122111"
      },
      {
        "amount": 1,
        "toUser": "R122111"
      },
      {
        "amount": 1,
        "toUser": "R122111"
      },
      {
        "amount": 1,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 2,
        "toUser": "R122111"
      },
      {
        "amount": 15,
        "toUser": "R122111"
      },
      {
        "amount": 15,
        "toUser": "R122111"
      }
    ]
  },
  "coins": 5,
  "inventory": [
    {
      "quantity": 2,
      "type": "book"
    },
    {
      "quantity": 4,
      "type": "cup"
    },
    {
      "quantity": 2,
      "type": "hoody"
    },
    {
      "quantity": 2,
      "type": "pen"
    },
    {
      "quantity": 3,
      "type": "powerbank"
    },
    {
      "quantity": 2,
      "type": "t-shirt"
    }
  ]
}
