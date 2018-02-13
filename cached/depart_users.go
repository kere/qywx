package cached

import (
	"errors"

	"github.com/kere/gno"
	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/depart"
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

// GetDepartUsersByID func
func GetDepartUsersByID(corpID int, agentName string, departID int) []depart.User {
	v := cachedDepartUsers.Get(corpID, agentName, departID)
	if v == nil {
		return nil
	}
	return v.([]depart.User)
}

// GetDepartUsers func
func GetDepartUsers(corpID int, agentName, departName string) []depart.User {
	items := GetDeparts(corpID, agentName, 0)
	departID := 0
	for _, v := range items {
		if departName == v.Name {
			departID = v.ID
		}
	}
	if departID == 0 {
		return nil
	}

	v := cachedDepartUsers.Get(corpID, agentName, departID)
	if v == nil {
		return nil
	}
	return v.([]depart.User)
}

// ClearDepart func
func ClearDepart() {
	cachedDepartUsers.ClearAll()
	cachedDeparts.ClearAll()
}

// Build func
func (t *DepartUsersMap) Build(args ...interface{}) (interface{}, int, error) {
	corpID := args[0].(int)
	agentName := args[1].(string)
	departID := args[2].(int)

	cp := corp.Get(corpID)
	if cp == nil {
		return nil, 0, errors.New("corp not found in departCached")
	}

	agent, err := cp.GetAgent(agentName)
	if err != nil {
		return nil, 0, err
	}

	token, err := agent.GetToken()
	if err != nil {
		return nil, 0, err
	}

	expires := gno.GetConfig().GetConf("data").DefaultInt("data_expires", 72000)

	usrs, err := depart.WxDepartUsers(departID, true, token)

	return usrs, expires, err
}
