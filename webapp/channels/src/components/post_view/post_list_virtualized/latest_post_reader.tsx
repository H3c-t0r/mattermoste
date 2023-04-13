// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useMemo} from 'react';
import {useSelector} from 'react-redux';

import {getPost} from 'mattermost-redux/selectors/entities/posts';

import {Post} from '@mattermost/types/posts';

import {GlobalState} from 'types/store';

import {getLatestPostId, usePostAriaLabel} from 'utils/post_utils';

interface Props {
    postIds?: string[];
}

const LatestPostReader = (props: Props): JSX.Element => {
    const {postIds} = props;
    const latestPostId = useMemo(() => getLatestPostId(postIds || []), [postIds]);
    const latestPost = useSelector<GlobalState, Post>((state) => getPost(state, latestPostId));

    const ariaLabel = usePostAriaLabel(latestPost);

    return (
        <span
            className='sr-only'
            aria-live='polite'
        >
            {ariaLabel}
        </span>
    );
};

export default LatestPostReader;
