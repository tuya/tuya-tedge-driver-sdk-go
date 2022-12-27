package cache

import (
	"sync"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
)

type DeviceCache struct {
	mutex     sync.RWMutex
	deviceMap map[string]*commons.DeviceInfo // key is device cid
}

func NewDeviceCache() *DeviceCache {
	devCache := &DeviceCache{
		deviceMap: make(map[string]*commons.DeviceInfo),
	}

	return devCache
}

// Add adds a new device to the cache. This method is used to populate the
// device cache with pre-existing or recently-added devices from Core Metadata.
func (dc *DeviceCache) Add(d commons.DeviceInfo) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	dc.deviceMap[d.Cid] = &d
}

// Update updates the device in the cache
func (dc *DeviceCache) Update(d commons.DeviceInfo) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	dc.deviceMap[d.Cid] = &d
}

// RemoveById removes the specified device by id from the cache.
func (dc *DeviceCache) RemoveById(id string) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	delete(dc.deviceMap, id)
}


// ById returns a device with the given device id.
func (dc *DeviceCache) ById(id string) (commons.DeviceInfo, bool) {
	dc.mutex.RLock()
	defer dc.mutex.RUnlock()

	d, ok := dc.deviceMap[id]
	if !ok {
		return commons.DeviceInfo{}, ok
	}
	return *d, ok
}

// All returns the current list of devices in the cache.
func (dc *DeviceCache) All() map[string]commons.DeviceInfo {
	dc.mutex.RLock()
	defer dc.mutex.RUnlock()

	dMap := make(map[string]commons.DeviceInfo, len(dc.deviceMap))
	for k, d := range dc.deviceMap {
		dMap[k] = *d
	}
	return dMap
}
