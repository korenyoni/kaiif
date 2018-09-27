Date: 2018-04-27
Title: Prometheus queries: negative lookahead
cat: ops

I've recently gotten stuck trying to use negative lookahead in Prometheus queries with the `~=` regex matcher:

### This query doesn't work

```
sum by (code,path)(rate(http_server_requests_total{path=~"^(/api/v1)(/api/v2)(/api/v3).*"}[5m])) * 60
```

In this example, we are trying to retrieve all sub-paths of `/api/v2` and `/api/v3` but not `/api/v1`.
However this is not possible in Prometheus since it uses the [re2](https://github.com/google/re2) library, which doesn't support
negative lookahead.

I was scratching my head for some time, but soon I realized Prometheus doesn't need negative-lookahead since it supports the `!~` (negative regex) matcher:


### This query works

```
sum by (code,path)(rate(http_server_requests_total{path=~"(/api/v2)(/api/v3).*", path!~"(/api/v1).*"}[5m])) * 60
```

By using both Prometheus's positive `=~` and negative `!~` regex matchers, we can get around [re2](https://github.com/google/re2)'s lack of negative lookahead.
