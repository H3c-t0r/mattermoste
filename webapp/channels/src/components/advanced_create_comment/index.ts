// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {ActionCreatorsMapObject, bindActionCreators, Dispatch} from 'redux';

import {GlobalState} from 'types/store/index.js';

import {ModalData} from 'types/actions.js';

import {ActionFunc, ActionResult} from 'mattermost-redux/types/actions.js';

import {PostDraft} from 'types/store/draft';

import {getCurrentUserId} from 'mattermost-redux/selectors/entities/common';

import {getConfig, getLicense} from 'mattermost-redux/selectors/entities/general';
import {haveIChannelPermission} from 'mattermost-redux/selectors/entities/roles';
import {getBool, isCustomGroupsEnabled} from 'mattermost-redux/selectors/entities/preferences';
import {getAllChannelStats, getChannelMemberCountsByGroup as selectChannelMemberCountsByGroup} from 'mattermost-redux/selectors/entities/channels';
import {makeGetMessageInHistoryItem} from 'mattermost-redux/selectors/entities/posts';
import {moveHistoryIndexBack, moveHistoryIndexForward, resetCreatePostRequest, resetHistoryIndex} from 'mattermost-redux/actions/posts';
import {getChannelTimezones} from 'mattermost-redux/actions/channels';
import {Permissions, Preferences, Posts} from 'mattermost-redux/constants';
import {getAssociatedGroupsForReferenceByMention} from 'mattermost-redux/selectors/entities/groups';
import {getCurrentTeamId} from 'mattermost-redux/selectors/entities/teams';
import {PreferenceType} from '@mattermost/types/preferences';
import {savePreferences} from 'mattermost-redux/actions/preferences';

import {AdvancedTextEditor, Constants, StoragePrefixes, UserStatuses} from 'utils/constants';
import {getCurrentLocale} from 'selectors/i18n';

import {
    clearCommentDraftUploads,
    makeOnEditLatestPost,
    updateCommentDraft,
    onSubmit,
} from 'actions/views/create_comment';
import {emitShortcutReactToLastPostFrom} from 'actions/post_actions';
import {getPostDraft, getIsRhsExpanded, getSelectedPostFocussedAt} from 'selectors/rhs';
import {showPreviewOnCreateComment} from 'selectors/views/textbox';
import {setShowPreviewOnCreateComment} from 'actions/views/textbox';
import {openModal} from 'actions/views/modals';

import AdvancedCreateComment from './advanced_create_comment';
import {getChannelMemberCountsFromMessage} from 'actions/channel_actions';
import {getStatusForUserId} from 'mattermost-redux/selectors/entities/users';

type OwnProps = {
    rootId: string;
    channelId: string;
    latestPostId: string;
};

function makeMapStateToProps() {
    const getMessageInHistoryItem = makeGetMessageInHistoryItem(Posts.MESSAGE_TYPES.COMMENT as 'comment');

    return (state: GlobalState, ownProps: OwnProps) => {
        const err = state.requests.posts.createPost.error || {};

        const draft = getPostDraft(state, StoragePrefixes.COMMENT_DRAFT, ownProps.rootId);
        const isRemoteDraft = state.views.drafts.remotes[`${StoragePrefixes.COMMENT_DRAFT}${ownProps.rootId}`] || false;

        const channelMembersCount = getAllChannelStats(state)[ownProps.channelId] ? getAllChannelStats(state)[ownProps.channelId].member_count : 1;
        const messageInHistory = getMessageInHistoryItem(state);

        const channel = state.entities.channels.channels[ownProps.channelId] || {};
        const config = getConfig(state);
        const license = getLicense(state);
        const currentUserId = getCurrentUserId(state);
        const userIsOutOfOffice = getStatusForUserId(state, currentUserId) === UserStatuses.OUT_OF_OFFICE;
        const enableConfirmNotificationsToChannel = config.EnableConfirmNotificationsToChannel === 'true';
        const isTimezoneEnabled = config.ExperimentalTimezone === 'true';
        const canPost = haveIChannelPermission(state, channel.team_id, channel.id, Permissions.CREATE_POST);
        const useChannelMentions = haveIChannelPermission(state, channel.team_id, channel.id, Permissions.USE_CHANNEL_MENTIONS);
        const isLDAPEnabled = license?.IsLicensed === 'true' && license?.LDAPGroups === 'true';
        const useCustomGroupMentions = isCustomGroupsEnabled(state) && haveIChannelPermission(state, channel.team_id, channel.id, Permissions.USE_GROUP_MENTIONS);
        const useLDAPGroupMentions = isLDAPEnabled && haveIChannelPermission(state, channel.team_id, channel.id, Permissions.USE_GROUP_MENTIONS);
        const channelMemberCountsByGroup = selectChannelMemberCountsByGroup(state, ownProps.channelId);
        const groupsWithAllowReference = useLDAPGroupMentions || useCustomGroupMentions ? getAssociatedGroupsForReferenceByMention(state, channel.team_id, channel.id) : null;
        const isFormattingBarHidden = getBool(state, Constants.Preferences.ADVANCED_TEXT_EDITOR, AdvancedTextEditor.COMMENT);
        const currentTeamId = getCurrentTeamId(state);

        return {
            currentTeamId,
            draft,
            isRemoteDraft,
            messageInHistory,
            channelMembersCount,
            currentUserId,
            isFormattingBarHidden,
            codeBlockOnCtrlEnter: getBool(state, Preferences.CATEGORY_ADVANCED_SETTINGS, 'code_block_ctrl_enter', true),
            ctrlSend: getBool(state, Preferences.CATEGORY_ADVANCED_SETTINGS, 'send_on_ctrl_enter'),
            createPostErrorId: err.server_error_id,
            enableConfirmNotificationsToChannel,
            locale: getCurrentLocale(state),
            rhsExpanded: getIsRhsExpanded(state),
            isTimezoneEnabled,
            selectedPostFocussedAt: getSelectedPostFocussedAt(state),
            canPost,
            useChannelMentions,
            shouldShowPreview: showPreviewOnCreateComment(state),
            groupsWithAllowReference,
            useLDAPGroupMentions,
            channelMemberCountsByGroup,
            useCustomGroupMentions,
            userIsOutOfOffice,
            channel,
        };
    };
}

function makeOnUpdateCommentDraft(channelId: string) {
    return (draft: PostDraft, save = false, instant = false) => updateCommentDraft({...draft, channelId}, save, instant);
}

type Actions = {
    clearCommentDraftUploads: () => void;
    onUpdateCommentDraft: (draft: PostDraft, save?: boolean) => void;
    onResetHistoryIndex: () => void;
    moveHistoryIndexBack: (index: string) => Promise<void>;
    moveHistoryIndexForward: (index: string) => Promise<void>;
    onEditLatestPost: () => ActionResult;
    resetCreatePostRequest: () => void;
    getChannelTimezones: (channelId: string) => Promise<ActionResult>;
    emitShortcutReactToLastPostFrom: (location: string) => void;
    setShowPreview: (showPreview: boolean) => void;
    getChannelMemberCountsFromMessage: (channelID: string, message: string) => void;
    openModal: <P>(modalData: ModalData<P>) => void;
    savePreferences: (userId: string, preferences: PreferenceType[]) => ActionResult;
    onSubmit: (draft: PostDraft, options: {ignoreSlash?: boolean}, latestPostId: string | undefined) => ActionResult;
};

function makeMapDispatchToProps() {
    let onUpdateCommentDraft: (draft: PostDraft, save?: boolean) => void;
    let onEditLatestPost: () => ActionFunc;

    function onResetHistoryIndex() {
        return resetHistoryIndex(Posts.MESSAGE_TYPES.COMMENT);
    }

    let rootId: string;
    let channelId: string;

    return (dispatch: Dispatch, ownProps: OwnProps) => {
        if (channelId !== ownProps.channelId) {
            onUpdateCommentDraft = makeOnUpdateCommentDraft(ownProps.channelId);
        }

        if (rootId !== ownProps.rootId) {
            onEditLatestPost = makeOnEditLatestPost(ownProps.rootId);
        }

        rootId = ownProps.rootId;
        channelId = ownProps.channelId;

        return bindActionCreators<ActionCreatorsMapObject<any>, Actions>(
            {
                clearCommentDraftUploads,
                onUpdateCommentDraft,
                onSubmit,
                onResetHistoryIndex,
                moveHistoryIndexBack,
                moveHistoryIndexForward,
                onEditLatestPost,
                resetCreatePostRequest,
                getChannelTimezones,
                emitShortcutReactToLastPostFrom,
                setShowPreview: setShowPreviewOnCreateComment,
                getChannelMemberCountsFromMessage,
                openModal,
                savePreferences,
            },
            dispatch,
        );
    };
}

export default connect(makeMapStateToProps, makeMapDispatchToProps, null, {forwardRef: true})(AdvancedCreateComment);
