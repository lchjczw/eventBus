package eventBus

type publisher struct{}

func (pub *publisher) Publish(topic string, events ...interface{}) {
	return
}

func (pub *publisher) PublishSync(topic string, events ...interface{}) {
	return
}
