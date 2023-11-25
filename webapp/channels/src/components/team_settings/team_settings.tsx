// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import type {Team} from '@mattermost/types/teams';

import AccessTab from './team_access_section';
import InfoTab from './team_info_section';

// todo sinan: why it is scrolled down in the tabs
type Props = {
    activeTab: string;
    closeModal: () => void;
    team?: Team;
};

const TeamSettings = ({
    activeTab = '',
    closeModal,
    team,
}: Props): JSX.Element | null => {
    if (!team) {
        return null;
    }

    let result;
    switch (activeTab) {
    case 'info':
        result = <InfoTab team={team}/>;
        break;
    case 'access':
        result = (
            <AccessTab
                team={team}
                closeModal={closeModal}
            />
        );
        break;
    default:
        result = (
            <div/>
        );
        break;
    }

    return result;
};

export default TeamSettings;
