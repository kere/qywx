package corp

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// corps Corporation map
var corps = make([]*Corporation, 0)

// Init 加载企业配置信息
// 可以配置多个企业
func Init(filename string) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(src, &corps)
	if err != nil {
		panic(err)
	}

	for _, c := range corps {
		// set corp
		for aName, agent := range c.AgentMap {
			agent.Name = aName
			agent.Corp = c
		}
	}

}

// Get get corp by name
func Get(i int) *Corporation {
	return corps[i]
}

// GetByID get corp by id
func GetByID(corpID string) (*Corporation, error) {
	for _, v := range corps {
		if v.ID == corpID {
			return v, nil
		}
	}
	return nil, errors.New("corp id is not found")
}

// GetByName get corp by name
func GetByName(corpName string) (*Corporation, error) {
	for _, v := range corps {
		if v.Name == corpName {
			return v, nil
		}
	}
	return nil, errors.New("corp name is not found")
}
