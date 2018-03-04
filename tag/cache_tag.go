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
func GetTags(corpIndex int) []Tag {
	v := cachedTags.Get(corpIndex)
	if v == nil {
		return nil
	}
	return v.([]Tag)
}

// Build func
func (t *TagMap) Build(args ...interface{}) (interface{}, int, error) {
	corpIndex := args[0].(int)

	cp := corp.Get(corpIndex)
	if cp == nil {
		return nil, 0, errors.New("corp not found in departCached")
	}

	token, err := cp.GetContactToken()
	if err != nil {
		return nil, 0, err
	}

	dat, err := WxTags(token)
	if err != nil {
		return nil, 0, err
	}

	return dat, util.Expires(), nil
}

// IsUserInTag 用户是否属于当前标签
func IsUserInTag(corpIndex int, userid, tagName string) bool {
	tagGet := GetTagUserData(corpIndex, tagName)
	if tagGet == nil {
		return false
	}

	for _, u := range tagGet.UserList {
		if u.UserID == userid {
			return true
		}
	}

	return false
}

// IsDepartInTag 用户是否属于当前标签
func IsDepartInTag(corpIndex, departID int, tagName string) bool {
	tagGet := GetTagUserData(corpIndex, tagName)
	if tagGet == nil {
		return false
	}

	for _, partyID := range tagGet.PartyList {
		if partyID == departID {
			return true
		}
	}

	return false
}
