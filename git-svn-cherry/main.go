package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var versionTarget = "1.17-1.17.5"
var versionRoot = "34065"

//@TEST
// 34191~34341,34413~34859,
//35142,35380,35387,35426,35584,35593,35596,35876,35955,36053                                                                                                                 //checkout
//36056
//36064,36085,36086,36838
//36849,37050,37073,37086
var versionAry = "" //cherry-pick
var workRoot = ""   //work dir

func main() {
	svnIDHashIDMap := map[string]string{}
	tmHashID := ""

	csvFile, err := os.Open("./commit.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			x := fmt.Sprintf("%s", err)
			if !strings.Contains(x, "wrong number of fields") {
				log.Fatal(err)
			}
			// fmt.Println(err)
		}

		for _, tmRecord := range record {
			if strings.Contains(tmRecord, "commit") {
				index := strings.Index(tmRecord, "commit")
				hashID := tmRecord[index+7 : len(tmRecord)]
				// fmt.Println(hashId)
				tmHashID = hashID
			}

			if strings.Contains(tmRecord, "git-svn-id") {
				index := strings.Index(tmRecord, "@")
				svnID := tmRecord[index+1 : index+6]
				svnIDHashIDMap[svnID] = tmHashID
				fmt.Println(svnID + "-" + tmHashID)
			}
		}
	}

	//access dir
	os.Chdir(workRoot)
	//checkout new branch
	execCommand("git", []string{"checkout", "-b", versionTarget, svnIDHashIDMap[versionRoot]})
	versionSplit := strings.Split(versionAry, ",")
	for _, s := range versionSplit {
		sAry := strings.Split(s, "~")
		fmt.Println(sAry)
		if len(sAry) == 2 {
			execCommand("git", []string{"cherry-pick", svnIDHashIDMap[sAry[0]] + "^.." + svnIDHashIDMap[sAry[1]]})
		} else {
			execCommand("git", []string{"cherry-pick", svnIDHashIDMap[sAry[0]]})
		}
	}
}

func execCommand(commandName string, params []string) bool {
	cmd := exec.Command(commandName, params...)
	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}

	cmd.Start()
	reader := bufio.NewReader(stdout)

	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	cmd.Wait()
	return true
}
