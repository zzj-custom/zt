package rmq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"sync"
)

type PushCallback func(ctx context.Context, result *primitive.SendResult, err error)

type Property struct {
	values map[string]string
	once   sync.Once
	mu     sync.Mutex
}

func (r *Property) Set(key, val string) {
	r.once.Do(func() {
		r.values = map[string]string{}
	})
	r.mu.Lock()
	defer func() {
		r.mu.Unlock()
	}()
	r.values[key] = val
}

func (r *Property) Remove(key string) {
	r.mu.Lock()
	defer func() {
		r.mu.Unlock()
	}()
	delete(r.values, key)
}

func (r *Property) Reset() {
	r.mu.Lock()
	defer func() {
		r.mu.Unlock()
	}()
	r.values = map[string]string{}
}

func (r *Property) Values() map[string]string {
	r.mu.Lock()
	defer func() {
		r.mu.Unlock()
	}()
	return r.values
}

type MessageExt struct {
	Topic      string
	Tag        string
	BizKey     []string
	DelayLevel int
	Properties *Property
	bodies     [][]byte
}

func (r *MessageExt) WithTag(tag string) {
	r.Tag = tag
}

func (r *MessageExt) AppendBody(body []byte) {
	r.bodies = append(r.bodies, body)
}

func (r *MessageExt) SetBodies(bodies [][]byte) {
	if len(bodies) == 0 {
		return
	}
	r.bodies = bodies
}

func (r *MessageExt) Reset() {
	r.bodies = [][]byte{}
}

func (r *MessageExt) Message() []*primitive.Message {
	var messages []*primitive.Message
	for _, body := range r.bodies {
		m := primitive.NewMessage(r.Topic, body)
		if r.Tag != "" {
			m.WithTag(r.Tag)
		}
		if len(r.BizKey) != 0 {
			m.WithKeys(r.BizKey)
		}
		if r.DelayLevel > 0 {
			m.WithDelayTimeLevel(r.DelayLevel)
		}
		if r.Properties != nil {
			m.WithProperties(r.Properties.Values())
		}
		messages = append(messages, m)
	}
	return messages
}

func NewMessage(topic string, body []byte) *MessageExt {
	m := &MessageExt{
		Topic:  topic,
		bodies: [][]byte{},
	}
	if body != nil {
		m.AppendBody(body)
	}
	return m
}
