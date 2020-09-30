package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"imageshow/models"
	"path"
)

// GetSpiderImages 获取没有审核过的图片
func GetSpiderImages() ([]models.SpiderImage, error) {
	var infos []models.SpiderImage
	sql := `SELECT id,url,title,md5 FROM t_images WHERE state=0 ORDER BY ID DESC LIMIT 10;`
	rows, err := spiderDB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var info models.SpiderImage
		err := rows.Scan(&info.Id, &info.Url, &info.Title, &info.MD5)
		if err != nil {
			continue
		}

		ext := path.Ext(info.Url)
		info.Url = fmt.Sprintf("./static/images/%c/%c/%s%s", info.MD5[0], info.MD5[1], info.MD5, ext)
		infos = append(infos, info)
	}

	return infos, nil
}

// GetSpiderImage 获取一张图片的信息
func GetSpiderImage(id int) (*models.SpiderImage, error) {
	sql := `SELECT id,url,title,md5 FROM t_images WHERE id=?;`
	info := &models.SpiderImage{}
	err := spiderDB.QueryRow(sql, id).Scan(&info.Id, &info.Url, &info.Title, &info.MD5)

	return info, err
}

// CheckImg 图片审核通过
func CheckImg(id int) error {
	strSQL := `UPDATE t_images SET state=1 WHERE id=? AND state=0;`
	_, err := spiderDB.Exec(strSQL, id)

	return err
}
