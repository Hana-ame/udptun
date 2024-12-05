package wsreverse

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	// 创建一个超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动一个goroutine执行任务
	go func() {
		// 模拟一个长时间的任务
		time.Sleep(10 * time.Second)
		t.Log("任务完成")
	}()

	// 等待任务完成或超时
	select {
	case <-ctx.Done():
		t.Log("超时终止")
	}

	time.Sleep(time.Minute)
}

func TestTimeoutTask(t *testing.T) {
	// 创建一个超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动一个goroutine执行任务
	go func() {
		task(ctx, t)
	}()

	// 等待任务完成或超时
	select {
	case <-ctx.Done():
		t.Log("超时终止")
	}

	time.Sleep(time.Minute)

}

func task(ctx context.Context, t *testing.T) {
	// 模拟一个长时间的任务
	for {
		select {
		case <-ctx.Done():
			t.Log("任务终止")
			return
		default:
			t.Log("任务正在执行")
			time.Sleep(1 * time.Second)
		}
	}
}

func TestConn(t *testing.T) {
	a, _ := net.Dial("tcp", "localhost:8080")
	a.Read(nil)
}

type a struct {
	int
}

func (a a) Print() {
	fmt.Println(a)
	fmt.Println(a.int)
}
func (a *a) Printp() {
	fmt.Println(a)
	fmt.Println(a.int)
}
func TestDefer(t *testing.T) {
	a := &a{1}
	func() {
		defer a.Print()  // {1}, 1
		defer a.Printp() // &{2}, 2

		a.int = 2

		a.Print()  // {2}, 2
		a.Printp() // &{2}, 2
	}()
}
