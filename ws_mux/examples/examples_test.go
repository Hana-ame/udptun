package examples

import "testing"

func TestResponse(t *testing.T) {
	// UploadFileBySha1sum("ws://127.0.0.1:8080/ws/server", nil)
	UploadFileBySha1sum("wss://file.moonchan.xyz/ws/server", nil)
}

func TestRequest(t *testing.T) {
	DownloadFileBySha1sum("ws://127.0.0.1:8080/ws", nil, "../source.txt", "destination.txt")
}
