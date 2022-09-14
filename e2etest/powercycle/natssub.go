package powercycle

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	metricspb "github.com/ci4rail/bogie-pdm/proto/go/metrics/v1"
	nats "github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

// NatsSubscriber is a wrapper for a NATS subscription
type NatsSubscriber struct {
	nc  *nats.Conn
	sub *nats.Subscription
}

type messageEnvelope struct {
	DataBase64 string
}

//type VeriferFunc func(subject string, m msg.Message) error

// NewNatsSubscriber creates a new NatsSubscriber for
// subject pattern provided as consumer on stream
func NewNatsSubscriber(natsURL string, credsFile string, subject string, consumer string, stream string) (*NatsSubscriber, error) {
	opts := []nats.Option{nats.UserCredentials(credsFile)}
	nc, err := nats.Connect(natsURL, opts...)
	if err != nil {
		return nil, err
	}
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, err
	}
	_, err = js.AddConsumer(stream, &nats.ConsumerConfig{
		Durable: consumer, AckPolicy: nats.AckExplicitPolicy, ReplayPolicy: nats.ReplayInstantPolicy, FilterSubject: subject,
	})
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("can't add consumer: %v", err)
	}
	sub, err := js.PullSubscribe(subject, consumer, nats.Bind(stream, consumer))
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("can't create subscription: %v", err)
	}
	return &NatsSubscriber{nc, sub}, nil
}

// Close closes the connection to the NATS server.
func (n *NatsSubscriber) Close() {
	n.nc.Close()
}

// WaitMsg checks for a message with a timestamp greater ts
func (n *NatsSubscriber) WaitMsg(timeout time.Duration, ts time.Time) (bool, error) {
	messages, err := n.sub.Fetch(10, nats.PullOpt(nats.MaxWait(timeout)))
	if err == nats.ErrTimeout {
		fmt.Print("Timeout reading stream\n")
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		for _, m := range messages {
			//fmt.Printf("Got message: %s\n", string(m.Data))
			var data messageEnvelope
			err := json.Unmarshal(m.Data, &data)
			if err != nil {
				return false, err
			}

			//fmt.Printf("json: %+v\n", data)
			err = m.Ack()
			if err != nil {
				return false, err
			}
			protoBytes, err := base64.StdEncoding.DecodeString(data.DataBase64)
			if err != nil {
				return false, err
			}
			metrics := &metricspb.Metrics{}
			err = proto.Unmarshal(protoBytes, metrics)
			if err != nil {
				return false, err
			}

			localLoc, err := time.LoadLocation("Local")
			if err != nil {
				return false, err
			}
			mts := metrics.Ts.AsTime()
			localMts := mts.In(localLoc)
			fmt.Printf("got messsage with ts: %+v %v\n", metrics, localMts)

			//verifier(m.Subject, data.Data)
			if localMts.Sub(ts) > 0 {
				return true, nil
			}
		}
		return false, nil
	}
}
