package cache

import (
	"sync"

	"github.com/tuya/tuya-tedge-driver-sdk-go/thingmodel"
)

type ThingModelProductCache struct {
	mu          sync.RWMutex
	productMap  map[string]*thingmodel.ThingModelProduct
	propertyMap map[string]map[string]thingmodel.PropertySpec // thing model property
	actionMap   map[string]map[string]thingmodel.Action       // thing model action
	eventMap    map[string]map[string]thingmodel.Event        // thing model event
}

func NewTyModelProduct() *ThingModelProductCache {
	return &ThingModelProductCache{
		productMap:  make(map[string]*thingmodel.ThingModelProduct),
		propertyMap: make(map[string]map[string]thingmodel.PropertySpec),
		actionMap:   make(map[string]map[string]thingmodel.Action),
		eventMap:    make(map[string]map[string]thingmodel.Event),
	}
}

func (tmpc *ThingModelProductCache) Add(p thingmodel.ThingModelProduct) {
	tmpc.mu.Lock()
	defer tmpc.mu.Unlock()

	tmpc.productMap[p.Id] = &p
	tmpc.propertyMap[p.Id] = propertySliceToMap(p.Property)
	tmpc.actionMap[p.Id] = actionSliceToMap(p.Action)
	tmpc.eventMap[p.Id] = eventSliceToMap(p.Event)
}

func (tmpc *ThingModelProductCache) ById(id string) (thingmodel.ThingModelProduct, bool) {
	tmpc.mu.RLock()
	defer tmpc.mu.RUnlock()

	p, ok := tmpc.productMap[id]
	if !ok {
		return thingmodel.ThingModelProduct{}, false
	}
	return *p, ok
}

func (tmpc *ThingModelProductCache) All() map[string]thingmodel.ThingModelProduct {
	tmpc.mu.RLock()
	defer tmpc.mu.RUnlock()

	ps := make(map[string]thingmodel.ThingModelProduct, len(tmpc.productMap))
	for k, p := range tmpc.productMap {
		ps[k] = *p
	}
	return ps
}

func (tmpc *ThingModelProductCache) Update(p thingmodel.ThingModelProduct) {
	tmpc.mu.Lock()
	defer tmpc.mu.Unlock()

	tmpc.removeById(p.Id)
	tmpc.productMap[p.Id] = &p
	tmpc.propertyMap[p.Id] = propertySliceToMap(p.Property)
	tmpc.actionMap[p.Id] = actionSliceToMap(p.Action)
	tmpc.eventMap[p.Id] = eventSliceToMap(p.Event)

}
func (tmpc *ThingModelProductCache) RemoveById(id string) {
	tmpc.mu.Lock()
	defer tmpc.mu.Unlock()

	tmpc.removeById(id)
}

func (tmpc *ThingModelProductCache) removeById(id string) {
	delete(tmpc.productMap, id)
	delete(tmpc.propertyMap, id)
	delete(tmpc.actionMap, id)
	delete(tmpc.eventMap, id)
}

func (tmpc *ThingModelProductCache) GetPropertySpecByCode(pid, code string) (thingmodel.PropertySpec, bool) {
	tmpc.mu.RLock()
	defer tmpc.mu.RUnlock()

	pm, ok := tmpc.propertyMap[pid]
	if !ok {
		return thingmodel.PropertySpec{}, false
	}

	ps, ok := pm[code]
	if !ok {
		return thingmodel.PropertySpec{}, false
	}

	return ps, ok
}

func (tmpc *ThingModelProductCache) GetActionSpecByCode(pid, code string) (thingmodel.Action, bool) {
	tmpc.mu.RLock()
	defer tmpc.mu.RUnlock()

	am, ok := tmpc.actionMap[pid]
	if !ok {
		return thingmodel.Action{}, false
	}

	a, ok := am[code]
	if !ok {
		return thingmodel.Action{}, false
	}

	return a, ok
}

func (tmpc *ThingModelProductCache) GetEventsByPid(pid string) (map[string]thingmodel.Event, bool) {
	tmpc.mu.RLock()
	defer tmpc.mu.RUnlock()

	e, ok := tmpc.eventMap[pid]
	if !ok {
		return map[string]thingmodel.Event{}, false
	}

	return e, ok
}

func propertySliceToMap(properties []thingmodel.PropertySpec) map[string]thingmodel.PropertySpec {
	pMap := make(map[string]thingmodel.PropertySpec, len(properties))
	for _, p := range properties {
		pMap[p.Code] = p
	}

	return pMap
}

func actionSliceToMap(actions []thingmodel.Action) map[string]thingmodel.Action {
	aMap := make(map[string]thingmodel.Action, len(actions))
	for _, a := range actions {
		aMap[a.Code] = a
	}

	return aMap
}

func eventSliceToMap(events []thingmodel.Event) map[string]thingmodel.Event {
	eMap := make(map[string]thingmodel.Event, len(events))
	for _, e := range events {
		eMap[e.Code] = e
	}

	return eMap
}
