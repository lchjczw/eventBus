package eventBus

import (
	mapSet "github.com/deckarep/golang-set"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12/core/memstore"
	"sync"
)

type EventBus interface {
	Publisher
	Subscriber
	Controller
}

type Controller interface {
	// 等待异步执行完成
	WaitAsync()
	// 设置同步订阅事务标记
	SetTransaction(topic string, tr bool)
	// 设置hook
	SetHook(topic string, hook Hook)
	// 关闭主题
	CloseTopic(topic string)
}

type Publisher interface {
	// 发布
	Publish(topic string, events ...interface{})
	// 同步发布
	PublishSync(topic string, events ...interface{}) error
	// 同步发布, 不等待异步调用完成
	PublishSyncNoWait(topic string, events ...interface{}) error
}

type Subscriber interface {
	// 同步订阅主题
	Register(topic string, callback Callback) error
	// 异步订阅主题
	RegisterAsync(topic string, callback Callback) error
	// 取消已订阅的主题
	UnSubscribe(topic string, callback Callback)
	// 取消所有已订阅的主题
	UnSubscribeAll(callback Callback)
}

// Callback 因为会导致重复订阅,所以必须用interface的形式
// type Callback = func(string, context.Context, ...interface{}) error
type Callback interface {
	Callback(topic string, ctx *memstore.Store, events ...interface{}) error
}

type Hook interface {
	Before(topic string, ctx *memstore.Store, events ...interface{}) error
	AfterSync(topic string, ctx *memstore.Store, events ...interface{})
	After(topic string, ctx *memstore.Store, events ...interface{})
	Error(topic string, ctx *memstore.Store, err error, events ...interface{})
}

type eventBus struct {
	topicMap sync.Map
	wg       sync.WaitGroup
	logger   *golog.Logger
	EventBus
}

type CallbackFunc = func(topic string, events ...interface{}) error

type topic struct {
	// 区分异步handlers和同步handlers
	//  异步handlers用set底层实现
	//  同步handlers用切片实现
	syncHandlers  []Callback
	asyncHandlers mapSet.Set
	transaction   bool
	//Before    BeforeCallback
	//AfterSync CycleCallback
	//After     CycleCallback
	//Error   ErrorCallback
	hook Hook
	sync.RWMutex
}

func NewBus(logger *golog.Logger) EventBus {
	return &eventBus{
		topicMap: sync.Map{},
		logger:   logger,
	}
}
