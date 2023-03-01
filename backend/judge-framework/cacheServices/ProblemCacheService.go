package cacheServices

import (
	"encoding/json"
	"fmt"
	"github.com/skmonir/mango-ui/backend/judge-framework/cache"
	"github.com/skmonir/mango-ui/backend/judge-framework/models"
)

func GetProblemListFromCache(platform string, cid string) []models.Problem {
	fmt.Println("Fetching problems from cache...")

	key := fmt.Sprintf("ParsedProblems:%v.%v", platform, cid)
	problemsObjStr := cache.GetGlobalCache().Get(key)

	var problems []models.Problem
	if problemsObjStr != "" {
		if err := json.Unmarshal([]byte(problemsObjStr), &problems); err != nil {
			fmt.Println(err)
			return problems
		}
		fmt.Println("Fetched problems from cache")
	}
	return problems
}

func SetProblemListIntoCache(platform string, cid string, problemList []models.Problem) {
	fmt.Println("Saving problems into cache...")

	data, err := json.MarshalIndent(problemList, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Saved problems into cache")

	key := fmt.Sprintf("ParsedProblems:%v.%v", platform, cid)
	cache.GetGlobalCache().Set(key, string(data))
}
