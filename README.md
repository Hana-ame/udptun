# udptun

a tiny tool for hole-punching

内网穿透用小玩具。自己写是因为搜索能力欠佳搜不到能用的……

总之nat类型按照从宽松到严格三种。

> Full Cone, Restricted Cone, Symmetric

最后一种因为 local addr 和 destination addr 对不上所以打不了洞的，请直接放弃。

Full Cone, 差不多就是公网IP，没打洞必要。所以就针对 Restricted Cone (有 Port Restricted Cone 类型但是意义不大，过程都一样的)

总之通过中介服务器你发一个空包我发一个空包就打通了，然后就能转发了…这个 repo 就是做了这件事

因为协议是udp所以还通过 [hysteri](https://github.com/HyNetwork/hysteria) 转了一下，能做代理。

其实v6更方便一点，没写不好意思。而且我家v6地址突然给我扬了，正好就不写了。

不过现在还有稀奇古怪的bug而且写成一坨屎山，不合理的地方也一堆，要重写懒得写凑合用吧。

稀奇古怪的bug比如运行三次会有一次不成功之类的……不过server端比较宽松，只能再苦一苦client了。

哦对了，校园网因为国内外路线的问题，连不了。可能真的还得整ipv6的适配……

[helper server 源代码](https://github.com/Hana-ame/Toys/blob/master/helper-server/server.go)

## build
```
go build .
```

## usage

> -s
when specified, it works as server. Otherwise it works as client.

> -a
when works as server, it refers to the address forward to.
when works as client, it's the address that listen connection.

> -h
the helper's host

> -p
the helper's path

helper should be the same, and only one server and one client is allowed.

## 

其实只要去[release](https://github.com/Hana-ame/udptun/releases)下载就好了，两个bat跑一下就行，记得改path，而且server最好自己做。[hysteri](https://github.com/HyNetwork/hysteria)在隔壁下载。
