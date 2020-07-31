# Go事件总线

#### 介绍
事件驱动golang event总线，适用于单进程代码快速解耦。

跨进程/集群(跨设备/容器)时应选择rabbit, rocket等开源MQ，而不是选择使用本项目。

#### 软件架构
使用`sync.map`与`github.com/deckarep/golang-set`底层完成

小范围读写锁确保发布时的性能

#### 使用说明

1. ###### 初始化实例

   ```go
   import "github.com/kataras/golog"
   
   EventBus := eventBus.NewBus(golog.Default)
   ```

2. ###### 实现callback接口

   ```go
   type subCallback struct {}
   
   var callback subCallback
   
   func (sub *subCallback) Callback(topic string, events ...interface{}) error {
   	config.Logger.Infof("event:%s%v", topic, events)
   	return nil
   }
   
   func initCallback() {
     callback = subCallback{}
   }
   
   ```

3. ###### 订阅主题

   ```go
   func sub() {
     // 同步订阅（收到事件后，按注册顺序同步执行callback）
     err := EventBus.Subscribe("topic:test", &callback)
     if err != nil {
       golog.Default.Error(err.Error())
     }
     // 异步订阅（收到事件后，为每个callback创建Goroutine，异步执行）
     err := EventBus.SubscribeAsync("topic:test", &callback)
     if err != nil {
       golog.Default.Error(err.Error())
     }
   }
   ```
   
4. ###### 发布主题

   ```go
   func pub() {
     // 异步发布，发布完成后不等待处理结果，立即返回
     EventBus.Publish("topic:test", "event:async")
     // 同步发布, 发布完成后等待同步与异步订阅全部完成后，返回同步订阅执行结果
     err := EventBus.PublishSync("topic:test", "event:sync")
     if err != nil {
       golog.Default.Error(err.Error())
     }
     // 同PublishSync，但不等待异步订阅完成
     err := EventBus.PublishSyncNoWait("topic:test", "event:syncNoWait")
     if err != nil {
       golog.Default.Error(err.Error())
     }
   }
   ```

5. ###### 取消订阅

   ```go
   func unSub() {
     // 取消单个主题的订阅
     EventBus.UnSubscribe("topic:test", &callback)
     // 取消全部订阅
     EventBus.UnSubscribeAll(&callback)
   }
   ```

6. ###### 关闭主题

   注意:  关闭主题只清空目标主题的订阅信息，再次发布与重新订阅该主题不受影响。

   ```go
   func closeTopic() {
     EventBus.CloseTopic("topic:test")
   }
   ```

#### demo

请参考`unitTest/pub_sub_test.go`

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request
