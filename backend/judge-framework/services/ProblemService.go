package services

import (
	"fmt"
	"strings"

	"github.com/skmonir/mango-ui/backend/judge-framework/cacheServices"
	"github.com/skmonir/mango-ui/backend/judge-framework/fileService"
	"github.com/skmonir/mango-ui/backend/judge-framework/models"
	"github.com/skmonir/mango-ui/backend/judge-framework/utils"
)

func GetProblem(platform string, cid string, label string) models.Problem {
	problems := GetProblemList(platform, cid)
	for _, prob := range problems {
		if prob.Label == label {
			return prob
		}
	}
	return models.Problem{}
}

func GetProblemListByUrl(url string) []models.Problem {
	platform, cid, pid := utils.ExtractInfoFromUrl(url)

	problems := GetProblemList(platform, cid)
	if len(problems) == 0 || pid == "" {
		return problems
	}

	var response []models.Problem
	for _, problem := range problems {
		if strings.Contains(problem.Url, url) {
			response = append(response, problem)
		}
	}
	return response
}

func GetProblemList(platform string, cid string) []models.Problem {
	fmt.Println("Fetching problems...")
	if problems := cacheServices.GetProblemListFromCache(platform, cid); len(problems) > 0 {
		return problems
	}
	problems := fileService.GetProblemListFromFile(platform, cid)
	if len(problems) > 0 {
		cacheServices.SetProblemListIntoCache(platform, cid, problems)
	}
	return problems
}

func SaveProblemList(problemList []models.Problem) {
	fmt.Println("Saving problems...")
	if len(problemList) > 0 {
		platform, cid := problemList[0].Platform, problemList[0].ContestId
		cacheServices.SetProblemListIntoCache(platform, cid, problemList)
		fileService.SaveProblemListIntoFile(platform, cid, problemList)
	}
}

func UpdateProblemList(newProblemList []models.Problem) {
	if len(newProblemList) > 0 {
		fmt.Println("Updating problems...")

		platform, cid := newProblemList[0].Platform, newProblemList[0].ContestId
		oldProblemList := GetProblemList(platform, cid)

		for _, newProblem := range newProblemList {
			isFound := false
			for i := 0; i < len(oldProblemList); i++ {
				if oldProblemList[i].Label == newProblem.Label {
					oldProblemList[i] = newProblem
					isFound = true
				}
			}
			if !isFound {
				oldProblemList = append(oldProblemList, newProblem)
			}
		}
		SaveProblemList(oldProblemList)
	}
}
