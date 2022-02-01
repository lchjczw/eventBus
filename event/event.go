package event

import (
	"errors"
	"gitee.com/super_step/eventBus"
	"github.com/kataras/golog"
	"strings"
	"sync"
)

var globalBus = eventBus.NewBus(golog.Default)
var root = NewRootEvent(globalBus)

func Root() Event {
	return root
}

type event struct {
	Topic    string
	Desc     string
	bus     eventBus.EventBus
	handler eventBus.Handler

	prent    *event
	children []*event
	sync.RWMutex
}

func (a *event) Check(topic string) error {
	if strings.ContainsAny(topic, "\\/") {
		return errors.New("topic不能包含\\/字符")
	}

	if a.IsExist(topic) {
		return errors.New("topic已经存在")
	}

	return nil
}

func (a *event) IsExist(topic string) bool {
	for i := range a.children {
		if topic == a.children[i].Topic {
			return true
		}
	}
	return false
}

func (a *event) Key() string {
	// 遍历到root,得到前缀
	var key string
	t := a
	for true {
		if t.Topic == "/" || t.Topic == "" {
			break
		}
		key = t.Topic + "/" + key
		if t.prent == nil {
			break
		}
		t = t.prent
	}
	return key
}
func (a *event) Event(topic, desc string) *event {
	if err := a.Check(topic); err != nil {
		panic(err)
	}

	t := &event{
		Topic:    topic,
		Desc:     desc,
		prent:    a,
		children: nil,
		bus:      a.bus,
	}
	a.Lock()
	a.children = append(a.children, t)
	a.Unlock()
	return t
}
func (a *event) PublishAsync(args ...interface{}) {
	a.bus.PublishAsync(a.Key(), args...)
}
func (a *event) PublishSync(args ...interface{}) error {
	return a.bus.PublishSync(a.Key(), args...)
}
func (a *event) PublishSyncNoWait(args ...interface{}) error {
	return a.bus.PublishSyncNoWait(a.Key(), args...)
}

// SubscribeSync 注册
// 给handler同时实现hook接口，则直接注入hook
func (a *event) SubscribeSync(handler eventBus.Handler) error {
	a.handler = handler
	err := a.bus.SubscribeSync(a.Key(), handler)
	if err != nil {
		return err
	}

	hook, ok := handler.(eventBus.Hook)
	if ok {
		a.bus.SetHook(a.Key(), hook)
	}

	return nil
}

func (a *event) SubscribeAsync(handler eventBus.Handler) error {
	a.handler = handler
	return a.bus.SubscribeAsync(a.Key(), handler)
}

type Event interface {
	Event(path, desc string) *event
	Key() string
	SubscribeSync(handler eventBus.Handler) error
	SubscribeAsync(handler eventBus.Handler) error
	PublishAsync(args ...interface{})
	PublishSync(args ...interface{}) error
	PublishSyncNoWait(args ...interface{}) error
}

func NewRootEvent(bus eventBus.EventBus) Event {
	if bus == nil {
		bus = globalBus
	}

	e := &event{
		Topic: "/",
		Desc:  "root节点",
		bus:   bus,
	}
	return e
}
