package services

import (
	"encoding/json"
	"fmt"
	"jing-sync/models"
	"jing-sync/services/db_services"
	"jing-sync/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

func NewOpenListClient(id string, db *gorm.DB) *OpenListClient {
	es := db_services.NewEngineService(db)
	engine, _ := es.GetByID(id)
	return &OpenListClient{
		Engine: engine,
		db:     db,
	}
}

type OpenListClient struct {
	Engine *models.Engine
	db     *gorm.DB
}

func (c *OpenListClient) GetChildPath(path string, speed int) (*utils.PageList[string], error) {
	res, err := c.GetChildPathRaw(path, speed)
	if err != nil {
		return nil, err
	}

	// 定义一个map来接收解析后的数据
	var r map[string]interface{}

	// 解析JSON字符串到map中
	e := json.Unmarshal(res, &r)
	if e != nil {
		return nil, e
	}

	return c.GetChildPathFormat(r)
}

func (c *OpenListClient) GetChildPathFormat(open_list_data map[string]interface{}) (*utils.PageList[string], error) {
	data := open_list_data["data"].(map[string]interface{})
	content_list := data["content"].([]interface{})
	fmt.Println(content_list)

	return nil, nil
}

func (c *OpenListClient) GetChildPathRaw(path string, speed int) ([]byte, error) {
	if speed == 2 {
		time.Sleep(1 * time.Second)
	}
	api := "/api/fs/list"
	data := map[string]interface{}{
		"path":     path,
		"refresh":  speed != 1,
		"page":     1,
		"per_page": 100,
	}

	res, err := c.Post(api, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *OpenListClient) Post(api string, data interface{}) ([]byte, error) {
	url := strings.Trim(c.Engine.Url, "/") + api
	token := c.Engine.Token

	ro := &utils.RequestOption{
		Headers: map[string]string{
			"Authorization": token,
		},
		Timeout: 30 * time.Second,
	}
	return utils.Post(url, data, ro)
}
