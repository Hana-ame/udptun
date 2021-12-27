./udptun.exe -mode "client" -l ":6000" -r "127.0.0.1:10000"
./udptun.exe -mode "server" -l ":10000" -r "127.0.0.1:4000"

go build -o u -mod vendor ./udptun &

./s -t "localhost:8080" -l ":4000"
./u -mode "server" -l "gcp" -r "127.0.0.1:4000"

./udptun.exe -mode "client" -l ":6000" -r "gcp"
./client.exe  -r "localhost:6000" -l ":8080"


Known issue:
when connecting, it seems that only one connection can work at the same time.
that may cause some problems. will be test at 211227 night.
