package jobs

/*
import (
	"fmt"

	"gopkg.in/robfig/cron.v2"
)

const DEFAULT_JOB_POOL_SIZE = 10

var (
	// Singleton instance of the underlying job scheduler.
	MainCron *cron.Cron

	// This limits the number of jobs allowed to run concurrently.
	workPermits chan struct{}

	// Is a single job allowed to run concurrently with itself?
	selfConcurrent bool
)

func main() {
	MainCron = cron.New()
	workPermits = make(chan struct{}, DEFAULT_JOB_POOL_SIZE)
	selfConcurrent = true
	MainCron.Start()
	fmt.Println("Go to /@jobs to see job status.")
}

*/
