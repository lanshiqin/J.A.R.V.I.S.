package orm

import (
	"github.com/jinzhu/gorm"
	"testing"
	"time"
)

const url = "jarvis_user:d82aad4c-5ad4-42dc-9b94-4e2fbb443ea0@/jarvis_db?charset=utf8&parseTime=True&loc=Local"

var mysql MySQL

func init() {
	mysql = MySQL{Url: url}
	mysql.InitMySQL()
}

func TestMySQL_InitMySQL(t *testing.T) {
	mysql = MySQL{Url: url}
	mysql.InitMySQL()
}

type UserInfo struct {
	gorm.Model
	Birthday time.Time
	Age      int
	Name     string `gorm:"not null"`
	Desc     string `gorm:"size:500"` // string默认长度为255, 使用tag设置为指定长度
}

// 指定映射的表名为user_info，默认表名为结构体名称加`s`后缀 user_infos
func (UserInfo) TableName() string {
	return "user_info"
}

func TestMySQL_CreateTable(t *testing.T) {
	// 检查模型`User`表是否存在
	isExist := mysql.DB.HasTable(&UserInfo{})
	if isExist {
		t.Log("表已经存在，删除表")
		mysql.DB.DropTable(&UserInfo{})
		t.Log("创建表")
		mysql.DB.CreateTable(&UserInfo{})
	} else {
		t.Log("表不存在，创建表")
		mysql.DB.CreateTable(&UserInfo{})
	}
}

// 插入数据
func TestMySQL_Insert(t *testing.T) {
	user := UserInfo{Name: "J.A.R.V.I.S.", Age: 1, Birthday: time.Now(), Desc: "A Personal Assistant"}
	mysql.DB.Create(&user)
	mysql.DB.Create(&UserInfo{Name: "Friday", Age: 0, Birthday: time.Now(), Desc: "AI"})
	mysql.DB.Create(&UserInfo{Name: "Tony Stark", Age: 35, Birthday: time.Now(), Desc: "STARK INDUSTRIES CEO"})
	mysql.DB.Create(&UserInfo{Name: "Iron Man", Age: 35, Birthday: time.Now(), Desc: "Super Hero"})
}

// 查询数据
func TestMySQL_Query(t *testing.T) {

	var firstUser UserInfo
	// 获取第一条记录, 按主键排序
	mysql.DB.First(&firstUser)
	//// SELECT * FROM user_info ORDER BY id LIMIT 1;
	t.Log(firstUser)

	var lastUser UserInfo
	// 获取最后一条记录, 按主键排序
	mysql.DB.Last(&lastUser)
	//// SELECT * FROM user_info ORDER BY id DESC LIMIT 1;
	t.Log(lastUser)

	var userList []UserInfo
	// 获取所有记录
	mysql.DB.Find(&userList)
	//// SELECT * FROM user_info;
	t.Log(userList)

	// 分页获取记录
	mysql.DB.Limit(2).Find(&userList)
	//// SELECT * FROM user_info limit 2;
	t.Log(userList)

	var userByIdObj UserInfo
	mysql.DB.First(&userByIdObj, 10)
	//// SELECT * FROM user_info where ID = 10;
	t.Log(userByIdObj)

}

// 查询数据 WHERE、IN、LIKE、AND
func TestMySQL_Query2(t *testing.T) {

	var firstUser UserInfo
	// 获取第一条记录, 按主键排序
	mysql.DB.Where("name = ?", "J.A.R.V.I.S.").First(&firstUser)
	//// SELECT * FROM user_info WHERE name = 'J.A.R.V.I.S.' ORDER BY id LIMIT 1;
	t.Log(firstUser)

	var userList []UserInfo
	// 获取所有匹配记录
	mysql.DB.Where("name = ?", "J.A.R.V.I.S.").Find(&userList)
	//// SELECT * FROM user_info WHERE name = 'J.A.R.V.I.S.';
	t.Log(userList)

	// IN
	mysql.DB.Where("name in (?)", []string{"Tony Stark", "Iron Man"}).Find(&userList)
	//// SELECT * FROM user_info WHERE name in ("Tony Stark","Iron Man");
	t.Log(userList)

	// LIKE
	mysql.DB.Where("name like ?", "%Stark%").Find(&userList)
	//// SELECT * FROM user_info WHERE name like "%Stark%";
	t.Log(userList)

	// AND
	mysql.DB.Where("name = ? and age >= ?", "Tony Stark", "22").Find(&userList)
	//// SELECT * FROM user_info WHERE name = "Tony Stark" and age >= "22" ;
	t.Log(userList)

}

func TestMySQL_Query3(t *testing.T) {

	var structByUser UserInfo
	// struct作为参数查询
	mysql.DB.Where(&UserInfo{Name: "J.A.R.V.I.S.", Age: 1}).First(&structByUser)
	//// SELECT * FROM user_info WHERE name = 'J.A.R.V.I.S.' and age = 1 ORDER BY id LIMIT 1;
	t.Log(structByUser)

	var mapByUser UserInfo
	// map作为参数查询
	mysql.DB.Where(map[string]interface{}{"name": "J.A.R.V.I.S.", "age": 1}).First(&mapByUser)
	//// SELECT * FROM user_info WHERE name = 'J.A.R.V.I.S.' and age = 1 ORDER BY id LIMIT 1;
	t.Log(mapByUser)
}

// 更新数据
func TestMySQL_Update(t *testing.T) {
	var firstUser UserInfo
	// 获取第一条记录, 按主键排序
	mysql.DB.First(&firstUser)
	//// SELECT * FROM user_info ORDER BY id LIMIT 1;
	t.Log(firstUser)

	// 修改数据
	firstUser.Name = "钢铁侠"
	// 更新数据
	mysql.DB.Save(&firstUser)
	t.Log(firstUser)

	// 根据条件更新
	mysql.DB.Model(&UserInfo{}).Where("name = ?", "J.A.R.V.I.S.").Update("age", 0)

}

// 删除数据
func TestMySQL_Delete(t *testing.T) {

	// 删除所有记录,如果模型有DeletedAt字段，它将自动获得软删除功能！ 那么在调用Delete时不会从数据库中永久删除，而是只将字段DeletedAt的值设置为当前时间
	mysql.DB.Delete(&UserInfo{})

	// 根据条件删除
	mysql.DB.Model(&UserInfo{}).Where("name = ?", "J.A.R.V.I.S.").Delete(&UserInfo{})

}

// 执行原生SQL
func TestMySQL_ExecSQL(t *testing.T) {

	mysql.DB.Exec("INSERT INTO user_info (`name`, `desc`) VALUES (?, ?)", "AI", "Desc Test")

	var userInfoList []UserInfo
	mysql.DB.Raw("SELECT * FROM user_info WHERE name = 'J.A.R.V.I.S.' ORDER BY id LIMIT 1;").Scan(&userInfoList)
	t.Log(userInfoList)
}

// 关闭连接
func TestMySQL_CloseMySQL(t *testing.T) {
	mysql.CloseMySQL()
}
