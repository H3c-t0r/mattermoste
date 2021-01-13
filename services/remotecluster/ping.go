// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package remotecluster

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
)

// pingLoop periodically sends a ping to all remote clusters.
func (rcs *Service) pingLoop(done <-chan struct{}) {
	pingChan := make(chan *model.RemoteCluster, MaxConcurrentSends*2)

	// create a thread pool to send pings concurrently to remotes.
	for i := 0; i < MaxConcurrentSends; i++ {
		go rcs.pingEmitter(pingChan, done)
	}

	go rcs.pingGenerator(pingChan, done)
}

func (rcs *Service) pingGenerator(pingChan chan *model.RemoteCluster, done <-chan struct{}) {
	defer close(pingChan)

	for {
		start := time.Now()

		// get all remotes, including any previously offline.
		remotes, err := rcs.server.GetStore().RemoteCluster().GetAll(model.RemoteClusterQueryFilter{})
		if err != nil {
			rcs.server.GetLogger().Log(mlog.LvlRemoteClusterServiceError, "Ping remote cluster failed (could not get list of remotes)", mlog.Err(err))
			select {
			case <-time.After(PingFreq):
				continue
			case <-done:
				return
			}
		}

		for _, rc := range remotes {
			if rc.SiteURL != "" { // filter out unconfirmed invites
				pingChan <- rc
			}
		}

		// try to maintain frequency
		elapsed := time.Since(start)
		if elapsed < PingFreq {
			sleep := time.Until(start.Add(PingFreq))
			select {
			case <-time.After(sleep):
			case <-done:
				return
			}
		}
	}
}

// pingEmitter pulls Remotes from the ping queue (pingChan) and pings them.
// Pinging a remote cannot take longer than PingTimeoutMillis.
func (rcs *Service) pingEmitter(pingChan <-chan *model.RemoteCluster, done <-chan struct{}) {
	for {
		select {
		case rc := <-pingChan:
			if rc == nil {
				return
			}
			if err := rcs.pingRemote(rc); err != nil {
				rcs.server.GetLogger().Log(mlog.LvlRemoteClusterServiceError, "Remote cluster ping failed", mlog.Err(err))
			}
		case <-done:
			return
		}
	}
}

func (rcs *Service) pingRemote(rc *model.RemoteCluster) error {
	frame, err := makePingFrame(rc)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s", rc.SiteURL, PingURL)

	resp, err := rcs.sendFrameToRemote(PingTimeout, frame, url)
	if err != nil {
		return err
	}

	ping := model.RemoteClusterPing{}
	err = json.Unmarshal(resp, &ping)
	if err != nil {
		return err
	}

	if metrics := rcs.server.GetMetrics(); metrics != nil {
		sentAt := time.Unix(0, ping.SentAt*int64(time.Millisecond))
		elapsed := time.Since(sentAt).Seconds()
		metrics.ObserveRemoteClusterPingDuration(rc.RemoteId, elapsed)

		// we approximate clock skew between remotes.
		skew := elapsed/2 - float64(ping.RecvAt-ping.SentAt)/1000
		metrics.ObserveRemoteClusterClockSkew(rc.RemoteId, skew)
	}

	rcs.server.GetLogger().Log(mlog.LvlRemoteClusterServiceDebug, "Remote cluster ping",
		mlog.String("remote", rc.DisplayName), mlog.Int64("SentAt", ping.SentAt), mlog.Int64("RecvAt", ping.RecvAt),
		mlog.Int64("Diff", ping.RecvAt-ping.SentAt))

	return nil
}

func makePingFrame(rc *model.RemoteCluster) (*model.RemoteClusterFrame, error) {
	ping := model.RemoteClusterPing{
		SentAt: model.GetMillis(),
	}
	pingRaw, err := json.Marshal(ping)
	if err != nil {
		return nil, err
	}

	msg := model.NewRemoteClusterMsg(PingTopic, pingRaw)

	frame := &model.RemoteClusterFrame{
		RemoteId: rc.RemoteId,
		Token:    rc.RemoteToken,
		Msg:      msg,
	}
	return frame, nil
}
