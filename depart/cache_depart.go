package depart

import (
	"errors"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/util"
)

// cachedDeparts
var (
	// cachedDeparts a
	cachedDeparts = newDepartsMap()
)

// DepartmentByName 通过部门名称获得部门信息
func DepartmentByName(name string) Department {
	all := GetDeparts(0, 0)
	for _, item := range all {
		if item.Name == name {
			return item
		}
	}

	return Department{}
}

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
func GetDeparts(corpIndex int, departID int) []Department {
	v := cachedDeparts.Get(corpIndex, departID)
	if v == nil {
		return nil
	}
	return v.([]Department)
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

	dat, err := WxDepartments(departID, token)
	if err != nil {
		return nil, 0, err
	}

	return dat, util.Expires(), nil
}
