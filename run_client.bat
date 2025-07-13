start release/udptun.exe -a "localhost:8080" -n client -h "http://localhost:8888/"
timeout /t 2
curl.exe --location --request POST "localhost:8080/server" --header "Content-Type: text/plain" --data-raw "127.0.0.1:4001"
start kcptun/client_windows_386.exe -r "127.0.0.1:4001" -l ":2222" -mode fast3 -nocomp -autoexpire 900 -sockbuf 16777217 -dscp 46
timeout /t 30
