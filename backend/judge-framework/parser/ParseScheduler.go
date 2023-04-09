package parser

import (
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/google/uuid"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/scheduler"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"strconv"
	"strings"
	"time"
)

const genericError = "Something went wrong!"

func ScheduleParse(url string) error {
	if checkIfParsingIsAlreadyScheduledForThisURL(url) {
		return errors.New("This URL is already scheduled for parsing.")
	}
	startTime := time.Time{}

	platform, cid, _ := utils.ExtractInfoFromUrl(url)
	if platform == "" || cid == "" {
		return errors.New("Error! Please check the url.")
	}

	if platform == "codeforces" && strings.Contains(url, "codeforces.com/contest") {
		err, body := getContestDetailPage("https://codeforces.com/contests")
		if err != nil {
			return err
		}
		startTime, err = getCodeforcesContestStartTime(body, cid, 2)
		if err != nil {
			return err
		}
	} else if platform == "codeforces" && strings.Contains(url, "codeforces.com/gym") {
		err, body := getContestDetailPage("https://codeforces.com/gyms")
		if err != nil {
			return err
		}
		startTime, err = getCodeforcesContestStartTime(body, cid, 1)
		if err != nil {
			return err
		}
	} else if platform == "atcoder" {
		err, body := getContestDetailPage(url)
		if err != nil {
			return err
		}
		startTime, err = getAtcoderContestStartTime(body)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Platform is not recognized")
	}

	fmt.Println(startTime)
	if err := scheduleTheParsing(url, startTime); err != nil {
		return err
	}

	return nil
}

func checkIfParsingIsAlreadyScheduledForThisURL(url string) bool {
	tasks := services.GetFutureParseScheduledTasks()
	for _, task := range tasks {
		if task.Url == url {
			return true
		}
	}
	return false
}

func getContestDetailPage(url string) (error, soup.Root) {
	html, err := utils.GetHtmlBody(url)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(genericError), soup.Root{}
	}

	return nil, soup.HTMLParse(html)
}

func getCodeforcesContestStartTime(body soup.Root, cid string, childId int) (time.Time, error) {
	defer utils.PanicRecovery()
	contestElement := body.Find("tr", "data-contestid", cid)
	if contestElement.Error != nil {
		return time.Time{}, errors.New("Schedule only upcoming contests")
	}
	timeDataElement := contestElement.FindAll("td")[childId].Find("a")
	if timeDataElement.Error != nil {
		return time.Time{}, errors.New("Schedule only upcoming contests")
	}
	datetimeUrl := strings.TrimSpace(timeDataElement.Attrs()["href"])

	queries := utils.ParseQueryMapFromUrl(datetimeUrl)
	day, _ := strconv.Atoi(queries.Get("day"))
	month, _ := strconv.Atoi(queries.Get("month"))
	year, _ := strconv.Atoi(queries.Get("year"))
	hour, _ := strconv.Atoi(queries.Get("hour"))
	min, _ := strconv.Atoi(queries.Get("min"))

	dateStr := fmt.Sprintf("%d-%02d-%02d %02d:%02d", year, month, day, hour, min)
	startTime, err := time.ParseInLocation("2006-01-02 15:04", dateStr, time.FixedZone("UTC+3", 3*60*60))
	if err != nil {
		logger.Error(err.Error())
		return time.Time{}, errors.New(genericError)
	}
	loc, err := time.LoadLocation("Local")
	if err != nil {
		logger.Error(err.Error())
		return time.Time{}, errors.New(genericError)
	}
	startTimeLocal := startTime.In(loc)

	if !utils.IsTimeInFuture(startTimeLocal) {
		return time.Time{}, errors.New("Schedule only upcoming contests")
	}

	return startTimeLocal, nil
}

func getAtcoderContestStartTime(body soup.Root) (time.Time, error) {
	defer utils.PanicRecovery()
	timeDataElement := body.Find("small", "class", "contest-duration").Find("a")
	datetimeUrl := strings.TrimSpace(timeDataElement.Attrs()["href"])
	queries := utils.ParseQueryMapFromUrl(datetimeUrl)
	iso := queries.Get("iso")
	dateStr := fmt.Sprintf("%v-%v-%v %v:%v", iso[0:4], iso[4:6], iso[6:8], iso[9:11], iso[11:13])
	startTime, err := time.ParseInLocation("2006-01-02 15:04", dateStr, time.FixedZone("UTC+9", 9*60*60))
	if err != nil {
		logger.Error(err.Error())
		return time.Time{}, errors.New(genericError)
	}
	loc, err := time.LoadLocation("Local")
	if err != nil {
		logger.Error(err.Error())
		return time.Time{}, errors.New(genericError)
	}
	startTimeLocal := startTime.In(loc)

	if !utils.IsTimeInFuture(startTimeLocal) {
		return time.Time{}, errors.New("Schedule only upcoming contests")
	}

	return startTimeLocal, nil
}

func scheduleTheParsing(url string, parsingTime time.Time) error {
	scheduleTask := models.ParseSchedulerTask{
		Id:        uuid.New().String(),
		Url:       url,
		StartTime: parsingTime,
		Stage:     "SCHEDULED",
	}

	if err := ScheduleTaskInScheduler(scheduleTask); err != nil {
		return errors.New(genericError)
	}
	services.AddParseScheduledTask(scheduleTask)
	fmt.Println("Parse scheduled successfully")
	return nil
}

func ScheduleTaskInScheduler(scheduleTask models.ParseSchedulerTask) error {
	return scheduler.ScheduleOneTimeTask(scheduleTask.Id, func() { go parseWithRetry(scheduleTask) }, scheduleTask.StartTime)
}

func parseWithRetry(scheduleTask models.ParseSchedulerTask) {
	isParsed := false
	for i := 0; i < 10; i++ {
		services.UpdateParseScheduledTask(scheduleTask.Id, "RUNNING")
		time.Sleep(15 * time.Second)
		problems := Parse(scheduleTask.Url)
		if len(problems) > 0 {
			isParsed = true
			scheduler.RemoveScheduledTask(scheduleTask.Id)
			services.UpdateParseScheduledTask(scheduleTask.Id, "COMPLETE")
			break
		}
	}
	if !isParsed {
		services.UpdateParseScheduledTask(scheduleTask.Id, "FAILED")
	}
}

func RemoveParseSchedule(taskId string) {
	scheduler.RemoveScheduledTask(taskId)
	services.RemoveParseScheduledTasksByIds([]string{taskId})
}
