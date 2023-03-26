package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	appealCount int
	mutex       sync.Mutex
	cookie      = "YOUR_BILIBILI_JCT_AND_SESSDATA"
)

func Appeal(aid string) {
	urlStr := "https://api.bilibili.com/x/web-interface/archive/appeal?jsonp=jsonp&csrf=097af8caab0ffc4e26f6be4b557ab972"
	queryParams := url.Values{}
	queryParams.Set("aid", aid)
	queryParams.Set("attach", "")
	queryParams.Set("desc", "标题")
	queryParams.Set("tid", "5")

	requestBody := bytes.NewBufferString(queryParams.Encode())
	//Create POST request

	req, err := http.NewRequest("POST", urlStr, requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error decoding response body:", err)
		return
	}
	if data["code"] == float64(0) {
		mutex.Lock()
		appealCount++
		mutex.Unlock()
	}

}

func SearchAndappeal(page int) {
	urlStr := "https://api.bilibili.com/x/web-interface/search/type"

	queryParams := url.Values{}
	queryParams.Set("search_type", "video")
	queryParams.Set("keyword", "收米直播 b775")
	queryParams.Set("page", strconv.Itoa(page))

	fullUrl := fmt.Sprintf("%s?%s", urlStr, queryParams.Encode())

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	//fmt.Println("Response Status:", resp.Status)
	//fmt.Println("Response Body:", string(body))

	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	data, ok := jsonData["data"].(map[string]interface{})
	if !ok {
		fmt.Println("'data' not found or is not an object")
		return
	}

	// 寻找'result'字段
	results, ok := data["result"].([]interface{})
	if !ok {
		fmt.Println("'result' not found or is not an array")
		return
	}

	for _, result := range results {
		res, ok := result.(map[string]interface{})
		if !ok {
			fmt.Println("Error: result is not a map")
			continue
		}
		title := res["title"].(string)
		//strings.Contains(title, "\uE708")
		if strings.Contains(title, "2023已更新") || strings.Contains(title, "网站") || strings.Contains(title, "体育") || strings.Contains(title, "app") {
			aidFloat, ok := res["aid"].(float64)
			if !ok {
				panic("aid is not of type float64")
			}
			aidString := fmt.Sprintf("%.0f", aidFloat)
			Appeal(aidString)
			fmt.Printf("Appeal Aid: %s sent, waiting 2 seconds...\n", aidString)
			time.Sleep(2 * time.Second)
		}

		//	}
		//}
	}
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 10) // 创建一个具有10个缓冲区的通道

	for i := 1; i < 50; i++ {
		wg.Add(1)
		ch <- struct{}{} // 向通道发送一个值，以便让另一个goroutine可以开始执行
		go func(i int) {
			defer wg.Done()
			SearchAndappeal(i)
			<-ch // 当goroutine完成时，从通道中接收一个值
		}(i)
	}

	wg.Wait()
	print("Video appealed: ", appealCount)
}
