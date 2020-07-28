package eventBus

type subscriber struct{}

// 同步订阅主题
func (sub *subscriber) Subscribe(topic string, callback func(topic string, events ...interface{}) error, err error) {
	return
}

// 同步订阅一次
func (sub *subscriber) SubscribeOnce(topic string, callback func(topic string, events ...interface{}) error, err error) {
	return
}

// 异步订阅主题
func (sub *subscriber) SubscribeAsync(topic string, callback func(topic string, events ...interface{}) error, err error) {
	return
}

// 异步订阅一次
func (sub *subscriber) SubscribeOnceAsync(topic string, callback func(topic string, events ...interface{}) error, err error) {
	return
}

// 取消订阅
func (sub *subscriber) UnSubscribe(topic string, callback func(topic string, events ...interface{}) error, err error) {
	return
}
