// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import PostStore from 'stores/post_store.jsx';

import {createPost} from 'actions/post_actions.jsx';

import Constants from 'utils/constants.jsx';

import {FormattedMessage} from 'react-intl';

import React from 'react';

export default class PendingPostOptions extends React.Component {
    constructor(props) {
        super(props);
        this.retryPost = this.retryPost.bind(this);
        this.cancelPost = this.cancelPost.bind(this);
        this.state = {};
    }
    retryPost(e) {
        e.preventDefault();

        var post = this.props.post;
        var doLoadPost = true;
        createPost(post, doLoadPost);

        post.state = Constants.POST_LOADING;
        PostStore.updatePendingPost(post);
        this.forceUpdate();
    }
    cancelPost(e) {
        e.preventDefault();

        var post = this.props.post;
        PostStore.removePendingPost(post.channel_id, post.pending_post_id);
        this.forceUpdate();
    }
    render() {
        return (<span className='pending-post-actions'>
            <a
                className='post-retry'
                href='#'
                onClick={this.retryPost}
            >
                <FormattedMessage
                    id='pending_post_actions.retry'
                    defaultMessage='Retry'
                />
            </a>
            {' - '}
            <a
                className='post-cancel'
                href='#'
                onClick={this.cancelPost}
            >
                <FormattedMessage
                    id='pending_post_actions.cancel'
                    defaultMessage='Cancel'
                />
            </a>
        </span>);
    }
}

PendingPostOptions.propTypes = {
    post: React.PropTypes.object
};
