package eventBus

type controller struct{}

func (ctr *controller) HasSubscribed(topic string) bool {
	return hasSubscribed(topic)
}

func (ctr *controller) WaitAsync() {
	return
}

func hasSubscribed(topic string) bool {
	return false
}

func hasTopic(topic string) bool {
	_, ok := bus.topicMap.Load(topic)
	return ok
}

func getTopic(topic string) (handlers []eventHandler) {
	res, _ := bus.topicMap.LoadOrStore(topic, handlers)
	return res.([]eventHandler)
}
