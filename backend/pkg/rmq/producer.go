package rmq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
)

type producer struct {
	options *ProducerOption
	rmq     rocketmq.Producer
}

func (r *producer) Push(messageExt *MessageExt, callback PushCallback) {
	messages := messageExt.Message()
	if len(messages) == 0 {
		return
	}
	result, err := r.rmq.SendSync(context.TODO(), messages...)
	ctx := context.WithValue(context.TODO(), "message", messageExt)
	callback(ctx, result, err)
}

func (r *producer) Start() error {
	return r.rmq.Start()
}

func (r *producer) Shutdown() error {
	return r.rmq.Shutdown()
}

func NewProducer(opts *ProducerOption) (*producer, error) {
	rmq, err := rocketmq.NewProducer(opts.Options()...)
	if err != nil {
		return nil, err
	}
	p := &producer{
		options: opts,
		rmq:     rmq,
	}
	err = rmq.Start()
	if err != nil {
		return nil, err
	}
	return p, nil
}
