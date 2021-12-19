package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"time"
)

type UserStatusResponse struct {
	Status string
	Result []*Submission
}
type Submission struct {
	Id, ContestId, CreationTimeSeconds, RelativeTimeSeconds int64
	Problem                                                 *Problem
	ProgrammingLanguage                                     string
	Verdict                                                 string // "OK" means the submission works
	Testset                                                 string
	PassedTestCount                                         int
	TimeConsumedMillis                                      int
	memoryConsumedBytes                                     int
}

type Problem struct {
	ContestId int
	Index     string
	Name      string
	Points    float64
	Tags      []string
}

func fetchSubmissions(handle string) *UserStatusResponse {

	url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=1000&lang=en", handle)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// var cfResp UserStatusResponse
	cfResp := &(UserStatusResponse{})
	if err := json.NewDecoder(resp.Body).Decode(cfResp); err != nil {
		log.Fatal(err)
	}
	if cfResp.Status != "OK" {
		log.Fatal("codeforces user.status API responded:", cfResp.Status)
	}
	return cfResp
}

func validateHandle(h string) bool {
	// TODO
	return true
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Please tell me the user handle :)")
	}
	handle := os.Args[1]
	if !validateHandle(handle) {
		log.Fatal(fmt.Sprintf("Something wrong with the user handle %s, please take a look :)", os.Args[1]))
	}
	res := fetchSubmissions(handle)

	solvedProblems, lines := map[string]bool{}, []string{}
	for _, v := range res.Result {
		id := fmt.Sprintf("%d/%s", v.Problem.ContestId, v.Problem.Index)
		if v.Verdict != "OK" || solvedProblems[id] {
			continue
		}

		y, m, d := time.Unix(v.CreationTimeSeconds, 0).Date()
		date := fmt.Sprintf("%d-%02d-%02d", y, m, d)
		lines = append(lines, fmt.Sprintf("%-9s%-50s%-12s", id, v.Problem.Name, date))

		solvedProblems[id] = true
	}

	fmt.Printf("%s solved %d problems.\n", handle, len(solvedProblems))
	fmt.Printf("%-9s%-50s%-12s\n", "ID", "Name", "Date")
	for _, l := range lines {
		fmt.Println(l)
	}
}
