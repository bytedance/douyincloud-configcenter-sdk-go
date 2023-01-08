# douyincloud-configcenter-sdk-go
本项目是抖音云配置中心的SDK插件，用以访问抖音云配置中心。

## 使用方法
### 初始化客户端
**方式一**：直接初始化
```go
// 默认初始化（默认轮询时间为60s，超时时间为5s）
sdkClient, err := base.Start()
if err != nil {
    panic(err)
}
```
**方式二**：设置轮询时间、超时时间的初始化
```go
config := base.NewClientConfig()
// 设置向配置中心发送请求的超时时间
config.SetTimeout(5 * time.Second)
// 设置轮询时间，最短为10s。设置参数小于10s则默认为10s
config.SetUpdateInterval(10 * time.Second)

sdkClient, err := base.StartWithConfig(config)
if err != nil {
    panic(err)
}
```
### 使用API调用配置中心
**根据key获取配置**
```go
sdkClient, err := base.Start()
value, err := sdkClient.Get("key_name")
```
**从云端获取配置刷新本地缓存**
```go
sdkClient, err := base.Start()
err := sdkClient.UpdateCache()
```
## 使用注意事项
- 使用前确保已在抖音云平台开启配置中心；
- 抖音云配置中心已通过专属的网络链路实现鉴权，您可以直接使用无需关心鉴权的具体逻辑；
- 由于抖音云配置中心有专属的鉴权逻辑，因此，您在开发调试时请使用抖音云的本地调试插件来连接抖音云配置中心。