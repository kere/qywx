package cached

import (
	"errors"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/depart"
)

// cachedDeparts
var (
	// cachedDeparts a
	cachedDeparts = newDepartsMap()
)

// DepartsMap class
type DepartsMap struct {
	cache.Map
}

func newDepartsMap() *DepartsMap {
	t := &DepartsMap{}
	t.Init(t)
	return t
}

// GetDeparts func
func GetDeparts(corpIndex int, departID int) []depart.Department {
	v := cachedDeparts.Get(corpIndex, departID)
	if v == nil {
		return nil
	}
	return v.([]depart.Department)
}

// Build func
func (t *DepartsMap) Build(args ...interface{}) (interface{}, int, error) {
	corpIndex := args[0].(int)
	departID := args[1].(int)

	cp := corp.Get(corpIndex)
	if cp == nil {
		return nil, 0, errors.New("corp not found in departCached")
	}

	token, err := cp.GetContactToken()
	if err != nil {
		return nil, 0, err
	}

	dat, err := depart.WxDepartments(departID, token)
	if err != nil {
		return nil, 0, err
	}

	return dat, Expires(), nil
}
