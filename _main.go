package main

import (
	"fmt"
	// "io/ioutil"
	"os"
	"bufio"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main()  {
	var proxies []string
	var proxy string
	// dat, err := ioutil.ReadFile("/Users/j/Desktop/config/proxy.txt")
	// check(err)
	// fmt.Print(string(dat))
	f, err := os.Open("/Users/j/Desktop/config/proxy.txt")
	check(err)
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		proxy = "http://" +fileScanner.Text()
		proxies = append(proxies,proxy)
    }

	defer f.Close()
	// b4, err := r4.Peek(5)
	// check(err)
	fmt.Println(proxies)
}