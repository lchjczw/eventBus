package event

import (
	"errors"
	"gitee.com/super_step/eventBus"
	"github.com/kataras/golog"
	"strings"
)

var globalBus = eventBus.NewBus(golog.Default)
var root = NewRootEvent(globalBus)

func Root() *Event {
	return root
}

type Event struct {
	Topic    string
	Desc     string
	callBack eventBus.Callback

	prent    *Event
	children []*Event
	bus      eventBus.EventBus
}

func (a *Event) Check(topic string) error {
	if strings.ContainsAny(topic, "\\/") {
		return errors.New("topic不能包含\\/字符")
	}

	if a.IsExist(topic) {
		return errors.New("topic已经存在")
	}

	return nil
}

func (a *Event) IsExist(topic string) bool {
	for i := range a.children {
		if topic == a.children[i].Topic {
			return true
		}
	}
	return false
}

func (a *Event) Key() string {
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
func (a *Event) Event(topic, desc string) *Event {
	if err := a.Check(topic); err != nil {
		panic(err)
	}

	t := &Event{
		Topic:    topic,
		Desc:     desc,
		prent:    a,
		children: nil,
		bus:      a.bus,
	}

	a.children = append(a.children, t)
	return t
}
func (a *Event) PublishAsync(args ...interface{}) {
	a.bus.PublishAsync(a.Key(), args...)
}
func (a *Event) PublishSync(args ...interface{}) error {
	return a.bus.PublishSync(a.Key(), args...)
}
func (a *Event) PublishSyncNoWait(args ...interface{}) error {
	return a.bus.PublishSyncNoWait(a.Key(), args...)
}

// SubscribeSync 注册
// 给callback同时实现hook接口，则直接注入hook
func (a *Event) SubscribeSync(callback eventBus.Callback) error {
	a.callBack = callback
	err := a.bus.SubscribeSync(a.Key(), callback)
	if err != nil {
		return err
	}

	hook, ok := callback.(eventBus.Hook)
	if ok {
		a.bus.SetHook(a.Key(), hook)
	}

	return nil
}

func (a *Event) SubscribeAsync(callback eventBus.Callback) error {
	a.callBack = callback
	return a.bus.SubscribeAsync(a.Key(), callback)
}

type EventInter interface {
	Event(path, desc string) *Event
	Key() string
	SubscribeSync(callback eventBus.Callback) error
	SubscribeAsync(callback eventBus.Callback) error
	PublishAsync(args ...interface{})
	PublishSync(args ...interface{}) error
	PublishSyncNoWait(args ...interface{}) error
}

func NewRootEvent(bus eventBus.EventBus) *Event {
	if bus == nil {
		bus = globalBus
	}

	e := &Event{
		Topic:    "/",
		Desc:     "root节点",
		bus:      bus,
	}
	return e
}
