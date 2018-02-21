package cached

import (
	"errors"

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

// GetTagUserData func
func GetTagUserData(corpIndex int, tagname string) *TagGet {
	v := cachedTagUsers.Get(corpIndex, tagname)
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
	tagname := args[1].(string)

	cp := corp.Get(corpIndex)
	if cp == nil {
		return nil, 0, errors.New("corp not found in departCached")
	}

	token, err := cp.GetContactToken()
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
			_, usrs, partyIds, err := tag.WxTagUsers(tags[i].ID, token)
			if err != nil {
				return nil, 0, err
			}

			departUsers := make([]depart.User, 0)
			for _, pid := range partyIds {
				if items := GetDepartUsersByID(corpIndex, pid); len(items) > 0 {
					departUsers = append(departUsers, items...)
				}

			}

			dat := &TagGet{UserList: usrs, PartyList: partyIds, DepartUserList: departUsers}

			return dat, Expires(), nil
		}
	}

	return nil, 0, nil
}
