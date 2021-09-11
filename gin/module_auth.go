package d_gin

import (
	"context"
	"errors"
	d "github.com/Etpmls/devtool"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var Auth auth

type auth struct {
	Optional optionalAuth
	enable   bool
}

type optionalAuth struct {
	AutoMigrate []interface{}
	SkipInsertingInitializationData bool	// 跳过插入初始化权限相关数据
	TokenExpirationTime time.Duration	// Token过期时间，单位s（秒）
}

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username string     `gorm:"unique;notnull"`
	Password string     `gorm:"notnull"`
	Avatar   Attachment `gorm:"polymorphic:Owner;polymorphicValue:user-avatar"`
	Roles    []Role     `gorm:"many2many:role_users;"`
}

type Role struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name string
	Remark string
	Users []User             `gorm:"many2many:role_users;"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name string
	Method string
	Path	string
	Remark string
	Roles []Role `gorm:"many2many:role_permissions;"`
}

type Attachment struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Path string	`gorm:"type:varchar(500)"`
	OwnerID uint
	OwnerType string
}

const (
	CacheMenuGetAll = "MenuGetAll"
	CacheUserGetCurrent = "UserGetCurrent"
)

func (this *auth) Init() error {
	// 检查数据库是否初始化
	if d.DatabaseClient == nil || !d.Database.GetEnabledStatus() {
		return errors.New("the database is not initialized")
	}
	// 检查日志是否初始化
	if !d.Log.GetEnabledStatus() {
		return errors.New("the log is not initialized")
	}
	// 检查验证器是否初始化
	if !d.Validator.GetEnabledStatus() {
		return errors.New("the validator is not initialized")
	}

	// 设置默认处理方式

	if this.Optional.TokenExpirationTime == 0 {
		this.Optional.TokenExpirationTime = 43200
	}

	err := d.DatabaseClient.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, &Role{}, &Permission{}, &Attachment{})
	if err != nil {
		return err
	}

	// 如果User表找不到数据，则自动填充
	if !this.Optional.SkipInsertingInitializationData {
		result := d.DatabaseClient.First(&User{})
		if result.RowsAffected == 0 {
			err = this.InsertBasicDataToDatabase()
			if err != nil {
				return err
			}
		}
	}

	this.enable = true
	return nil
}

// 获取启动的状态
func (this *auth) GetEnabledStatus() bool {
	return this.enable
}

// 无验证码模式用户登录
type UserLoginRequest struct {
	Username string `json:"username" validate:"required,max=255"`
	Password string `json:"password" validate:"required,max=255"`
}
func (this *auth) UserLogin(c *gin.Context, json *UserLoginRequest, translator ut.Translator) (string, error) {
	err := Validate(c, json, translator)
	if err != nil {
		return "", err
	}

	usrID, usrName, err := this.UserVerify(json.Username, json.Password)
	if err != nil {
		return "", err
	}

	//JWT
	token, err := this.TokenGenerate(usrID, usrName)
	if err != nil {
		return "", err
	}

	return token, nil
}

// 创建用户
type UserCreateRequest struct {
	ID        uint `json:"-"`
	CreatedAt time.Time	`json:"-"`
	UpdatedAt time.Time	`json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	Username string `json:"username" validate:"required,max=50"`
	Password string `json:"password" validate:"required,max=50"`
	Roles []Role `gorm:"many2many:role_users" json:"roles"`
}
func (this *auth) UserCreate(c *gin.Context, json *UserCreateRequest, translator ut.Translator) error {
	err := Validate(c, json, translator)
	if err != nil {
		return err
	}

	// ValidatorClient Username unique
	var count int64
	if err := d.DatabaseClient.Model(&User{}).Where("username = ?", json.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("user already exists")
	}

	// 格式转化
	type User UserCreateRequest
	u := User(*json)

	// 创建数据
	err = d.DatabaseClient.Transaction(func(tx *gorm.DB) error {

		// Bcrypt Password
		u.Password, err = d.BcryptPassword(u.Password)
		if err != nil {
			logrus.Error(err)
			return err
		}

		result := tx.Create(&u)
		if result.Error != nil {
			logrus.Error(err)
			return result.Error
		}

		return nil
	})

	// 更改原值
	*json = UserCreateRequest(u)

	return err
}

// 修改用户
type UserEditRequest struct {
	ID uint `json:"id" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"max=50"`
	Roles []Role `gorm:"many2many:role_users" json:"roles"`
}
func (this *auth) UserEdit(c *gin.Context, json *UserEditRequest, translator ut.Translator) error {
	err := Validate(c, json, translator)
	if err != nil {
		return err
	}

	// Find User
	var form User
	d.DatabaseClient.First(&form, json.ID)

	// If user set new password
	if len(json.Password) > 0 {
		form.Password, err = d.BcryptPassword(json.Password)
		if err != nil {
			return err
		}
	}

	form.Username = json.Username	// Username

	err = d.DatabaseClient.Transaction(func(tx *gorm.DB) error {
		// 替换关联
		err = tx.Model(&User{ID:form.ID}).Association("Roles").Replace(json.Roles)
		if err != nil {
			logrus.Error(err)
			return err
		}
		// 修改数据
		result := tx.Save(&form)
		if result.Error != nil {
			logrus.Error(err)
			return result.Error
		}

		return nil
	})

	if d.Cache.GetEnabledStatus() {
		d.CacheClient.HDel(context.Background(), CacheUserGetCurrent, strconv.Itoa(int(json.ID))).Err()
	}

	return err

}

// Delete users (allow multiple deletions at the same time)
// 删除用户（允许同时删除多个）
type UserDeleteRequest struct {
	Users []User `json:"users" validate:"required,min=1"`
}
func (this *auth) UserDelete(c *gin.Context, json *UserDeleteRequest, translator ut.Translator) error {
	err := Validate(c, json, translator)
	if err != nil {
		return err
	}

	var ids []uint
	for _, v := range json.Users {
		ids = append(ids, v.ID)
	}

	// Find if admin is included in ids
	// 查找ids中是否包含admin
	b := this.CheckIfOneIsIncludeInIds(ids)
	if b {
		return errors.New("can not include administrator")
	}

	err = d.DatabaseClient.Transaction(func(tx *gorm.DB) error {
		var u []User
		tx.Where("id IN ?", ids).Find(&u)

		// 删除用户
		result := tx.Delete(&u)
		if result.Error != nil {
			logrus.Error(result.Error)
			return result.Error
		}

		// 删除关联
		err = tx.Model(&u).Association("Roles").Clear()
		if err != nil {
			logrus.Error(err)
			return err
		}

		return nil
	})

	if d.Cache.GetEnabledStatus() {
		var tmp []string
		for _, v := range ids {
			tmp = append(tmp, strconv.Itoa(int(v)))
		}
		d.CacheClient.HDel(context.Background(), CacheUserGetCurrent, strings.Join(tmp, " ")).Err()
	}

	return err
}

// Update user information
// 更新用户信息
type UserUpdateInformationRequest struct {
	ID        uint	`json:"-"`
	CreatedAt time.Time	`json:"-"`
	UpdatedAt time.Time	`json:"-"`
	DeletedAt gorm.DeletedAt	`json:"-"`
	Password string `json:"password" validate:"omitempty,min=6,max=50"`
	Avatar Attachment	`gorm:"polymorphic:Owner;polymorphicValue:user-avatar" json:"avatar"`
}
func (this *auth) UserUpdateInformation(c *gin.Context, json *UserUpdateInformationRequest, translator ut.Translator) error {
	err := Validate(c, json, translator)
	if err != nil {
		return err
	}

	// Get current user id
	id, err := this.GetUserIdByRequest(c)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = d.DatabaseClient.Transaction(func(tx *gorm.DB) error {
		// 如果表单包含缩略图，
		if len(json.Avatar.Path) > 0 {
			// 1.删除同名缓存
			result := tx.Unscoped().Where("path = ?", json.Avatar.Path).Delete(Attachment{})
			if result.Error != nil {
				logrus.Error(result.Error)
				return result.Error
			}
		}

		// 2.删除历史avatar
		var old Attachment
		result2 := tx.Where("owner_id = ?", id).Where("owner_type = ?", "user-avatar").First(&old)
		// 如果找到记录则删除
		if result2.RowsAffected > 0 {
			// 根据Path删除附件
			err := this.AttachmentBatchDelete([]string{old.Path})
			if err != nil {
				logrus.Error(err)
				return err
			}
		}
		// 3.新增avatar
		err := tx.Model(&User{ID: id}).Association("Avatar").Replace(&Attachment{Path:json.Avatar.Path})
		if err != nil {
			logrus.Error(err)
			return err
		}

		// 4.更新
		// Update password if exists
		if len(json.Password) > 0 {
			json.Password, err = d.BcryptPassword(json.Password)
		}

		result := tx.Debug().Model(&User{ID: id}).Updates(json)
		if result.Error != nil {
			logrus.Error(result.Error)
			return result.Error
		}
		return nil
	})

	if d.Cache.GetEnabledStatus() {
		d.CacheClient.HDel(context.Background(), CacheUserGetCurrent, strconv.Itoa(int(id))).Err()
	}

	return err
}

// 插入初始化数据到数据库
func (this *auth)InsertBasicDataToDatabase() error {
	err := d.DatabaseClient.Transaction(func(tx *gorm.DB) error {
		// Create Role
		role := Role{
			Name:        "Administrator",
			Remark: "System Administrator",
		}
		if err := d.DatabaseClient.Debug().Create(&role).Error; err != nil {
			return err
		}


		// Create User
		user := User{
			Username: "admin",
			Password: "$2a$10$yNoJrsN7mrtHzUyvm6s8KOwHrnkkGmqcRJvcieQKItIfQNwyzqfMy",
			Roles: []Role{
				{
					ID:1,
				},
			},
		}
		if err := d.DatabaseClient.Debug().Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&user).Error; err != nil {
			return err
		}

		// Create Permission
		permission := []Permission{
			{
				Name: "View User",
				Method: "GET",
				Path: "/api/*/user/getAll",
				Remark: "View user list",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "Create User",
				Method: "POST",
				Path: "/api/*/user/create",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "Edit User",
				Method: "PUT",
				Path: "/api/*/user/edit",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "Delete User",
				Method: "DELETE",
				Path: "/api/*/user/delete",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "View Role",
				Method: "GET",
				Path: "/api/*/role/getAll",
				Remark: "View role list",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "Create Role",
				Method: "POST",
				Path: "/api/*/role/create",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "Edit Role",
				Method: "PUT",
				Path: "/api/*/role/edit",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "Delete Role",
				Method: "DELETE",
				Path: "/api/*/role/delete",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "View Permission",
				Method: "GET",
				Path: "/api/*/permission/getAll",
				Remark: "View permission list",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "Create Permission",
				Method: "POST",
				Path: "/api/*/permission/create",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "Edit Permission",
				Method: "PUT",
				Path: "/api/*/permission/edit",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "Delete Permission",
				Method: "DELETE",
				Path: "/api/*/permission/delete",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "Create/Edit Menu",
				Method: "POST",
				Path: "/api/*/menu/create",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "Clear Cache",
				Method: "GET",
				Path: "/api/*/setting/clearCache",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},
			{
				Name: "Disk Cleanup",
				Method: "GET",
				Path: "/api/*/setting/diskCleanup",
				Roles: []Role{
					{
						ID:1,
					},
				},
			},

		}
		if err := d.DatabaseClient.Debug().Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&permission).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

// Get token by header Or query
// 从header或query中获取token
func (this *auth) GetTokenByRequest(c *gin.Context) (token string, err error) {
	// Get Query Token
	token, b := c.GetQuery("token")
	if b {
		return token, err
	}

	// Get Header Token
	token = c.GetHeader("X-Token")
	if len(token) != 0 {
		return token, err
	}

	token = c.GetHeader("Token")
	if len(token) != 0 {
		return token, err
	}

	logrus.Error("token acquisition failed")
	return token, errors.New("token acquisition failed")
}
// Only check if the token exists
// 只检查token是否存在
func (this *auth) MiddlewareCheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Get Token
		// 获取token
		token, err := this.GetTokenByRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Get Claims
		// 获取Claims
		_, err = d.Token.Parse(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Next()
		return
	}
}
// Check whether the token exists, check whether the user's role has permissions
// 检查token是否存在，检查用户所在角色是否拥有权限
func (this *auth) MiddlewareCheckPermission() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Get Token
		// 获取token
		token, err := this.GetTokenByRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Get Claims
		// 获取Claims
		tk, err := d.Token.Parse(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 判断所属角色是否有相应的权限
		if claims,ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
			if userId, ok := claims["jti"].(string); ok {
				b, err := this.permissionCheck(c, userId)
				if err == nil && b {
					c.Next()
					return
				}
			}
		}

		// 没权限就是401
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No Permission"})
		c.Abort()
		return
	}
}
// 权限检查逻辑
func (this *auth) permissionCheck(c *gin.Context, idStr string) (b bool, err error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return b, err
	}

	// 1.获取用户ID
	var u User
	d.DatabaseClient.Preload("Roles").First(&u, id)
	var ids []uint
	for _, v := range u.Roles {
		// 如果为管理员组
		if v.ID == 1 {
			return true, nil
		}
		ids = append(ids, v.ID)
	}
	// 获取角色相关权限
	var r []Role
	d.DatabaseClient.Preload("Permissions").Where(ids).Find(&r)

	// 获取当前URL Path
	tmpUri, err := url.Parse(c.Request.RequestURI)
	if err != nil {
		return b, err
	}
	uri := tmpUri.Path

	// Determine whether there is a request permission
	// 判断是否有请求权限
	for _, v := range r {
		for _, subv := range v.Permissions {

			// define an empty slice
			// 定义一个空切片
			var mtd = []string{}
			mtd = strings.Split(subv.Method, ",")

			// Path comparison
			// 路径对比
			b, _ := filepath.Match(subv.Path, uri)
			if b {

				// Method comparison
				// 方法对比
				for _, mtdv := range mtd {
					// If it is ALL, return the permission verification success directly
					// 如果是ALL直接返回权限验证成功
					if strings.ToUpper(mtdv) == "ALL" {
						return true, nil
					}
					// If the method is the same as the current request, return the verification success
					// 如果与当前请求方法相同，返回验证成功
					if mtdv == c.Request.Method {
						return true, nil
					}
				}

			}
		}
	}

	return false, err
}

// Check if 1(Admin User/Admin Role) is included in ids
// 查看ids中是否包含admin
func (this *auth) CheckIfOneIsIncludeInIds(ids []uint) bool {
	for _, v := range ids {
		if v == 1 {
			return true
		}
	}

	return false
}

// 根据请求信息获取用户id
func (this *auth) GetUserIdByRequest(c *gin.Context) (uint, error) {
	token, err := this.GetTokenByRequest(c)
	if err != nil {
		return 0, err
	}
	id, err := d.Token.GetJwtIdByToken(token)
	if err != nil {
		return 0, err
	}
	return uint(id), err
}

// Verify user logic
// 验证用户逻辑
func (this *auth) UserVerify(username string, password string) (id uint, unm string, err error) {
	//Search User
	var user User
	d.DatabaseClient.Where("username = ?", username).First(&user)
	if !(user.ID > 0) {
		return id, unm, errors.New("The username does not exist!")
	}

	//Password is wrong
	b, err := d.VerifyPassword(password, user.Password)
	if err != nil || !b {
		return id, unm, errors.New("Verification failed or wrong password!")
	}

	return user.ID, user.Username, err
}

func (this *auth) TokenGenerate(userId uint, username string) (string, error) {
	return d.Token.Create(&jwt.StandardClaims{
		Id: strconv.Itoa(int(userId)),	// 用户ID
		ExpiresAt: time.Now().Add(time.Second * this.Optional.TokenExpirationTime).Unix(),	// 过期时间 - 12个小时
		Subject:    username,	// 发行者
	})
}

// Validate Path is a file in storage/upload
// 验证路径是否在storage/upload中
func (this *auth) AttachmentValidatePath(path string) error {
	const upload_path = "storage/upload/"
	// 截取前十五个字符判断和Path是否相同
	if len(path) <= len(upload_path) || !strings.Contains(path[:len(upload_path)], upload_path) {
		logrus.Error("illegal request path")
		return  errors.New("Illegal request path!")
	}
	f, err := os.Stat(path)
	if err != nil {
		logrus.Error(err)
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	// 如果文件是目录
	if f.IsDir() {
		logrus.Error("cannot delete directory")
		return errors.New("Cannot delete directory!")
	}
	return nil
}

// Batch delete any type of files in storage/upload/
// 批量删除storage/upload/中的任何类型文件
func (this *auth) AttachmentBatchDelete(s []string) (err error) {
	// 第一遍遍历先验证
	for _, v := range s {
		// Validate If a File
		err = this.AttachmentValidatePath(v)
		if err != nil {
			return err
		}
	}
	// 第二遍遍历再删除
	for _, v := range s {
		// Delete Image
		_ = os.Remove(v)
	}

	// Delete Database
	if err = d.DatabaseClient.Unscoped().Where("path IN (?)", s).Delete(Attachment{}).Error; err != nil {
		logrus.Error(err)
		return err
	}

	return err
}


