package dao

// 使用 gorm 框架，https://gorm.io/zh_CN/docs/index.html
import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
)

var (
	DB   *gorm.DB
	err  error
	user []User
)

type User struct {
	ID    uint   `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Name  string `gorm:"type:varchar(32) not null;default:''"`
	Email string `gorm:"type:varchar(64) not null;default:''"`
	Age   uint   `gorm:"type:int(8) not null;default: 18"`
}

func init() {
	// 手动创建数据库 gorm
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: "root:admin@123@tcp(172.16.4.121:4017)/gorm?charset=utf8&parseTime=True&loc=Local", // DSN data source name
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})

	if err != nil {
		//println("建立数据库连接失败：" + err.Error())
		return
	}
	sqlDB, _ := DB.DB()

	//设置数据库连接池参数
	sqlDB.SetMaxOpenConns(100) //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)  //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
	// 创建表
	DB.Migrator().DropTable(&User{})
	DB.Migrator().CreateTable(&User{})

	//var users = []User{{ID: 1, Name: "Lily", Email: "lily@gail.com", Age: 18}, {ID: 2, Name: "Tom", Email: "Tom@gail.com", Age: 19}, {ID: 3, Name: "Jane", Email: "Jane@gail.com", Age: 20}, {ID: 4, Name: "Lisa", Email: "lisa@gail.com", Age: 22}}
	//DB.Create(&users)
}

func GetDB() *gorm.DB {
	return DB
}

func InsertData(chWrite chan bool) {
	db := GetDB()
	defer close(chWrite)
	u := [4]User{{ID: 1, Name: "Lily", Email: "lily@gail.com", Age: 18},
		{ID: 2, Name: "Tom", Email: "Tom@gail.com", Age: 19},
		{ID: 3, Name: "Jane", Email: "Jane@gail.com", Age: 20},
		{ID: 4, Name: "Lisa", Email: "lisa@gail.com", Age: 22}}

	if err := db.Create(&u).Error; err != nil {
		//fmt.Println(err)
		chWrite <- false
		return
	} else {
		chWrite <- true
		return
	}

	//db.Model(&User{}).Create([]map[string]interface{}{
	//	{"Name": "Lily", "Email": "lily@gail.com", "Age": 18},
	//	{"Name": "Tom", "Email": "Tom@gail.com", "Age": 19},
	//	{"Name": "Jane", "Email": "Jane@gail.com", "Age": 20},
	//	{"Name": "Lisa", "Email": "lisa@gail.com", "Age": 22},
	//})
}
func SelectData(chRead chan []User) {
	db := GetDB()
	defer close(chRead)

	result := db.Find(&user)
	if result.Error != nil {
		return
	} else {
		fmt.Printf("user:%#v\n", user)
		chRead <- user
		return
	}

}

func Write2file(ch chan []User) {

	file, err := os.Create("./userinfo.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	result := <-ch
	//fmt.Println(user)
	for _, data := range result {
		splitLine := fmt.Sprintf("%d,%s,%s,%d\n", data.ID, data.Name, data.Email, data.Age)
		_, err := file.WriteString(splitLine)
		if err != nil {
			break
		}
	}

	//defer close(ch)
	defer file.Close()
	return

}
