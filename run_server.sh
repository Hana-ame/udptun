release/udptun -d "localhost:4000" -n server -h "http://localhost:8888/"
kcptun/server_linux_amd64 -t "127.0.0.1:22" -l ":4000" -mode fast3 -nocomp -sockbuf 16777217 -dscp 46