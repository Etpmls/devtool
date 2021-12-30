package d_gin

import (
	"context"
	"encoding/json"
	d "github.com/etpmls/devtool"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

const (
	CacheMenuGet = "MenuGet"
)

// Create Menu
// 创建菜单
type MenuCreateRequest struct {
	Menu string `json:"menu" binding:"required"`
}
func MenuCreate(c *gin.Context, translator ut.Translator) error {
	var j MenuCreateRequest
	err := Validate(c, &j, translator)
	if err != nil {
		logrus.Error(err)
		return err
	}

	// 创建菜单
	err = d.MenuCreate(j.Menu)
	if err != nil {
		return err
	}

	// Delete Cache
	// 删除缓存
	if d.Cache.GetEnabledStatus() {
		d.CacheClient.Del(context.Background(), CacheMenuGet).Err()
	}

	return nil
}

// Get all menu
// 获取全部菜单
func MenuGet() (interface{}, error) {
	if d.Cache.GetEnabledStatus() {
		return menuGetCache()
	} else {
		return menuGetNoCache()
	}
}
func menuGetCache() (interface{}, error) {
	// Get the menu from cache
	// 从缓存中获取menu
	ctx, err := d.CacheClient.Get(context.Background(), CacheMenuGet).Result()
	if err != nil {
		if err == redis.Nil {
			return menuGetNoCache()
		}
		return nil, err
	}

	var j interface{}
	err = json.Unmarshal([]byte(ctx), &j)
	if err != nil {
		_ = d.CacheClient.Del(context.Background(), CacheMenuGet).Err()
		return nil, err
	}
	return j, nil
}
func menuGetNoCache() (interface{}, error) {
	ctx, err := d.MenuGet()
	if err != nil {
		return nil, err
	}

	// Save menu
	// 储存菜单
	if d.Cache.GetEnabledStatus() {
		_ = d.CacheClient.Set(context.Background(), CacheMenuGet, string(ctx), 0).Err()
	}

	var j interface{}
	err = json.Unmarshal(ctx, &j)
	if err != nil {
		return nil, err
	}

	return j, nil
}