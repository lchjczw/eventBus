package eventBus

import "sync"

type EventBus interface {
	Publisher
	Subscriber
	Controller
}

type Controller interface {
	HasSubscribed(topic string) bool
	WaitAsync()
}

type Publisher interface {
	Publish(topic string, events ...interface{})
	PublishSync(topic string, events ...interface{})
}

type Subscriber interface {
	// 同步订阅主题
	Subscribe(topic string, callback CallbackFunc, err error)
	// 同步订阅一次
	SubscribeOnce(topic string, callback CallbackFunc, err error)
	// 异步订阅主题
	SubscribeAsync(topic string, callback CallbackFunc, err error)
	// 异步订阅一次
	SubscribeOnceAsync(topic string, callback CallbackFunc, err error)
	// 取消订阅
	UnSubscribe(topic string, callback CallbackFunc, err error)
}

type eventBus struct {
	topicMap sync.Map
	wg       sync.WaitGroup
	EventBus
}

type topic struct {
	// todo 这里应该直接区分异步handlers和同步handlers
	//  异步handlers用set底层实现
	//  同步handlers用切片实现
	handlers []eventHandler
	lock     sync.RWMutex
}

type eventHandler struct {
	callback CallbackFunc
	flagOnce bool
	async    bool
}

var bus eventBus

func GetBus() EventBus {
	return &bus
}

func init() {
	bus = eventBus{
		topicMap: sync.Map{},
	}
}
