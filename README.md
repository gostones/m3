Mirr
-----

An IPFS DApp that mirrors the old web, a.k.a,  the world wide web.
it has a built-in forward proxy and load balancer to distribute the load to multiple peers.


### Usage

```
$ ./mirr -port 18080
```

### Credits

https://github.com/voldyman/GoLoadBalance

https://github.com/kintoandar/fwd

https://github.com/elazarl/goproxy

<!-- https://github.com/FelisCatus/SwitchyOmega -->
<!-- https://github.com/PuerkitoBio/gocrawl -->
<!-- https://github.com/gocolly/colly -->

### License

Mirr is released under MIT license

Author: Qiang Li <liqiang@gmail.com>

<!--
https://docs.ipfs.io/reference/api/http/

curl "http://localhost:5001/api/v0/swarm/addrs/local?id=<value>"
curl http://127.0.0.1:5001/api/v0/swarm/peers

curl "http://localhost:5001/api/v0/p2p/stream/dial?arg=<Peer>&arg=<Protocol>&arg=<BindAddress>"

-->