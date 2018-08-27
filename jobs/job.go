// The MIT License (MIT)
//
// Copyright (C) 2012-2016 The Revel Framework Authors.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
