package main

import (
	"download/downloader"
	"fmt"
)

func main() {
    var url string = "https://dl.softmgr.qq.com/original/im/QQ9.5.0.27852.exe"
    
    httpDownload := downloader.New(url, 4)
    fmt.Printf("Bool:%v\nContent:%d\n", httpDownload.AcceptRanges(), httpDownload.ContentLength())

    httpDownload.Download()
}