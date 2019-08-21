package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)
//need file "localsVariable.txt" "excludeFiles.txt"
var originDir = "rootDir"
var localFile = "localsVariable.txt"
var formDir = ""
var toDir = ""

func main() {
	fmt.Println(os.Args)
	// useDir := "use"
	// bakDir := "bak"
	param := "no"
	if len(os.Args) > 1 {
		param = os.Args[1]
	}

	cpIndex := -1
	if param == "o2b" {
		cpIndex = 3
	} else if param == "o2u" {
		cpIndex = 2
	} else if param == "b2o" {
		cpIndex = 1
	} else if param == "u2o" {
		cpIndex = 0
	} else {
		// var err error
		readLine(localFile, func(line []byte) {
			str := string(line)
			cpIndex, _ = strconv.Atoi(str)
		})
	}

	// fmt.Printf("cpIndex:%d\n", cpIndex)
	if cpIndex == 3 { //origin to bak
		toDir = path.Join(originDir, "excludeFiles/bak/")
		formDir = originDir
	} else if cpIndex == 2 { //origin to use
		toDir = path.Join(originDir, "excludeFiles/use/")
		formDir = originDir
	} else if cpIndex == 1 { //bak to origin
		formDir = path.Join(originDir, "excludeFiles/bak/")
		toDir = originDir
		var data = []byte("0")
		saveFile(localFile, data)
	} else if cpIndex == 0 { //use to origin
		formDir = path.Join(originDir, "excludeFiles/use/")
		toDir = originDir
		var data = []byte("1")
		saveFile(localFile, data)
	} else {
		return
	}

	readLine("excludeFiles.txt", processLine)
}

func processLine(line []byte) {
	str := string(line)
	str = strings.Replace(str, "\n", "", -1)

	checkAndMakeDir(str)
	strSrc := path.Join(formDir, str)
	strDst := path.Join(toDir, str)
	copyFile(strSrc, strDst)
}

func readLine(filePth string, hookfn func([]byte)) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		hookfn(line)    //放在错误处理前面，即使发生错误，也会处理已经读取到的数据。
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func copyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		fmt.Println("src Open---", err)
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		fmt.Println("dst Open---", err)
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func checkAndMakeDir(targetDir string) {
	sDir, _ := path.Split(targetDir)
	isDirExit, err := pathExists(path.Join(toDir, sDir))
	if isDirExit == false && err == nil {
		createDir(path.Join(toDir, sDir))
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func createDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	return err
}

func saveFile(filename string, body []byte) error {
	return ioutil.WriteFile(filename, body, 0600)
}
