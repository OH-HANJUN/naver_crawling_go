package scraper

import (
	"fmt"
	"sort"
	"strings"
	"encoding/json"
	"unicode/utf8"
	"sync"

	"naver/domain"
	"github.com/gocolly/colly"
	"github.com/mitchellh/mapstructure"
)

type CatalogData struct {
	Props struct {
		PageProps struct {
			InitialState struct {
				Catalog struct {
					Products []interface{} `json:"products"`
				} `json:"catalog"`
			} `json:"initialState"`
		} `json:"pageProps"`
	} `json:"props"`
}

type ShoppingData struct {
	Props struct {
		PageProps struct {
			InitialState struct {
				Products struct {
					List []interface{} `json:"list"`
				} `json:"products"`
			} `json:"initialState"`
		} `json:"pageProps"`
	} `json:"props"`

}

	
type CatalogReturn struct {
	PcPrice     int    
	MallName    string 
	ProductName string 
}

type ShoppingItem struct {
	LowPrice int    `json:"lowPrice"`
	MallName string `json:"mallName"`
}
type LowMallLists struct {
	Name string `json:"name"`
}

type ShoppingItems struct {
	Item struct {
		LowPrice int    `json:"lowPrice"`
		MallName string `json:"mallName"`
		LowMallList []LowMallLists  `json:"lowMallList"`
	} `json:"item"`
}
// type ChType[T domain.NaverData | []CatalogReturn ] chan T


func GetData(wg *sync.WaitGroup, ch chan domain.NaverData){
	for data := range ch{
		fmt.Println(data)
		var shoppingData ShoppingItem
		ch := make( chan interface{} )
		go catalog_test(ch)
		ch <- data
		if strings.HasPrefix(data.NaverShoppingURL, "http"){
			shoppingData = naverShopping(data)
			fmt.Println(shoppingData,"shoppping")
		}
		fmt.Println(<-ch ,"testtttt")
		close(ch)
		
		// if len(data_list) > 0{
		// 	sort.Slice(data_list, func(i,j int) bool {
		// 		return data_list[i].PcPrice < data_list[j].PcPrice})
		// }
			
		// testChan := make(chan chan interface{} )
		
		// close(ch)
		// testChan2 <- data
		// testChan <- testChan2
		// close(tes?tChan2) 
		
		
		
		wg.Done()
	}	
	
}

func catalog_test(ch chan interface{} ){
	var resultArray [3]CatalogReturn
	for d := range ch {
		data := d.(domain.NaverData)
		if data.NaverCatalogID1.Valid && utf8.RuneCountInString(data.NaverCatalogID1.StringVal) >3 {
			id1 :=catalog("https://search.shopping.naver.com/catalog/"+data.NaverCatalogID1.StringVal)
			fmt.Println(id1,"#1")
			resultArray[0] = id1
		}
		if data.NaverCatalogID2.Valid && utf8.RuneCountInString(data.NaverCatalogID2.StringVal) >3 {
			id2 :=catalog("https://search.shopping.naver.com/catalog/"+data.NaverCatalogID2.StringVal)
			fmt.Println("#2",id2)
			resultArray[1] = id2
		}
		if data.NaverCatalogID3.Valid && utf8.RuneCountInString(data.NaverCatalogID3.StringVal) >3 {
			id3 :=catalog("https://search.shopping.naver.com/catalog/"+data.NaverCatalogID3.StringVal)
			fmt.Println("#2",id3)
			resultArray[2] = id3
		}
		ch <- resultArray
	}
 }



func catalog(url string) CatalogReturn{
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"))
	var data_list  []CatalogReturn
	c.OnHTML("script[id='__NEXT_DATA__']", func(e *colly.HTMLElement) {
			containerScript := e.Text
			JSONparse := CatalogData{}
			json.Unmarshal([]byte(containerScript), &JSONparse)
			products := JSONparse.Props.PageProps.InitialState.Catalog.Products
			
			for idx, product := range products {
				if idx >= 5{
					break
				}
				productsPage, ok  := product.(map[string]interface{})["productsPage"]
				if !ok || productsPage == nil{
					continue 
				}
				
				products2, ok := productsPage.(map[string]interface{})["products"].([]interface{})
				if !ok || len(products2) ==0{
					continue 
				}
				for _, temp := range products2{
	
					temp2, _  := temp.(map[string]interface{})
					temp_data := &CatalogReturn{}
					if err := mapstructure.Decode(temp2,&temp_data); err != nil {
						fmt.Println(err)
					}
					data_list = append(data_list,*temp_data)
	
				}
							
			}
			sort.Slice(data_list, func(i,j int) bool {
				return data_list[i].PcPrice < data_list[j].PcPrice})
			})
		c.Visit(url)
		
		return data_list[0]
	
}

func naverShopping(data domain.NaverData) ShoppingItem{
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"))
	var shoppingArray  []ShoppingItem
	c.OnHTML("script[id='__NEXT_DATA__']", func(e *colly.HTMLElement) {
		containerScript := e.Text
		JSONparse := ShoppingData{}
		json.Unmarshal([]byte(containerScript), &JSONparse)
		list := JSONparse.Props.PageProps.InitialState.Products.List
		
		
		for idx, item := range list{
			if idx >= 5{
				break
			}
			temp := &ShoppingItems{}
			mapstructure.Decode(item,&temp)
			mallName := temp.Item.MallName
			if len(temp.Item.LowMallList) > 0{
				mallName = temp.Item.LowMallList[0].Name
			}

			temp_data := ShoppingItem{LowPrice : temp.Item.LowPrice, MallName: mallName}
			shoppingArray = append(shoppingArray,temp_data)
		}
		if len(shoppingArray) >0{
			sort.Slice(shoppingArray, func(i,j int) bool {
				return shoppingArray[i].LowPrice < shoppingArray[j].LowPrice})
		}else{
			shoppingArray = append(shoppingArray,ShoppingItem{LowPrice : 999999, MallName: "없음"})
		}
		
		
	})
	c.Visit(data.NaverShoppingURL)
	return shoppingArray[0]

}