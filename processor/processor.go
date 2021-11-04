package processor

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"sync"
)

var waitGroups = make(map[string]*sync.WaitGroup, 0)
var commentaryBuffer = make(map[string]chan *domain.Commentary)
var metadataBuffer = make(map[string]chan *domain.Metadata)

func Process(urls []string) {
	wg := common.WaitGroup(len(urls))
	for _, url := range urls {
		go async(url, wg)
	}
	wg.Wait()
}

func async(url string, waitGroup *sync.WaitGroup) {
	waitGroups[url] = common.WaitGroup(3)
	commentaryBuffer[url] = make(chan *domain.Commentary)
	metadataBuffer[url] = make(chan *domain.Metadata)

	go produceMetadata(url)
	go produceCommentary(url)
	// TODO: consumer does not need be a goroutine if it implements a infinite loop, unless we want extra process after that.
	go consume(url)

	waitGroups[url].Wait()
	waitGroup.Done()
}
