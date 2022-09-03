# udptun

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
