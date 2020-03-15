// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package audit

import (
	"fmt"
	"sort"

	"github.com/wiggin77/logr"
	"github.com/wiggin77/logr/format"
)

type Level logr.Level

type Audit struct {
	lgr    *logr.Logr
	logger logr.Logger

	// OnQueueFull is called on an attempt to add an audit record to a full queue.
	// On return the calling goroutine will block until the audit record can be added.
	OnQueueFull func(qname string, maxQueueSize int)

	// OnError is called when an error occurs while writing an audit record.
	OnError func(err error)
}

func (a *Audit) Init(maxQueueSize int) {
	a.lgr = &logr.Logr{MaxQueueSize: maxQueueSize}
	a.logger = a.lgr.NewLogger()

	a.lgr.OnQueueFull = a.onQueueFull
	a.lgr.OnTargetQueueFull = a.onTargetQueueFull
	a.lgr.OnLoggerError = a.onLoggerError
}

// MakeFilter creates a filter which only allows the specified audit levels to be output.
func (a *Audit) MakeFilter(level ...Level) *logr.CustomFilter {
	filter := &logr.CustomFilter{}
	for _, l := range level {
		filter.Add(logr.Level(l))
	}
	return filter
}

// MakeJSONFormatter creates a formatter that outputs JSON suitable for audit records.
func (a *Audit) MakeJSONFormatter() *format.JSON {
	f := &format.JSON{
		DisableTimestamp:  true,
		DisableMsg:        true,
		DisableStacktrace: true,
		DisableLevel:      true,
		ContextSorter:     sortAuditFields,
	}
	return f
}

// LogRecord emits an audit record with complete info.
func (a *Audit) LogRecord(level Level, rec Record) {
	flds := logr.Fields{}
	flds[KeyAPIPath] = rec.APIPath
	flds[KeyEvent] = rec.Event
	flds[KeyStatus] = rec.Status
	flds[KeyUserID] = rec.UserID
	flds[KeySessionID] = rec.SessionID
	flds[KeyClient] = rec.Client
	flds[KeyIPAddress] = rec.IPAddress

	for k, v := range rec.Meta {
		flds[k] = v
	}

	l := a.logger.WithFields(flds)
	l.Log(logr.Level(level))
}

// Log emits an audit record based on minimum required info.
func (a *Audit) Log(level Level, path string, evt string, status string, userID string, sessionID string, meta Meta) {
	a.LogRecord(level, Record{
		APIPath:   path,
		Event:     evt,
		Status:    status,
		UserID:    userID,
		SessionID: sessionID,
		Meta:      meta,
	})
}

// AddTarget adds a Logr target to the list of targets each audit record will be output to.
func (a *Audit) AddTarget(target logr.Target) {
	a.lgr.AddTarget(target)
}

// Shutdown cleanly stops the audit engine after making best efforts to flush all targets.
func (a *Audit) Shutdown() {
	err := a.lgr.Shutdown()
	if err != nil {
		a.onLoggerError(err)
	}
}

func (a *Audit) onQueueFull(rec *logr.LogRec, maxQueueSize int) bool {
	if a.OnQueueFull != nil {
		a.OnQueueFull("main", maxQueueSize)
	}
	// block until record can be added.
	return false
}

func (a *Audit) onTargetQueueFull(target logr.Target, rec *logr.LogRec, maxQueueSize int) bool {
	if a.OnQueueFull != nil {
		a.OnQueueFull(fmt.Sprintf("%v", target), maxQueueSize)
	}
	// block until record can be added.
	return false
}

func (a *Audit) onLoggerError(err error) {
	if a.OnError != nil {
		a.OnError(err)
	}
}

// sortAuditFields sorts the context fields of an audit record such that some fields
// are prepended in order, some are appended in order, and the rest are sorted alphabetically.
// This is done to make reading the records easier since common fields will appear in the same order.
func sortAuditFields(fields logr.Fields) []format.ContextField {
	prependKeys := []string{KeyEvent, KeyStatus, KeyUserID, KeySessionID, KeyIPAddress}
	appendKeys := []string{KeyClusterID, KeyClient}

	// sort the remaining.
	keys := make([]string, 0, len(fields))
	for k := range fields {
		if !findIn(k, prependKeys, appendKeys) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	allKeys := append(prependKeys, keys...)
	allKeys = append(allKeys, appendKeys...)

	cfs := make([]format.ContextField, 0, len(allKeys))
	for _, k := range allKeys {
		cfs = append(cfs, format.ContextField{Key: k, Val: fields[k]})
	}
	return cfs
}

func findIn(s string, arrs ...[]string) bool {
	for _, list := range arrs {
		for _, key := range list {
			if s == key {
				return true
			}
		}
	}
	return false
}
