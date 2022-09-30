package errcheck

import "log"

// 捕获错误
func Check(e error) {
	if e != nil {
		log.Println(e)
		panic(e)
	}
}