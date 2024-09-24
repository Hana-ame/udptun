
# usage

```sh
# as server
# ./udptun.exe -d "[destination addr:port]" -n [name] -h "[help server]" -m "[udp4|udp6]"
./udptun.exe -d "localhost:4000" -n server -h "http://localhost:8888/"

# as client
# ./udptun.exe -a "[control addr:port]" -h "[help server]" -m "[udp4|udp6]"
./udptun.exe -a "localhost:8080" -n client -h "http://localhost:8888/"

# as help server
./udptun.exe -isHelpServer -h "127.0.0.1:8888"
```

create a connection

```sh
# curl --location --request POST '[help server]/[server name]' \
# --header 'Content-Type: text/plain' \
# --data-raw '[local addr:port]'

curl.exe --location --request POST 'localhost:8080/server' --header 'Content-Type: text/plain' --data-raw '127.0.0.1:4001'
```

## kcptun
```sh
./client_windows_386.exe -r "127.0.0.1:4001" -l ":2222" -mode fast3 -nocomp -autoexpire 900 -sockbuf 16777217 -dscp 46
./server_windows_386.exe -t "127.0.0.1:22" -l ":4000" -mode fast3 -nocomp -sockbuf 16777217 -dscp 46
```


# dev
git submodule add -b go git@github.com:Hana-ame/Tools.git
git submodule add  git@github.com:Hana-ame/tcp-over-websocket.git


悟了，走proxy巨慢无比，走直连还有穿透问题。基本没法搞。


# tcp-over-websocket
总觉得太难看了，自己写也写出来了
但是自己写速度不达标

server
		go startWsServer(listenHostPort, isSsl, sslCrt, sslKey)


client_port -------------------- server_port
client ---------- websocket ----------server