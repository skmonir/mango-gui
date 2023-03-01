package config

//
//func GetSourceDirPath() string {
//	config := GetJudgeConfig()
//	return filepath.Join(config.WorkspaceDirectory, strings.ToLower(config.), config.CurrentContestId, config.SrcDir)
//}
//
//func GetSourceFilePathWithExt(problemId string) string {
//	config := GetJudgeConfig()
//	return filepath.Join(config.GetSourceDirPath(), problemId+".cpp")
//}
//
//func GetSourceFilePathWithoutExt(problemId string) string {
//	config := GetJudgeConfig()
//	return filepath.Join(config.GetSourceDirPath(), problemId)
//}
//
//func GetTestcaseDirPath() string {
//	config := GetJudgeConfig()
//	return filepath.Join(config.Workspace, strings.ToLower(config.OJ), config.CurrentContestId, config.TestDir)
//}
//
//func GetTestcaseFilePath(problemId string) string {
//	config := GetJudgeConfig()
//	return filepath.Join(config.GetTestcaseDirPath(), problemId+".json")
//}
//
//func ResolveTescasePath(problemId string) error {
//	config := GetJudgeConfig()
//	testCaseDirPath := config.GetTestcaseDirPath()
//
//	if err := utils.CreateFile(testCaseDirPath, problemId+".json"); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func GetProblemInfo(problemId string) (models.Problem, error) {
//	config := GetJudgeConfig()
//	var problemInfo models.Problem
//	testpath := config.GetTestcaseFilePath(problemId)
//
//	if !utils.IsFileExist(testpath) {
//		return problemInfo, errors.New("problem not found")
//	}
//
//	data, err := ioutil.ReadFile(testpath)
//	if err != nil {
//		return problemInfo, err
//	}
//
//	err = json.Unmarshal(data, &problemInfo)
//	if err != nil {
//		return problemInfo, err
//	}
//
//	return problemInfo, nil
//}
//
//func GetProblemIdListForTester() []string {
//	config := GetJudgeConfig()
//	files, err := ioutil.ReadDir(config.GetTestcaseDirPath())
//	if err != nil {
//		return []string{}
//	}
//
//	problemIdList := []string{}
//	for _, fileService := range files {
//		filename := strings.Split(fileService.Lang(), ".")[0]
//		problemIdList = append(problemIdList, filename)
//	}
//	return problemIdList
//}
