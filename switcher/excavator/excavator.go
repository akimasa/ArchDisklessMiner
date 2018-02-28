package excavator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var execCmd *exec.Cmd

func MakeJSON(algo string, port int, wallet string) {
	tmpl := `[
		{"time":0,"commands":[
			{"id":1,"method":"algorithm.add","params":["%s","%s.jp.nicehash.com:%d","%s.testrig"]}
		]},
		{"time":3,"commands":[
			%s
		]},
		{"time":10,"loop":10,"commands":[
			{"id":1,"method":"worker.print.speed","params":["0"]},
			{"id":1,"method":"algorithm.print.speeds","params":["0"]}
		]}
		]`
	jsonfile, _ := os.OpenFile("excavator-run.json", os.O_CREATE|os.O_WRONLY, 0666)
	jsonfile.Truncate(0)
	jsonfile.WriteString(fmt.Sprintf(tmpl, algo, algo, port, wallet, GetWorkerTemplate()))
	jsonfile.Close()
}

// BenchAlgo ... only support equihash now
func BenchAlgo(algo string) {
	BenchmarkAlgo(algo)
}

// currently only support equihash
func BenchmarkAlgo(algo string) (value float32) {
	execCmd = exec.Command("./excavator", "-p", "4000")
	stdout, err := execCmd.StdoutPipe()

	execCmd.Start()
	time.Sleep(time.Second * 5)

	defer execCmd.Process.Kill()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(stdout)

	// total speed: 390.490956 H/s
	re := regexp.MustCompile(`speed: (?P<rate>\d*.\d*?) H/s`)
	receiver := make(chan string)
	go func() {
		hr := ""
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Println(text)
			match := re.FindStringSubmatch(text)
			if match != nil {
				fmt.Println(match[1])
				hr = match[1]
				receiver <- hr
				return
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

	}()

	cmd := fmt.Sprintf(`{"id":1,"method":"algorithm.add","params":["%s","benchmark","test.testrig"]}`, algo)
	ret, err := GetJSON(cmd)
	if err != nil {
		panic(err)
	}
	wt := GetWorkerTemplate()
	wts := strings.Split(wt, ",\n")
	for _, v := range wts {
		ret, err := GetJSON(v)
		if err != nil {
			panic(err)
		}
		fmt.Println(ret)
	}
	time.Sleep(time.Second * 10)
	ret, err = GetJSON(`{"id":1,"method":"algorithm.print.speeds","params":["0"]}`)
	if err != nil {
		panic(err)
	}
	fmt.Println(ret)

	for {
		select {
		case recv := <-receiver:
			fmt.Println("recv", recv)
			hashrate, err := strconv.ParseFloat(recv, 32)
			if err != nil {
				fmt.Println("err", err)
				return 0
			}
			return float32(hashrate)
		case <-time.After(time.Second * 1):
			fmt.Println("timeout")
			return 0
		}
	}

}
func getDeviceList() string {

	StartExcavator()
	defer StopExcavator()
	ret, err := GetJSON(`{"id":1,"method":"device.list","params":[]}`)
	if err != nil {
		panic(err)
	}
	fmt.Println(ret)
	// ret := `{"devices":[{"device_id":0,"name":"GeForce GTX 1070","gpgpu_type":1,"subvendor":"10de","details":{"cuda_id":0,"sm_major":6,"sm_minor":1,"bus_id":9}}],"id":1,"error":null}`
	// ret := `{"devices":[{"device_id":0},{"device_id":1},{"device_id":2}]}`
	dec := json.NewDecoder(strings.NewReader(ret))
	var m interface{}
	err = dec.Decode(&m)
	if err != nil {
		panic(err)
	}
	jsonstr := ""
	// fmt.Println(m.(map[string]interface{})["devices"].([]interface{})[0].(map[string]interface{})["device_id"])
	for _, v := range m.(map[string]interface{})["devices"].([]interface{}) {
		deviceid := int(v.(map[string]interface{})["device_id"].(float64))
		// fmt.Println(k, v, deviceid)
		jsonstr += fmt.Sprintf(`{"id":1,"method":"worker.add","params":["0","%d"]},`, deviceid)
		jsonstr += "\n"
	}

	return jsonstr[0 : len(jsonstr)-2]
}
func GetWorkerTemplate() string {
	wtbyte, err := ioutil.ReadFile("workertemplate.txt")
	if os.IsNotExist(err) {
		devlist := getDeviceList()
		wtfile, err := os.OpenFile("workertemplate.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}
		wtfile.WriteString(devlist)
		return devlist
	} else {
		return string(wtbyte)
	}

}
func StartExcavator() {
	execCmd = exec.Command("./excavator", "-p", "4000")
	execCmd.Stderr = os.Stderr
	execCmd.Stdout = os.Stdout
	execCmd.Start()
	time.Sleep(time.Second * 5)
}
func GetJSON(cmd string) (res string, err error) {
	conn, err := net.Dial("tcp", "127.0.0.1:4000")
	if err != nil {
		return "", err
	}
	fmt.Fprintf(conn, "%s\n", cmd)
	ret, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	return ret, nil
}
func StopExcavator() {
	execCmd.Process.Kill()
}
func main() {

}
