package sns

type SNS struct {
	Topics        []Topic        `tree:"topics"`
	Subscriptions []Subscription `tree:"subscriptions"`
}

func (s *SNS) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range s.Topics {
		s.Topics[i].Relationships.Subscriptions = nil
	}

	// link subscriptions to topics
	for i, sub := range s.Subscriptions {
		for j := range s.Topics {
			if sub.TopicARN.Equal(s.Topics[j].ID) {
				s.Topics[j].Relationships.Subscriptions = append(s.Topics[j].Relationships.Subscriptions, &s.Subscriptions[i])
				break
			}
		}
	}
}
