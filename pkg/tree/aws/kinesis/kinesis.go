package kinesis

import "strings"

type Kinesis struct {
	Streams         []Stream                 `tree:"streams"`
	DeliveryStreams []FirehoseDeliveryStream `tree:"delivery_streams"`
}

func (k *Kinesis) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range k.DeliveryStreams {
		k.DeliveryStreams[i].Relationships.KinesisStream = nil
	}

	// link firehose delivery streams to kinesis streams
	for i, deliveryStream := range k.DeliveryStreams {
		if deliveryStream.KinesisStreamARN.IsEmpty() {
			continue
		}
		for j := range k.Streams {
			if strings.Contains(deliveryStream.KinesisStreamARN.Value(), k.Streams[j].ID) {
				k.DeliveryStreams[i].Relationships.KinesisStream = &k.Streams[j]
				break
			}
		}
	}
}
