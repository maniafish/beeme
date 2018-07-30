package bg

import (
	"beeme/util/mylog"
	"os"
	"testing"
	"time"
)

func testJob() {}

func TestStatus(t *testing.T) {
	if len(Status()) != 0 {
		t.Errorf("Status() length != 0")
	}

	log := mylog.New("INFO", os.Stdout, 0)

	Run(log, "test1", testJob)
	Run(log, "test2", testJob)

	if len(Status()) != 2 {
		t.Errorf("Status() length != 2, %v", Status())
	}

	StopAll(log, time.Second)

	time.Sleep(1 * time.Second)
	if len(Status()) != 0 {
		t.Errorf("Status() length != 0, %v", Status())
	}

	if !Stopped() {
		t.Errorf("not stopped")
	}

}
