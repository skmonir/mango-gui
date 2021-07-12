package taskParser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/skmonir/mango-gui/utils"
)

func (parser CodeforcesParser) GetContestType(contestId string) string {
	id, err := strconv.ParseInt(contestId, 10, 30)
	if err != nil {
		if len(contestId) > 5 {
			return "gym"
		}
		return "contest"
	}
	if id > 100000 {
		return "gym"
	}
	return "contest"
}

func (parser CodeforcesParser) GetContestUrl(contestId string) string {
	contestType := parser.GetContestType(contestId)
	return fmt.Sprintf("https://codeforces.com/%v/%v", contestType, contestId)
}

func (parser CodeforcesParser) GetProblemUrl(contestId string, problemId string) string {
	contestType := parser.GetContestType(contestId)
	return fmt.Sprintf("https://codeforces.com/%v/%v/problem/%v", contestType, contestId, problemId)
}

func (parser CodeforcesParser) ParseContestAndProblemId(cmd string) (string, string, error) {
	cmd = strings.TrimSpace(cmd)
	if len(cmd) == 0 {
		return "", "", errors.New("command is not valid")
	}

	ptr, sz := 0, len(cmd)

	contestId := ""
	for ptr < sz && utils.IsDigit(rune(cmd[ptr])) {
		contestId += string(cmd[ptr])
		ptr++
	}

	problemId := ""
	for ptr < sz && rune(cmd[ptr]) != ' ' {
		problemId += string(cmd[ptr])
		ptr++
	}

	return contestId, problemId, nil
}

func (parser CodeforcesParser) GetProblemList(body []byte) []string {
	stat_regexp := regexp.MustCompile(`class="problems"[\s\S]+?</tr>([\s\S]+?)</table>`)
	stat_body := stat_regexp.FindSubmatch(body)
	if stat_body == nil {
		return []string{}
	}
	probs_table := stat_body[1]

	row_idx_regex := regexp.MustCompile(`<tr[\s\S]*?>`)
	row_idxs := row_idx_regex.FindAllIndex(probs_table, -1)
	if row_idxs == nil {
		return []string{}
	}

	row_idxs = append(row_idxs, []int{0, len(probs_table)})
	id_td_regex := regexp.MustCompile(`<td[\s\S]*?>`)
	prob_id_regex := regexp.MustCompile(`<a[\s\S]*?>([\s\S]*)</a>`)

	problem_ids := make([]string, len(row_idxs)-1)
	for i := 1; i < len(row_idxs); i++ {
		current_row := probs_table[row_idxs[i-1][0]:row_idxs[i][1]]
		td_idxs := id_td_regex.FindAllIndex(current_row, -1)

		current_prob_elem := current_row[td_idxs[0][0]:td_idxs[1][1]]
		id := prob_id_regex.FindSubmatch(current_prob_elem)
		if id != nil {
			problem_ids[i-1] = strings.TrimSpace(string(id[1]))
		} else {
			problem_ids[i-1] = "$"
		}
	}

	return problem_ids
}

func (parser CodeforcesParser) GetProblemConstraints(body []byte) (int64, uint64) {
	trg := regexp.MustCompile(`class="time-limit"[\s\S]*?([\d]+) seconds`)
	mrg := regexp.MustCompile(`class="memory-limit"[\s\S]*?([\d]+) megabytes`)
	a := trg.FindSubmatch(body)
	b := mrg.FindSubmatch(body)

	var timeLimit int64 = 2      // default time-limit
	var memoryLimit uint64 = 512 // default memory-limit

	if len(a) > 0 {
		TL, err := strconv.Atoi(strings.TrimSpace(string(FilterHtml(a[1]))))
		if err == nil {
			timeLimit = int64(TL)
		}
	}
	if len(b) > 0 {
		ML, err := strconv.Atoi(strings.TrimSpace(string(FilterHtml(b[1]))))
		if err == nil {
			memoryLimit = uint64(ML)
		}
	}

	return timeLimit, memoryLimit
}

func (parser CodeforcesParser) GetProblemSamples(body []byte) ([][]byte, [][]byte, error) {
	irg := regexp.MustCompile(`class="input"[\s\S]*?<pre>([\s\S]*?)</pre>`)
	org := regexp.MustCompile(`class="output"[\s\S]*?<pre>([\s\S]*?)</pre>`)
	a := irg.FindAllSubmatch(body, -1)
	b := org.FindAllSubmatch(body, -1)
	if a == nil || b == nil || len(a) != len(b) {
		return nil, nil, fmt.Errorf("cannot parse samples")
	}

	var input [][]byte
	var output [][]byte
	for i := 0; i < len(a); i++ {
		input = append(input, FilterHtml(a[i][1]))
		output = append(output, FilterHtml(b[i][1]))
	}
	return input, output, nil
}

func (parser CodeforcesParser) GetProblemName(body []byte) string {
	name_body_regex := regexp.MustCompile(`class="title"([\s\S]*?)class="time-limit"`)
	name_body := name_body_regex.FindSubmatch(body)
	if name_body == nil {
		return ""
	}

	name_regex := regexp.MustCompile(`>([\s\S]*?)</div>[\s\S]*?`)
	name := name_regex.FindSubmatch(name_body[1])
	if name == nil {
		return ""
	}
	return strings.TrimSpace(string(name[1]))
}
