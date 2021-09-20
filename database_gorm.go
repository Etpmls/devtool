package d

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)


var (
	Database       database
	DatabaseClient *gorm.DB
)

type database struct {
	Host string
	Port string
	DBName string
	Username string
	Password string
	Optional optionalGorm
	enable bool
}

type optionalGorm struct {
	TablePrefix    string
	AutoMigrate    []interface{}
	Charset        string
	DoNotParseTime bool
	Location       string
	FuzzySearch 	string
}

// 初始化数据库
// https://gorm.io/docs/connecting_to_the_database.html
func (this *database) Init() error {
	if this.Optional.Charset == "" {
		this.Optional.Charset = "utf8mb4"
	}

	var parseTime string
	if this.Optional.DoNotParseTime {
		parseTime = "False"
	} else {
		parseTime = "True"
	}

	if this.Optional.Location == "" {
		this.Optional.Location = "Local"
	}

	if this.Optional.FuzzySearch == "" {
		this.Optional.FuzzySearch = "LIKE"
	}

	var err error
	dsn := this.Username + ":"+ this.Password + "@tcp(" + this.Host + ":" + this.Port + ")/" + this.DBName + "?charset=" + this.Optional.Charset + "&parseTime=" + parseTime + "&loc=" + this.Optional.Location
	DatabaseClient, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   this.Optional.TablePrefix,
		},
	})
	if err != nil {
		return err
	}
	this.enable = true
	return nil
}

// 获取启动的状态
func (this *database) GetEnabledStatus() bool {
	return this.enable
}

// 创建数据库表
// https://gorm.io/docs/migration.html
func (this *database) Migration() error {
	// Add table suffix when creating tables
	err := DatabaseClient.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(this.Optional.AutoMigrate...)
	if err != nil {
		return err
	}

	return nil
}