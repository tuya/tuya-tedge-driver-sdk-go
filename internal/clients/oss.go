package clients

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/utils"
)

type OSSClientI interface {
	UploadFile(ctx context.Context, filename string, content []byte) (string, string, error)
}

func NewOSSClient(cfg *common.ResultNode, logger commons.TedgeLogger) (OSSClientI, error) {
	switch cfg.Provider {
	case "oss":
		return newCOSClient(cfg, logger)
	default:
		return nil, errors.New("unsupported oss provider")
	}
}

type COSClient struct {
	*cos.Client
	cfg    *common.ResultNode
	logger commons.TedgeLogger
}

func newCOSClient(cfg *common.ResultNode, logger commons.TedgeLogger) (*COSClient, error) {
	rawurl := fmt.Sprintf("https://%s.%s", cfg.Bucket, cfg.Endpoint) //COS不需要region，试出来的
	u, _ := url.Parse(rawurl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     cfg.AK,
			SecretKey:    cfg.SK,
			SessionToken: cfg.Token, //使用临时密钥
		},
	})
	return &COSClient{
		Client: client,
		cfg:    cfg,
		logger: logger,
	}, nil
}

// 返回:{"bucket":"ty-cn-storage30","objectKey":"/88012e-34125598-194c9d2a7048680e/v","secretKey":"c86440d1912a413d8d90d13097391159","expireTime":1551083405}
func (c *COSClient) UploadFile(ctx context.Context, filename string, content []byte) (string, string, error) {
	crypted, key, err := prepareContent(content)
	if err != nil {
		c.logger.Errorf("prepare file content error: %s", err)
		return "", "", err
	}

	opt := &cos.ObjectPutOptions{}
	name := c.cfg.PathConfig.Detect[1:] + "/" + filename
	if _, err = c.Object.Put(ctx, name, bytes.NewReader(crypted), opt); err != nil {
		c.logger.Errorf("upload file error: %s", err)
		return "", "", err
	}
	imageInfo, err := json.Marshal(common.UploadResopnse{
		Bucket:     c.cfg.Bucket,
		ObjectKey:  name,
		SecretKey:  key,
		ExpireTime: time.Now().Unix() + c.cfg.LifeCycle,
	})
	fileId := strconv.FormatInt(time.Now().UnixNano()/1e6, 10) + "_" + filename
	return fileId, string(imageInfo), err
}

// 返回值：加密后的内容，aes加密key（hex encode to string），error
func prepareContent(content []byte) ([]byte, string, error) {
	var (
		err     error
		crypted []byte
		version [4]byte
		reserve [40]byte
		length  = make([]byte, 4)
	)

	uid, _ := uuid.New().MarshalBinary()
	str := hex.EncodeToString(uid)
	key := []byte(str[:16])
	iv := []byte(str[16:])
	if crypted, err = utils.AESCBCEncrypt(content, key, iv); err != nil {
		return []byte{}, "", err
	}

	var buff bytes.Buffer
	version[3] = '1'
	buff.Write(version[:])
	buff.Write(iv)
	binary.LittleEndian.PutUint32(length, uint32(len(content)))
	buff.Write(length)
	buff.Write(reserve[:])
	buff.Write(crypted)
	return buff.Bytes(), string(key), nil
}
