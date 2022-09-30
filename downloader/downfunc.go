package downloader

import (
	"download/errcheck"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

// 初始化下载器
func New(url string, numThreads int) *HttpDownloader {
	var urlSplits []string = strings.Split(url, "/")
	var filename string = urlSplits[len(urlSplits)-1]

	res, err := http.Head(url)
	errcheck.Check(err)

	acceptRange := false
	if len(res.Header["Accept-Ranges"]) != 0 && res.Header["Accept-Ranges"][0] == "bytes" {
		acceptRange = true
	}

	httpDownload := HttpDownloader{
		url:           url,
		contentLength: int(res.ContentLength),
		numThreads:    numThreads,
		filename:      filename,
		acceptRanges:  acceptRange,
	}

	return &httpDownload
}

// 下载综合调度
func (h *HttpDownloader) Download() {
	f, err := os.Create(h.filename)
	errcheck.Check(err)
	defer f.Close()

	if h.acceptRanges == false {
		fmt.Println("该文件不支持多线程下载，单线程下载中：")
		resp, err := http.Get(h.url)
		errcheck.Check(err)
		savefile(h.filename, 0, resp)
	} else {
		var wg sync.WaitGroup
		for _, ranges := range h.Split() {
			fmt.Printf("多线程下载中:%d-%d\n", ranges[0], ranges[1])
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				h.download(start, end)
			}(ranges[0], ranges[1])
		}
		wg.Wait()
	}
}

// 多线程下载
func (h *HttpDownloader) download(start, end int) {
	req, err := http.NewRequest("GET", h.url, nil)
	errcheck.Check(err)
	req.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", start, end))
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	errcheck.Check(err)
	defer resp.Body.Close()

	savefile(h.filename, int64(start), resp)
}

// 下载文件分段
func (h *HttpDownloader) Split() [][]int {
	ranges := [][]int{}
	blockSize := h.contentLength / h.numThreads
	for i := 0; i < h.numThreads; i++ {
		var start int = i * blockSize
		var end int = (i+1)*blockSize - 1
		if i == h.numThreads-1 {
			end = h.contentLength - 1
		}
		ranges = append(ranges, []int{start, end})
	}
	return ranges
}

// 保存文件
func savefile(filename string, offset int64, resp *http.Response) {
	f, err := os.OpenFile(filename, os.O_WRONLY, 0660)
	errcheck.Check(err)
	f.Seek(offset, 0)
	defer f.Close()

	content, err := ioutil.ReadAll(resp.Body)
	errcheck.Check(err)
	f.Write(content)
}
