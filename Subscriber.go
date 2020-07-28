package eventBus

type subscriber struct{}

type CallbackFunc = func(topic string, events ...interface{}) error

// 同步订阅主题
func (sub *subscriber) Subscribe(topic string, callback CallbackFunc) {
	sub.subscribe(topic, callback, false, false)
}

// 同步订阅一次
func (sub *subscriber) SubscribeOnce(topic string, callback CallbackFunc) {
	sub.subscribe(topic, callback, true, false)
}

// 异步订阅主题
func (sub *subscriber) SubscribeAsync(topic string, callback CallbackFunc) {
	sub.subscribe(topic, callback, false, true)
}

// 异步订阅一次
func (sub *subscriber) SubscribeOnceAsync(topic string, callback CallbackFunc) {
	sub.subscribe(topic, callback, true, true)
}

// 取消订阅
func (sub *subscriber) UnSubscribe(topic string, callback CallbackFunc) {
	Topic := getTopic(topic)
	Topic.lock.Lock()

}

func (sub *subscriber) subscribe(topic string, callback CallbackFunc, once, async bool) {
	Topic := getTopic(topic)
	Topic.lock.Lock()
	Topic.handlers = append(Topic.handlers, eventHandler{
		callback: callback,
		flagOnce: once,
		async:    async,
	})
	Topic.lock.Unlock()
}
