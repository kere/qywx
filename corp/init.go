package corp

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/kere/gno/db"
)

var (
	//TableCorp table
	TableCorp = "corps"
	//TableAgent table
	TableAgent = "agents"

	// corps Corporation map
	corps = make([]*Corporation, 0)

	isLoaded = false
)

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

	isLoaded = true
}

// Get get corp by name
func Get(i int) *Corporation {
	return corps[i]
}

// GetByID get corp by id
func GetByID(corpid string) (*Corporation, error) {
	if isLoaded {
		for _, v := range corps {
			if v.Corpid == corpid {
				return v, nil
			}
		}
		return nil, errors.New("corp id is not found")
	}

	return LoadByID(corpid)
}

// GetByName get corp by name
func GetByName(corpName string) (*Corporation, error) {
	if isLoaded {
		for _, v := range corps {
			if v.Name == corpName {
				return v, nil
			}
		}
		return nil, errors.New("corp name is not found")
	}

	return LoadByName(corpName)
}

// LoadByID get corp by id
func LoadByID(corpid string) (*Corporation, error) {
	q := db.NewQueryBuilder(TableCorp).Cache()
	row, err := q.Where("corpid=?", corpid).QueryOne()
	if err != nil {
		return nil, err
	}

	if row.IsEmpty() {
		return nil, errors.New("query corp is empty by " + corpid)
	}

	return buildCorp(row)
}

// LoadByName get corp by name
func LoadByName(corpName string) (*Corporation, error) {
	q := db.NewQueryBuilder(TableCorp).Cache()
	row, err := q.Where("name=?", corpName).QueryOne()
	if err != nil {
		return nil, err
	}

	if row.IsEmpty() {
		return nil, errors.New("query corp is empty by " + corpName)
	}

	return buildCorp(row)
}

// buildCorp get corp by id
func buildCorp(row db.DataRow) (*Corporation, error) {
	c := &Corporation{}

	c.ID = row.Int("id")
	c.Corpid = row.String("corpid")
	c.Name = row.String("name")
	c.Title = row.String("title")

	c.ContactsToken = row.String("c_token")
	c.ContactsSecret = row.String("c_secret")
	c.ContactsAesKey = row.String("c_aeskey")

	rows, _ := db.NewQueryBuilder(TableAgent).Cache().Where("corp_id=?", c.ID).Query()

	c.AgentMap = make(map[string]*Agent, 0)
	var name string
	var agent *Agent
	for _, row = range rows {
		name = row.String("name")
		agent = &Agent{
			Corp:      c,
			ID:        row.Int("id"),
			Agentid:   row.Int("agentid"),
			Name:      row.String("name"),
			Secret:    row.String("secret"),
			MsgToken:  row.String("msgtoken"),
			MsgAesKey: row.String("msgaeskey"),
		}
		c.AgentMap[name] = agent
	}

	return c, nil
}
