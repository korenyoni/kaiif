Date: 2018-05-06
Title: Promethues: default values
cat: ops

If you have queries that may not return values for some parameters, e.g:

```
sum by (code,path)(rate(http_server_requests_total{path="/api/obscure_path/"}[5m])) * 60
```

You can OR this time-series vector with a time-series vector with all values 0:

```
sum by (code,path)(rate(http_server_requests_total{path="/api/obscure_path/"}[5m])) * 60 OR on() vector(0)
```

This is useful if this obscure metric will otherwise cause an alarm to have a status of "NO DATA".
