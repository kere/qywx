package cached

import (
	"errors"

	"github.com/kere/gno"
	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/tag"
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
func GetTags(corpID int, agentName string) []tag.Tag {
	v := cachedTags.Get(corpID, agentName)
	if v == nil {
		return nil
	}
	return v.([]tag.Tag)
}

// Build func
func (t *TagMap) Build(args ...interface{}) (interface{}, int, error) {
	corpID := args[0].(int)
	agentName := args[1].(string)

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

	dat, err := tag.WxTags(token)
	if err != nil {
		return nil, 0, err
	}

	expires := gno.GetConfig().GetConf("data").DefaultInt("data_expires", 72000)

	return dat, expires, nil
}
