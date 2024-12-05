package mymux

import (
	"testing"
)

// TestCtrlFrame 测试 MyFrame 结构的各个方法。
func TestCtrlFrame(t *testing.T) {
	// 创建一个示例 CtrlFrame 数据
	frame := Frame{0x01, 0x02, 0x03, 0x04, 0x7F, 0x00, 0x01}

	// 测试 Source 方法
	if got := frame.Source(); got != 0x0102 {
		t.Errorf("Source() = %v; want 0x0102", got) // 期望结果为 0x0102
	}

	// 测试 Destination 方法
	if got := frame.Destination(); got != 0x0304 {
		t.Errorf("Destination() = %v; want 0x0304", got) // 期望结果为 0x0304
	}

	// 测试 Port 方法
	if got := frame.Port(); got != 0x7F {
		t.Errorf("Port() = %v; want 0x7F", got) // 期望结果为 0x7F
	}

	// 测试 SetSource 方法
	frame.SetSource(0x1234)
	if got := frame.Source(); got != 0x1234 {
		t.Errorf("SetSource() = %v; want 0x1234", got) // 期望结果为 0x1234
	}

	// 测试 SetDestination 方法
	frame.SetDestination(0x5678)
	if got := frame.Destination(); got != 0x5678 {
		t.Errorf("SetDestination() = %v; want 0x5678", got) // 期望结果为 0x5678
	}

	// 测试 SetPort 方法
	frame.SetPort(0x7E)
	if got := frame.Port(); got != 0x7E {
		t.Errorf("SetPort() = %v; want 0x7E", got) // 期望结果为 0x7E
	}

	// 测试 SetCommand 方法
	frame.SetCommand(Data) // 设置为 Data 类型
	if got := frame[5]; got != byte(Data) {
		t.Errorf("SetCommand() = %v; want %v", got, byte(Data)) // 期望结果为 Data 的字节值
	}

	// 测试 SetSequenceNumber 方法
	frame.SetSequenceNumber(0x02)
	if got := frame.SequenceNumber(); got != 0x02 {
		t.Errorf("SetSequenceNumber() = %v; want 0x02", got) // 期望结果为 0x02
	}
}
