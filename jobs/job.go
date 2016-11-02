// jobs 包来源于 https://github.com/revel/modules/jobs ,移除了 revel 环境依赖
package jobs

import (
	"log"
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"gopkg.in/robfig/cron.v2"
)

type Job struct {
	Name    string
	inner   cron.Job
	status  uint32
	running sync.Mutex
}

var (
	SelfConcurrent bool
	JobPoolSize    int
	MainCron       *cron.Cron
	workPermits    chan struct{}
)

const UNNAMED = "(unnamed)"
const DefaultJobPoolSize = 50

func init() {
	if JobPoolSize == 0 {
		JobPoolSize = DefaultJobPoolSize
	}
	workPermits = make(chan struct{}, JobPoolSize)
	MainCron = cron.New()
	MainCron.Start()
}

func New(job cron.Job) *Job {
	name := reflect.TypeOf(job).Name()
	if name == "Func" {
		name = UNNAMED
	}
	return &Job{
		Name:  name,
		inner: job,
	}
}

func (j *Job) Status() string {
	if atomic.LoadUint32(&j.status) > 0 {
		return "RUNNING"
	}
	return "IDLE"
}

func (j *Job) Run() {
	// If the job panics, just print a stack trace.
	// Don't let the whole process die.
	defer func() {
		if err := recover(); err != nil {
			log.Print(err, "\n", string(debug.Stack()))
		}
	}()

	if !SelfConcurrent {
		if atomic.LoadUint32(&j.status) > 0 {
			return
		}
		j.running.Lock()
		defer j.running.Unlock()
	}
	
	/*
	if !SelfConcurrent {
		j.running.Lock()
		defer j.running.Unlock()
	}*/
	

	if workPermits != nil {
		workPermits <- struct{}{}
		defer func() { <-workPermits }()
	}

	atomic.StoreUint32(&j.status, 1)
	defer atomic.StoreUint32(&j.status, 0)

	j.inner.Run()
}
