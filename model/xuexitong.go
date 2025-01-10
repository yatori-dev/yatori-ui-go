package model

type XueXiTong struct {
	Id         uint   `gorm:"primaryKey" json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	CourseId   string `json:"courseId"`
	CourseName string `json:"courseName"`
	Status     int    `json:"status"` // 0 未开始 1 队列中 2 进行中 3 已完成 4 已失败 5 已暂停
}

type AddXueXiTongDto struct {
	XueXiTong
}

type ListXueXiTongDto struct {
	XueXiTong
	Page int
	Size int
}
