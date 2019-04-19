// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/mattermost/mattermost-server/model"
import store "github.com/mattermost/mattermost-server/store"

// PluginStore is an autogenerated mock type for the PluginStore type
type PluginStore struct {
	mock.Mock
}

// CompareAndUpdate provides a mock function with given fields: keyVal, oldValue
func (_m *PluginStore) CompareAndUpdate(keyVal *model.PluginKeyValue, oldValue []byte) (bool, *model.AppError) {
	ret := _m.Called(keyVal, oldValue)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*model.PluginKeyValue, []byte) bool); ok {
		r0 = rf(keyVal, oldValue)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.PluginKeyValue, []byte) *model.AppError); ok {
		r1 = rf(keyVal, oldValue)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// Delete provides a mock function with given fields: pluginId, key
func (_m *PluginStore) Delete(pluginId string, key string) store.StoreChannel {
	ret := _m.Called(pluginId, key)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string, string) store.StoreChannel); ok {
		r0 = rf(pluginId, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// DeleteAllExpired provides a mock function with given fields:
func (_m *PluginStore) DeleteAllExpired() store.StoreChannel {
	ret := _m.Called()

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func() store.StoreChannel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// DeleteAllForPlugin provides a mock function with given fields: PluginId
func (_m *PluginStore) DeleteAllForPlugin(PluginId string) store.StoreChannel {
	ret := _m.Called(PluginId)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string) store.StoreChannel); ok {
		r0 = rf(PluginId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// Get provides a mock function with given fields: pluginId, key
func (_m *PluginStore) Get(pluginId string, key string) store.StoreChannel {
	ret := _m.Called(pluginId, key)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string, string) store.StoreChannel); ok {
		r0 = rf(pluginId, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// List provides a mock function with given fields: pluginId, page, perPage
func (_m *PluginStore) List(pluginId string, page int, perPage int) store.StoreChannel {
	ret := _m.Called(pluginId, page, perPage)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string, int, int) store.StoreChannel); ok {
		r0 = rf(pluginId, page, perPage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// SaveOrUpdate provides a mock function with given fields: keyVal
func (_m *PluginStore) SaveOrUpdate(keyVal *model.PluginKeyValue) store.StoreChannel {
	ret := _m.Called(keyVal)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(*model.PluginKeyValue) store.StoreChannel); ok {
		r0 = rf(keyVal)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}
