// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import type {ConnectedProps} from 'react-redux';
import {bindActionCreators} from 'redux';
import type {ActionCreatorsMapObject, Dispatch} from 'redux';

import type {Team} from '@mattermost/types/teams';

import {getTeam, patchTeam, removeTeamIcon, setTeamIcon} from 'mattermost-redux/actions/teams';
import {getConfig} from 'mattermost-redux/selectors/entities/general';
import type {ActionResult, GenericAction} from 'mattermost-redux/types/actions';

import type {GlobalState} from 'types/store/index';

import TeamInfoTab from './team_info_tab';

export type OwnProps = {
    team?: Team & { last_team_icon_update?: number };
    hasChanges: boolean;
    hasChangeTabError: boolean;
    setHasChanges: (hasChanges: boolean) => void;
    setHasChangeTabError: (hasChangesError: boolean) => void;
    closeModal: () => void;
    collapseModal: () => void;
};

function mapStateToProps(state: GlobalState) {
    const config = getConfig(state);
    const maxFileSize = parseInt(config.MaxFileSize ?? '', 10);

    return {
        maxFileSize,
    };
}

type Actions = {
    getTeam: (teamId: string) => Promise<ActionResult>;
    patchTeam: (team: Partial<Team>) => Promise<ActionResult>;
    removeTeamIcon: (teamId: string) => Promise<ActionResult>;
    setTeamIcon: (teamId: string, teamIconFile: File) => Promise<ActionResult>;
};

function mapDispatchToProps(dispatch: Dispatch<GenericAction>) {
    return {
        actions: bindActionCreators<ActionCreatorsMapObject, Actions>({
            getTeam,
            patchTeam,
            removeTeamIcon,
            setTeamIcon,
        }, dispatch),
    };
}

const connector = connect(mapStateToProps, mapDispatchToProps);

export type PropsFromRedux = ConnectedProps<typeof connector>;

export default connector(TeamInfoTab);
