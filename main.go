package main

import (
	"fmt"
	"context"
	"os"
	"sync"
	"time"
	// "reflect"
	// "strings"
	// "bufio"

	"naver/database"
	"naver/scraper"
	"naver/domain"
)

// func getProxies() []string{
// 	var proxies []string
// 	var proxy string

// 	f, _ := os.Open("/Users/j/Desktop/config/proxy.txt")
// 	// check(err)
// 	fileScanner := bufio.NewScanner(f)
// 	fileScanner.Split(bufio.ScanLines)
// 	for fileScanner.Scan() {
// 		proxy = "http://" +fileScanner.Text()
// 		proxies = append(proxies,proxy)
//     }

// 	defer f.Close()
// 	return proxies
// }


func main()  {
	ctx := context.Background()
	client, err := database.GetClient(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()
	query:= client.Query(`
	WITH A AS (
		SELECT
			ROW_NUMBER() OVER (PARTITION BY product_id ORDER BY parttimer_searched_dt DESC) AS row_num,
			* EXCEPT(queenit_final_price)
		FROM` +"`damoa-fb351.naver_shopping.price_search_result_raw`" + `
	)
	SELECT
		* EXCEPT(row_num)
	FROM
		A
	WHERE
		row_num = 1
		AND naver_min_price_parttimer IS NOT NULL
	ORDER BY
		parttimer_searched_dt DESC
		LIMIT 30
	`)
	rows, err := query.Read(ctx)
	if err != nil {
		// TODO: Handle error.
		fmt.Println("errr")
	}
	data,err := database.GetData(os.Stdout, rows)
	if err != nil {
		fmt.Println(err)
	}
	// proxies := getProxies()
	ch1 := make(chan domain.NaverData)
	var wg sync.WaitGroup
	startTime := time.Now()
	for idx, temp := range data{
		if idx <5{
			fmt.Printf("#%d 진행중",idx)
			wg.Add(1)
			go scraper.GetData(&wg, ch1)
			ch1 <- temp
			
		}else{
			break
		}	
		
	}
	
	wg.Wait()
	fmt.Println("end", time.Since(startTime))
}