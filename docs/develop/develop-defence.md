# 周界设备接入DP点数据格式定义
在此之前，请务必阅读[DP模型驱动开发指南](./developdp.md)、以及[设备接入补充说明文档](./develop-addtion.md)
- 如果云端下发的 DP 指令中，带有 sn 字段，则驱动程序在处理完指令后，要向云端上报指令执行成功或失败
- 上报指令执行成功，调用Tedge API：`CmdRespSuccess`，携带 sn 字段
- 上报指令执行失败，调用Tedge API：`CmdRespFail`，携带 sn 字段

## 周界设备 DP 点消息格式定义
### DPId:2
- 云端下发消音指令：DP Id 2：(注，所有云端下发的消息，都经由 handleCommands 接口），其将 CommandRequest 中的 data 字段反序列化后，会得到如下数据结构，根据实际情况，有的字段可能没有)
```
//SilenceBellDp 消音DP，执行完成必须向云端上报指令执行成功或失败
type SilenceBellDp struct {
	Sn uint64 `json:"sn"`
}
```

### DPId:3
- 云端下发子系统布撤防指令：DP Id 3：
```
//IssueDp 子系统布撤防，执行完成必须向云端上报指令执行成功或失败
type IssueDp struct {
	Status uint8 `json:"status"`
	Sn     int64 `json:"sn"`
}
```

### DPId:5
- 云端下发旁路指令：DP Id 5
```
//ByPassDP 旁路，执行完成必须向云端上报指令执行成功或失败
type ByPassDP struct {
	Status uint8    `json:"status"`
	Points []string `json:"points"`
	Sn     int64    `json:"sn"`
}
```

### DPId:6
- 设备同步指令
```
//SyncHostDp 同步消息，执行完成必须向云端上报指令执行成功或失败
type SyncHostDp struct {
	Sn int64 `json:"sn"`
}
```

### DPId:1(主机设备)
- 向云端上报主机状态变化
```
#status = 0x09 //布防
#status = 0x04 //撤防
{
    "status": int64 `json:"status"`
}
```

### DPId:1(防区设备)
- 向云端上报防区状态变化
```
#status = 5 //旁路开启
#status = 3 //旁路关闭
#status = 4 //告警
{
    "contactId": string `json:"contactId"`
    "status": int64 `json:"contactId"`
    "eventId": string `json:"eventId"`
}
```
