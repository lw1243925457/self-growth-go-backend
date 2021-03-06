package job

import (
	"encoding/json"
	"fmt"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
	"go.mongodb.org/mongo-driver/bson"
	v1 "seltGrowth/internal/api/v1"
	"seltGrowth/internal/api/v1/game_text_auto"
	"seltGrowth/internal/pkg/service/statistics"
	"time"
)

type DayAchievementCal interface {
	Cal() error
}

type dayAchievementCal struct {
	dayService statistics.DayStatisticsService
}

func NewDayAchievementCal() DayAchievementCal {
	return &dayAchievementCal{
		dayService: statistics.NewDayStatisticsService(),
	}
}

func (d *dayAchievementCal) Cal() error {
	log.Infof("每日成就统计转换开始：%s", time.Now())

	var users []v1.User
	err := mgm.Coll(&v1.User{}).SimpleFind(&users, bson.M{})
	if err != nil {
		log.Error(err)
		return err
	}

	yesterday := time.Now().AddDate(0, 0, -1)
	for _, user := range users {
		dayStatistics, err := d.dayService.DayStatistics(yesterday, user.Email, true, false)
		if err != nil {
			log.Error(err)
			return err
		}

		log.Infof("连续打卡奖励倍率为：learn -- %d, running -- %d, sleep -- %d, improve -- %d", user.LearnPersistDay, user.RunningPersistDay, user.SleepPersistDay, user.ImprovePersistDay)
		achievement := game_text_auto.NewDayAchievement(dayStatistics, user.LearnPersistDay, user.RunningPersistDay, user.SleepPersistDay, user.ImprovePersistDay)
		err = mgm.Coll(&game_text_auto.DayAchievement{}).Create(achievement)
		if err != nil {
			log.Error(err)
			return err
		}

		s, err := json.MarshalIndent(achievement, "", "    ")
		log.Infof("昨日成就：：%s", string(s))

		if achievement.Spirit > 0 {
			user.SleepPersistDay = user.SleepPersistDay + 1
		}
		if achievement.Strength > 0 {
			user.RunningPersistDay = user.RunningPersistDay + 1
		}
		if achievement.Reiki > 0 {
			user.LearnPersistDay = user.LearnPersistDay + 1
		}
		if achievement.Spirit > 0 && achievement.Strength > 0 && achievement.Reiki > 0 {
			user.ImprovePersistDay = user.ImprovePersistDay + 1
		}
		err = mgm.Coll(&user).Update(&user)
		if err != nil {
			log.Error(err)
		}

		msg := model.NewMessage("AT_CrUMMtfpshG6oo8dyAzWRrK4gl5PNwm6").SetContent(string(s)).AddUId("UID_clQRP3Pig7dxpSqGGICcrBy4KPTB")
		msgArr, err := wxpusher.SendMessage(msg)
		fmt.Println(msgArr, err)
	}
	return nil
}
