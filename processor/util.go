package processor

import (
	"fmt"
	"github.com/margostino/anfield/domain"
)

func printCommentary(h2h string, commentary *domain.Commentary) {
	if isTimedComment(commentary) {
		fmt.Printf("[%s] # %s - %s\n", h2h, commentary.Time, commentary.Comment)
	} else {
		fmt.Printf("[%s] # %s\n", h2h, commentary.Comment)
	}
}

func end(commentary *domain.Commentary) bool {
	return commentary == nil || (commentary.Time == "end")
}

func notStarted(commentary *domain.Commentary) bool {
	return commentary.Time == "not-started"
}

func isTimedComment(commentary *domain.Commentary) bool {
	return commentary != nil && commentary.Time != "" && commentary.Comment != ""
}
