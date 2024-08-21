package notifier

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type DelayNotify struct {
	Notify
	mu          sync.Mutex
	done        chan struct{}
	wg          sync.WaitGroup
	messages    chan string
	maxMessages int           // 消息收集器的最大容量
	collectFreq time.Duration // 收集频率
}

// NewDelayNotifier 创建一个新的延迟通知实例
func NewDelayNotifier(webhookURL string, textBodyFormat string) Notifier {
	delayNotify := &DelayNotify{
		Notify: Notify{
			client:         http.DefaultClient,
			webhookURL:     webhookURL,
			method:         http.MethodPost, // 默认方法为 POST
			headers:        map[string]string{"Content-Type": "application/json"},
			textBodyFormat: textBodyFormat,
		},
		mu:          sync.Mutex{},
		done:        make(chan struct{}),
		wg:          sync.WaitGroup{},
		messages:    make(chan string, 1000),
		maxMessages: 5,
		collectFreq: time.Duration(500) * time.Millisecond,
	}
	delayNotify.wg.Add(1)
	go delayNotify.collectAndSendMessages()
	return delayNotify
}

func (n *DelayNotify) SendMessage(ctx context.Context, message string) error {
	n.messages <- message
	return nil
}

// collectAndSendMessages 负责定期合并并发送消息
func (n *DelayNotify) collectAndSendMessages() {
	defer n.wg.Done()
	ticker := time.NewTicker(n.collectFreq)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n.mu.Lock()
			n.sendCollectedMessages()
			n.mu.Unlock()
		case <-n.done:
			// 在关闭时处理所有剩余的消息
			n.mu.Lock()
			n.sendCollectedMessages()
			n.mu.Unlock()
			return
		}
	}
}

// sendCollectedMessages 合并并发送收集到的消息
func (n *DelayNotify) sendCollectedMessages() {
	// 如果有消息则合并并发送
	if len(n.messages) == 0 {
		return
	}
	var collectedMessages []string
	//for len(collectedMessages) < n.maxMessages {
	//	select {
	//	case msg := <-n.messages:
	//		collectedMessages = append(collectedMessages, msg)
	//	default:
	//		break
	//	}
	//}
	for {
		select {
		case msg := <-n.messages:
			collectedMessages = append(collectedMessages, msg)
		default:
			goto sendMessages
		}
	}

sendMessages:
	if len(collectedMessages) == 0 {
		return
	}

	err := n.Notify.SendMessage(context.Background(), strings.Join(collectedMessages, "\\n===\\n"))
	if err != nil {
		log.Println(err)
	}

}

// Close 关闭通知器并清理资源
func (n *DelayNotify) Close() error {
	close(n.done)
	n.wg.Wait()
	close(n.messages)
	return nil
}
