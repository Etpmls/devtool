package d

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)


var (
	Database database
	DB *gorm.DB
)

type database struct {
	Host string
	Port string
	DBName string
	Username string
	Password string
	Optional optionalGorm
}

type optionalGorm struct {
	TablePrefix string
	AutoMigrate []interface{}
	Charset string
	ParseTime bool
	Location string
}

// 初始化数据库
// https://gorm.io/docs/connecting_to_the_database.html
func (this *database) Init() error {
	if this.Optional.Charset == "" {
		this.Optional.Charset = "utf8mb4"
	}

	var parseTime string
	if this.Optional.ParseTime {
		parseTime = "True"
	} else {
		parseTime = "False"
	}

	if this.Optional.Location == "" {
		this.Optional.Location = "Local"
	}

	var err error
	dsn := this.Username + ":"+ this.Password + "@tcp(" + this.Host + ":" + this.Port + ")/" + this.DBName + "?charset=" + this.Optional.Charset + "&parseTime=" + parseTime + "&loc=" + this.Optional.Location
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   this.Optional.TablePrefix,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// 创建数据库表
// https://gorm.io/docs/migration.html
func (this *database) Migration() error {
	// Add table suffix when creating tables
	err := DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(this.Optional.AutoMigrate...)
	if err != nil {
		return err
	}

	return nil
}