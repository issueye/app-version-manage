package main

import (
	"fmt"

	bootstrap "github.com/issueye/app-version-manage/bootstrap"

	"github.com/asticode/go-astilectron"
)

var isMax = false

type Message struct {
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

func sendMessageTest(w *astilectron.Window) {

	// 发送消息，收不到返回数据
	bootstrap.SendMessage(w, "test", "test data", func(m *bootstrap.MessageIn) {})
}

// handleMessages handles messages
func handleMessages(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	if m.Name == "hello" {
		// 发送测试
		sendMessageTest(w)

		return "world", nil
	}

	// 最小化
	if m.Name == "minWindow" {
		fmt.Println("最小化窗口")
		w.Minimize()
		return "", nil
	}

	// 全屏
	if m.Name == "screen" {
		fmt.Println("isMax", isMax)
		if isMax {
			isMax = !isMax
			w.Unmaximize()
			return "恢复窗口", nil
		} else {
			isMax = !isMax
			w.Maximize()
			return "最大化窗口", nil
		}
	}

	if m.Name == "close-app" {
		w.Close()
		return "", nil
	}
	return nil, nil
}
