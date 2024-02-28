// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getPost} from 'mattermost-redux/selectors/entities/posts';

import {openModal} from 'actions/views/modals';

import MarkdownImage from './markdown_image';

function mapStateToProps(state, ownProps) {
    const post = getPost(state, ownProps.postId);
    const isUnsafeLinksPost = post?.props?.unsafe_links === 'true';

    return {
        isUnsafeLinksPost,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            openModal,
        }, dispatch),
    };
}

const connector = connect(mapStateToProps, mapDispatchToProps);

export default connector(MarkdownImage);
