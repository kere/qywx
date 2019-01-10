package corp

import (
	"encoding/json"
	"errors"
	"fmt"
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

// GetByCorpid get corp by id
func GetByCorpid(corpid string) (*Corporation, error) {
	if isLoaded {
		for _, v := range corps {
			if v.Corpid == corpid {
				return v, nil
			}
		}
		return nil, errors.New("corp id is not found")
	}

	return LoadByCorpid(corpid)
}

// GetByID get corp by id
func GetByID(corpID int) (*Corporation, error) {
	if isLoaded {
		for _, v := range corps {
			if v.ID == corpID {
				return v, nil
			}
		}
		return nil, errors.New("corp id is not found")
	}

	return LoadByID(corpID)
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
func LoadByID(corpID int) (*Corporation, error) {
	q := db.NewQueryBuilder(TableCorp)
	q.Cache()
	row, err := q.Where("id=?", corpID).QueryOne()
	if err != nil {
		return nil, err
	}

	if row.IsEmpty() {
		return nil, errors.New("query corp is empty by corp " + fmt.Sprint(corpID))
	}

	return buildCorp(row)
}

// LoadByCorpid get corp by id
func LoadByCorpid(corpid string) (*Corporation, error) {
	q := db.NewQueryBuilder(TableCorp)
	q.Cache()
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
	q := db.NewQueryBuilder(TableCorp)
	q.Cache()
	row, err := q.Where("name=?", corpName).QueryOne()
	if err != nil {
		return nil, err
	}

	if row.IsEmpty() {
		return nil, errors.New("query corp is empty by " + corpName)
	}

	return buildCorp(row)
}

// GetAgentByID get agent by id
func GetAgentByID(agentID int) (*Agent, error) {
	q := db.NewQueryBuilder(TableAgent)
	q.Where("id=?", agentID).Cache()
	row, err := q.QueryOne()
	if err != nil {
		return nil, err
	}
	if row.IsEmpty() {
		return nil, errors.New("GetAgentByID: not found")
	}

	c, _ := GetByID(row.Int("corp_id"))
	if c == nil {
		return nil, errors.New("GetAgentByID: not found corp")
	}

	a := c.GetAgentByID(agentID)
	if a == nil {
		return nil, errors.New("GetAgentByID: not found agent " + fmt.Sprint(agentID))
	}
	return a, nil
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
	c.MainURL = row.String("main_url")

	q := db.NewQueryBuilder(TableAgent)
	rows, _ := q.Cache().Where("corp_id=?", c.ID).Query()

	c.AgentMap = make(map[string]*Agent, 0)
	var name string
	var agent *Agent
	for _, row = range rows {
		name = row.String("name")
		agent = row2agent(c, row)

		if agent.ParentID > 0 {
			agent.Parent = row2agent(c, rows.Search("id", agent.ParentID))
		}

		c.AgentMap[name] = agent
	}

	return c, nil
}

func row2agent(c *Corporation, row db.DataRow) *Agent {
	if row.IsEmpty() {
		return nil
	}
	return &Agent{
		Corp:      c,
		ID:        row.Int("id"),
		Agentid:   row.Int("agentid"),
		Name:      row.String("name"),
		Secret:    row.String("secret"),
		ParentID:  row.IntDefault("parent_id", 0),
		MsgToken:  row.String("msgtoken"),
		MsgAesKey: row.String("msgaeskey"),
	}
}
