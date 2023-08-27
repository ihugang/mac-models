package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type MacModel struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}

func GetMacModel() string {
	return "Mac"
}

// 获取最新的JSON文件
func GenerateModelsJson() {
	url := "https://everymac.com/systems/by_capability/mac-specs-by-machine-model-machine-id.html"

	// html := getUrl(url)
	// fmt.Println(html)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	// 设置User-Agent来模拟Chrome浏览器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")

	// 设置一些其他常见请求头
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	array := []MacModel{}

	doc.Find("#contentcenter_specs_externalnav_wrapper").Each(func(i int, s *goquery.Selection) {
		line := s.Text()
		lines := strings.Split(line, "\n")
		if len(lines) > 3 {
			fmt.Println(fmt.Sprintf("%s => %s", lines[3], lines[2]))
			item := MacModel{}
			item.Name = lines[2]
			item.Identifier = lines[3]
			array = append(array, item)
		}
	})

	//转换为JSON
	jsonBytes, err := json.Marshal(array)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	fileName := "mac-models.json"
	err = os.WriteFile(fileName, jsonBytes, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("done.")
}

func GetMacInfo() string {
	out, err := exec.Command("sysctl", "-n", "hw.model").Output()
	if err != nil {
		fmt.Printf("Failed to get model: %v\n", err)
		return ""
	}
	identifier := strings.TrimSpace(string(out))

	models := []MacModel{}

	jsonFile := "mac-models.json"
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return ""
	}
	err = json.Unmarshal(data, &models)
	if err != nil {
		fmt.Printf("Failed to unmarshal file: %v\n", err)
		return ""
	}

	for _, model := range models {
		if model.Identifier == identifier {
			return model.Name
		}
	}

	fmt.Println("Not found.")
	return ""
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	jsonFile := "mac-models.json"
	if !fileExists(jsonFile) {
		GenerateModelsJson()
	}
	var modelName = GetMacInfo()
	fmt.Println("Current Mac model: " + modelName)
}
