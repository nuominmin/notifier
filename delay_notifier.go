package notifier

import (
	"context"
	"log"
	"math"
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

type delayNotify struct {
	Notifier
	sep         string // 分隔符
	mu          sync.Mutex
	done        chan struct{}
	wg          sync.WaitGroup
	messages    chan string
	maxMessages int           // 每次合并发送的最大消息数
	collectFreq time.Duration // 收集频率
}

// NewDelayNotifier 创建一个新的延迟通知实例
func NewDelayNotifier(webhookUrlFormat string, textBodyFormat string, tokens ...string) DelayNotifier {
	dn := &delayNotify{
		sep:         "\\n===\\n",
		Notifier:    NewNotifier(webhookUrlFormat, textBodyFormat, tokens...),
		mu:          sync.Mutex{},
		done:        make(chan struct{}),
		wg:          sync.WaitGroup{},
		messages:    make(chan string, 1000),
		maxMessages: 5,
		collectFreq: time.Duration(500) * time.Millisecond,
	}
	dn.wg.Add(1)
	go dn.collectAndSendMessages()
	return dn
}

// SetSep 设置分隔符
func (n *delayNotify) SetSep(sep string) DelayNotifier {
	n.sep = sep
	return n
}

// SetCollectFreq 设置收集频率
func (n *delayNotify) SetCollectFreq(freq time.Duration) DelayNotifier {
	n.collectFreq = freq
	return n
}

// SetMaxMessages 设置每次合并发送的最大消息数
func (n *delayNotify) SetMaxMessages(maxMessages int) DelayNotifier {
	n.maxMessages = maxMessages
	return n
}

func (n *delayNotify) SendMessage(ctx context.Context, message string) error {
	n.messages <- message
	return nil
}

// collectAndSendMessages 负责定期合并并发送消息
func (n *delayNotify) collectAndSendMessages() {
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
func (n *delayNotify) sendCollectedMessages() {
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
		default:
			goto sendMessages
		}
	}

sendMessages:
	collectedMessagesLen := len(collectedMessages)
	if collectedMessagesLen == 0 {
		return
	}

	group := n.groupInMax(make([]int, collectedMessagesLen), n.maxMessages)
	for i := 0; i < len(group); i++ {
		start := i * n.maxMessages
		end := start + n.maxMessages
		if end > collectedMessagesLen {
			end = collectedMessagesLen
		}
		err := n.Notifier.SendMessage(context.Background(), strings.Join(collectedMessages[start:end], n.sep))
		if err != nil {
			log.Println("send message error:", err.Error())
		}
	}
}

func (n *delayNotify) groupInMax(array []int, max int) (res [][]int) {
	length := len(array)

	if length <= 0 {
		return [][]int{}
	}

	// 小于等于 max 直接返回原数据
	if length <= max || max == 0 {
		return [][]int{array}
	}

	// 分组数量 = 长度除以最大数
	groupNum := int(math.Ceil(float64(length) / float64(max)))

	// 开始和结束的指针位置
	var start, end int
	for i := 1; i <= groupNum; i++ {
		end = i * max
		if end > length {
			end = length
		}
		res = append(res, array[start:end])
		start = end
	}
	return res
}

// Close 关闭通知器并清理资源
func (n *delayNotify) Close() error {
	close(n.done)
	n.wg.Wait()
	close(n.messages)
	return nil
}
