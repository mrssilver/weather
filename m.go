package main
//weatherapi set
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"text/tabwriter"
	"time"
)

var	apiKey = ""
// 定义5天天气预报API响应结构体
type WeatherForecastResponse struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`
	Forecast struct {
		ForecastDay []struct {
			Date string `json:"date"`
			Day  struct {
				MaxTempC      float64 `json:"maxtemp_c"`
				MaxTempF      float64 `json:"maxtemp_f"`
				MinTempC      float64 `json:"mintemp_c"`
				MinTempF      float64 `json:"mintemp_f"`
				AvgTempC      float64 `json:"avgtemp_c"`
				AvgTempF      float64 `json:"avgtemp_f"`
				Condition     struct {
					Text string `json:"text"`
				} `json:"condition"`
				MaxWindKph  float64 `json:"maxwind_kph"`
				TotalPrecip float64 `json:"totalprecip_mm"`
				AvgHumidity float64 `json:"avghumidity"`
				UV          float64 `json:"uv"`
			} `json:"day"`
		} `json:"forecastday"`
	} `json:"forecast"`
}


func main() {

	// 定义城市数组
	cities := []string{"Shanghai(china)", "Hangzhou(china)", "Jiujiang(china)"}
	
	// 获取5天天气预报
	for _, city := range cities {
		fmt.Printf("\n正在获取 %s 的天气预报...\n", city)
		weatherData, err := getWeatherForecast(apiKey, city)
		if err != nil {
			fmt.Printf("获取 %s 天气数据失败: %v\n", city, err)
			continue
		}
		
		// 打印天气信息表格
		printWeatherForecastTable(weatherData, city)
		
		// 添加延迟以避免API速率限制
		time.Sleep(1 * time.Second)
	}
}

// 获取5天天气预报
func getWeatherForecast(apiKey, city string) (*WeatherForecastResponse, error) {
	// 构建API请求URL
	baseURL := "http://api.weatherapi.com/v1/forecast.json"
	params := url.Values{}
	params.Add("key", apiKey)
	params.Add("q", city)
	params.Add("days", "5")
	params.Add("aqi", "no")
	params.Add("alerts", "no")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 发送HTTP请求
	response, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求错误: %s, 响应内容: %s", response.Status, string(body))
	}

	// 解析JSON响应
	var weatherData WeatherForecastResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return &weatherData, nil
}

// 打印5天天气预报表格
func printWeatherForecastTable(weatherData *WeatherForecastResponse, city string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	
	fmt.Fprintf(w, "\n%s (%s, %s) 5天天气预报\n", 
		weatherData.Location.Name, 
		weatherData.Location.Region, 
		weatherData.Location.Country)
	fmt.Fprintf(w, "查询时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "==============\t==============\t==============\t==============\t==============\t==============\t==============\n")
	fmt.Fprintf(w, "日期\t\t最高温度\t最低温度\t平均温度\t天气状况\t降水量(mm)\t湿度\n")
	fmt.Fprintf(w, "==============\t==============\t==============\t==============\t==============\t==============\t==============\n")
	
	for _, day := range weatherData.Forecast.ForecastDay {
		// 将日期格式化为更易读的形式
		date, err := time.Parse("2006-01-02", day.Date)
		if err != nil {
			fmt.Fprintf(w, "%s\t", day.Date)
		} else {
			fmt.Fprintf(w, "%s\t", date.Format("01/02"))
		}
		
		fmt.Fprintf(w, "%.1f°C\t", day.Day.MaxTempC)
		fmt.Fprintf(w, "%.1f°C\t", day.Day.MinTempC)
		fmt.Fprintf(w, "%.1f°C\t", day.Day.AvgTempC)
		fmt.Fprintf(w, "%s\t", day.Day.Condition.Text)
		fmt.Fprintf(w, "%.1f\t", day.Day.TotalPrecip)
		fmt.Fprintf(w, "%.0f%%\n", day.Day.AvgHumidity)
	}
	
	fmt.Fprintf(w, "==============\t==============\t==============\t==============\t==============\t==============\t==============\n")
	w.Flush()
}

