// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import LocalizedIcon from 'components/localized_icon';

import {t} from 'utils/i18n';

type Props = {
    text: React.ReactNode;
    style?: React.CSSProperties;
}

const LoadingSpinner: React.FC<Props> = React.memo(({text, style}) => {
    return (
        <span
            id='loadingSpinner'
            className={'LoadingSpinner' + (text ? ' with-text' : '')}
            style={style}
            data-testid='loadingSpinner'
        >
            <LocalizedIcon
                className='fa fa-spinner fa-fw fa-pulse spinner'
                component='span'
                title={{id: t('generic_icons.loading'), defaultMessage: 'Loading Icon'}}
            />
            {text}
        </span>
    );
});

LoadingSpinner.defaultProps = {
    text: null,
};

export default LoadingSpinner;
