package processor

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/context"
	"sync"
)

var done = make(chan bool)

// Process TODO: spawn one process per URL
func Process(urls []string) {
	wg := common.WaitGroup(len(urls))
	for _, url := range urls {
		go async(url, wg)
	}
	wg.Wait()
}

func async(url string, waitGroup *sync.WaitGroup) {
	commentaryUrl := url + context.Config().Commentary.Params
	go metadata(url)
	go commentary(commentaryUrl)
	go consume()
	<-done
	waitGroup.Done()
}
