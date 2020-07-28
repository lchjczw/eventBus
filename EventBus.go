package eventBus

import "sync"

import "github.com/deckarep/golang-set"

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
	Subscribe(topic string, callback func(topic string, events ...interface{}) error, err error)
	// 同步订阅一次
	SubscribeOnce(topic string, callback func(topic string, events ...interface{}) error, err error)
	// 异步订阅主题
	SubscribeAsync(topic string, callback func(topic string, events ...interface{}) error, err error)
	// 异步订阅一次
	SubscribeOnceAsync(topic string, callback func(topic string, events ...interface{}) error, err error)
	// 取消订阅
	UnSubscribe(topic string, callback func(topic string, events ...interface{}) error, err error)
}

type eventBus struct {
	topicMap sync.Map
	wg       sync.WaitGroup
	EventBus
}

type eventHandler struct {
	callback func(topic string, events ...interface{}) error
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
