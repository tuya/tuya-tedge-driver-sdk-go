package driversvc

import (
	"context"
	"errors"
	"time"

	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"google.golang.org/grpc/status"
)

/////////////////////////////////////////////////////////////////////////////////////////////
// 新的KV方法，获取到的KV存储支持云端备份
func (bds *BaseService) QueryBackupKV(prefix string) (map[string][]byte, error) {
	keys, err := bds.GetBackupKVKeys(prefix)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return map[string][]byte{}, nil
	}

	return bds.GetBackupKV(keys)
}

// 新的KV方法，根据前缀获取key，获取到的KV存储支持云端备份
func (bds *BaseService) GetBackupKVKeys(prefix string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	keys, err := bds.resourceCli.DriverStorageClient.GetKeysV2(ctx, &proto.GetPrefixReq{
		DriverServiceId: bds.GetGatewayId() + "/" + bds.GetServiceId(),
		Prefix:          prefix,
	})
	if err == nil && keys != nil {
		return keys.Key, nil
	} else {
		return []string{}, err
	}
}

// 新的KV方法，根据key获取KV，获取到的KV存储支持云端备份
func (bds *BaseService) GetBackupKV(keys []string) (map[string][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	var req = proto.GetReq{
		DriverServiceId: bds.GetGatewayId() + "/" + bds.GetServiceId(),
		Keys:            keys,
	}

	if resp, err := bds.resourceCli.DriverStorageClient.GetV2(ctx, &req); err != nil {
		return nil, errors.New(status.Convert(err).Message())
	} else {
		kvs := make(map[string][]byte, len(resp.GetKvs()))
		for _, value := range resp.GetKvs() {
			kvs[value.GetKey()] = value.GetValue()
		}
		return kvs, nil
	}
}

// 新的KV方法，更新KV，KV支持云端备份
func (bds *BaseService) PutBackupKV(kvs map[string][]byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var kv []*proto.KV
	for k, v := range kvs {
		kv = append(kv, &proto.KV{
			Key:   k,
			Value: v,
		})
	}
	var req = proto.PutReq{
		DriverServiceId: bds.GetGatewayId() + "/" + bds.GetServiceId(),
		Data:            kv,
	}

	if _, err := bds.resourceCli.DriverStorageClient.PutV2(ctx, &req); err != nil {
		return errors.New(status.Convert(err).Message())
	}
	return nil
}

// 新的KV方法，删除KV，KV支持云端备份
func (bds *BaseService) DelBackupKV(keys []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	var req = proto.DeleteReq{
		DriverServiceId: bds.GetGatewayId() + "/" + bds.GetServiceId(),
		Keys:            keys,
	}

	if _, err := bds.resourceCli.DriverStorageClient.DeleteV2(ctx, &req); err != nil {
		return errors.New(status.Convert(err).Message())
	}
	return nil
}

func (bds *BaseService) GetKV(keys []string) (map[string][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var req = proto.GetReq{
		DriverServiceId: bds.GetGatewayId() + "/" + bds.GetServiceId(),
		Keys:            keys,
	}
	if resp, err := bds.resourceCli.DriverStorageClient.Get(ctx, &req); err != nil {
		return nil, errors.New(status.Convert(err).Message())
	} else {
		kvs := make(map[string][]byte, len(resp.GetKvs()))
		for _, value := range resp.GetKvs() {
			kvs[value.GetKey()] = value.GetValue()
		}
		return kvs, nil
	}
}

func (bds *BaseService) PutKv(kvs map[string][]byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var kv []*proto.KV
	for k, v := range kvs {
		kv = append(kv, &proto.KV{
			Key:   k,
			Value: v,
		})
	}
	var req = proto.PutReq{
		DriverServiceId: bds.GetGatewayId() + "/" + bds.GetServiceId(),
		Data:            kv,
	}

	if _, err := bds.resourceCli.DriverStorageClient.Put(ctx, &req); err != nil {
		return errors.New(status.Convert(err).Message())
	}
	return nil
}

func (bds *BaseService) DeleteKV(keys []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	var req = proto.DeleteReq{
		DriverServiceId: bds.GetGatewayId() + "/" + bds.GetServiceId(),
		Keys:            keys,
	}

	if _, err := bds.resourceCli.DriverStorageClient.Delete(ctx, &req); err != nil {
		return errors.New(status.Convert(err).Message())
	}
	return nil
}

// Deprecated:GetAllValue 获取所有驱动存储的自定义内容，不支持云端备份
func (bds *BaseService) getAllKV() (map[string][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := bds.resourceCli.DriverStorageClient.All(ctx, &proto.AllReq{
		DriverServiceId: bds.GetGatewayId() + "/" + bds.GetServiceId(),
	})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}

	kvs := make(map[string][]byte, len(resp.Kvs))
	for _, v := range resp.GetKvs() {
		kvs[v.GetKey()] = v.GetValue()
	}

	return kvs, nil
}

func (bds *BaseService) GetKVKeys(prefix string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	keys, err := bds.resourceCli.DriverStorageClient.GetKeys(ctx, &proto.GetPrefixReq{
		DriverServiceId: bds.GetGatewayId() + "/" + bds.GetServiceId(),
		Prefix:          prefix,
	})
	if err == nil && keys != nil {
		return keys.Key, nil
	} else {
		return []string{}, err
	}
}

func (bds *BaseService) QueryKV(prefix string) (map[string][]byte, error) {
	keys, err := bds.GetKVKeys(prefix)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return map[string][]byte{}, nil
	}

	return bds.GetKV(keys)
}
