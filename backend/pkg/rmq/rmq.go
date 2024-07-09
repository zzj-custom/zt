package rmq

import (
	rmqConsumer "github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	rmqProducer "github.com/apache/rocketmq-client-go/v2/producer"
	"time"
)

type Credentials struct {
	AccessKey     string `toml:"access_key"`
	SecretKey     string `toml:"secret_key"`
	SecurityToken string `toml:"security_token"`
}

func (r *Credentials) AsPrimitive() primitive.Credentials {
	return primitive.Credentials{
		AccessKey:     r.AccessKey,
		SecretKey:     r.SecretKey,
		SecurityToken: r.SecurityToken,
	}
}

type ConsumerOption struct {
	Name     string   `toml:"name"`
	Endpoint []string `toml:"endpoint"`
	Instance string   `toml:"instance"`

	Retries          int          `toml:"retries"`
	Group            string       `toml:"group"`
	PullBatchSize    int          `toml:"pull_batch_size"`
	ConsumeBatchSize int          `toml:"consume_batch_size"`
	Order            bool         `toml:"order"`
	Namespace        string       `toml:"namespace"`
	VIPChannel       bool         `toml:"vip_channel"`
	FromWhere        int          `toml:"from_where"`
	ConsumeTimestamp string       `toml:"consume_timestamp"`
	Credentials      *Credentials `toml:"credentials"`
}

func (r *ConsumerOption) Options() []rmqConsumer.Option {
	opts := []rmqConsumer.Option{
		rmqConsumer.WithNameServer(r.Endpoint),
	}
	if r.Group != "" {
		opts = append(opts, rmqConsumer.WithGroupName(r.Group))
	}
	if r.Retries >= 0 {
		opts = append(opts, rmqConsumer.WithRetry(r.Retries))
	}
	if r.PullBatchSize > 0 {
		opts = append(opts, rmqConsumer.WithPullBatchSize(int32(r.PullBatchSize)))
	}
	if r.ConsumeBatchSize > 0 {
		opts = append(opts, rmqConsumer.WithConsumeMessageBatchMaxSize(r.ConsumeBatchSize))
	}
	if r.Order {
		opts = append(opts, rmqConsumer.WithConsumerOrder(true))
	}
	if r.Namespace != "" {
		opts = append(opts, rmqConsumer.WithNamespace(r.Namespace))
	}
	if r.VIPChannel {
		opts = append(opts, rmqConsumer.WithVIPChannel(r.VIPChannel))
	}
	if r.FromWhere >= 0 {
		opts = append(opts, rmqConsumer.WithConsumeFromWhere(rmqConsumer.ConsumeFromWhere(r.FromWhere)))
		if r.FromWhere == int(rmqConsumer.ConsumeFromTimestamp) {
			consumeTimestamp := r.ConsumeTimestamp
			if r.ConsumeTimestamp == "" {
				consumeTimestamp = "19700101000000"
			}
			opts = append(opts, rmqConsumer.WithConsumeTimestamp(consumeTimestamp))
		}
	}
	if r.Credentials != nil {
		opts = append(opts, rmqConsumer.WithCredentials(r.Credentials.AsPrimitive()))
	}
	if r.Instance != "" {
		opts = append(opts, rmqConsumer.WithInstance(r.Instance))
	}
	return opts
}

type ProducerOption struct {
	Name     string   `toml:"name"`
	Endpoint []string `toml:"endpoint"`
	Instance string   `toml:"instance"`

	Retries     int          `toml:"retries"`
	Group       string       `toml:"group"`
	Namespace   string       `toml:"namespace"`
	VIPChannel  bool         `toml:"vip_channel"`
	Credentials *Credentials `toml:"credentials"`
	SendTimeout int          `toml:"send_timeout"`
}

func (r *ProducerOption) Options() []rmqProducer.Option {
	var opts []rmqProducer.Option
	opts = append(opts, rmqProducer.WithNameServer(r.Endpoint))
	if r.Group != "" {
		opts = append(opts, rmqProducer.WithGroupName(r.Group))
	}
	if r.Instance != "" {
		opts = append(opts, rmqProducer.WithInstanceName(r.Instance))
	}
	if r.Namespace != "" {
		opts = append(opts, rmqProducer.WithNamespace(r.Namespace))
	}
	if r.Retries > 0 {
		opts = append(opts, rmqProducer.WithRetry(r.Retries))
	}
	if r.Credentials != nil {
		opts = append(opts, rmqProducer.WithCredentials(r.Credentials.AsPrimitive()))
	}
	if r.VIPChannel {
		opts = append(opts, rmqProducer.WithVIPChannel(true))
	}
	if r.SendTimeout > 0 {
		opts = append(opts, rmqProducer.WithSendMsgTimeout(time.Duration(r.SendTimeout)*time.Second))
	}
	return opts
}
