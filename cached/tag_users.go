package cached

import (
	"errors"

	"github.com/kere/gno"
	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/depart"
	"github.com/kere/qywx/tag"
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
	UserList       []tag.User    `json:"userlist"`
	PartyList      []int         `json:"partylist"`
	DepartUserList []depart.User `json:"depart_userlist"`
}

func newTagUsersMap() *TagUsersMap {
	t := &TagUsersMap{}
	t.Init(t)
	return t
}

// GetTagUsers func
func GetTagUsers(corpID int, agentName, tagname string) *TagGet {
	v := cachedTagUsers.Get(corpID, agentName, tagname)
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
	corpID := args[0].(int)
	agentName := args[1].(string)
	tagname := args[2].(string)

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

	tags := GetTags(corpID, agentName)
	l := len(tags)
	if l == 0 {
		return nil, 0, err
	}

	expires := gno.GetConfig().GetConf("data").DefaultInt("data_expires", 72000)

	for i := 0; i < l; i++ {
		if tags[i].Name == tagname {
			_, usrs, partyIds, err := tag.WxTagUsers(tags[i].ID, token)
			if err != nil {
				return nil, 0, err
			}

			departUsers := make([]depart.User, 0)
			for _, pid := range partyIds {
				if items := GetDepartUsersByID(corpID, agentName, pid); len(items) > 0 {
					departUsers = append(departUsers, items...)
				}

			}

			dat := &TagGet{UserList: usrs, PartyList: partyIds, DepartUserList: departUsers}

			return dat, expires, nil
		}
	}

	return nil, 0, nil
}
