// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package audit

import "github.com/wiggin77/logr"

const (
	KeyID        = "id"
	KeyAPIPath   = "api_path"
	KeyEvent     = "event"
	KeyUserID    = "user_id"
	KeySessionID = "session_id"
	KeyClient    = "client"
	KeyIPAddress = "ip_address"
)

var (
	// IDGenerator creates a new unique id for audit records.
	// Reassign to generate custom ids.
	IDGenerator func() string = newID
)

// Meta represents metadata that can be added to a audit record as name/value pairs.
type Meta map[string]interface{}

// Record provides a consistent set of fields used for all audit logging.
type Record struct {
	ID        string
	APIPath   string
	Event     string
	UserID    string
	SessionID string
	Client    string
	IPAddress string
	Meta      Meta
}

// Log emits an audit record with complete info.
func LogRecord(level Level, rec Record) {
	flds := logr.Fields{}
	flds[KeyID] = rec.ID
	flds[KeyAPIPath] = rec.APIPath
	flds[KeyEvent] = rec.Event
	flds[KeyUserID] = rec.UserID
	flds[KeySessionID] = rec.SessionID
	flds[KeyClient] = rec.Client
	flds[KeyIPAddress] = rec.IPAddress

	for k, v := range rec.Meta {
		flds[k] = v
	}

	l := logger.WithFields(flds)
	l.Log(logr.Level(level))
}

// Log emits an audit record based on minimum required info.
func Log(level Level, path string, evt string, userID string, sessionID string, meta Meta) {
	LogRecord(level, Record{
		ID:        IDGenerator(),
		APIPath:   path,
		Event:     evt,
		UserID:    userID,
		SessionID: sessionID,
		Meta:      meta,
	})
}
