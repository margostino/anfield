package processor

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"sync"
)

var asyncPool = make(map[string]chan bool, 0)
var commentaryBuffer = make(map[string]chan *domain.Commentary)
var metadataBuffer = make(map[string]chan *domain.Metadata)

// Process TODO: spawn one process per URL
func Process(urls []string) {
	wg := common.WaitGroup(len(urls))
	for _, url := range urls {
		go async(url, wg)
	}
	wg.Wait()
}

func async(url string, waitGroup *sync.WaitGroup) {
	var done = make(chan bool)
	asyncPool[url] = done
	commentaryBuffer[url] = make(chan *domain.Commentary)
	metadataBuffer[url] = make(chan *domain.Metadata)

	go metadata(url)
	go commentary(url)
	go consume(url)

	<-done
	waitGroup.Done()
}
