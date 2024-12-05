package examples

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// mux version. notused.

// func DownloadFileBySha1sum(url string, requestHeader http.Header, sha1sum string, filePath string) {
// 	// 建立 WebSocket 连接
// 	wsc, _, err := websocket.DefaultDialer.Dial(url, requestHeader)
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 	}
// 	defer wsc.Close()

// 	mux := wsmux.NewWsMux(wsc, wsmux.MuxSeqClient)
// go mux.ReadDaemon(wsc)
// 	mxc, _ := mux.Dial()
// 	buf := make([]byte, 1024) // 1KB缓冲区
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		log.Println("打开文件失败:", err)
// 		return
// 	}
// 	defer file.Close()

// 	io.CopyBuffer(file, mxc, buf)
// }

func DownloadFileBySha1sum(url string, requestHeader http.Header, sha1sum string, filePath string) {
	// 建立 WebSocket 连接
	ws, _, err := websocket.DefaultDialer.Dial(url, requestHeader)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	file, err := os.Create(filePath)
	if err != nil {
		log.Println("打开文件失败:", err)
		return
	}
	defer file.Close()

	ws.WriteMessage(websocket.BinaryMessage, []byte(sha1sum))
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		file.Write(msg)
	}
}
