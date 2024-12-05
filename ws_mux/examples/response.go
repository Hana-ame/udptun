package examples

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	db "github.com/Hana-ame/udptun/Tools/db_filehash"
	wsmux "github.com/Hana-ame/udptun/Tools/ws_mux"

	"github.com/gorilla/websocket"
)

func UploadFileBySha1sum(url string, requestHeader http.Header) {
	// 建立 WebSocket 连接
	wsc, _, err := websocket.DefaultDialer.Dial(url, requestHeader)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer wsc.Close()

	mux := wsmux.NewWsMux(wsc, wsmux.MuxSeqClient)
	go func() {
		for {
			mxc := mux.Accept()
			log.Println("Accepted:", mxc)
			go handleRequestFileBySha1sum(mxc)
		}
	}()

	mux.ReadDaemon(wsc)
	log.Println("exit")
}

func handleRequestFileBySha1sum(muc *wsmux.WsMuxConn) {
	defer muc.Close()

	pkg := muc.ReadPackage()
	sha1sum := string(pkg.Message)

	filepath, err := getPathFromSha1sum(sha1sum)
	if err != nil {
		log.Println(filepath, err)
		return
	}

	log.Println(filepath)

	// buf := make([]byte, 1024) // 1KB缓冲区
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	defer file.Close()

	// 这段代码用了之后文件尾有问题。
	// closed := false
	// stuck := make(chan struct{}, 1)
	// go func() {
	// 	for !closed {
	// 		buf := make([]byte, 1024*32)
	// 		n, err := file.Read(buf)
	// 		if err != nil {
	// 			stuck <- struct{}{}
	// 			break
	// 		}
	// 		if _, err := muc.Write(buf[:n]); err != nil {
	// 			stuck <- struct{}{}
	// 			break
	// 		}
	// 	}
	// }()

	// select {
	// case <-stuck:
	// case <-muc.ReadChan:
	// }

	// 不能加，加了会导致文件卡死在最后0B，还会断流
	// time.Sleep(time.Minute * 2)
	// closed = true

	buf := make([]byte, 1024)
	n, err := io.CopyBuffer(muc, file, buf)

	log.Println("handleRequestFileBySha1sum", n, err)
	// log.Println("handleRequestFileBySha1sum")

}

// このまま
func getPathFromSha1sum(sha1sum string) (string, error) {
	return db.ReadPathByHash(sha1sum)
}
