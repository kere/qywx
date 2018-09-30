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
	t.Init(t, 0)
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
func (t *DepartUsersMap) Build(args ...interface{}) (interface{}, error) {
	corpID := args[0].(int)
	departID := args[1].(int)

	cp := corp.Get(corpID)
	if cp == nil {
		return nil, errors.New("corp not found in departCached")
	}

	token, err := cp.GetContactToken()
	if err != nil {
		return nil, err
	}

	usrs, err := WxDepartUsers(departID, true, token)

	t.SetExpires(util.Expires())
	return usrs, err
}
