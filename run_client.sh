release/udptun -a "localhost:8080" -n client -h "http://localhost:8888/" &
sleep 2
curl --location --request POST "localhost:8080/server" --header "Content-Type: text/plain" --data-raw "127.0.0.1:4001"
kcptun/client_linux_amd64 -r "127.0.0.1:4001" -l ":2222" -mode fast3 -nocomp -autoexpire 900 -sockbuf 16777217 -dscp 46

