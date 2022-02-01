package eventBus

import (
	"gitee.com/super_step/eventBus/pkg/memstore"
	mapSet "github.com/deckarep/golang-set"
	"github.com/kataras/golog"
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
	// 异步发布
	PublishAsync(topic string, args ...interface{})
	// 同步发布
	PublishSync(topic string, args ...interface{}) error
	// 同步发布, 不等待异步调用完成
	PublishSyncNoWait(topic string, args ...interface{}) error
}

type Subscriber interface {
	// 同步订阅主题
	SubscribeSync(topic string, handler Handler) error
	// 异步订阅主题
	SubscribeAsync(topic string, handler Handler) error
	// 取消已订阅的主题
	UnSubscribe(topic string, handler Handler)
	// 取消所有已订阅的主题
	UnSubscribeAll(handler Handler)
}

// Handler 因为会导致重复订阅,所以必须用interface的形式
// type Handler = func(string, context.Context, ...interface{}) error
type Handler interface {
	Handler(topic string, ctx *memstore.Store, args ...interface{}) error
}

type Hook interface {
	// 发布前回调
	Before(topic string, ctx *memstore.Store, args ...interface{}) error
	// 同步完成时回调
	AfterSync(topic string, ctx *memstore.Store, args ...interface{})
	// 全部完成时回调
	After(topic string, ctx *memstore.Store, args ...interface{})
	// 错误时回调
	Error(topic string, ctx *memstore.Store, err error, args ...interface{})
}

type eventBus struct {
	topicMap sync.Map
	wg       sync.WaitGroup
	logger   *golog.Logger
	EventBus
}

type HandlerFunc = func(topic string, args ...interface{}) error

type topic struct {
	// 区分异步handlers和同步handlers
	//  异步handlers用set底层实现
	//  同步handlers用切片实现
	syncHandlers  []Handler
	asyncHandlers mapSet.Set
	transaction   bool

	hook Hook
	sync.RWMutex
}

func NewBus(logger *golog.Logger) EventBus {
	return &eventBus{
		topicMap: sync.Map{},
		logger:   logger,
	}
}
