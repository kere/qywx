package cached

import (
	"errors"

	"github.com/kere/gno"
	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/users"
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
func GetDeparts(corpID string, agentID, departID int) []users.Department {
	v := cachedDeparts.Get(corpID, agentID, departID)
	if v == nil {
		return nil
	}
	return v.([]users.Department)
}

// Build func
func (t *DepartsMap) Build(args ...interface{}) (interface{}, int, error) {
	corpID := args[0].(string)
	agentID := args[1].(int)
	departID := args[2].(int)

	cp := corp.GetByID(corpID)
	if cp == nil {
		return nil, 0, errors.New("corp not found in departCached")
	}

	agent := cp.GetAgentByID(agentID)
	if agent == nil {
		return nil, 0, errors.New("agent not found in departCached")
	}

	token, err := agent.GetToken()
	if err != nil {
		return nil, 0, err
	}

	dat, err := users.WxDepartments(departID, token)
	if err != nil {
		return nil, 0, err
	}

	expires := gno.GetConfig().GetConf("data").DefaultInt("data_expires", 72000)

	return dat, expires, nil
}
