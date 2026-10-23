package services

import (
	"encoding/json"
	"fmt"
	"jing-sync/internal/models"
	"jing-sync/internal/services/db_services"
	"jing-sync/internal/utils"
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

func (s *OpenListService) GetOpenListInfo(engine_id string, path string) (*utils.PageList[string], error) {
	return NewOpenListClient(engine_id, s.db).GetChildPath(path, 0)
}

// FileInfo 文件与目录元数据
type FileInfo struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"is_dir"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

type ChildPathRawInfo struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
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

type CommonResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
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

// TestConnection 测试引擎连通性
func (c *OpenListClient) TestConnection() error {
	if c.Engine == nil || c.Engine.Url == "" {
		return fmt.Errorf("engine url is empty")
	}

	// 优先调用 /api/me 校验
	res, err := c.Get("/api/me", nil)
	if err == nil {
		var commonResp CommonResponse
		if err := json.Unmarshal(res, &commonResp); err == nil && commonResp.Code == 200 {
			return nil
		}
	}

	// 回退调用 /api/fs/list 根路径校验连通性
	data := map[string]interface{}{
		"path":     "/",
		"refresh":  false,
		"page":     1,
		"per_page": 1,
	}
	res, err = c.Post("/api/fs/list", data)
	if err != nil {
		return err
	}

	var listResp CommonResponse
	if err := json.Unmarshal(res, &listResp); err != nil {
		return err
	}
	if listResp.Code != 200 {
		return fmt.Errorf("connection failed: %s", listResp.Message)
	}

	return nil
}

// ListFilesDetailed 获取包含文件大小、类型的完整列表
func (c *OpenListClient) ListFilesDetailed(path string, speed int) ([]FileInfo, error) {
	res, err := c.GetChildPathRaw(path, speed)
	if err != nil {
		return nil, err
	}

	var r ChildPathRawResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return nil, err
	}
	if r.Code != 200 {
		return nil, fmt.Errorf("get file list failed: %s", r.Message)
	}

	files := make([]FileInfo, 0, len(r.Data.Content))
	for _, item := range r.Data.Content {
		files = append(files, FileInfo{
			Name:     item.Name,
			Size:     item.Size,
			IsDir:    item.IsDir,
			Created:  item.Created,
			Modified: item.Modified,
		})
	}
	return files, nil
}

// CopyFile 调用 OpenList/AList v3 /api/fs/copy 复制文件
func (c *OpenListClient) CopyFile(srcDir, dstDir, name string) (string, error) {
	data := map[string]interface{}{
		"src_dir": srcDir,
		"dst_dir": dstDir,
		"names":   []string{name},
	}
	res, err := c.Post("/api/fs/copy", data)
	if err != nil {
		return "", err
	}

	var resp CommonResponse
	if err := json.Unmarshal(res, &resp); err != nil {
		return "", err
	}
	if resp.Code != 200 {
		return "", fmt.Errorf("copy failed: %s", resp.Message)
	}

	taskId := ""
	if resp.Data != nil {
		if tid, ok := resp.Data["task_id"].(string); ok {
			taskId = tid
		}
	}
	return taskId, nil
}

// MoveFile 调用 OpenList/AList v3 /api/fs/move 移动文件
func (c *OpenListClient) MoveFile(srcDir, dstDir, name string) (string, error) {
	data := map[string]interface{}{
		"src_dir": srcDir,
		"dst_dir": dstDir,
		"names":   []string{name},
	}
	res, err := c.Post("/api/fs/move", data)
	if err != nil {
		return "", err
	}

	var resp CommonResponse
	if err := json.Unmarshal(res, &resp); err != nil {
		return "", err
	}
	if resp.Code != 200 {
		return "", fmt.Errorf("move failed: %s", resp.Message)
	}

	taskId := ""
	if resp.Data != nil {
		if tid, ok := resp.Data["task_id"].(string); ok {
			taskId = tid
		}
	}
	return taskId, nil
}

// RemoveFile 调用 OpenList/AList v3 /api/fs/remove 删除文件
func (c *OpenListClient) RemoveFile(dir, name string) error {
	data := map[string]interface{}{
		"dir":   dir,
		"names": []string{name},
	}
	res, err := c.Post("/api/fs/remove", data)
	if err != nil {
		return err
	}

	var resp CommonResponse
	if err := json.Unmarshal(res, &resp); err != nil {
		return err
	}
	if resp.Code != 200 {
		return fmt.Errorf("remove failed: %s", resp.Message)
	}
	return nil
}

func (c *OpenListClient) GetChildPath(path string, speed int) (*utils.PageList[string], error) {
	res, err := c.GetChildPathRaw(path, speed)
	if err != nil {
		return nil, err
	}

	var r ChildPathRawResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return nil, err
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

func (c *OpenListClient) Get(api string, params map[string]string) ([]byte, error) {
	url := strings.Trim(c.Engine.Url, "/") + api
	token := c.Engine.Token

	ro := &utils.RequestOption{
		Headers: map[string]string{
			"Authorization": token,
		},
		Timeout: 30 * time.Second,
	}
	return utils.Get(url, ro)
}
