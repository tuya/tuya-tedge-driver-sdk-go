package common

import "fmt"

type UploadTokenInfo struct {
	UploadUrl string `json:"uploadUrl"`
	Token     string `json:"token"`
}

// UploadTokenResp 对应的消息回复
type UploadTokenResp struct {
	Result  UploadTokenInfo `json:"result"`
	T       int64           `json:"t"`
	Success bool            `json:"success"`
	Status  string          `json:"status"`
}

// UploadImageResp 上传图片对应的消息回复
type UploadImageResp struct {
	Success bool          `json:"success"`
	Msg     string        `json:"msg,omitempty"`
	T       int64         `json:"t"`
	Result  ImageFileInfo `json:"result"`
}

type ImageFileInfo struct {
	TmpFileId string `json:"tmp_file_id"`
}

type OssStruct struct {
	Result  []ResultNode `json:"result"`
	Success bool         `json:"success"`
	E       bool         `json:"e"`
	T       int64        `json:"t"`
}

type ResultNode struct {
	PathConfig struct {
		Log    string `json:"log"`
		Detect string `json:"detect"`
	} `json:"pathConfig"` // 网关devID
	AK         string `json:"ak"`
	SK         string `json:"sk"`
	Token      string `json:"token"`
	Bucket     string `json:"bucket"`
	LifeCycle  int64  `json:"lifeCycle"`
	Endpoint   string `json:"endpoint"`
	Provider   string `json:"provider"`
	Region     string `json:"region"`
	Expiration string `json:"expiration"`
	Id         string `json:"id"`
}

func (rn ResultNode) String() string {
	return fmt.Sprintf("id: %s, expiration: %s, provider: %s, endpoint: %s", rn.Id, rn.Expiration, rn.Provider, rn.Endpoint)
}

type UploadResopnse struct {
	Bucket     string `json:"bucket"`
	ObjectKey  string `json:"objectKey"`
	SecretKey  string `json:"secretKey"`
	ExpireTime int64  `json:"expireTime"`
}
