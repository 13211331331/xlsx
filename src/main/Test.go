package main

import (
	"bar"
	"time"
	//"sync"
	"fmt"
	"myutil"
)
func main() {


	fmt.Println(myutil.CountToExcel(1))



	b1 := bar.New("多线程进度条").NewBar("export ... ", 1000)



	//var wg2 sync.WaitGroup
	//wg2.Add(1)
	go func() {
		//defer wg2.Done()
		for i := 0; i < 1000; i++ {
			b1.AddNumber(1)

			time.Sleep(time.Second / 2000)
		}
	}()


	//wg2.Wait()

	fmt.Println("3333333333333")
}