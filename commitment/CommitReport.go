package commitment

import "fmt"

// A CommitReport exposes commit properties used for reporting to users
type CommitReport struct {
	Message, URL, Date string
}

// Print outputs a summary to the console for this CommitReport
func (c *CommitReport) Print() {
	fmt.Println("Commit message: ", c.Message)
	fmt.Println("Commit date: ", c.Date)
	fmt.Println("View changes: ", c.URL)
}
