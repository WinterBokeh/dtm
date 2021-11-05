package test

import (
	"time"

	"github.com/yedf/dtm/common"
	"github.com/yedf/dtm/dtmcli/dtmimp"
	"github.com/yedf/dtm/dtmsvr"
)

var config = common.DtmConfig

func dbGet() *common.DB {
	return common.DbGet(config.DB)
}

// waitTransProcessed only for test usage. wait for transaction processed once
func waitTransProcessed(gid string) {
	dtmimp.Logf("waiting for gid %s", gid)
	select {
	case id := <-dtmsvr.TransProcessedTestChan:
		for id != gid {
			dtmimp.LogRedf("-------id %s not match gid %s", id, gid)
			id = <-dtmsvr.TransProcessedTestChan
		}
		dtmimp.Logf("finish for gid %s", gid)
	case <-time.After(time.Duration(time.Second * 3)):
		dtmimp.LogFatalf("Wait Trans timeout")
	}
}

func cronTransOnce() {
	gid := dtmsvr.CronTransOnce()
	if dtmsvr.TransProcessedTestChan != nil && gid != "" {
		waitTransProcessed(gid)
	}
}

var e2p = dtmimp.E2P

// TransGlobal alias
type TransGlobal = dtmsvr.TransGlobal

// TransBranch alias
type TransBranch = dtmsvr.TransBranch

func cronTransOnceForwardNow(seconds int) {
	old := dtmsvr.NowForwardDuration
	dtmsvr.NowForwardDuration = time.Duration(seconds) * time.Second
	cronTransOnce()
	dtmsvr.NowForwardDuration = old
}

func cronTransOnceForwardCron(seconds int) {
	old := dtmsvr.CronForwardDuration
	dtmsvr.CronForwardDuration = time.Duration(seconds) * time.Second
	cronTransOnce()
	dtmsvr.CronForwardDuration = old
}
