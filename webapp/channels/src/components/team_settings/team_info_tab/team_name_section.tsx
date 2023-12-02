// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import type {ChangeEvent} from 'react';
import {useIntl} from 'react-intl';

import type {Team} from '@mattermost/types/teams';

import Input from 'components/widgets/inputs/input/input';
import BaseSettingItem, {type BaseSettingItemProps} from 'components/widgets/modals/components/base_setting_item';

import Constants from 'utils/constants';

type Props = {
    handleNameChanges: (name: string) => void;
    name?: Team['display_name'];
    clientError: BaseSettingItemProps['error'];
};

const TeamNameSection = (props: Props) => {
    const {formatMessage} = useIntl();

    const updateName = (e: ChangeEvent<HTMLInputElement>) => props.handleNameChanges(e.target.value);

    const nameSectionInput = (
        <Input
            id='teamName'
            className='form-control'
            type='text'
            maxLength={Constants.MAX_TEAMNAME_LENGTH}
            onChange={updateName}
            value={props.name}
            label={formatMessage({id: 'general_tab.teamName', defaultMessage: 'Team Name'})}
        />
    );

    return (
        <BaseSettingItem
            title={{id: 'general_tab.teamInfo', defaultMessage: 'Team info'}}
            description={{id: 'general_tab.teamNameInfo', defaultMessage: 'This name will appear on your sign-in screen and at the top of the left sidebar.'}}
            content={nameSectionInput}
            error={props.clientError}
        />
    );
};

export default TeamNameSection;
