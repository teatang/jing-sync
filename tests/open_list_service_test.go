package tests

import (
	"jing-sync/services"
	"jing-sync/boot"
	"testing"
)


func TestOpenListClientPost(t *testing.T) {
	ol := services.NewOpenListClient("1", boot.GetDB())
	data := map[string]interface{}{
		"path": "/",
		"refresh": true,
		"page": 1,
		"per_page": 100,
	}

	res, _ := ol.Post("/api/fs/list", data)

	t.Error(string(res))
}