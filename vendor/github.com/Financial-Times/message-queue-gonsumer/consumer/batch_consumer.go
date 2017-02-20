package consumer

import "net/http"

// NewBatchedQueueConsumer returns a new instance of a QueueConsumer that handles batches of messages
func NewBatchedQueueConsumer(config QueueConfig, handler func(m []Message), client http.Client) QueueConsumer {
	offset := "largest"
	if len(config.Offset) > 0 {
		offset = config.Offset
	}
	queue := &defaultQueueCaller{
		addrs:            config.Addrs,
		group:            config.Group,
		topic:            config.Topic,
		offset:           offset,
		autoCommitEnable: config.AutoCommitEnable,
		caller:           defaultHTTPCaller{config.Queue, config.AuthorizationKey, client},
	}
	return &DefaultQueueConsumer{config, queue, nil, make(chan bool, 1), BatchedMessageProcessor{handler}}
}

//BatchedMessageProcessor process messages in batches
type BatchedMessageProcessor struct {
	handler func(m []Message)
}

func (b BatchedMessageProcessor) consume(msgs ...Message) {
	if len(msgs) > 0 {
		b.handler(msgs)
	}
}