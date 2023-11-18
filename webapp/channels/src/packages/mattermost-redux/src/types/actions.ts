// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AnyAction} from 'redux';
import type {ThunkAction} from 'redux-thunk';

import type {GlobalState} from '@mattermost/types/store';

export type Thunk<ReturnType, State = GlobalState> = ThunkAction<ReturnType, State, never, AnyAction>;

export type ActionFunc<Data = any, Error = any> = Thunk<Promise<ActionResult<Data, Error>> | ActionResult<Data, Error>>;
export type ActionFuncPromiseDoNotUse<Data = any, Error = any> = Thunk<Promise<ActionResult<Data, Error>>>;

export type DispatchFunc<ReturnType = any, State = GlobalState> = Parameters<Thunk<ReturnType, State>>[0];
export type GetStateFunc<State = GlobalState> = Parameters<Thunk<any, State>>[1];

export type GenericAction = AnyAction;

export type Action = GenericAction | Thunk<any>;

export type ActionResult<Data = any, Error = any> = {
    data?: Data;
    error?: Error;
};

