package parser

import (
	"context"
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/procyon-projects/chrono"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"github.com/skmonir/mango-gui/backend/socket"
	"strconv"
	"strings"
	"time"
)

const genericError = "Something went wrong!"

func ScheduleParse(url string) (time.Time, error) {
	startTime := time.Time{}

	platform, cid, _ := utils.ExtractInfoFromUrl(url)
	if platform == "" || cid == "" {
		return startTime, errors.New("Error! Please check the url.")
	}

	if platform == "codeforces" && strings.Contains(url, "codeforces.com/contest") {
		err, body := getContestDetailPage("https://codeforces.com/contests")
		if err != nil {
			return time.Time{}, err
		}
		startTime, err = getCodeforcesContestStartTime(body, cid, 2)
		if err != nil {
			return startTime, err
		}
	} else if platform == "codeforces" && strings.Contains(url, "codeforces.com/gym") {
		err, body := getContestDetailPage("https://codeforces.com/gyms")
		if err != nil {
			return time.Time{}, err
		}
		startTime, err = getCodeforcesContestStartTime(body, cid, 1)
		if err != nil {
			return startTime, err
		}
	} else if platform == "atcoder" {
		err, body := getContestDetailPage(url)
		if err != nil {
			return startTime, err
		}
		startTime, err = getAtcoderContestStartTime(body)
		if err != nil {
			return startTime, err
		}
	} else {
		return startTime, errors.New("Platform is not recognized")
	}

	fmt.Println(startTime)
	if err := scheduleTheParsing(url, startTime); err != nil {
		return time.Time{}, err
	}

	return startTime, nil
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
	timeDataElement := body.Find("tr", "data-contestid", cid).FindAll("td")[childId].Find("a")
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
	// 2023 04 08 T 21 00
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
	taskScheduler := chrono.NewDefaultTaskScheduler()
	_, err := taskScheduler.Schedule(func(ctx context.Context) {
		go parseWithRetry(url)
	}, chrono.WithTime(parsingTime))
	if err != nil {
		logger.Error(err.Error())
		return errors.New(genericError)
	}
	fmt.Println("Parse scheduled successfully")
	return nil
}

func parseWithRetry(url string) {
	for i := 0; i < 10; i++ {
		time.Sleep(15 * time.Second)
		socket.PublishStatusMessage("parse_schedule_event", "running", "info")
		probs := Parse(url)
		if len(probs) > 0 {
			socket.PublishStatusMessage("parse_schedule_event", "done", "info")
			break
		}
	}
}
