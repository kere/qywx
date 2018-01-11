package corp

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// corpMap Corporation map
var corpMap = make(map[string]*Corporation, 0)

// Init 加载企业配置信息
// 可以配置多个企业
func Init(filename string) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(src, &corpMap)
	if err != nil {
		panic(err)
	}

	for _, c := range corpMap {
		// set corp
		for _, agent := range c.AgentMap {
			agent.Corp = c
		}
	}

}

// GetByName get corp by name
func GetByName(name string) (*Corporation, error) {
	if c, isok := corpMap[name]; isok {
		return c, nil
	}

	return nil, errors.New("corp name not found")
}
