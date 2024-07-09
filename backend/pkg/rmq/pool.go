package rmq

import (
	"log/slog"
	"sync"
)

type RMQConfig struct {
	Consumers []*ConsumerOption `toml:"consumers"`
	Producers []*ProducerOption `toml:"producers"`
}

var (
	consumerPool sync.Map
	producerPool sync.Map
	once         sync.Once
)

func Init(config RMQConfig) error {
	var err error
	once.Do(func() {
		for _, opts := range config.Consumers {
			c, err := NewConsumer(opts)
			if err != nil {
				return
			}
			consumerPool.Store(opts.Name, c)
		}
		for _, opts := range config.Producers {
			p, err := NewProducer(opts)
			if err != nil {
				return
			}
			producerPool.Store(opts.Name, p)
		}
	})
	return err
}

func Shutdown() {
	consumerPool.Range(func(name, v any) bool {
		c, ok := v.(*consumer)
		if !ok {
			return true
		}
		_ = c.Shutdown()
		slog.With("name", name).Info("consumer已关闭")
		return true
	})
	producerPool.Range(func(name, v any) bool {
		p, ok := v.(*producer)
		if !ok {
			return true
		}
		_ = p.Shutdown()
		slog.With("name", name).Info("producer已关闭")
		return true
	})

}

func Producer(name string) (*producer, bool) {
	v, ok := producerPool.Load(name)
	if !ok {
		return nil, false
	}
	p, ok := v.(*producer)
	if !ok {
		return nil, false
	}
	return p, true
}

func Consumer(name string) (*consumer, bool) {
	v, ok := consumerPool.Load(name)
	if !ok {
		return nil, false
	}
	c, ok := v.(*consumer)
	if !ok {
		return nil, false
	}
	return c, true
}
