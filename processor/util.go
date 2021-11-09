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
	if commentary == nil || (commentary.Time == "end" && commentary.Comment == "end") {
		return true
	}
	return false
}

func isTimedComment(commentary *domain.Commentary) bool {
	if commentary != nil {
		if commentary.Time != "" && commentary.Comment != "" {
			return true
		}
	}
	return false
}
