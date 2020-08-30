package eventBus

func (bus *eventBus) SetCycleBefore(topic string, callback CycleCallback) {
	Topic := bus.getTopic(topic)
	Topic.beforeCallback = callback
}

func (bus *eventBus) SetCycleAfterSync(topic string, callback CycleCallback) {
	Topic := bus.getTopic(topic)
	Topic.afterSyncCallback = callback
}

func (bus *eventBus) SetCycleAfterAll(topic string, callback CycleCallback) {
	Topic := bus.getTopic(topic)
	Topic.afterCallback = callback
}
