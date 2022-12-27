package commons

import "github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Payload []byte `json:"payload"`
}

type AppDriverReq struct {
	Header  *Header
	Payload []byte
}

type Header struct {
	Tag    string
	From   string
	Option map[string]string
}

func (req *AppDriverReq) ToRpc() *proto.Data {
	return &proto.Data{
		Name: "",
		Header: &proto.Header{
			Tag:    req.Header.Tag,
			From:   req.Header.From,
			Option: req.Header.Option,
		},
		Payload: req.Payload,
	}
}


//////////////////////////////////////////////////////////////////////////////////////
type ProxyInfo struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func DefaultProxyInfo(port int) ProxyInfo {
	return ProxyInfo{
		Host: "127.0.0.1",
		Port: port,
	}
}
