package common

import (
	"path/filepath"
	"runtime"

	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	Resource   = "Resource"
	MQTTBroker = "MQTTBroker"
)

const ContentMaxLen = 0xFF

var EmptyPb = new(emptypb.Empty)

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

var basepath string

func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	return filepath.Join(basepath, rel)
}
