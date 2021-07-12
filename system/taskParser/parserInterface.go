package taskParser

type Parser interface {
	GetContestType(string) string
	GetContestUrl(string) string
	GetProblemUrl(string, string) string
	GetProblemName([]byte) string
	GetProblemList([]byte) []string                       // problem id
	GetProblemConstraints([]byte) (int64, uint64)         // timeLimit, memoryLimit
	GetProblemSamples([]byte) ([][]byte, [][]byte, error) // input, output, error
	ParseContestAndProblemId(string) (string, string, error)
}
