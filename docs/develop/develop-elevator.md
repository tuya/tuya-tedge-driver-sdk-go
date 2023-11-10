# 梯控设备接入DP点数据格式定义
- 在此之前，请务必阅读[DP模型驱动开发指南](./developdp.md)、以及[设备接入补充说明文档](./develop-addtion.md)
- 如果云端下发的 DP 指令中，带有 sn 字段，则驱动程序在处理完指令后，要向云端上报指令执行成功或失败
- 上报指令执行成功，调用Tedge API：`CmdRespSuccess`，携带 sn 字段
- 上报指令执行失败，调用Tedge API：`CmdRespFail`，携带 sn 字段

## 梯控设备 DP 点消息格式定义
### DPId:12
- 通行事件上报(驱动向云端上报)：上报数据类型为字符串类型，将如下结构体序列化成 json 字符串后上报；
```
//说明：注解 `json:"leftNum"` 表示 json 中的 key 为 leftNum；
type PassRecord struct {
    LeftNum    int    `json:"leftNum"`              //非必填, 电梯井号
    Uid        string `json:"uid"`                  //必填, 涂鸦人员id，不超过32个字符
    StartFloor int    `json:"startFloor,omitempty"` //非必填, 开始楼层
    EndFloor   int    `json:"endFloor,omitempty"`   //非必填, 结束楼层
    T          int64  `json:"t"`                    //必填, 通行时间 时间戳、毫秒
    Way        int    `json:"way"`                  //必填，通行方式 1(卡)/2(密码)/3(⼆维码)/4(⼈脸)/5(指纹)/9(远程)/10(其它)
    ImageId    string `json:"imageId,omitempty"`    //非必填, 通行图片id
    Temp       string `json:"temp,omitempty"`       //非必填, 温度
    IsSuccess  uint   `json:"isSuccess"`            //必填, 0 失败 1 成功
    Message    string `json:"message,omitempty"`    //非必填, 消息提示
}
```

### DPId:2
- 云端下发呼梯消息：DP Id 2：(注，所有云端下发的消息，都经由 handleCommands 接口），其将 CommandRequest 中的 data 字段反序列化后，会得到如下数据结构，根据实际情况，有的字段可能没有)
- 执行完成必须向云端上报指令执行成功或失败
```
type elevatorControlDp struct {
    Type       string `json:"type"`       //必填，类型,in:内呼,out:外呼,unify:内外呼结合
    StartFloor int    `json:"startFloor"` //非必填，起始楼层
    EndFloor   int    `json:"endFloor"`   //非必填，起始楼层
    Sn         int64  `json:"sn"`         //消息序列号
    Direction  uint   `json:"direction,omitempty"` ////非必填，方向
}

CommandRequest 示例：
{
    "protocol":5,
    "t":1696904311,
    "s":1696904311,
    "data":{
        "2":"{\"type\":\"add\",\"startFloor\":1,\"endFloor\":10,\"sn\":1696904311}"
    }
}
```

### DPId:3
- 云端下发设置通行消息(人脸信息)：DP Id 3：新增、更新、删除
- 执行完成必须向云端上报指令执行成功或失败
```
type faceDp struct {
    Data faceDpData `json:"data"`
    Mode string     `json:"mode"` //ec_faceInfo
    Sn   int64      `json:"sn"`   //消息序列号
    Type string     `json:"type"` //add 新增人脸；update 更新；del 删除
}

type faceDpData struct {
    FaceId string `json:"faceId"` 			//人脸id
    Uid    string `json:"uid"`    			//涂鸦用户id
    Url    string `json:"url,omitempty"`//照片 url
}
```

### DPId:4
- 云端下发设置用户消息(人员信息)：DP Id 4：新增、更新、删除；
- 说明：当 Type 为 add/update 时，CommandRequest 中的 data 会同时存在 DP Id 5，表示该人员拥有权限的楼层列表；
- 执行完成必须向云端上报指令执行成功或失败
```
type personDp struct {
    Data personDpData `json:"data"`
    Mode string       `json:"mode"`  //ec_userInfo
    Sn   int64        `json:"sn"`
    Type string       `json:"type"`  //add 新增人员；update 更新；del 删除
}

type personDpData struct {
    BeginTime   uint64 `json:"beginTime,omitempty"` //非必填，有效期开始时间，13位时间戳
    EndTime     uint64 `json:"endTime,omitempty"`   //非必填，有效期结束时间，13位时间戳
    IdCard      string `json:"idCard,omitempty"`    //非必填，证件id，不超过32个字符
    Name        string `json:"name,omitempty"`      //用户名，不超过64个字符
    Phone       string `json:"phone,omitempty"`     //非必填，手机号，11位数字
    Uid         string `json:"uid"`                 //涂鸦人员id，不超过32个字符
    SecretKey   string `json:"secretKey,omitempty"` 
    RefreshTime uint64 `json:"refreshTime,omitempty"`
}
```

### DPId:5
- 云端下发设置楼层权限组消息：DP Id 5：(参考 DP Id 4)
```
type floorAuthDp struct {
    FloorAuthList []int `json:"floorAuthList,omitempty"`
}
```

### DPId:6
- 向云端上报电梯运行状态消息：DP Id 6
```
type statusReport struct {
    LeftNum   int    `json:"leftNum"`    //必填, 电梯井号
    CurFloor  int    `json:"curFloor"`   //非必填, 当前楼层,物理楼层,从1开始
    Status    int     `json:"status"`    //非必填, 电梯状态, 0(停止)/1(上行中)/2(下行中)/3(故障)
    Direction int    `json:"direction"`  //非必填, 电梯运行方向, 0(上行)/1下行
}
```


### DPId:13
- 云端下发二维码消息: DP Id 13
- 执行完成必须向云端上报指令执行成功或失败ß
```
type qrcodeDp struct {
    Sn   int64  `json:"sn"`   //序列号
    Type string `json:"type"` //操作类型：add/update/del，分别表示二新增二维码、更新、和删除
    Data struct {
        UID    string `json:"uid"`    //涂鸦用户 id
        Qrcode string `json:"qrcode"` //二维码
    } `json:"data"`
}
```
