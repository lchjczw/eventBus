package unitTest

import (
	"errors"
	"fmt"
	"gitee.com/super_step/eventBus"
	"testing"
)

const (
	testTopic = "testTopic"
)

var testBus eventBus.EventBus

func syncCallbackError(topic string, events ...interface{}) error {
	fmt.Println(fmt.Sprintf("sync# %s: %v", topic, events))
	return errors.New("test")
}

func syncCallback(topic string, events ...interface{}) error {
	fmt.Println(fmt.Sprintf("sync# %s: %v", topic, events))
	return nil
}

func syncCallback1(topic string, events ...interface{}) error {
	fmt.Println(fmt.Sprintf("sync# %s: %v", topic, events))
	return nil
}

func asyncCallback(topic string, events ...interface{}) error {
	fmt.Println(fmt.Sprintf("async# %s: %v", topic, events))
	return nil
}

func TestSub(t *testing.T) {
	err := testBus.Subscribe(testTopic, syncCallbackError)
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = testBus.Subscribe(testTopic, syncCallbackError)
	if err == nil {
		t.Error("重复订阅")
		return
	}

	err = testBus.Subscribe(testTopic, syncCallback)
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = testBus.Subscribe(testTopic, syncCallback)
	if err == nil {
		t.Error("重复订阅")
		return
	}

	err = testBus.Subscribe(testTopic, syncCallback1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = testBus.Subscribe(testTopic, syncCallback1)
	if err == nil {
		t.Error("重复订阅")
		return
	}
}

func TestSubAsync(t *testing.T) {
	err := testBus.SubscribeAsync(testTopic, asyncCallback)
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = testBus.SubscribeAsync(testTopic, asyncCallback)
	if err == nil {
		t.Error("重复订阅")
		return
	}
}

func TestPublish(t *testing.T) {
	testBus.SetTransaction(testTopic, false)
	testBus.Publish(testTopic, "not tr")
	testBus.SetTransaction(testTopic, true)
	testBus.Publish(testTopic, "tr")
}

func TestPublishSync(t *testing.T) {
	testBus.PublishSync(testTopic, "PublishSync")
}

func TestUnSubAsync(t *testing.T) {
	testBus.UnSubscribe(testTopic, asyncCallback)
	testBus.Publish(testTopic, "UnSubAsync")
}

func TestUnSubSync(t *testing.T) {
	testBus.UnSubscribe(testTopic, syncCallback1)
	testBus.Publish(testTopic, "UnSubSync", "callback1")
	testBus.WaitAsync()
	testBus.UnSubscribe(testTopic, syncCallback)
	testBus.Publish(testTopic, "UnSubSync", "callback")
	testBus.WaitAsync()
	testBus.UnSubscribe(testTopic, syncCallbackError)
	testBus.PublishSync(testTopic, "UnSubSync", "callbackError")
}

func TestWait(t *testing.T) {
	testBus.WaitAsync()
}

func TestMain(m *testing.M) {
	testBus = eventBus.NewBus()
	m.Run()
}
