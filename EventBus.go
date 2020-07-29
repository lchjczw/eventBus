package eventBus

import (
	mapset "github.com/deckarep/golang-set"
	"sync"
)

type EventBus interface {
	Publisher
	Subscriber
	Controller
}

type Controller interface {
	HasSubscribed(topic string) bool
	WaitAsync()
	SetTransaction(topic string, tr bool) bool
}

type Publisher interface {
	Publish(topic string, events ...interface{})
	PublishSync(topic string, events ...interface{})
}

type Subscriber interface {
	// 同步订阅主题
	Subscribe(topic string, callback CallbackFunc) error
	// 同步订阅一次
	// SubscribeOnce(topic string, callback CallbackFunc, err error)
	// 异步订阅主题
	SubscribeAsync(topic string, callback CallbackFunc) error
	// 异步订阅一次
	// SubscribeOnceAsync(topic string, callback CallbackFunc, err error)
	// 取消订阅
	UnSubscribe(topic string, callback CallbackFunc)
}

type eventBus struct {
	topicMap sync.Map
	wg       sync.WaitGroup
	EventBus
}

type CallbackFunc = func(topic string, events ...interface{}) error

type topic struct {
	// 区分异步handlers和同步handlers
	//  异步handlers用set底层实现
	//  同步handlers用切片实现
	syncHandlers  []CallbackFunc
	asyncHandlers mapset.Set
	transaction   bool
	sync.RWMutex
}

func NewBus() EventBus {
	return &eventBus{
		topicMap: sync.Map{},
	}
}
