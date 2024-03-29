# 设备接入指南补充说明

## 一、复习一下驱动开发流程
- 0.开发流程已经在[DP模型驱动开发指南](./developdp.md)介绍过
- 1.创建驱动服务，并实现SDK中定义的驱动接口 `DPModelDriver`(使用 sdk api)
- 2.实现 `驱动程序的职责` 中列出的 4 大类基本功能
    - 2.1.发现真实的子设备，并将真实设备上报给Tedge：激活子设备(使用 sdk api)
    - 2.2.定时获取子设备的状态，并上报子设备状态：在线或离线(使用 sdk api)
    - 2.3.根据实际情况，获取子设备各项功能数据，并上报到云端：DP点值(使用 sdk api)
    - 2.4.接收云端的设备指令，解析指令并将指令转化成设备可接收的格式，最后将指令发送到真实子设备；有些指令可能需要向云端上报执行结果，具体依赖所接入的设备类型；
    - 说明：如上所列 4 项属于驱动程序需要实现的核心功能，其中 1、2 比较简单；3、4 是驱动程序的主要代码量； 
- 3.将驱动打包成一个 docker 镜像，并推送到镜像仓库
- 4.在Tedge Web通过自定义方式，将开发好的驱动程序安装到Tedge中

## 二、一些细节方面的问题
- 1.设备发现
    - 驱动程序要实现的第一项核心功能就是要连接第三方系统，并通过第三方系统提供的 SDK 或 API 获取子设备列表，然后将子设备激活
    - 第三方系统的连接信息，一般配置在"专家模式配置"中，驱动读取配置后连接第三方系统
    - 激活子设备时，参数`Cid`、`ProductId`、`Name` 这三个参数是必填的
    - `Cid`: 整个网关下必须唯一，一般尽量使用第三方 API 返回的设备唯一编号；`ProductId`：由对接的涂鸦同学提供；`Name`: 若无法获取设备名，可设置成 Cid；
    - 驱动程序每次重新启动后，首先要获取子设备列表，然后通过 Tedge API `GetActiveDevices` 或 `GetDeviceById` 判断该子设备是否已经激活过；若未激活过，则当然新设备激活；否则直接在驱动程序中创建设备索引即可；
- 2.设备状态
    - 驱动必须定时上报子设备状态，通过 Tedge API `ReportDeviceStatus`
    - 若第三方系统或设备提供回调的形式通知设备在线状态变更，那么驱动程序应当监听回调，并在收到设备状态变更时立即上报设备状态
    - 若第三方系统或设备提供查询设备在线状态的功能，则驱动程序可以定时查询并上报，一般每隔一分钟查询一次比较合适
    - 若第三方系统无任何形式获取设备在线状态，则只要能连接成功，则驱动程序可将该设备一直置为在线状态
- 3.专家模式配置，格式为 yaml 格式，在驱动中通过 Tedge API `GetCustomConfig()` 获取配置
    ![新增镜像仓库](../images/专家配置.png)
- 4.驱动处理云端下发的DP指令
    - 云端下发到驱动程序的指令，全部都通过接口 `HandleCommands` 下发，完整的消息参考参数 `CommandRequest`
    - 特别注意事项：驱动程序在收到 HandleCommands 消息后，应当将消息发送异步队列中进行处理，不能在 HandleCommands 中进行阻塞处理！！！
    - CommandRequest 中 Protocol 固定为 5，T 为时间戳，S 暂时无用；
    - 开发同学需要关心的是 CommandRequest 中的 Data，参考示例如下：其中 dps 中的 key 就是 dpId，value 就是 dpValue；
```json
{
    "Protocol":5,
    "T":1672036365,
    "Data":{
        "cid":"test_cid_0001",
        "dps":{
            "101": 100,
            "102": false,
            "103": "xxxx"
        }
    }
}
```
- 5.驱动向云端上报设备DP消息
    - 开发过程中，必须先将 productId 填加到 Tedge 中，并且点击同步产品信息(即 DP 列表)，具体操作参考[边缘网关操作介绍](./tedgeweb.md)`## 四、产品和子设备`
    - 在驱动程序中使用 Tedge API `ReportWithDPData` 向涂鸦云上报设备DP消息，注意虽然该接口参数为列表，但是目前只支持一个元素；
    - 参数 DPId 必须在 productId 中的 DP 列表中，否则边缘网关会拒绝该消息
    - 具体示例：请参考驱动程序 Demo
    - 常用基本类型：value(即int64)、 bool、enum、string
    - 字符串类型，某些 DP 点数据实际是json字符串，这类消息驱动程序要将消息序列化成 json 字符串再上报
    -  DP 消息上报时机：
        - 如果第三方系统或设备提供的 API 主动推送设备信息，则驱动程序可以收到相应信息后向涂鸦云端上报
        - 如果第三方系统或设备提供提供回调机制，则同上
        - 否则，驱动程序应当定时去查询设备运行状态，并将查到的信息上报到涂鸦云
