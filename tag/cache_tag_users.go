package tag

import (
	"errors"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/util"
)

// CachedTag
var (
	// cachedTag a
	cachedTagUsers = newTagUsersMap()
)

// TagUsersMap class
type TagUsersMap struct {
	cache.Map
}

// TagGet class
type TagGet struct {
	UserList  []User `json:"userlist"`
	PartyList []int  `json:"partylist"`
	// DepartUserList []depart.User `json:"depart_userlist"`
}

func newTagUsersMap() *TagUsersMap {
	t := &TagUsersMap{}
	t.Init(t, 0)
	return t
}

// GetTagUserData func
func GetTagUserData(corpid, agentName, tagname string) *TagGet {
	v := cachedTagUsers.Get(corpid, agentName, tagname)
	if v == nil {
		return nil
	}
	return v.(*TagGet)
}

// ClearTag func
func ClearTag() {
	cachedTagUsers.ClearAll()
	cachedTags.ClearAll()
}

// Build func
func (t *TagUsersMap) Build(args ...interface{}) (interface{}, error) {
	corpid := args[0].(string)
	agentName := args[1].(string)
	tagname := args[2].(string)

	cp, err := corp.GetByCorpid(corpid)
	if err != nil {
		return nil, errors.New("corp name not found")
	}
	if cp == nil {
		return nil, errors.New("corp not found in departCached")
	}
	agent, err := cp.GetAgent(agentName)
	var token string
	if err != nil {
		token, err = cp.GetContactToken()
	} else {
		token, err = agent.GetToken()
	}

	if err != nil {
		return nil, err
	}

	tags := GetTags(corpid)
	l := len(tags)
	if l == 0 {
		return nil, err
	}

	for i := 0; i < l; i++ {
		if tags[i].Name == tagname {
			_, usrs, partyIds, err := WxTagUsers(tags[i].ID, token)
			if err != nil {
				return nil, err
			}

			dat := &TagGet{UserList: usrs, PartyList: partyIds}

			if t.GetExpires() == 0 {
				t.SetExpires(util.Expires())
			}
			return dat, nil
		}
	}

	return nil, nil
}
