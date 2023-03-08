package fileService

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
)

func CreateSourceFiles(problems []models.Problem) {
	judgeConfig := config.GetJudgeConfigFromCache()

	for _, prob := range problems {
		saveSourceIntoFile(judgeConfig, prob)
	}
}

func saveSourceIntoFile(judgeConfig *config.JudgeConfig, problem models.Problem) {
	folderPath := fmt.Sprintf("%v/%v/%v/source", strings.TrimRight(judgeConfig.WorkspaceDirectory, "/"), problem.Platform, problem.ContestId)
	fileName := problem.Label + judgeConfig.ActiveLanguage.FileExtension
	filePath := filepath.Join(folderPath, fileName)

	if !utils.IsFileExist(filePath) {
		template := getTemplate(judgeConfig, problem)
		utils.WriteFileContent(folderPath, fileName, []byte(template))
	}
}

func getTemplate(judgeConfig *config.JudgeConfig, meta models.Problem) string {
	header := getTemplateHeader(judgeConfig.Author, strings.ToUpper(meta.Label)+" - "+meta.Name)
	body := ""

	if len(judgeConfig.ActiveLanguage.TemplatePath) > 0 {
		body = utils.ReadFileContent(judgeConfig.ActiveLanguage.TemplatePath, 123456, 123456)
		if len(body) == 0 {
			body = getGenericTemplateBody()
		}
	} else {
		body = getGenericTemplateBody()
	}

	return header + body
}

func getTemplateHeader(author string, problemName string) string {
	header := "/**\n"
	header += fmt.Sprintf(" *     author:  %v\n", author)
	header += fmt.Sprintf(" *    created:  %v\n", time.Now().Local().Format("2-Jan-2006 15:04:05"))
	header += fmt.Sprintf(" *    problem:  %v\n", problemName)
	header += "**/\n\n"

	return header
}

func getGenericTemplateBody() string {
	body := ""

	body += "#include <bits/stdc++.h>\n"
	body += "\n"
	body += "using namespace std;\n"
	body += "\n"
	body += "const int N = 1e5 + 7;\n"
	body += "const int INF = 1e9 + 7;\n"
	body += "const int MOD = 1e9 + 7;\n"
	body += "\n"
	body += "\n"
	body += "int solver() {\n"
	body += "\t// your code goes here\n"
	body += "\treturn 0;\n"
	body += "}\n"
	body += "\n"
	body += "\n"
	body += "int main() {\n"
	body += "\tios::sync_with_stdio(0), cin.tie(0);\n"
	body += "\tint tt = 1;\n"
	body += "\t// cin >> tt;\n"
	body += "\tfor (int t = 1; t <= tt; ++t) {\n"
	body += "\t\t// cout << \"Case \" << t << \": \";\n"
	body += "\t\tsolver();\n"
	body += "\t}\n"
	body += "\treturn 0;\n"
	body += "}"

	return body
}

func OpenSourceByPath(filePath string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", filePath).Run()
	case "windows":
		err = exec.Command("cmd", fmt.Sprintf("/C start %v", filePath)).Run()
	case "darwin":
		err = exec.Command("open", filePath).Run()
	default:
		fmt.Println("unsupported os")
	}
	if err != nil {
		fmt.Println(err)
	}
}

func OpenSourceByMetadata(platform string, cid string, label string) {
	filePath := getSourceFilePath(platform, cid, label)
	if utils.IsFileExist(filePath) {
		OpenSourceByPath(filePath)
	}
}

func GetCodeByMetadata(platform string, cid string, label string) string {
	filePath := getSourceFilePath(platform, cid, label)
	if utils.IsFileExist(filePath) {
		return utils.ReadFileContent(filePath, 123456, 123456)
	}
	return ""
}

func getSourceFilePath(platform string, cid string, label string) string {
	judgeConfig := config.GetJudgeConfigFromCache()

	folderPath := fmt.Sprintf("%v/%v/%v/source", strings.TrimRight(judgeConfig.WorkspaceDirectory, "/"), platform, cid)
	fileName := label + judgeConfig.ActiveLanguage.FileExtension
	filePath := filepath.Join(folderPath, fileName)

	return filePath
}

func GetSourceBinaryPath(platform string, cid string, label string) string {
	judgeConfig := config.GetJudgeConfigFromCache()

	folderPath := fmt.Sprintf("%v/%v/%v/source", strings.TrimRight(judgeConfig.WorkspaceDirectory, "/"), platform, cid)
	binaryPath := filepath.Join(folderPath, label)

	return binaryPath
}
