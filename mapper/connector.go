package mapper

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
	//a       = 1
	//b       = []int{1, 2}
	mappers = []interface{}{&Comment{}, &Follower{}, &Like{}, &Message{}, &User{}, &Video{}}
)

// InitDBConnector 连接数据库的方法, 在该方法外定义了数据库连接对象"DBConn", 本方法需要在项目启动时运行, 以避免连接数据库失败
// 方法仅有error返回值, 当err不为空时, 需要中止项目运行
func InitDBConnector() (err error) {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:root@tcp(127.0.0.1:3306)/tiktok?charset=utf8mb4&parseTime=True&loc=Local"
	//var err error
	DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
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
