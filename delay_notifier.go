package notifier

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

// DelayNotifier 消息通知接口
type DelayNotifier interface {
	Notifier
	SetCollectFreq(freq time.Duration) DelayNotifier
	SetSep(sep string) DelayNotifier
	SetMaxMessages(maxMessages int) DelayNotifier
	Close() error
}

type DelayNotify struct {
	sep         string // 分隔符
	notifier    Notifier
	mu          sync.Mutex
	done        chan struct{}
	wg          sync.WaitGroup
	messages    chan string
	maxMessages int           // 每次合并发送的最大消息数
	collectFreq time.Duration // 收集频率
}

// NewDelayNotifier 创建一个新的延迟通知实例
func NewDelayNotifier(webhookURL string, textBodyFormat string) DelayNotifier {
	delayNotify := &DelayNotify{
		sep:         "\\n===\\n",
		notifier:    NewNotifier(webhookURL, textBodyFormat),
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

// SetSep 设置分隔符
func (n *DelayNotify) SetSep(sep string) DelayNotifier {
	n.sep = sep
	return n
}

// SetCollectFreq 设置收集频率
func (n *DelayNotify) SetCollectFreq(freq time.Duration) DelayNotifier {
	n.collectFreq = freq
	return n
}

// SetMaxMessages 设置每次合并发送的最大消息数
func (n *DelayNotify) SetMaxMessages(maxMessages int) DelayNotifier {
	n.maxMessages = maxMessages
	return n
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
			n.sendCollectedMessages()
		case <-n.done:
			// 在关闭时处理所有剩余的消息
			n.sendCollectedMessages()
			return
		}
	}
}

// sendCollectedMessages 合并并发送收集到的消息
func (n *DelayNotify) sendCollectedMessages() {
	n.mu.Lock()
	defer n.mu.Unlock()

	// 如果有消息则合并并发送
	if len(n.messages) == 0 {
		return
	}

	mapMessage := make(map[string]struct{})
	collectedMessages := make([]string, 0)
	for {
		select {
		case msg := <-n.messages:
			if _, ok := mapMessage[msg]; !ok {
				collectedMessages = append(collectedMessages, msg)
				mapMessage[msg] = struct{}{}
			}
			if len(collectedMessages) >= n.maxMessages {
				goto sendMessages
			}
		default:
			goto sendMessages
		}
	}

sendMessages:
	if len(collectedMessages) == 0 {
		return
	}

	err := n.notifier.SendMessage(context.Background(), strings.Join(collectedMessages, n.sep))
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
