// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package pglayer

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/store/storetest"
)

func TestSessionStore(t *testing.T) {
	StoreTest(t, storetest.TestSessionStore)
}
