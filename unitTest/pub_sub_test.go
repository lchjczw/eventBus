package unitTest

import (
	"errors"
	"fmt"
	"gitee.com/super_step/eventBus"
	"testing"
)

const (
	syncTopic  = "syncTopic"
	asyncTopic = "asyncTopic"
)

var (
	testBus   eventBus.EventBus
	callbacks []eventBus.Callback
)

type callback struct {
	Name      string
	Err       bool
	Recursion bool
}

func (callback *callback) Callback(topic string, events ...interface{}) error {
	if len(events) == 0 {
		fmt.Println(fmt.Sprintf("%s# %s: %v", callback.Name, topic, "Recursioned"))
	} else {
		fmt.Println(fmt.Sprintf("%s# %s: %v", callback.Name, topic, events))
		if callback.Recursion {
			testBus.Publish(topic)
		}
	}
	if callback.Err {
		return errors.New(callback.Name)
	}
	return nil
}

func TestSub(t *testing.T) {
	for index, callback := range callbacks {
		err := testBus.Subscribe(syncTopic, callback)
		if err != nil {
			t.Error(index, err.Error())
			return
		}
		err = testBus.Subscribe(syncTopic, callback)
		if err == nil {
			t.Error("重复订阅未报告异常")
			return
		}
	}
}

func TestSubAsync(t *testing.T) {
	for index, callback := range callbacks {
		err := testBus.SubscribeAsync(asyncTopic, callback)
		if err != nil {
			t.Error(index, err.Error())
			return
		}
		err = testBus.SubscribeAsync(asyncTopic, callback)
		if err == nil {
			t.Error("重复订阅未报告异常")
			return
		}
	}
}

func TestPublish(t *testing.T) {
	testBus.SetTransaction(syncTopic, true)
	testBus.Publish(syncTopic, "publish test", "transaction")
	testBus.WaitAsync()
	testBus.SetTransaction(syncTopic, false)
	testBus.Publish(syncTopic, "publish test")
	testBus.Publish(asyncTopic, "publish test")
	testBus.WaitAsync()
}

func TestPublishSync(t *testing.T) {
	testBus.PublishSync(syncTopic, "PublishSync")
	testBus.PublishSync(asyncTopic, "PublishSync")
}

func TestUnSubAsync(t *testing.T) {
	for index, callback := range callbacks {
		testBus.Publish(syncTopic, "Sync", "UnSub", index)
		testBus.WaitAsync()
		testBus.UnSubscribe(syncTopic, callback)
	}
}

func TestUnSubSync(t *testing.T) {
	for index, callback := range callbacks {
		testBus.Publish(asyncTopic, "Async", "UnSub", index)
		testBus.WaitAsync()
		testBus.UnSubscribe(asyncTopic, callback)
	}
}

func TestUnSubAll(t *testing.T) {
	TestSub(t)
	TestSubAsync(t)
	length := len(callbacks)
	for i := 1; i <= length; i++ {
		testBus.Publish(asyncTopic, "Async", "UnSubAll", i)
		testBus.Publish(syncTopic, "Sync", "UnSubAll", i)
		testBus.WaitAsync()
		testBus.UnSubscribeAll(callbacks[length-i])
	}

	for index, callback := range callbacks {
		testBus.Publish(asyncTopic, "Async", "UnSubAll", index)
		testBus.Publish(syncTopic, "Sync", "UnSubAll", index)
		testBus.WaitAsync()
		testBus.UnSubscribeAll(callback)
	}
}

func TestRecursionSub(t *testing.T) {
	callbackRecursion := callback{
		Name:      "Recursion",
		Err:       false,
		Recursion: true,
	}
	err := testBus.Subscribe(syncTopic, &callbackRecursion)
	if err != nil {
		t.Error("Recursion", err.Error())
		return
	}
	testBus.Publish(syncTopic, "Recursion")
	testBus.WaitAsync()
}

func TestMain(m *testing.M) {
	testBus = eventBus.NewBus()
	for i := 0; i < 4; i++ {
		callbackStruck := newCallback(fmt.Sprintf("%d#", i), i == 2)
		callbacks = append(callbacks, callbackStruck)
	}
	m.Run()
	testBus.WaitAsync()
}

func newCallback(name string, err bool) eventBus.Callback {
	callbackStruck := callback{
		Name: name,
		Err:  err,
	}
	return &callbackStruck
}
