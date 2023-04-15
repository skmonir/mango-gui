package parser

import (
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/google/uuid"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/oj/client"
	"github.com/skmonir/mango-gui/backend/judge-framework/scheduler"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const genericError = "Something went wrong!"

func ScheduleParse(url string) error {
	if checkIfParsingIsAlreadyScheduledForThisURL(url) {
		return errors.New("This URL is already scheduled for parsing.")
	}

	startTime, err := fetchStartTime(url)
	if err != nil {
		return err
	}

	if !utils.IsTimeInFuture(startTime) {
		return errors.New("Schedule only upcoming contests")
	}

	fmt.Println(startTime)
	if err := scheduleTheParsing(url, "SCHEDULED", startTime); err != nil {
		return err
	}

	return nil
}

func fetchStartTime(url string) (time.Time, error) {
	startTime := time.Time{}

	platform, cid, _, _ := utils.ExtractInfoFromUrl(url)
	if platform == "" || cid == "" {
		return time.Time{}, errors.New("Error! Please check the url.")
	}

	_, httpClient := client.GetHttpClientByPlatform(platform)

	if platform == "codeforces" && strings.Contains(url, "codeforces.com/contest") {
		err, body := getContestDetailPage(httpClient, "https://codeforces.com/contests?complete=true")
		if err != nil {
			return time.Time{}, err
		}
		startTime, err = getCodeforcesContestStartTime(body, cid, 2)
		if err != nil {
			return time.Time{}, err
		}
	} else if platform == "codeforces" && strings.Contains(url, "codeforces.com/gym") {
		err, body := getContestDetailPage(httpClient, "https://codeforces.com/gyms")
		if err != nil {
			return time.Time{}, err
		}
		startTime, err = getCodeforcesContestStartTime(body, cid, 1)
		if err != nil {
			return time.Time{}, err
		}
	} else if platform == "atcoder" {
		err, body := getContestDetailPage(httpClient, url)
		if err != nil {
			return time.Time{}, err
		}
		startTime, err = getAtcoderContestStartTime(body)
		if err != nil {
			return time.Time{}, err
		}
	} else {
		return time.Time{}, errors.New("Platform is not recognized")
	}
	return startTime, nil
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

func getContestDetailPage(httpClient *http.Client, url string) (error, soup.Root) {
	html, err := utils.GetBody(httpClient, url)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(genericError), soup.Root{}
	}

	return nil, soup.HTMLParse(string(html))
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

	return startTimeLocal, nil
}

func scheduleTheParsing(url, stage string, parsingTime time.Time) error {
	scheduleTask := models.ParseSchedulerTask{
		Id:        uuid.New().String(),
		Url:       url,
		StartTime: parsingTime,
		Stage:     stage,
	}

	if err := ScheduleTaskInScheduler(scheduleTask); err != nil {
		return errors.New(genericError)
	}
	services.AddParseScheduledTask(scheduleTask)
	fmt.Println("Parse scheduled successfully")
	return nil
}

func ScheduleTaskInScheduler(scheduleTask models.ParseSchedulerTask) error {
	if !utils.IsTimeInFuture(scheduleTask.StartTime) {
		return errors.New("Schedule only upcoming contests")
	}
	return scheduler.ScheduleOneTimeTask(scheduleTask.Id, func() { go parseWithRetry(scheduleTask) }, scheduleTask.StartTime)
}

func parseWithRetry(scheduleTask models.ParseSchedulerTask) {
	if timeChanged, newStartTime := checkIfScheduledTimeChanged(scheduleTask.Url, scheduleTask.StartTime); timeChanged {
		RemoveParseSchedule(scheduleTask.Id)
		_ = scheduleTheParsing(scheduleTask.Url, "RE_SCHEDULED", newStartTime)
		return
	}
	isParsed := false
	services.UpdateParseScheduledTaskStage(scheduleTask.Id, "RUNNING")
	for i := 0; i < 10; i++ {
		problems := Parse(scheduleTask.Url)
		if len(problems) > 0 {
			isParsed = true
			scheduler.RemoveScheduledTask(scheduleTask.Id)
			services.UpdateParseScheduledTaskStage(scheduleTask.Id, "COMPLETE")
			break
		}
		time.Sleep(15 * time.Second)
	}
	if !isParsed {
		services.UpdateParseScheduledTaskStage(scheduleTask.Id, "FAILED")
	}
}

func checkIfScheduledTimeChanged(url string, prevStartTime time.Time) (bool, time.Time) {
	newStartTime, err := fetchStartTime(url)
	if err != nil || newStartTime == prevStartTime {
		return false, time.Time{}
	}
	return true, newStartTime
}

func RemoveParseSchedule(taskId string) {
	scheduler.RemoveScheduledTask(taskId)
	services.RemoveParseScheduledTasksByIds([]string{taskId})
}
