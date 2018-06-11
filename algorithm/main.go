package main

import (
	"log"
	// "go-algorithms/sort/bubble"

	"golibrary/algorithm/utils"
	"golibrary/algorithm/quick"
)

func main() {
	list := utils.GetArrayOfSize(10)
	log.Println(list)
	quick.HoareSort(list, 0, len(list)-1)
	// quick.LomutoSort(list, 0, len(list)-1)
	// quick.Sort(list)
	// bubble.Sort(list)
	log.Println(list)

}
