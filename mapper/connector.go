package mapper

import (
	"fmt"
	"os"
	"tiktok/model"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DBConn   *gorm.DB
	dbConfig *DBConfig
	mappers  = []interface{}{&model.Comment{}, &model.Follower{}, &model.Like{}, &model.Message{}, &model.User{}, &model.Video{}}
)

type DBConfig struct {
	Host        string `yaml:"Host"`
	UserName    string `yaml:"Username"`
	PassWord    string `yaml:"Password"`
	DBName      string `yaml:"DBName"`
	Port        string `yaml:"Port"`
	MaxOpenConn int    `yaml:"MaxOpenConn"`
	MaxIdleConn int    `yaml:"MaxIdleConn"`
}

// 从数据库配置文件中获取DSN
func getDSN() string {
	yamlFile, err := os.ReadFile("./config/db.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &dbConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.UserName,
		dbConfig.PassWord,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
	return dsn
}

// InitDBConnector 连接数据库的方法, 在该方法外定义了数据库连接对象"DBConn", 本方法需要在项目启动时运行, 以避免连接数据库失败
// 方法仅有error返回值, 当err不为空时, 需要中止项目运行
func InitDBConnector() (err error) {
	dsn := getDSN()
	DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 配置连接池
	sqlDB, _ := DBConn.DB()
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConn)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)

	return nil
}

// 从数据库配置文件中获取DSN
func getDSNSupportTest() string {
	//yamlFile, err := os.ReadFile("./config/db.yaml")
	yamlFile, err := os.ReadFile("E:\\OneDrive\\MyCode\\Go\\TikTok\\config\\db.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &dbConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.UserName,
		dbConfig.PassWord,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
	return dsn
}

func InitDBConnectorSupportTest() (err error) {
	dsn := getDSNSupportTest()
	DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 配置连接池
	sqlDB, _ := DBConn.DB()
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConn)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)

	return nil
}

func createDateTable(t interface{}) error {
	if !DBConn.Migrator().HasTable(t) {
		err := DBConn.AutoMigrate(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateAllTable() (err error) {
	for _, mapper := range mappers {
		err = createDateTable(mapper)
		if err != nil {
			//这里需要打一下日志, 白天需要查看一下有哪些合适的日志库
			return err
		}
	}
	return nil
}
