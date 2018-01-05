package corp

import (
	"encoding/json"
	"io/ioutil"
)

// Corp current Corporation
var Corp *Corporation

// Init 加载企业配置信息
// 可以配置多个企业
func Init(filename string) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(src, &Corp)
	if err != nil {
		panic(err)
	}

	// set corp
	for _, agent := range Corp.AgentMap {
		agent.Corp = Corp
	}
}
