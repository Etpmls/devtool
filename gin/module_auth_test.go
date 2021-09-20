package d_gin_test

import (
	"context"
	"fmt"
	d "github.com/Etpmls/devtool"
	d_gin "github.com/Etpmls/devtool/gin"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	// 初始化数据库
	var (
		Host = "127.0.0.1"
		Port = "3306"
		DBName = "test"
		Username = "root"
		Password = "123456"
		TablePrefix = "ts_"
	)
	d.Database.Host = Host
	d.Database.Port = Port
	d.Database.DBName = DBName
	d.Database.Username = Username
	d.Database.Password = Password
	d.Database.Optional.TablePrefix = TablePrefix
	err := d.Database.Init()
	if err != nil {
		t.Fatal(err)
		return
	}
	// 初始化日志
	d.Log.Init()
	// 初始化验证器
	err = d.Validator.Init()
	if err != nil {
		t.Fatal(err)
		return
	}
	// 初始化Auth权限模块
	err = d_gin.Auth.Init()
	if err != nil {
		t.Fatal(err)
		return
	}
	// 启动gin
	r := gin.Default()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	r.POST("/user/login", func(c *gin.Context) {
		var u d_gin.UserLoginRequest
		token, err := d_gin.Auth.UserLogin(c, &u, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//Return Token
		resData := make(map[string]string)
		resData["token"] = token
		c.JSON(http.StatusOK, resData)
		fmt.Println(resData)
		QuitTest(srv)
	})
	r.GET("/user/getCurrent", func(c *gin.Context) {
		u, err := d_gin.Auth.UserGetCurrent(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, u)
		fmt.Println(u)
		QuitTest(srv)
	})
	r.POST("/user/create", func(c *gin.Context) {
		var u d_gin.UserCreateRequest
		err := d_gin.Auth.UserCreate(c, &u, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		QuitTest(srv)
	})
	r.PUT("/user/edit", func(c *gin.Context) {
		var u d_gin.UserEditRequest
		err := d_gin.Auth.UserEdit(c, &u, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		QuitTest(srv)
	})
	r.DELETE("/user/delete", func(c *gin.Context) {
		var u d_gin.UserDeleteRequest
		err := d_gin.Auth.UserDelete(c, &u, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		QuitTest(srv)
	})
	r.PUT("/user/updateInformation", func(c *gin.Context) {
		var u d_gin.UserUpdateInformationRequest
		err := d_gin.Auth.UserUpdateInformation(c, &u, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		QuitTest(srv)
	})
	r.POST("/role/create", func(c *gin.Context) {
		var r d_gin.RoleCreateRequest
		err := d_gin.Auth.RoleCreate(c, &r, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		QuitTest(srv)
	})
	r.PUT("/role/edit", func(c *gin.Context) {
		var r d_gin.RoleEditRequest
		err := d_gin.Auth.RoleEdit(c, &r, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		QuitTest(srv)
	})
	r.DELETE("/role/delete", func(c *gin.Context) {
		var r d_gin.RoleDeleteRequest
		err := d_gin.Auth.RoleDelete(c, &r, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		QuitTest(srv)
	})
	r.GET("/role/get", func(c *gin.Context) {
		data, count := d_gin.Auth.RoleGet(c)
		c.JSON(http.StatusOK, gin.H{"data": data, "count": count})
		QuitTest(srv)
	})
	r.POST("/permission/create", func(c *gin.Context) {
		var p d_gin.PermissionCreateRequest
		err := d_gin.Auth.PermissionCreate(c, &p, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		QuitTest(srv)
	})
	r.PUT("/permission/edit", func(c *gin.Context) {
		var p d_gin.PermissionEditRequest
		err := d_gin.Auth.PermissionEdit(c, &p, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		QuitTest(srv)
	})
	r.DELETE("/permission/delete", func(c *gin.Context) {
		var p d_gin.PermissionDeleteRequest
		err := d_gin.Auth.PermissionDelete(c, &p, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		QuitTest(srv)
	})
	r.GET("/permission/get", func(c *gin.Context) {
		data, count := d_gin.Auth.PermissionGet(c)
		c.JSON(http.StatusOK, gin.H{"data": data, "count": count})
		QuitTest(srv)
	})
	srv.ListenAndServe()
	// r.Run() // listen and serve on 0.0.0.0:8080
}

func QuitTest(srv *http.Server)  {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	return
}