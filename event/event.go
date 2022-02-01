package event

import (
	"errors"
	"fmt"
	"gitee.com/super_step/eventBus"
	"github.com/kataras/golog"
	"path"
	"strings"
)

var globalBus eventBus.EventBus
var root *Event

func Root() *Event {
	return root
}

func init() {

	logger := golog.New()
	globalBus = eventBus.NewBus(logger)
	root = &Event{
		Topic: "/",
		Desc:  "root节点",
		bus:   globalBus,
	}
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
	return path.Clean(key)
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
func (a *Event) Publish(args ...interface{}) {
	a.bus.Publish(a.Key(), args...)
}
func (a *Event) PublishSync(args ...interface{}) error {
	return a.bus.PublishSync(a.Key(), args...)
}
func (a *Event) PublishSyncNoWait(args ...interface{}) error {
	return a.bus.PublishSyncNoWait(a.Key(), args...)
}
// Register 注册
func (a *Event) Register(callback eventBus.Callback) error {
	a.callBack = callback
	return a.bus.Register(a.Key(), callback)
}

func (a *Event) RegisterAsync(callback eventBus.Callback) error {
	a.callBack = callback
	return a.bus.RegisterAsync(a.Key(), callback)
}

type EventInter interface {
	Event(path, desc string) *Event
	Key() string
	Register(callback eventBus.Callback) error
	RegisterAsync(callback eventBus.Callback) error
	Publish(args ...interface{})
	PublishSync(args ...interface{}) error
	PublishSyncNoWait(args ...interface{}) error
}

func Init() {
	sys := root.Event(`system`, `系统级事件`)
	fmt.Println(sys)
	sale := root.Event(`sale`, `销售相关事件`)
	fmt.Println(sale)
}

func NewRootEvent(bus eventBus.EventBus, topic, desc string) *Event {
	if bus == nil {
		bus = globalBus
	}

	e := &Event{
		Topic:    topic,
		Desc:     desc,
		callBack: nil,
		children: nil,
		bus:      bus,
	}
	return e
}
