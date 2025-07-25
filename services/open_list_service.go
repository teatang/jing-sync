package services

import (
	"gorm.io/gorm"
	"jing-sync/models"
	"jing-sync/services/db_services"
	"jing-sync/utils"
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

// func (c *OpenListClient) GetChildPath(path string) (utils.PageList[string], error) {

// }

// func (c *OpenListClient) GetChildPathRaw(path string) (utils.PageList[string], error) {

// }

func (c *OpenListClient) Post(api string, data interface{}) ([]byte, error) {
	url := c.Engine.Url + api
	token := c.Engine.Token

	ro := &utils.RequestOption{
		Headers: map[string]string{
			"Authorization": token,
		},
		Timeout: 30,
	}
	return utils.Post(url, data, ro)
}
