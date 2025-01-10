package service

import (
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"github.com/yatori-dev/yatori-go-core/aggregation/xuexitong"
	"github.com/yatori-dev/yatori-go-core/api/entity"
	xuexitongApi "github.com/yatori-dev/yatori-go-core/api/xuexitong"
	"log"
	"strconv"
	"time"
	"yatory-gui-wails3/config"
	"yatory-gui-wails3/model"
	"yatory-gui-wails3/utils"
)

type Xuexitong struct {
	_map map[string][]model.XueXiTong
}

// NewXuexitong 初始化，并异步启动数据库中未完成的任务
func NewXuexitong() *Xuexitong {
	// TODO: 启动数据库中未完成的任务
	c := &Xuexitong{
		_map: make(map[string][]model.XueXiTong),
	}
	// 获取所有未完成的任务
	list := c.findAllTaskNotFinish()
	// 创建任务队列
	for _, task := range list {
		c.joinPlay(task.Username, &task)
	}
	return c
}

// QueryCourse 查询课程
func (c *Xuexitong) QueryCourse(dto *model.AddXueXiTongDto) *model.Result {
	// 参数校验
	if dto.Username == "" || dto.Password == "" {
		return &model.Result{
			Code:    400,
			Data:    nil,
			Message: "参数错误",
		}
	}
	cache := &xuexitongApi.XueXiTUserCache{Name: dto.Username, Password: dto.Password}
	// TODO 处理登录异常（等完善前后端通讯）
	err := xuexitong.XueXiTLoginAction(cache) // 登录
	utils.OkOrPanic(err)
	// TODO 处理获取课程异常
	action, err := xuexitong.XueXiTPullCourseAction(cache)
	utils.OkOrPanic(err)
	return &model.Result{
		Code:    200,
		Data:    action,
		Message: "success",
	}
}

// Add 添加任务
func (c *Xuexitong) Add(dto *model.AddXueXiTongDto) {
	println("添加学习通任务")
	data := &model.XueXiTong{
		CourseId:   dto.CourseId,
		CourseName: dto.CourseName,
		Username:   dto.Username,
		Password:   dto.Password,
		Status:     0,
	}
	config.DB.Create(data)
	// 启动新的任务
	log.Println("添加学习通任务")
	c.joinPlay(dto.Username, data)
}

// List 分页获取任务列表
func (c *Xuexitong) List(dto *model.ListXueXiTongDto) *model.Result {
	utils.PageFormat(dto)
	var list []model.XueXiTong
	queryCondition := config.DB.Model(&model.XueXiTong{}).
		Where("username like ?", "%"+dto.Username+"%")
	var total int64 = 0
	queryCondition.Count(&total)
	queryCondition.Omit("password").Offset((dto.Page - 1) * dto.Size).Limit(dto.Size)
	queryCondition.Find(&list)
	return &model.Result{
		Code:    200,
		Data:    list,
		Message: "success",
		Total:   total,
	}
}

// Delete 删除任务
func (c *Xuexitong) Delete(id uint) {
	// Todo 删除任务
	config.DB.Delete(&model.XueXiTong{}, id)
}

// Pause 暂停任务
func (c *Xuexitong) Pause(id uint) {
	// Todo 暂停任务
	c.updateStatus(id, 0)
}

func (c *Xuexitong) findAllTaskNotFinish() []model.XueXiTong {
	var list []model.XueXiTong
	config.DB.Model(&model.XueXiTong{}).Where("status in ?", []int{0, 1, 2}).Find(&list)
	return list
}

// updateStatus 更新任务状态
func (c *Xuexitong) updateStatus(id uint, status int) {
	config.DB.Model(&model.XueXiTong{}).Where("id = ?", id).Update("status", status)
}

// joinPlay 加入启动任务
func (c *Xuexitong) joinPlay(queueId string, task *model.XueXiTong) {
	// 拦截异常
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	// 更新当前task为队列中
	c.updateStatus(task.Id, 1)
	// 不存在队列或者队列为空，直接新建队列
	if c._map[queueId] == nil || len(c._map[queueId]) == 0 {
		fmt.Println("创建队列" + queueId)
		c._map[queueId] = make([]model.XueXiTong, 0)
		c._map[queueId] = append(c._map[queueId], *task)
		go c.startPlay(queueId)
	} else {
		fmt.Println("加入队列" + queueId)
		c._map[queueId] = append(c._map[queueId], *task)
	}
}

// startPlay 开始任务(异步线程)
func (c *Xuexitong) startPlay(queueId string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	for {
		// 获取任务
		fmt.Println("开始任务" + queueId)
		if len(c._map[queueId]) == 0 {
			// 队列为空，删除队列
			delete(c._map, queueId)
			break
		}
		task := c._map[queueId][0]
		// 更新任务状态为进行中
		c.updateStatus(task.Id, 2)
		fmt.Println("更新状态完成")
		time.Sleep(10000)
		// todo 对这个任务进行刷课
		cache := &xuexitongApi.XueXiTUserCache{Name: task.Username, Password: task.Password}
		err := xuexitong.XueXiTLoginAction(cache) // 登录
		utils.OkOrPanic(err, "登录失败")

		// TODO 处理获取课程异常
		courseList, err := xuexitong.XueXiTPullCourseAction(cache)
		utils.OkOrPanic(err, "获取课程失败")
		userId, _ := strconv.Atoi(cache.UserID)
		log.Println("用户ID:", cache.UserID)
		// 获取指定课程
		log.Println("获取课程数量:", len(courseList))
		var course xuexitong.XueXiTCourse
		for _, _course := range courseList {
			if _course.CourseID == task.CourseId {
				course = _course
				break
			}
		}
		if course.CourseID == "" {
			log.Println("没有获取到对应课程")
			//TODO 没有查询到课程处理
		}
		log.Printf("刷课信息：[%+v]", course)
		key, _ := strconv.Atoi(course.Key)
		courseId, _ := strconv.Atoi(course.CourseID)
		chapterList, _, _ := xuexitong.PullCourseChapterAction(cache, course.Cpi, key) //获取对应章节信息

		var nodes []int
		for _, item := range chapterList.Knowledge {
			nodes = append(nodes, item.ID)
		}
		// 检测节点完成情况
		// TODO: 获取节点完成情况
		pointList, err := xuexitong.ChapterFetchPointAction(cache, nodes, &chapterList, key, userId, course.Cpi, courseId)
		log.Printf("获取章节信息：%+v\r\n", pointList)
		utils.OkOrPanic(err, "获取章节信息出错")
		var isFinished = func(index int) bool {
			if index < 0 || index >= len(pointList.Knowledge) {
				return false
			}
			i := pointList.Knowledge[index]
			return i.PointTotal >= 0 && i.PointTotal == i.PointFinished
		}
		for index, _ := range nodes {
			if isFinished(index) { //如果完成了的那么直接跳过
				continue
			}
			_, fetchCards, err := xuexitong.ChapterFetchCardsAction(cache, &chapterList, nodes, index, courseId, key, course.Cpi)
			if err != nil {
				utils.OkOrPanic(err)
			}
			videoDTOs, workDTOs, documentDTOs := entity.ParsePointDto(fetchCards)
			if videoDTOs == nil && workDTOs == nil && documentDTOs == nil {
				log.Println("没有可学习的内容")
			}
			log.Println("待完成的视频任务：", videoDTOs)
			log.Println("待完成的文档任务：", documentDTOs)

			// 视屏类型
			if videoDTOs != nil {
				for _, videoDTO := range videoDTOs {
					card, err := xuexitong.PageMobileChapterCardAction(
						cache, key, courseId, videoDTO.KnowledgeID, videoDTO.CardIndex, course.Cpi)
					if err != nil {
						log.Fatal(err)
					}
					videoDTO.AttachmentsDetection(card)
					ExecuteVideo(cache, &videoDTO)
					time.Sleep(5 * time.Second)
				}
			}
			// 文档类型
			if documentDTOs != nil {
				for _, documentDTO := range documentDTOs {
					card, err := xuexitong.PageMobileChapterCardAction(
						cache, key, courseId, documentDTO.KnowledgeID, documentDTO.CardIndex, course.Cpi)
					if err != nil {
						utils.OkOrPanic(err)
					}
					documentDTO.AttachmentsDetection(card)
					//point.ExecuteDocument(userCache, &documentDTO)
					ExecuteDocument(cache, &documentDTO)
					time.Sleep(5 * time.Second)
				}
			}

		}

		// 如果完成了，更新状态为已完成
		c.updateStatus(task.Id, 3)
		c._map[queueId] = c._map[queueId][1:]
	}

}

// ExecuteVideo 常规刷视频逻辑
func ExecuteVideo(cache *xuexitongApi.XueXiTUserCache, p *entity.PointVideoDto) {
	if state, _ := xuexitong.VideoDtoFetchAction(cache, p); state {
		var playingTime = p.PlayTime
		for {
			// TODO 检测当前任务是否被暂停 或 被删除
			if p.Duration-playingTime >= 58 {
				playReport, err := cache.VideoDtoPlayReport(p, playingTime, 0, 8, nil)
				if gojsonq.New().JSONString(playReport).Find("isPassed") == nil || err != nil {
					//lg.Print(lg.INFO, `[`, cache.Name, `] `, lg.BoldRed, "提交学时接口访问异常，返回信息：", playReport, err.Error())
					// TODO 处理信息
				}
				if gojsonq.New().JSONString(playReport).Find("isPassed").(bool) == true { //看完了，则直接退出
					//lg.Print(lg.INFO, "[", lg.Green, cache.Name, lg.Default, "] ", " 【", p.Title, "】 >>> ", "提交状态：", lg.Green, strconv.FormatBool(gojsonq.New().JSONString(playReport).Find("isPassed").(bool)), lg.Default, " ", "观看时间：", strconv.Itoa(p.Duration)+"/"+strconv.Itoa(p.Duration), " ", "观看进度：", fmt.Sprintf("%.2f", float32(p.Duration)/float32(p.Duration)*100), "%")
					// TODO 处理信息
					break
				}
				// TODO 处理信息
				//lg.Print(lg.INFO, "[", lg.Green, cache.Name, lg.Default, "] ", " 【", p.Title, "】 >>> ", "提交状态：", lg.Green, lg.Green, strconv.FormatBool(gojsonq.New().JSONString(playReport).Find("isPassed").(bool)), lg.Default, " ", "观看时间：", strconv.Itoa(playingTime)+"/"+strconv.Itoa(p.Duration), " ", "观看进度：", fmt.Sprintf("%.2f", float32(playingTime)/float32(p.Duration)*100), "%")
				playingTime = playingTime + 58
				time.Sleep(58 * time.Second)
			} else if p.Duration-playingTime < 58 {
				playReport, err := cache.VideoDtoPlayReport(p, p.Duration, 2, 8, nil)
				if gojsonq.New().JSONString(playReport).Find("isPassed") == nil || err != nil {
					// TODO 处理信息
					//lg.Print(lg.INFO, `[`, cache.Name, `] `, lg.BoldRed, "提交学时接口访问异常，返回信息：", playReport, err.Error())
				}
				if gojsonq.New().JSONString(playReport).Find("isPassed").(bool) == true { //看完了，则直接退出
					// TODO 处理信息
					//lg.Print(lg.INFO, "[", lg.Green, cache.Name, lg.Default, "] ", " 【", p.Title, "】 >>> ", "提交状态：", lg.Green, lg.Green, strconv.FormatBool(gojsonq.New().JSONString(playReport).Find("isPassed").(bool)), lg.Default, " ", "观看时间：", strconv.Itoa(p.Duration)+"/"+strconv.Itoa(p.Duration), " ", "观看进度：", fmt.Sprintf("%.2f", float32(p.Duration)/float32(p.Duration)*100), "%")
					break
				}
				// TODO 处理信息
				//lg.Print(lg.INFO, "[", lg.Green, cache.Name, lg.Default, "] ", " 【", p.Title, "】 >>> ", "提交状态：", lg.Green, lg.Green, strconv.FormatBool(gojsonq.New().JSONString(playReport).Find("isPassed").(bool)), lg.Default, " ", "观看时间：", strconv.Itoa(p.Duration)+"/"+strconv.Itoa(p.Duration), " ", "观看进度：", fmt.Sprintf("%.2f", float32(p.Duration)/float32(p.Duration)*100), "%")
				time.Sleep(time.Duration(p.Duration-playingTime) * time.Second)
			}
		}
	} else {
		log.Fatal("视频解析失败")
	}
}

// ExecuteDocument 常规刷文档逻辑
func ExecuteDocument(cache *xuexitongApi.XueXiTUserCache, p *entity.PointDocumentDto) {
	report, err := cache.DocumentDtoReadingReport(p)
	if gojsonq.New().JSONString(report).Find("status") == nil || err != nil {
		//lg.Print(lg.INFO, `[`, cache.Name, `] `, lg.BoldRed, "提交学时接口访问异常，返回信息：", report, err.Error())
		// TODO 处理异常
		log.Fatalln(err)
	}
	if gojsonq.New().JSONString(report).Find("status").(bool) {
		//// TODO 处理信息
		//lg.Print(lg.INFO, "[", lg.Green, cache.Name, lg.Default, "] ", " 【", p.Title, "】 >>> ", "文档阅览状态：", lg.Green, lg.Green, strconv.FormatBool(gojsonq.New().JSONString(report).Find("status").(bool)), lg.Default, " ")
	}
}
