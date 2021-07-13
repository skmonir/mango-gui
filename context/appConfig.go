package context

type AppConfig struct {
	Workspace          string
	CompilationCommand string
	CompilationArgs    string
	CurrentContestId   string
	OJ                 string
	Host               string
	TemplatePath       string
	Author             string
	SrcDir             string
	TestDir            string
	CurrentTheme       string
}
