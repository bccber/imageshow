package models

// User 用户信息
type User struct {
	Id           int64
	Name         string
	Password     string
	Remark       string
	Created_Time int64
}

// Image 图片信息
type Image struct {
	Id            int64
	MD5           string
	Url           string
	Title         string
	Comment_Count int
	Like_Count    int
	Created_Time  int64
}

// ImageList 图片列表
type ImageList []Image

func (img ImageList) Len() int      { return len(img) }
func (img ImageList) Swap(i, j int) { img[i], img[j] = img[j], img[i] }

// 按名Id排序
func (img ImageList) Less(i, j int) bool { return img[i].Id > img[j].Id }

// Comment 评论信息
type Comment struct {
	Id           int64
	UId          int64
	ImgId        int64
	UserName     string
	Content      string
	Created_Time int64
}

// SpiderImage 爬虫库的图片信息
type SpiderImage struct {
	Id    int
	MD5   string
	Url   string
	Title string
	State int
}
