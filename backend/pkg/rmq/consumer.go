package rmq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	rmqConsumer "github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type consumer struct {
	options *ConsumerOption
	rmq     rocketmq.PushConsumer
}

func (r *consumer) Subscribe(
	topic string,
	tag string,
	f func(ctx context.Context, ext ...*primitive.MessageExt) (rmqConsumer.ConsumeResult, error),
) error {
	selector := rmqConsumer.MessageSelector{}
	if tag != "" {
		selector.Type = rmqConsumer.TAG
		selector.Expression = tag
	}
	return r.rmq.Subscribe(topic, selector, f)
}

func (r *consumer) Shutdown() error {
	return r.rmq.Shutdown()
}

func NewConsumer(opt *ConsumerOption) (*consumer, error) {
	pushConsumer, err := rocketmq.NewPushConsumer(opt.Options()...)
	if err != nil {
		return nil, err
	}
	err = pushConsumer.Start()
	if err != nil {
		return nil, err
	}
	return &consumer{
		options: opt,
		rmq:     pushConsumer,
	}, nil
}
