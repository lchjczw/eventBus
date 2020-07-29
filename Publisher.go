package eventBus

import mapset "github.com/deckarep/golang-set"

// 异步发布
func (bus *eventBus) Publish(topic string, events ...interface{}) {
	bus.wg.Add(1)
	go bus.publish(topic, events...)
	return
}

// 同步发布
func (bus *eventBus) PublishSync(topic string, events ...interface{}) {
	bus.publishSync(topic, events...)
	return
}

func (bus *eventBus) publish(topic string, events ...interface{}) {
	bus.publishSync(topic, events...)
	bus.wg.Done()
	return
}

func (bus *eventBus) publishSync(topic string, events ...interface{}) {
	Topic := bus.getTopic(topic)
	bus.callAsync(topic, events, Topic.asyncHandlers)
	bus.callSync(topic, events, Topic)
}

// 执行同步订阅回调
func (bus *eventBus) callSync(topic string, events []interface{}, Topic *topic) {
	syncHandlers := make([]CallbackFunc, 0)
	Topic.RLock()
	if len(Topic.syncHandlers) > 0 {
		_ = clone(syncHandlers, Topic.syncHandlers)
	}
	Topic.RUnlock()
	for _, syncHandler := range syncHandlers {
		err := syncHandler(topic, events...)
		if err != nil {
			if Topic.transaction {
				return
			}
		}
	}
}

// 执行异步订阅回调
func (bus *eventBus) callAsync(topic string, events []interface{}, asyncHandlers mapset.Set) {
	for _, asyncHandler := range asyncHandlers.ToSlice() {
		callback, ok := asyncHandler.(CallbackFunc)
		if ok {
			go func() {
				_ = callback(topic, events...)
			}()
		}
	}
}
