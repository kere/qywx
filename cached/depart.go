package cached

import (
	"errors"

	"github.com/kere/gno"
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
func GetDeparts(corpID int, agentName string, departID int) []depart.Department {
	v := cachedDeparts.Get(corpID, agentName, departID)
	if v == nil {
		return nil
	}
	return v.([]depart.Department)
}

// Build func
func (t *DepartsMap) Build(args ...interface{}) (interface{}, int, error) {
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

	dat, err := depart.WxDepartments(departID, token)
	if err != nil {
		return nil, 0, err
	}

	expires := gno.GetConfig().GetConf("data").DefaultInt("data_expires", 72000)

	return dat, expires, nil
}
