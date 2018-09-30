package depart

import (
	"errors"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/util"
)

var (
	cachedDepartSimpleUsers = newDepartSimpleUsersMap()
)

// DepartSimpleUsersMap class
type DepartSimpleUsersMap struct {
	cache.Map
}

func newDepartSimpleUsersMap() *DepartSimpleUsersMap {
	t := &DepartSimpleUsersMap{}
	t.Init(t, 0)
	return t
}

// GetDepartUsersByID func
func GetDepartUsersByID(corpIndex int, departID int) []User {
	v := cachedDepartSimpleUsers.Get(corpIndex, departID)
	if v == nil {
		return nil
	}
	return v.([]User)
}

// GetDepartSimpleUsers func
func GetDepartSimpleUsers(corpIndex int, departName string) []User {
	items := GetDeparts(corpIndex, 0)
	departID := 0
	for _, v := range items {
		if departName == v.Name {
			departID = v.ID
		}
	}
	if departID == 0 {
		return nil
	}

	v := cachedDepartSimpleUsers.Get(corpIndex, departID)
	if v == nil {
		return nil
	}
	return v.([]User)
}

// Build func
func (t *DepartSimpleUsersMap) Build(args ...interface{}) (interface{}, error) {
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

	usrs, err := WxDepartSimpleUsers(departID, true, token)

	t.SetExpires(util.Expires())
	return usrs, err
}
