package cached

import (
	"errors"

	"github.com/kere/gno"
	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/users"
)

// CachedTag
var (
	// cachedTag a
	cachedTags = newTagMap()
)

// TagMap class
type TagMap struct {
	cache.Map
}

func newTagMap() *TagMap {
	t := &TagMap{}
	t.Init(t)
	return t
}

// GetTags func
func GetTags(corpID string, agentID int) []users.Tag {
	v := cachedTags.Get(corpID, agentID)
	if v == nil {
		return nil
	}
	return v.([]users.Tag)
}

// Build func
func (t *TagMap) Build(args ...interface{}) (interface{}, int, error) {
	corpID := args[0].(string)
	agentID := args[1].(int)

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

	dat, err := users.WxTags(token)
	if err != nil {
		return nil, 0, err
	}

	expires := gno.GetConfig().GetConf("data").DefaultInt("data_expires", 72000)

	return dat, expires, nil
}
