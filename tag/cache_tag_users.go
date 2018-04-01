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
	t.Init(t)
	return t
}

// GetTagUserData func
func GetTagUserData(corpIndex int, agentName, tagname string) *TagGet {
	v := cachedTagUsers.Get(corpIndex, agentName, tagname)
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
func (t *TagUsersMap) Build(args ...interface{}) (interface{}, int, error) {
	corpIndex := args[0].(int)
	agentName := args[1].(string)
	tagname := args[2].(string)

	cp := corp.Get(corpIndex)
	if cp == nil {
		return nil, 0, errors.New("corp not found in departCached")
	}
	agent, err := cp.GetAgent(agentName)
	var token string
	if err != nil {
		token, err = cp.GetContactToken()
	} else {
		token, err = agent.GetToken()
	}

	if err != nil {
		return nil, 0, err
	}

	tags := GetTags(corpIndex)
	l := len(tags)
	if l == 0 {
		return nil, 0, err
	}

	for i := 0; i < l; i++ {
		if tags[i].Name == tagname {
			_, usrs, partyIds, err := WxTagUsers(tags[i].ID, token)
			if err != nil {
				return nil, 0, err
			}

			// departUsers := make([]depart.User, 0)
			// for _, pid := range partyIds {
			// 	if items := depart.GetDepartUsersByID(corpIndex, pid); len(items) > 0 {
			// 		departUsers = append(departUsers, items...)
			// 	}
			//
			// }

			dat := &TagGet{UserList: usrs, PartyList: partyIds}

			return dat, util.Expires(), nil
		}
	}

	return nil, 0, nil
}
