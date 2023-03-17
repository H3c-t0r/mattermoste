// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"database/sql"
	"os"
	"testing"

	"github.com/mattermost/mattermost-server/v6/server/boards/model"
	"github.com/mattermost/mattermost-server/v6/server/boards/services/store"
	"github.com/mattermost/mattermost-server/v6/server/channels/store/storetest"
	"github.com/mattermost/mattermost-server/v6/server/platform/shared/mlog"
	"github.com/mgdelacroix/foundation"
	"github.com/stretchr/testify/require"
)

type storeType struct {
	Name       string
	ConnString string
	Store      store.Store
	Logger     *mlog.Logger
}

func NewStoreType(t *testing.T, name string, driver string, skipMigrations bool) *storeType {
	settings := storetest.MakeSqlSettings(driver, false)
	require.NotNil(t, settings.DataSource)
	connectionString := *settings.DataSource

	logger := mlog.CreateConsoleTestLogger(false, mlog.LvlDebug)

	sqlDB, err := sql.Open(driver, connectionString)
	require.NoError(t, err)
	err = sqlDB.Ping()
	require.NoError(t, err)

	storeParams := Params{
		DBType:           driver,
		ConnectionString: connectionString,
		SkipMigrations:   skipMigrations,
		TablePrefix:      "focalboard_",
		Logger:           logger,
		DB:               sqlDB,
		IsPlugin:         false, // ToDo: to be removed
	}
	store, err := New(storeParams)
	require.NoError(t, err)

	return &storeType{name, connectionString, store, logger}
}

func initStores(t *testing.T, skipMigrations bool) []*storeType {
	var storeTypes []*storeType

	if os.Getenv("IS_CI") == "true" {
		switch os.Getenv("MM_SQLSETTINGS_DRIVERNAME") {
		case "mysql":
			storeTypes = append(storeTypes, NewStoreType(t, "MySQL", model.MysqlDBType, skipMigrations))
		case "postgres":
			storeTypes = append(storeTypes, NewStoreType(t, "PostgreSQL", model.PostgresDBType, skipMigrations))
		default:
			t.Errorf(
				"Invalid value %q for MM_SQLSETTINGS_DRIVERNAME when IS_CI=true",
				os.Getenv("MM_SQLSETTINGS_DRIVERNAME"),
			)
		}
	} else {
		storeTypes = append(storeTypes,
			NewStoreType(t, "PostgreSQL", model.PostgresDBType, skipMigrations),
			NewStoreType(t, "MySQL", model.MysqlDBType, skipMigrations),
		)
	}

	return storeTypes
}

func RunStoreTests(t *testing.T, f func(*testing.T, store.Store)) {
	storeTypes := initStores(t, false)

	for _, st := range storeTypes {
		st := st
		t.Run(st.Name, func(t *testing.T) {
			f(t, st.Store)
		})
		require.NoError(t, st.Store.Shutdown())
		require.NoError(t, st.Logger.Shutdown())
	}
}

func RunStoreTestsWithSqlStore(t *testing.T, f func(*testing.T, *SQLStore)) {
	storeTypes := initStores(t, false)

	for _, st := range storeTypes {
		st := st
		sqlstore := st.Store.(*SQLStore)
		t.Run(st.Name, func(t *testing.T) {
			f(t, sqlstore)
		})
		require.NoError(t, st.Store.Shutdown())
		require.NoError(t, st.Logger.Shutdown())
	}
}

func RunStoreTestsWithFoundation(t *testing.T, f func(*testing.T, *foundation.Foundation)) {
	storeTypes := initStores(t, true)

	for _, st := range storeTypes {
		st := st
		t.Run(st.Name, func(t *testing.T) {
			sqlstore := st.Store.(*SQLStore)
			f(t, foundation.New(t, NewBoardsMigrator(sqlstore)))
		})
		require.NoError(t, st.Store.Shutdown())
		require.NoError(t, st.Logger.Shutdown())
	}
}
