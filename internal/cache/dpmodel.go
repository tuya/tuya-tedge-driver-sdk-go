package cache

import (
	"sync"

	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
)

type DPModelProductCache struct {
	mutex      sync.RWMutex
	productMap map[string]*dpmodel.DPModelProduct // key is product id
	dpMap      map[string]map[string]dpmodel.DP
}

func NewDPProductCache() *DPModelProductCache {
	return &DPModelProductCache{
		productMap: make(map[string]*dpmodel.DPModelProduct),
		dpMap:      make(map[string]map[string]dpmodel.DP),
	}
}

func (pc *DPModelProductCache) Add(p dpmodel.DPModelProduct) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	pc.productMap[p.Id] = &p
	pc.dpMap[p.Id] = DpmodelSliceToMap(p.Dps)
}

func (pc *DPModelProductCache) ById(id string) (dpmodel.DPModelProduct, bool) {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	p, ok := pc.productMap[id]
	if !ok {
		return dpmodel.DPModelProduct{}, ok
	}
	return *p, ok
}

func (pc *DPModelProductCache) All() map[string]dpmodel.DPModelProduct {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	ps := make(map[string]dpmodel.DPModelProduct, len(pc.productMap))
	for k, p := range pc.productMap {
		ps[k] = *p
	}
	return ps
}


func DpmodelSliceToMap(resources []dpmodel.DP) map[string]dpmodel.DP {
	drMap := make(map[string]dpmodel.DP, len(resources))
	for _, dr := range resources {
		drMap[dr.Id] = dr
	}
	return drMap
}

func (pc *DPModelProductCache) Update(p dpmodel.DPModelProduct) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	pc.removeById(p.Id)
	pc.productMap[p.Id] = &p
	pc.dpMap[p.Id] = DpmodelSliceToMap(p.Dps)
}

func (pc *DPModelProductCache) RemoveById(id string) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	pc.removeById(id)
}

func (pc *DPModelProductCache) removeById(id string) {
	delete(pc.productMap, id)
	delete(pc.dpMap, id)
}

func (pc *DPModelProductCache) Dp(productId, dpId string) (dpmodel.DP, bool) {
	var (
		ok  bool
		drs map[string]dpmodel.DP
		dr  dpmodel.DP
	)
	pc.mutex.RLock()
	pc.mutex.RUnlock()

	if drs, ok = pc.dpMap[productId]; !ok {
		return dr, ok
	}
	dr, ok = drs[dpId]
	return dr, ok
}
