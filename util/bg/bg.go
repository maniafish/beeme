package bg

import (
	"beeme/util/mylog"
	"context"
	"runtime/debug"
	"sort"
	"sync"
	"time"
)

// bg 实现管理后台goroutine

var ctxGlobal, cancelGlobal = context.WithCancel(context.Background())
var wg = &sync.WaitGroup{}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
// code from https://stackoverflow.com/questions/32840687/timeout-for-waitgroup-wait
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

// StopAll stop all jobs
func StopAll(log mylog.MyLogger, timeout time.Duration) {
	jobsMu.Lock()
	stopped = true
	jobsMu.Unlock()

	cancelGlobal()

	if waitTimeout(wg, timeout) {
		log.Warnf("can not stop after %d second", timeout/time.Second)
	}
}

// Stopped check all jobs stopped
func Stopped() bool {
	return stopped
}

// JobFunc job function
type JobFunc func()

var (
	jobsMu  sync.Mutex
	jobs    = make(map[*Job]struct{})
	stopped bool
)

// Job a background job
type Job struct {
	ctx    context.Context
	cancel context.CancelFunc
	name   string
	fn     JobFunc

	startTime time.Time
	err       error

	log mylog.MyLogger
	// track
}

// JobStatus means a job status
type JobStatus struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"starttime"`
}

func newJob(ctx context.Context, name string, fn JobFunc, log mylog.MyLogger) *Job {
	job := &Job{
		name: name,
		fn:   fn,

		startTime: time.Now(),

		log: log,
	}

	if ctx == nil {
		job.ctx = ctxGlobal
		job.cancel = cancelGlobal
	} else {
		job.ctx, job.cancel = context.WithCancel(ctx)
	}

	return job
}

func (j *Job) run() {
	j.log.Infof("job: %s start runing at %v", j.name, j.startTime)
	register(j)
	wg.Add(1)

	go func() {

		defer func() {
			if err := recover(); err != nil {
				j.log.Infof("job: %s panic, err: %v, panic stack: %s",
					j.name, err, debug.Stack())
			}

			wg.Done()
			unregister(j)
			j.log.Infof("job: %s stop runing at %v", j.name, time.Now())
		}()

		for {
			select {
			case <-j.ctx.Done():
				j.log.Infof("job: %s ctx done: %v", j.name, j.ctx.Err())
				return
			default:
				j.fn()
			}
		}
	}()
}

// Err return job error
func (j *Job) Err() error {
	err := j.ctx.Err()
	if err == nil {
		err = j.err
	}
	return err
}

type byName []JobStatus

func (a byName) Len() int {
	return len(a)
}

func (a byName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

// Run run a background job
func Run(log mylog.MyLogger, name string, fn JobFunc) *Job {

	job := newJob(nil, name, fn, log)
	job.run()

	return job
}

func register(job *Job) {
	jobsMu.Lock()
	jobs[job] = struct{}{}
	if stopped {
		panic("when called StopAll(), can't call Run()")
	}
	jobsMu.Unlock()
}

func unregister(job *Job) {
	jobsMu.Lock()
	delete(jobs, job)
	if len(jobs) == 0 {
		stopped = true
	}
	jobsMu.Unlock()
}

// Status get jobs status
func Status() []JobStatus {
	jobsMu.Lock()
	defer jobsMu.Unlock()

	js := make([]JobStatus, 0, 16)
	for k := range jobs {
		js = append(js, JobStatus{
			Name:      k.name,
			StartTime: k.startTime,
		})
	}

	sort.Sort(byName(js))
	return js
}

// SleepWithWaitCtx sleep for a while with waiting context canceled, return true if context canceled
func SleepWithWaitCtx(t time.Duration) bool {
	select {
	case <-ctxGlobal.Done():
		return true
	case <-time.After(t):
		return false
	}
}
