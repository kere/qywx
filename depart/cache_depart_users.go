package depart

import (
	"errors"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/users"
	"github.com/kere/qywx/util"
)

// CachedDepart
var (
	// cachedDepart a
	cachedDepartUsers = newDepartUsersMap()
)

// DepartUsersMap class
type DepartUsersMap struct {
	cache.Map
}

func newDepartUsersMap() *DepartUsersMap {
	t := &DepartUsersMap{}
	t.Init(t)
	return t
}

// GetDepartUsers func
func GetDepartUsers(corpIndex int, departName string) []users.UserDetail {
	dp := DepartmentByName(departName)
	departID := dp.ID
	if departID == 0 {
		return nil
	}

	v := cachedDepartUsers.Get(corpIndex, departID)
	if v == nil {
		return nil
	}
	return v.([]users.UserDetail)
}

// ClearDepart func
func ClearDepart() {
	cachedDepartUsers.ClearAll()
	cachedDepartSimpleUsers.ClearAll()
	cachedDeparts.ClearAll()
}

// Build func
func (t *DepartUsersMap) Build(args ...interface{}) (interface{}, int, error) {
	corpID := args[0].(int)
	departID := args[1].(int)

	cp := corp.Get(corpID)
	if cp == nil {
		return nil, 0, errors.New("corp not found in departCached")
	}

	token, err := cp.GetContactToken()
	if err != nil {
		return nil, 0, err
	}

	usrs, err := WxDepartUsers(departID, true, token)

	return usrs, util.Expires(), err
}
