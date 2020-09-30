package db

import (
	"errors"
	"fmt"
	"imageshow/models"
	"sort"
)

// AppendImage 增加图片
func AppendImage(image models.Image) error {
	// 通过id定位数据库，这里使用主库
	// 因为id里有Url的DNA，
	// 所以，其实是通过Url定位数据库
	db, err := getMasterDB(image.Id)
	if err != nil {
		return err
	}

	// 判断图片是否存在
	count := 0
	strSQL := `SELECT COUNT(0) FROM t_images WHERE md5=?;`
	db.QueryRow(strSQL, image.MD5).Scan(&count)
	if count > 0 {
		return errors.New("记录已存在")
	}

	// 增加图片
	strSQL = `INSERT INTO t_images(id,md5,title,url,comment_count,like_count,created_time) VALUE(?,?,?,?,?,?,?);`
	_, err = db.Exec(strSQL, image.Id, image.MD5, image.Title, image.Url, image.Comment_Count, image.Like_Count, image.Created_Time)

	return err
}

// GetImage 获取一张图片的信息
func GetImage(id int64) (*models.Image, error) {
	// 通过id定位数据库，这里使用从库
	db, err := getSlaveDB(id)
	if err != nil {
		return nil, err
	}

	img := &models.Image{}
	strSQL := "SELECT id,url,title,md5,comment_count,like_count,created_time FROM t_images WHERE id=?"
	err = db.QueryRow(strSQL, id).Scan(&img.Id, &img.Url, &img.Title, &img.MD5, &img.Comment_Count, &img.Like_Count, &img.Created_Time)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// GetImages 图片分页，这里使用“禁止跳页法”分页
func GetImages(action string, lastMaxId, lastMinId int64) (models.ImageList, error) {
	var list models.ImageList

	pageSize := 10
	strSQL := fmt.Sprintf("SELECT id,url,title,md5,comment_count,like_count,created_time FROM "+
		"t_images ORDER BY id DESC LIMIT %d;", pageSize)
	if action == "prev" {
		strSQL = fmt.Sprintf("SELECT id,url,title,md5,comment_count,like_count,created_time FROM "+
			"t_images WHERE id>%d ORDER BY id LIMIT %d;", lastMaxId, pageSize)
	} else if action == "next" {
		strSQL = fmt.Sprintf("SELECT id,url,title,md5,comment_count,like_count,created_time FROM "+
			"t_images WHERE id<%d ORDER BY id DESC LIMIT %d;", lastMinId, pageSize)
	}

	for _, db := range slaveDB {
		rows, err := db.Query(strSQL)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var info models.Image
			rows.Scan(&info.Id, &info.Url, &info.Title, &info.MD5, &info.Comment_Count, &info.Like_Count, &info.Created_Time)

			list = append(list, info)
		}
		rows.Close()
	}

	// 重新排序，取pageSize条
	sort.Stable(list)
	if len(list) < pageSize {
		pageSize = len(list)
	}

	if action == "prev" {
		return list[len(list)-pageSize:], nil
	} else {
		return list[:pageSize], nil
	}
}

// UpdateLike 更新点赞次数
func UpdateLike(uid, imgid int64) error {
	if uid <= 0 || imgid <= 0 {
		return errors.New("参数不正确")
	}

	// 通过imgid的基因分库
	// 把点赞的数据和图片分布到同一个数据库
	// 仍然可以使用mysql的事务
	db, err := getMasterDB(imgid)
	if err != nil {
		return err
	}

	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	count := 0
	err = tx.QueryRow("SELECT COUNT(0) FROM t_likes WHERE uid=? AND imgid=?;", uid, imgid).Scan(&count)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 如果没有点赞过，插入到t_likes表，并且t_images表加1
	strSQL1 := "INSERT INTO t_likes(uid,imgid) Values(?,?);"
	strSQL2 := "UPDATE t_images SET like_count = like_count + ? WHERE id=?;"
	total := 1
	if count > 0 {
		// 如果已经点赞过，删除t_likes表，并且t_images表减1
		strSQL1 = "DELETE FROM t_likes WHERE uid=? AND imgid=?;"
		total = -1
	}
	_, err = tx.Exec(strSQL1, uid, imgid)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(strSQL2, total, imgid)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}

	return err
}

// GetComments 获取图片的所有评论，因为评论和图片分配在同一个数据库，可以使用传统的分页方法
func GetComments(imgid int64, pageIndex int) ([]models.Comment, error) {
	db, err := getSlaveDB(imgid)
	if err != nil {
		return nil, err
	}

	var list []models.Comment

	pageSize := 10
	strSQL := "SELECT id,uid,imgid,username,content,created_time FROM t_comments WHERE imgid=? ORDER BY id DESC LIMIT ?,?;"
	rows, err := db.Query(strSQL, imgid, (pageIndex-1)*pageSize, pageSize)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var info models.Comment
		rows.Scan(&info.Id, &info.UId, &info.ImgId, &info.UserName, &info.Content, &info.Created_Time)

		list = append(list, info)
	}
	rows.Close()

	return list, nil
}

// AppendComment 增加图片的评论
func AppendComment(comment models.Comment) error {
	// 通过id定位数据库，这里使用主库
	// 因为id里有imgid的DNA，
	// 所以，其实是通过imgid定位数据库
	db, err := getMasterDB(comment.Id)
	if err != nil {
		return err
	}

	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 增加评论
	strSQL := `INSERT INTO t_comments(id,uid,imgid,username,content,created_time) VALUE(?,?,?,?,?,?);`
	_, err = tx.Exec(strSQL, comment.Id, comment.UId, comment.ImgId, comment.UserName, comment.Content, comment.Created_Time)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 评论次数加1
	strSQL = `UPDATE t_images SET comment_count = comment_count + 1 WHERE id=?;`
	_, err = tx.Exec(strSQL, comment.ImgId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return err
}
