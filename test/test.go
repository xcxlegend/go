package main

/*
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type LineData struct {
	Event      string            `json:"event"`
	OtherEvent string            `json:"otherEvent"`
	Data       map[string]string `json:"data"`
}

type ReduceMap map[string]int64

func main() {
	const BASE_DIR = "/data/aws-kinesis-agent-user/statis/"
	var ws = new(sync.WaitGroup)
	files, _ := ioutil.ReadDir(BASE_DIR)
	fmt.Println("file numbers:", len(files))
	var count = 0
	var mapC = make(chan ReduceMap, len(files))
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "statis") && strings.HasSuffix(f.Name(), ".log") {
			ws.Add(1)
			count++
			go MapReduce(BASE_DIR+f.Name(), mapC, ws)
			// fmt.Println(BASE_DIR + f.Name())
		}
	}
	var total = ReduceMap{}
	ws.Add(1)
	go func() {
		var i = 0
		fmt.Println("total count:", count)
		for m := range mapC {
			i++
			fmt.Println("file over:", i)
			for k, v := range m {
				total[k] += v
			}
			if i >= count {
				break
			}
		}
		ws.Done()
	}()
	ws.Wait()
	// fmt.Println(total)
	for k, v := range total {
		fmt.Printf("%s = %d\r\n", k, v)
	}
}

func MapReduce(file string, c chan ReduceMap, ws *sync.WaitGroup) {
	defer ws.Done()
	var m = make(ReduceMap)
	f, err := os.Open(file)
	if err != nil {
		c <- m
		return
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil || io.EOF == err {
			break
		}
		var data = new(LineData)
		err = json.Unmarshal([]byte(line), data)
		if err != nil {
			continue
		}
		if data.Event == "sys" && data.OtherEvent == "api" {
			if api, ok := data.Data["api"]; ok {
				apis := strings.SplitN(api, "?", -1)
				m[apis[0]]++
			}
		}
	}
	c <- m
}
*/
