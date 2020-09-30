package db

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"imageshow/models"
	"imageshow/utils"
	"strings"
)

// GetUser 获取用户信息
func GetUser(name string) (*models.User, error) {
	// 通过name定位数据库，从库
	name = strings.ToLower(name)
	nameCRC32 := int64(utils.CRC32(name))
	db, err := getSlaveDB(nameCRC32)
	if err != nil {
		return nil, err
	}

	// 增加用户
	user := &models.User{}
	strSQL := `SELECT id,name,password,remark,created_time FROM t_users WHERE name=?;`
	err = db.QueryRow(strSQL, name).Scan(&user.Id, &user.Name, &user.Password, &user.Remark, &user.Created_Time)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AppendUser 增加用户
func AppendUser(user models.User) error {
	// 通过id定位数据库，主库
	// 因为id里有name的DNA，
	// 所以，其实是通过name定位数据库
	db, err := getMasterDB(user.Id)
	if err != nil {
		return err
	}

	// 判断用户名是否存在
	strSQL := `SELECT count(0) FROM t_users WHERE name=?;`
	count := 0
	db.QueryRow(strSQL, user.Name).Scan(&count)
	if count > 0 {
		return errors.New("记录已存在")
	}

	// 增加用户
	strSQL = `INSERT INTO t_users(id,name,password,remark,created_time) VALUE(?,?,?,?,?);`
	_, err = db.Exec(strSQL, user.Id, user.Name, user.Password, user.Remark, user.Created_Time)

	return err
}
