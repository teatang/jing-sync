package services

import (
	"encoding/json"
	"jing-sync/models"
	"jing-sync/services/db_services"
	"jing-sync/utils"

	"strings"
	"time"

	"gorm.io/gorm"
)

type OpenListService struct {
    db *gorm.DB
}

func NewOpenListService(db *gorm.DB) *OpenListService {
	return &OpenListService{db: db}
}

func (s *OpenListService)GetOpenListInfo(engine_id string, path string) (*utils.PageList[string], error) {
	return NewOpenListClient(engine_id, s.db).GetChildPath(path, 0)
}

type ChildPathRawInfo struct {
	Name     string `json:"name"`
	IsDir    bool   `json:"is_dir"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

type ChildPathRawResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Total   int64  `json:"total"`
	Data    struct {
		Content []ChildPathRawInfo `json:"content"`
	} `json:"data"`
}

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

	var r ChildPathRawResponse
	// 解析JSON字符串到ChildPathRawResponse中
	e := json.Unmarshal(res, &r)
	if e != nil {
		return nil, e
	}

	return c.GetChildPathFormat(r)
}

func (c *OpenListClient) GetChildPathFormat(open_list_data ChildPathRawResponse) (*utils.PageList[string], error) {
	var path_list []string
	for _, c := range open_list_data.Data.Content {
		if !c.IsDir {
			continue
		}

		path_list = append(path_list, c.Name)
	}

	return &utils.PageList[string]{
		List: path_list,
		Pagination: utils.PageInfo{
			Page:  1,
			Size:  100,
			Total: open_list_data.Total,
		},
	}, nil
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
