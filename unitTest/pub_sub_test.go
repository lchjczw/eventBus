package unitTest

import (
	"errors"
	"fmt"
	"gitee.com/super_step/eventBus"
	"gitee.com/super_step/go_utils/logger"
	"gitee.com/super_step/go_utils/myError"
	"github.com/kataras/golog"
	"github.com/kataras/pio"
	"os"
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
		golog.Default.Infof("%s# %s: %v", callback.Name, topic, "Recursioned")
	} else {
		golog.Default.Infof("%s# %s: %v", callback.Name, topic, events)
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
		if err != nil {
			err = myError.Warp(err, "%d#", index)
			golog.Default.Error(err.Error())
		} else {
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
	err := testBus.PublishSync(syncTopic, "PublishSync")
	if err != nil {
		golog.Default.Error(err.Error())
	}
	err = testBus.PublishSyncNoWait(syncTopic, "PublishSync")
	if err != nil {
		golog.Default.Error(err.Error())
	}
	err = testBus.PublishSync(asyncTopic, "PublishSync")
	if err != nil {
		golog.Default.Error(err.Error())
	}
	err = testBus.PublishSyncNoWait(asyncTopic, "PublishSync")
	if err != nil {
		golog.Default.Error(err.Error())
	}
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
	golog.Default.Printer = pio.NewPrinter("", os.Stdout).
		EnableDirectOutput().Hijack(logger.Hijacker).SetSync(false)
	golog.Default.SetTimeFormat("2006-01-02 15:04:05")
	golog.Default.SetLevel("debug")
	testBus = eventBus.NewBus(golog.Default)
	for i := 0; i < 4; i++ {
		callbackStruck := newCallback(fmt.Sprintf("%d", i), i == 2)
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
