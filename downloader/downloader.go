package downloader

const (
	userAgent = `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36`
)

type HttpDownloader struct {
	url           string
	filename      string
	contentLength int
	acceptRanges  bool // 是否支持断点续传
	numThreads    int  // 同时下载线程数
}

func (h *HttpDownloader) Url() string {
	return h.url
}

func (h *HttpDownloader) Filename() string {
	return h.filename
}

func (h *HttpDownloader) ContentLength() int {
	return h.contentLength
}

func (h *HttpDownloader) AcceptRanges() bool {
	return h.acceptRanges
}

func (h *HttpDownloader) NumThreads() int {
	return h.numThreads
}