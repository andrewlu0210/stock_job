package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/andrewlu0210/stock"

	"github.com/spf13/viper"
)

const (
	layout = "20060102"
	limit  = 366 //一次下載天數最多365天
)

var (
	DBHost    string
	DBName    string
	DBAccount string
	DBPasswd  string
	CSVRoot   string //下載之CVS檔案存放目錄
	startDate string
	count     int
)

func main() {
	viperDir := "/home/andrew/stock"
	initViper(viperDir)
	if len(os.Args) > 1 {
		startDate = checkStartDate(os.Args[1])
		if len(os.Args) > 2 {
			count = getCount(os.Args[2])
		} else {
			count = 1
		}
	} else {
		startDate = checkStartDate("")
		count = 1
	}

	fmt.Printf("準備下載Stock CSV檔，並存入MongoDB，從%s開始，共%d筆\n", startDate, count)

	stock.SetMongo(DBHost, DBName, DBAccount, DBPasswd)
	stock.ConnectDb()
	defer stock.DisconnectDb()
	stockDownloader := stock.GetPriceService().GetDownloader(CSVRoot)

	//prepare date string
	date, _ := time.Parse(layout, startDate)
	saveCnt := 0
	for i := 0; i < count; i++ {
		is_saved, is_download := stockDownloader.DownloadStockPrice(date.Format(layout))
		if is_download {
			//sleep and download again
			time.Sleep(1500 * time.Millisecond)
		}
		if is_saved {
			saveCnt += 1
		}
		date = date.Add(24 * time.Hour)

	} // end for loop
	fmt.Printf("總共存入%d天資料\n", saveCnt)
}

func initViper(viperConfigDir string) {
	viper.SetConfigName("stock")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if viperConfigDir != "" {
		viper.AddConfigPath(viperConfigDir)
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	DBHost = viper.GetString("stock_db.host")
	DBName = viper.GetString("stock_db.db_name")
	DBAccount = viper.GetString("stock_db.db_account")
	DBPasswd = viper.GetString("stock_db.db_passwd")
	CSVRoot = viper.GetString("csvRoot")
}

func checkStartDate(dateStr string) string {
	now := time.Now()
	date, err := time.Parse(layout, dateStr)
	if err != nil || date.After(now) {
		return now.Format(layout)
	}
	return dateStr
}

func getCount(cntStr string) int {
	val, err := strconv.ParseInt(cntStr, 10, 32)
	if err != nil {
		return 1
	}
	if val <= 0 {
		return 1
	} else if val > limit {
		return limit
	}
	return int(val)
}
