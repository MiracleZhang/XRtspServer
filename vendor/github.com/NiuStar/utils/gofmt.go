package utils

import (
	"bufio"
	"fmt"

	"io"
	"os/exec"
	//"os"
)



func Gofmt(filename ,path string) (bool,string) {
	{
		return execCommand("cmd", "/c", "gofmt",  "-w", filename, path)
	}
}

func execCommand(commandName string, params ...string) (bool, string) {
	cmd := exec.Command(commandName, params...)

	//显示运行的命令
	fmt.Println(cmd.Args)

	/*stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println("err: ",err)
		return false,err.Error()
	}*/
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("err2", err)
		return false, err.Error()
	}

	cmd.Start()
	//fmt.Println("stdout:",getOutput(stdout))
	//fmt.Println("stderr:",getOutput(stderr))
	//go io.Copy(os.Stdout, stdout)
	//go io.Copy(os.Stderr, stderr)
	result_err := getOutput(stderr)
	cmd.Wait()
	return true, result_err
}

func getOutput(out io.ReadCloser) string {
	reader := bufio.NewReader(out)
	var result string
	//实时循环读取输出流中的一行内容
	for {
		line, ok, err2 := reader.ReadLine()
		fmt.Println("OK:", ok)
		if err2 != nil || io.EOF == err2 {
			fmt.Println("err2:", err2)
			break
		}
		result += string(line) + "\r\n"
		//fmt.Println("line:",string(line))
	}
	return result
}
