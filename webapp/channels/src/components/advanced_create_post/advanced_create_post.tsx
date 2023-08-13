// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/* eslint-disable max-lines */

import React from 'react';

import {Posts} from 'mattermost-redux/constants';
import {sortFileInfos} from 'mattermost-redux/utils/file_utils';
import {ActionResult} from 'mattermost-redux/types/actions';

import {Channel, ChannelMemberCountsByGroup} from '@mattermost/types/channels';
import {Post, PostPriority, PostPriorityMetadata} from '@mattermost/types/posts';
import {PreferenceType} from '@mattermost/types/preferences';
import {ServerError} from '@mattermost/types/errors';
import {CommandArgs} from '@mattermost/types/integrations';
import {Group, GroupSource} from '@mattermost/types/groups';
import {FileInfo} from '@mattermost/types/files';
import {Emoji} from '@mattermost/types/emojis';

import * as GlobalActions from 'actions/global_actions';
import Constants, {
    StoragePrefixes,
    ModalIdentifiers,
    Locations,
    A11yClassNames,
    Preferences,
    AdvancedTextEditor as AdvancedTextEditorConst,
} from 'utils/constants';
import * as Keyboard from 'utils/keyboard';
import {
    specialMentionsInText,
    postMessageOnKeyPress,
    shouldFocusMainTextbox,
    isErrorInvalidSlashCommand,
    splitMessageBasedOnCaretPosition,
    groupsMentionedInText,
    mentionsMinusSpecialMentionsInText,
    hasRequestedPersistentNotifications,
    isStatusSlashCommand,
    extractCommand,
} from 'utils/post_utils';
import * as UserAgent from 'utils/user_agent';
import * as Utils from 'utils/utils';
import EmojiMap from 'utils/emoji_map';

import NotifyConfirmModal from 'components/notify_confirm_modal';
import EditChannelHeaderModal from 'components/edit_channel_header_modal';
import EditChannelPurposeModal from 'components/edit_channel_purpose_modal';
import {FileUpload as FileUploadClass} from 'components/file_upload/file_upload';
import ResetStatusModal from 'components/reset_status_modal';
import TextboxClass from 'components/textbox/textbox';
import PostPriorityPickerOverlay from 'components/post_priority/post_priority_picker_overlay';
import PersistNotificationConfirmModal from 'components/persist_notification_confirm_modal';

import {PostDraft} from 'types/store/draft';
import {ModalData} from 'types/actions';

import PriorityLabels from './priority_labels';
import Foo from 'components/advanced_text_editor/foo';

const KeyCodes = Constants.KeyCodes;

type TextboxElement = HTMLInputElement | HTMLTextAreaElement;

type Props = {

    // ref passed from channelView for EmojiPickerOverlay
    getChannelView?: () => void;

    // Data used in notifying user for @all and @channel
    currentChannelMembersCount: number;

    // Data used in multiple places of the component
    currentChannel: Channel;

    //Data used for DM prewritten messages
    currentChannelTeammateUsername?: string;

    //Data used in executing commands for channel actions passed down to client4 function
    currentTeamId: string;

    //Data used for posting message
    currentUserId: string;

    //Force message submission on CTRL/CMD + ENTER
    codeBlockOnCtrlEnter?: boolean;

    //Flag used for handling submit
    ctrlSend?: boolean;

    //Flag used for adding a class center to Postbox based on user pref
    fullWidthTextBox?: boolean;

    // Data used for deciding if tutorial tip is to be shown
    showSendTutorialTip: boolean;

    // Data used populating message state when triggered by shortcuts
    messageInHistoryItem?: string;

    // Data used for populating message state from previous draft
    draft: PostDraft;

    // Data used for knowing if the draft came from a WS event
    isRemoteDraft: boolean;

    // Data used dispatching handleViewAction ex: edit post
    latestReplyablePostId?: string;
    locale: string;

    // Data used for calling edit of post
    currentUsersLatestPost?: Post | null;

    //Whether to check with the user before notifying the whole channel.
    enableConfirmNotificationsToChannel: boolean;

    emojiMap: EmojiMap;

    //Whether to display a confirmation modal to reset status.
    userIsOutOfOffice: boolean;
    rhsExpanded: boolean;

    //If RHS open
    rhsOpen: boolean;

    //To check if the timezones are enable on the server.
    isTimezoneEnabled: boolean;

    canPost: boolean;

    //To determine if the current user can send special channel mentions
    useChannelMentions: boolean;

    //Should preview be showed
    shouldShowPreview: boolean;

    isFormattingBarHidden: boolean;

    isPostPriorityEnabled: boolean;

    actions: {

        //Set show preview for textbox
        setShowPreview: (showPreview: boolean) => void;

        // func called after message submit.
        addMessageIntoHistory: (message: string) => void;

        // func called for navigation through messages by Up arrow
        moveHistoryIndexBack: (index: string) => Promise<void>;

        // func called for navigation through messages by Down arrow
        moveHistoryIndexForward: (index: string) => Promise<void>;

        // func called for adding a reaction
        addReaction: (postId: string, emojiName: string) => void;

        // func called for posting message
        onSubmitPost: (post: Post, fileInfos: FileInfo[]) => void;

        // func called for removing a reaction
        removeReaction: (postId: string, emojiName: string) => void;

        // func called on load of component to clear drafts
        clearDraftUploads: () => void;

        //hooks called before a message is sent to the server
        runMessageWillBePostedHooks: (originalPost: Post) => ActionResult;

        //hooks called before a slash command is sent to the server
        runSlashCommandWillBePostedHooks: (originalMessage: string, originalArgs: CommandArgs) => ActionResult;

        // func called for setting drafts
        setDraft: (name: string, value: PostDraft | null, draftChannelId: string, save?: boolean, instant?: boolean) => void;

        // func called for editing posts
        setEditingPost: (postId?: string, refocusId?: string, title?: string, isRHS?: boolean) => void;

        // func called for opening the last replayable post in the RHS
        selectPostFromRightHandSideSearchByPostId: (postId: string) => void;

        //Function to open a modal
        openModal: <P>(modalData: ModalData<P>) => void;

        executeCommand: (message: string, args: CommandArgs) => ActionResult;

        //Function to get the users timezones in the channel
        getChannelTimezones: (channelId: string) => ActionResult;
        scrollPostListToBottom: () => void;

        //Function to set or unset emoji picker for last message
        emitShortcutReactToLastPostFrom: (emittedFrom: string) => void;

        getChannelMemberCountsFromMessage: (channelId: string, message: string) => void;

        //Function used to advance the tutorial forward
        savePreferences: (userId: string, preferences: PreferenceType[]) => ActionResult;
        onSubmit: (draft: PostDraft, options: {ignoreSlash?: boolean}, latestPostId?: string) => ActionResult;
    };

    groupsWithAllowReference: Map<string, Group> | null;
    channelMemberCountsByGroup: ChannelMemberCountsByGroup;
    useLDAPGroupMentions: boolean;
    useCustomGroupMentions: boolean;
}

type State = {
    message: string;
    caretPosition: number;
    showEmojiPicker: boolean;
    renderScrollbar: boolean;
    scrollbarWidth: number;
    currentChannel: Channel;
    errorClass: string | null;
    serverError: (ServerError & {submittedMessage?: string}) | null;
    postError?: React.ReactNode;
    showFormat: boolean;
    isFormattingBarHidden: boolean;
    showPostPriorityPicker: boolean;
};

class AdvancedCreatePost extends React.PureComponent<Props, State> {
    static defaultProps = {
        latestReplyablePostId: '',
    };

    private lastBlurAt = 0;
    private lastChannelSwitchAt = 0;
    private draftsForChannel: {[channelID: string]: PostDraft | null} = {};
    private lastOrientation?: string;
    private isDraftSubmitting = false;

    private textboxRef: React.RefObject<TextboxClass>;
    private fileUploadRef: React.RefObject<FileUploadClass>;

    static getDerivedStateFromProps(props: Props, state: State): Partial<State> {
        let updatedState: Partial<State> = {
            currentChannel: props.currentChannel,
        };
        if (
            props.currentChannel.id !== state.currentChannel.id ||
            (props.isRemoteDraft && props.draft.message !== state.message)
        ) {
            updatedState = {
                ...updatedState,
                message: props.draft.message,
                serverError: null,
            };
        }
        return updatedState;
    }

    constructor(props: Props) {
        super(props);
        this.state = {
            message: props.draft.message,
            caretPosition: props.draft.message.length,
            showEmojiPicker: false,
            renderScrollbar: false,
            scrollbarWidth: 0,
            currentChannel: props.currentChannel,
            errorClass: null,
            serverError: null,
            showFormat: false,
            isFormattingBarHidden: props.isFormattingBarHidden,
            showPostPriorityPicker: false,
        };

        this.textboxRef = React.createRef<TextboxClass>();
        this.fileUploadRef = React.createRef<FileUploadClass>();
    }

    componentDidMount() {
        const {actions} = this.props;
        this.onOrientationChange();
        actions.setShowPreview(false);
        actions.clearDraftUploads();
        this.focusTextbox();
        document.addEventListener('keydown', this.documentKeyHandler);
        this.setOrientationListeners();
        this.getChannelMemberCountsByGroup();
    }

    componentDidUpdate(prevProps: Props, prevState: State) {
        const {currentChannel, actions} = this.props;
        if (prevProps.currentChannel.id !== currentChannel.id) {
            this.lastChannelSwitchAt = Date.now();
            this.focusTextbox();
            this.saveDraft(prevProps);
            this.getChannelMemberCountsByGroup();
        }

        if (currentChannel.id !== prevProps.currentChannel.id) {
            actions.setShowPreview(false);
        }

        // Focus on textbox when emoji picker is closed
        if (prevState.showEmojiPicker && !this.state.showEmojiPicker) {
            this.focusTextbox();
        }

        // Focus on textbox when returned from preview mode
        if (prevProps.shouldShowPreview && !this.props.shouldShowPreview) {
            this.focusTextbox();
        }
    }

    componentWillUnmount() {
        document.removeEventListener('keydown', this.documentKeyHandler);
        this.removeOrientationListeners();
        this.saveDraft();
    }

    getChannelMemberCountsByGroup = () => {
        this.props.actions.getChannelMemberCountsFromMessage(this.props.currentChannel.id, this.state.message);
    };

    saveDraft = (props = this.props) => {
        if (props.currentChannel) {
            const channelId = props.currentChannel.id;
            props.actions.setDraft(StoragePrefixes.DRAFT + channelId, this.draftsForChannel[channelId], channelId, true);
        }
    };

    setShowPreview = (newPreviewValue: boolean) => {
        this.props.actions.setShowPreview(newPreviewValue);
    };

    setOrientationListeners = () => {
        if (window.screen.orientation && 'onchange' in window.screen.orientation) {
            window.screen.orientation.addEventListener('change', this.onOrientationChange);
        } else if ('onorientationchange' in window) {
            window.addEventListener('orientationchange', this.onOrientationChange);
        }
    };

    removeOrientationListeners = () => {
        if (window.screen.orientation && 'onchange' in window.screen.orientation) {
            window.screen.orientation.removeEventListener('change', this.onOrientationChange);
        } else if ('onorientationchange' in window) {
            window.removeEventListener('orientationchange', this.onOrientationChange);
        }
    };

    onOrientationChange = () => {
        if (!UserAgent.isIosWeb()) {
            return;
        }

        const LANDSCAPE_ANGLE = 90;
        let orientation = 'portrait';
        if (window.orientation) {
            orientation = Math.abs(window.orientation as number) === LANDSCAPE_ANGLE ? 'landscape' : 'portrait';
        }

        if (window.screen.orientation) {
            orientation = window.screen.orientation.type.split('-')[0];
        }

        if (
            this.lastOrientation &&
            orientation !== this.lastOrientation &&
            (document.activeElement || {}).id === 'post_textbox'
        ) {
            this.textboxRef.current?.blur();
        }

        this.lastOrientation = orientation;
    };

    handlePostError = (postError: React.ReactNode) => {
        if (this.state.postError !== postError) {
            this.setState({postError});
        }
    };

    toggleEmojiPicker = (e?: React.MouseEvent<HTMLButtonElement, MouseEvent>): void => {
        e?.stopPropagation();
        this.setState({showEmojiPicker: !this.state.showEmojiPicker});
    };

    hideEmojiPicker = () => {
        this.handleEmojiClose();
    };

    shouldEnableAddButton = () => {
        const message = this.state.message;
        const fileInfos = this.props.draft.fileInfos;
        const serverError = this.state.serverError;
        if (message.trim().length !== 0 || fileInfos.length !== 0) {
            return true;
        }
        return isErrorInvalidSlashCommand(serverError);
    };

    doSubmit = async (e?: React.FormEvent) => {
        const channelId = this.props.currentChannel.id;
        const draft = {...this.props.draft, message: this.state.message};
        const serverError = this.state.serverError;
        const message = draft.message;
        const latestPost = this.props.latestReplyablePostId;
        const scrollPostListToBottom = this.props.actions.scrollPostListToBottom;
        if (e) {
            e.preventDefault();
        }

        const enableAddButton = this.shouldEnableAddButton();

        if (!enableAddButton) {
            this.isDraftSubmitting = false;
            return;
        }

        if (draft.uploadsInProgress.length > 0) {
            this.isDraftSubmitting = false;
            return;
        }

        if (this.state.postError) {
            this.setState({errorClass: 'animation--highlight'});
            setTimeout(() => {
                this.setState({errorClass: null});
            }, Constants.ANIMATION_TIMEOUT);
            this.isDraftSubmitting = false;
            return;
        }

        const fasterThanHumanWillClick = 150;
        const forceFocus = Date.now() - this.lastBlurAt < fasterThanHumanWillClick;
        this.focusTextbox(forceFocus);

        let ignoreSlash = false;
        if (serverError && isErrorInvalidSlashCommand(serverError) && serverError.submittedMessage === message) {
            ignoreSlash = true;
        }

        const options = {ignoreSlash};

        const res = await this.props.actions.onSubmit(draft, options, latestPost);
        if (res.error) {
            const err = res.error;
            err.submittedMessage = draft.message;
            this.setState({
                serverError: err,
                message,
            });
            this.isDraftSubmitting = false;
            return;
        }

        this.setState({message: ''});
        this.setState({
            serverError: null,
            postError: null,
            showFormat: false,
        });

        scrollPostListToBottom();
        this.isDraftSubmitting = false;
        this.draftsForChannel[channelId] = null;
    };

    handleNotifyAllConfirmation = () => {
        this.doSubmit();
    };

    showNotifyAllModal = (mentions: string[], channelTimezoneCount: number, memberNotifyCount: number) => {
        this.props.actions.openModal({
            modalId: ModalIdentifiers.NOTIFY_CONFIRM_MODAL,
            dialogType: NotifyConfirmModal,
            dialogProps: {
                mentions,
                channelTimezoneCount,
                memberNotifyCount,
                onConfirm: () => this.handleNotifyAllConfirmation(),
                onExited: () => {
                    this.isDraftSubmitting = false;
                },
            },
        });
    };

    showPersistNotificationModal = (message: string, specialMentions: {[key: string]: boolean}, channelType: Channel['type']) => {
        this.props.actions.openModal({
            modalId: ModalIdentifiers.PERSIST_NOTIFICATION_CONFIRM_MODAL,
            dialogType: PersistNotificationConfirmModal,
            dialogProps: {
                currentChannelTeammateUsername: this.props.currentChannelTeammateUsername,
                specialMentions,
                channelType,
                message,
                onConfirm: this.handleNotifyAllConfirmation,
            },
        });
    };

    handleSubmit = async (e: React.FormEvent) => {
        const message = this.state.message;
        const channelMembersCount = this.props.currentChannelMembersCount;
        const channelId = this.props.currentChannel.id;
        const getChannelTimezones = this.props.actions.getChannelTimezones;
        const openModal = this.props.actions.openModal;

        e.preventDefault();
        this.setShowPreview(false);
        this.isDraftSubmitting = true;

        const {
            enableConfirmNotificationsToChannel,
            useChannelMentions,
            isTimezoneEnabled,
            groupsWithAllowReference,
            channelMemberCountsByGroup,
            useLDAPGroupMentions,
            useCustomGroupMentions,
            userIsOutOfOffice,
            currentChannel: updateChannel,
        } = this.props;

        const notificationsToChannel = enableConfirmNotificationsToChannel && useChannelMentions;
        let memberNotifyCount = 0;
        let channelTimezoneCount = 0;
        let mentions: string[] = [];

        const specialMentions = specialMentionsInText(message);
        const hasSpecialMentions = Object.values(specialMentions).includes(true);

        if (enableConfirmNotificationsToChannel && !hasSpecialMentions && (useLDAPGroupMentions || useCustomGroupMentions)) {
            // Groups mentioned in users text
            const mentionGroups = groupsMentionedInText(message, groupsWithAllowReference);
            if (mentionGroups.length > 0) {
                mentionGroups.
                    forEach((group) => {
                        if (group.source === GroupSource.Ldap && !useLDAPGroupMentions) {
                            return;
                        }
                        if (group.source === GroupSource.Custom && !useCustomGroupMentions) {
                            return;
                        }
                        const mappedValue = channelMemberCountsByGroup[group.id];
                        if (mappedValue && mappedValue.channel_member_count > Constants.NOTIFY_ALL_MEMBERS && mappedValue.channel_member_count > memberNotifyCount) {
                            memberNotifyCount = mappedValue.channel_member_count;
                            channelTimezoneCount = mappedValue.channel_member_timezones_count;
                        }
                        mentions.push(`@${group.name}`);
                    });
                mentions = [...new Set(mentions)];
            }
        }

        if (notificationsToChannel && channelMembersCount > Constants.NOTIFY_ALL_MEMBERS && hasSpecialMentions) {
            memberNotifyCount = channelMembersCount - 1;

            for (const k in specialMentions) {
                if (specialMentions[k]) {
                    mentions.push('@' + k);
                }
            }

            if (isTimezoneEnabled) {
                const {data} = await getChannelTimezones(channelId);
                channelTimezoneCount = data ? data.length : 0;
            }
        }

        if (
            this.props.isPostPriorityEnabled &&
            hasRequestedPersistentNotifications(this.props.draft?.metadata?.priority)
        ) {
            this.showPersistNotificationModal(this.state.message, specialMentions, updateChannel.type);
            this.isDraftSubmitting = false;
            return;
        }

        if (memberNotifyCount > 0) {
            this.showNotifyAllModal(mentions, channelTimezoneCount, memberNotifyCount);
            return;
        }

        const status = extractCommand(message);
        if (userIsOutOfOffice && isStatusSlashCommand(status)) {
            const resetStatusModalData = {
                modalId: ModalIdentifiers.RESET_STATUS,
                dialogType: ResetStatusModal,
                dialogProps: {newStatus: status},
            };

            openModal(resetStatusModalData);

            this.resetMessage();
            this.isDraftSubmitting = false;
            return;
        }

        if (message.trimEnd() === '/header') {
            const editChannelHeaderModalData = {
                modalId: ModalIdentifiers.EDIT_CHANNEL_HEADER,
                dialogType: EditChannelHeaderModal,
                dialogProps: {channel: updateChannel},
            };

            openModal(editChannelHeaderModalData);

            this.resetMessage();
            this.isDraftSubmitting = false;
            return;
        }

        const isDirectOrGroup =
            updateChannel.type === Constants.DM_CHANNEL || updateChannel.type === Constants.GM_CHANNEL;
        if (!isDirectOrGroup && message.trimEnd() === '/purpose') {
            const editChannelPurposeModalData = {
                modalId: ModalIdentifiers.EDIT_CHANNEL_PURPOSE,
                dialogType: EditChannelPurposeModal,
                dialogProps: {channel: updateChannel},
            };

            openModal(editChannelPurposeModalData);

            this.resetMessage();
            this.isDraftSubmitting = false;
            return;
        }

        await this.doSubmit(e);
    };

    resetMessage() {
        if (this.state.message !== '') {
            this.setState({message: ''});
        }
    }

    focusTextbox = (keepFocus = false) => {
        const postTextboxDisabled = !this.props.canPost;
        if (this.textboxRef.current && postTextboxDisabled) {
            this.textboxRef.current.blur(); // Fixes Firefox bug which causes keyboard shortcuts to be ignored (MM-22482)
            return;
        }
        if (this.textboxRef.current && (keepFocus || !UserAgent.isMobile())) {
            this.textboxRef.current.focus();
        }
    };

    postMsgKeyPress = (e: React.KeyboardEvent<TextboxElement>) => {
        const {ctrlSend, codeBlockOnCtrlEnter} = this.props;

        const {allowSending, withClosedCodeBlock, ignoreKeyPress, message} = postMessageOnKeyPress(
            e,
            this.state.message,
            Boolean(ctrlSend),
            Boolean(codeBlockOnCtrlEnter),
            Date.now(),
            this.lastChannelSwitchAt,
            this.state.caretPosition,
        ) as {
            allowSending: boolean;
            withClosedCodeBlock?: boolean;
            ignoreKeyPress?: boolean;
            message?: string;
        };

        if (ignoreKeyPress) {
            e.preventDefault();
            e.stopPropagation();
            return;
        }

        if (allowSending && this.isValidPersistentNotifications()) {
            if (e.persist) {
                e.persist();
            }
            if (this.textboxRef.current) {
                this.isDraftSubmitting = true;
                this.textboxRef.current.blur();
            }

            if (withClosedCodeBlock && message) {
                this.setState({message}, () => this.handleSubmit(e));
            } else {
                this.handleSubmit(e);
            }

            this.setShowPreview(false);
        }

        this.emitTypingEvent();
    };

    emitTypingEvent = () => {
        const channelId = this.props.currentChannel.id;
        GlobalActions.emitLocalUserTypingEvent(channelId, '');
    };

    handleChange = (e: React.ChangeEvent<TextboxElement>) => {
        const message = e.target.value;

        let serverError = this.state.serverError;
        if (isErrorInvalidSlashCommand(serverError)) {
            serverError = null;
        }

        this.setState({
            message,
            serverError,
        });

        const draft = {
            ...this.props.draft,
            message,
        };

        this.handleDraftChange(draft);
    };

    handleDraftChange = (draft: PostDraft, instant = false) => {
        const channelId = this.props.currentChannel.id;
        this.props.actions.setDraft(StoragePrefixes.DRAFT + channelId, draft, channelId, false, instant);
        this.draftsForChannel[channelId] = draft;
    };

    handleFileUploadChange = () => {
        this.focusTextbox();
    };

    handleUploadStart = (clientIds: string[], channelId: string) => {
        const uploadsInProgress = [...this.props.draft.uploadsInProgress, ...clientIds];

        const draft = {
            ...this.props.draft,
            uploadsInProgress,
        };

        this.props.actions.setDraft(StoragePrefixes.DRAFT + channelId, draft, channelId);
        this.draftsForChannel[channelId] = draft;

        // this is a bit redundant with the code that sets focus when the file input is clicked,
        // but this also resets the focus after a drag and drop
        this.focusTextbox();
    };

    handleFileUploadComplete = (fileInfos: FileInfo[], clientIds: string[], channelId: string) => {
        const draft = {...this.draftsForChannel[channelId]!};

        // remove each finished file from uploads
        for (let i = 0; i < clientIds.length; i++) {
            if (draft.uploadsInProgress) {
                const index = draft.uploadsInProgress.indexOf(clientIds[i]);

                if (index !== -1) {
                    draft.uploadsInProgress = draft.uploadsInProgress.filter((item, itemIndex) => index !== itemIndex);
                }
            }
        }

        if (draft.fileInfos) {
            draft.fileInfos = sortFileInfos(draft.fileInfos.concat(fileInfos), this.props.locale);
        }

        this.handleDraftChange(draft, true);
    };

    handleUploadError = (uploadError: string | ServerError | null, clientId?: string, channelId?: string) => {
        if (clientId && channelId) {
            const draft = {...this.draftsForChannel[channelId]!};

            if (draft.uploadsInProgress) {
                const index = draft.uploadsInProgress.indexOf(clientId);

                if (index !== -1) {
                    const uploadsInProgress = draft.uploadsInProgress.filter((item, itemIndex) => index !== itemIndex);
                    const modifiedDraft = {
                        ...draft,
                        uploadsInProgress,
                    };
                    this.props.actions.setDraft(StoragePrefixes.DRAFT + channelId, modifiedDraft, channelId);
                    this.draftsForChannel[channelId] = modifiedDraft;
                }
            }
        }

        if (typeof uploadError === 'string') {
            if (uploadError.length !== 0) {
                this.setState({serverError: new Error(uploadError)});
            }
        } else {
            this.setState({serverError: uploadError});
        }
    };

    removePreview = (id: string) => {
        const draft = {...this.props.draft};
        const fileInfos = [...draft.fileInfos];
        const uploadsInProgress = [...draft.uploadsInProgress];
        const channelId = this.props.currentChannel.id;

        // Clear previous errors
        this.setState({serverError: null});

        // id can either be the id of an uploaded file or the client id of an in progress upload
        let index = draft.fileInfos.findIndex((info) => info.id === id);
        if (index === -1) {
            index = draft.uploadsInProgress.indexOf(id);

            if (index !== -1) {
                uploadsInProgress.splice(index, 1);

                if (this.fileUploadRef.current) {
                    this.fileUploadRef.current.cancelUpload(id);
                }
            }
        } else {
            fileInfos.splice(index, 1);
        }

        const modifiedDraft = {
            ...draft,
            fileInfos,
            uploadsInProgress,
        };

        this.props.actions.setDraft(StoragePrefixes.DRAFT + channelId, modifiedDraft, channelId, false);
        this.draftsForChannel[channelId] = modifiedDraft;

        this.handleFileUploadChange();
    };

    focusTextboxIfNecessary = (e: KeyboardEvent) => {
        // Focus should go to the RHS when it is expanded
        if (this.props.rhsExpanded) {
            return;
        }

        // Hacky fix to avoid cursor jumping textbox sometimes
        if (this.props.rhsOpen && document.activeElement?.tagName === 'BODY') {
            return;
        }

        // Bit of a hack to not steal focus from the channel switch modal if it's open
        // This is a special case as the channel switch modal does not enforce focus like
        // most modals do
        if (document.getElementsByClassName('channel-switch-modal').length) {
            return;
        }

        if (shouldFocusMainTextbox(e, document.activeElement)) {
            this.focusTextbox();
        }
    };

    documentKeyHandler = (e: KeyboardEvent) => {
        const ctrlOrMetaKeyPressed = e.ctrlKey || e.metaKey;
        const lastMessageReactionKeyCombo = ctrlOrMetaKeyPressed && e.shiftKey && Keyboard.isKeyPressed(e, KeyCodes.BACK_SLASH);
        if (lastMessageReactionKeyCombo) {
            this.reactToLastMessage(e);
            return;
        }

        this.focusTextboxIfNecessary(e);
    };

    fillMessageFromHistory() {
        const lastMessage = this.props.messageInHistoryItem;
        this.setState({
            message: lastMessage || '',
        });
    }

    handleMouseUpKeyUp = (e: React.MouseEvent | React.KeyboardEvent) => {
        this.setState({
            caretPosition: (e.target as HTMLInputElement).selectionStart || 0,
        });
    };

    editLastPost = (e: React.KeyboardEvent) => {
        e.preventDefault();

        const lastPost = this.props.currentUsersLatestPost;
        if (!lastPost) {
            return;
        }

        let type;
        if (lastPost.root_id && lastPost.root_id.length > 0) {
            type = Utils.localizeMessage('create_post.comment', Posts.MESSAGE_TYPES.COMMENT);
        } else {
            type = Utils.localizeMessage('create_post.post', Posts.MESSAGE_TYPES.POST);
        }
        if (this.textboxRef.current) {
            this.textboxRef.current.blur();
        }
        this.props.actions.setEditingPost(lastPost.id, 'post_textbox', type);
    };

    replyToLastPost = (e: React.KeyboardEvent) => {
        e.preventDefault();
        const latestReplyablePostId = this.props.latestReplyablePostId;
        const replyBox = document.getElementById('reply_textbox');
        if (replyBox) {
            replyBox.focus();
        }
        if (latestReplyablePostId) {
            this.props.actions.selectPostFromRightHandSideSearchByPostId(latestReplyablePostId);
        }
    };

    loadPrevMessage = (e: React.KeyboardEvent) => {
        e.preventDefault();
        this.props.actions.moveHistoryIndexBack(Posts.MESSAGE_TYPES.POST).then(() => this.fillMessageFromHistory());
    };

    loadNextMessage = (e: React.KeyboardEvent) => {
        e.preventDefault();
        this.props.actions.moveHistoryIndexForward(Posts.MESSAGE_TYPES.POST).then(() => this.fillMessageFromHistory());
    };

    reactToLastMessage = (e: KeyboardEvent) => {
        e.preventDefault();

        const {rhsExpanded, actions: {emitShortcutReactToLastPostFrom}} = this.props;
        const noModalsAreOpen = document.getElementsByClassName(A11yClassNames.MODAL).length === 0;
        const noPopupsDropdownsAreOpen = document.getElementsByClassName(A11yClassNames.POPUP).length === 0;

        // Block keyboard shortcut react to last message when :
        // - RHS is completely expanded
        // - Any dropdown/popups are open
        // - Any modals are open
        if (!rhsExpanded && noModalsAreOpen && noPopupsDropdownsAreOpen) {
            emitShortcutReactToLastPostFrom(Locations.CENTER);
        }
    };

    handleBlur = () => {
        if (!this.isDraftSubmitting) {
            this.saveDraft();
        }

        this.lastBlurAt = Date.now();
    };

    handleEmojiClose = () => {
        this.setState({showEmojiPicker: false});
    };

    setMessageAndCaretPosition = (newMessage: string, newCaretPosition: number) => {
        const textbox = this.textboxRef.current?.getInputBox();

        this.setState({
            message: newMessage,
            caretPosition: newCaretPosition,
        }, () => {
            Utils.setCaretPosition(textbox, newCaretPosition);

            const draft = {
                ...this.props.draft,
                message: this.state.message,
            };

            this.handleDraftChange(draft);
        });
    };

    prefillMessage = (message: string, shouldFocus?: boolean) => {
        this.setMessageAndCaretPosition(message, message.length);

        if (shouldFocus) {
            const inputBox = this.textboxRef.current?.getInputBox();
            if (inputBox) {
                // programmatic click needed to close the create post tip
                inputBox.click();
            }
            this.focusTextbox(true);
        }
    };

    handleEmojiClick = (emoji: Emoji) => {
        const emojiAlias = ('short_names' in emoji && emoji.short_names && emoji.short_names[0]) || emoji.name;

        if (!emojiAlias) {
            //Oops.. There went something wrong
            return;
        }

        if (this.state.message === '') {
            const newMessage = ':' + emojiAlias + ': ';
            this.setMessageAndCaretPosition(newMessage, newMessage.length);
        } else {
            const {message} = this.state;
            const {firstPiece, lastPiece} = splitMessageBasedOnCaretPosition(this.state.caretPosition, message);

            // check whether the first piece of the message is empty when cursor is placed at beginning of message and avoid adding an empty string at the beginning of the message
            const newMessage =
                firstPiece === '' ? `:${emojiAlias}: ${lastPiece}` : `${firstPiece} :${emojiAlias}: ${lastPiece}`;

            const newCaretPosition =
                firstPiece === '' ? `:${emojiAlias}: `.length : `${firstPiece} :${emojiAlias}: `.length;
            this.setMessageAndCaretPosition(newMessage, newCaretPosition);
        }

        this.handleEmojiClose();
    };

    handleGifClick = (gif: string) => {
        if (this.state.message === '') {
            this.setState({message: gif});
        } else {
            const newMessage = (/\s+$/).test(this.state.message) ? this.state.message + gif : this.state.message + ' ' + gif;
            this.setState({message: newMessage});

            const draft = {
                ...this.props.draft,
                message: newMessage,
            };

            this.handleDraftChange(draft);
        }
        this.handleEmojiClose();
    };

    toggleAdvanceTextEditor = () => {
        this.setState({
            isFormattingBarHidden:
                !this.state.isFormattingBarHidden,
        });
        this.props.actions.savePreferences(this.props.currentUserId, [{
            category: Preferences.ADVANCED_TEXT_EDITOR,
            user_id: this.props.currentUserId,
            name: AdvancedTextEditorConst.POST,
            value: String(!this.state.isFormattingBarHidden),
        }]);
    };

    handleRemovePriority = () => {
        this.handlePostPriorityApply();
    };

    handlePostPriorityApply = (settings?: PostPriorityMetadata) => {
        const updatedDraft = {
            ...this.props.draft,
        };

        if (settings?.priority || settings?.requested_ack) {
            updatedDraft.metadata = {
                priority: {
                    ...settings,
                    priority: settings!.priority || '',
                    requested_ack: settings!.requested_ack,
                },
            };
        } else {
            updatedDraft.metadata = {};
        }

        this.handleDraftChange(updatedDraft, true);
        this.focusTextbox();
    };

    handlePostPriorityHide = () => {
        this.focusTextbox(true);
    };

    hasPrioritySet = () => {
        return (
            this.props.isPostPriorityEnabled &&
            this.props.draft.metadata?.priority && (
                this.props.draft.metadata.priority.priority ||
                this.props.draft.metadata.priority.requested_ack
            )
        );
    };

    isValidPersistentNotifications = (): boolean => {
        if (!this.hasPrioritySet()) {
            return true;
        }

        const {currentChannel} = this.props;
        const {priority, persistent_notifications: persistentNotifications} = this.props.draft.metadata!.priority!;
        if (priority !== PostPriority.URGENT || !persistentNotifications) {
            return true;
        }

        if (currentChannel.type === Constants.DM_CHANNEL) {
            return true;
        }

        if (this.hasSpecialMentions()) {
            return false;
        }

        const mentions = mentionsMinusSpecialMentionsInText(this.state.message);

        return mentions.length > 0;
    };

    getSpecialMentions = (): {[key: string]: boolean} => {
        return specialMentionsInText(this.state.message);
    };

    hasSpecialMentions = (): boolean => {
        return Object.values(this.getSpecialMentions()).includes(true);
    };

    onMessageChange = (message: string, callback?: (() => void) | undefined) => {
        this.handleDraftChange({
            ...this.props.draft,
            message,
        });
        this.setState({message}, callback);
    };

    render() {
        const {draft, canPost} = this.props;

        let centerClass = '';
        if (!this.props.fullWidthTextBox) {
            centerClass = 'center';
        }

        if (!this.props.currentChannel || !this.props.currentChannel.id) {
            return null;
        }

        return (
            <Foo
                location={Locations.CENTER}
                textboxRef={this.textboxRef}
                currentUserId={this.props.currentUserId}
                message={this.state.message}
                showEmojiPicker={this.state.showEmojiPicker}
                currentChannel={this.state.currentChannel}
                channelId={this.props.currentChannel.id}
                postId={''}
                errorClass={this.state.errorClass}
                serverError={this.state.serverError}
                isFormattingBarHidden={this.state.isFormattingBarHidden}
                draft={draft}
                showSendTutorialTip={this.props.showSendTutorialTip}
                handleSubmit={this.handleSubmit}
                removePreview={this.removePreview}
                setShowPreview={this.setShowPreview}
                shouldShowPreview={this.props.shouldShowPreview}
                canPost={canPost}
                useChannelMentions={this.props.useChannelMentions}
                handleBlur={this.handleBlur}
                postError={this.state.postError}
                handlePostError={this.handlePostError}
                emitTypingEvent={this.emitTypingEvent}
                handleMouseUpKeyUp={this.handleMouseUpKeyUp}
                onKeyPress={this.postMsgKeyPress}
                handleChange={this.handleChange}
                toggleEmojiPicker={this.toggleEmojiPicker}
                handleGifClick={this.handleGifClick}
                handleEmojiClick={this.handleEmojiClick}
                hideEmojiPicker={this.hideEmojiPicker}
                toggleAdvanceTextEditor={this.toggleAdvanceTextEditor}
                handleUploadError={this.handleUploadError}
                handleFileUploadComplete={this.handleFileUploadComplete}
                handleUploadStart={this.handleUploadStart}
                handleFileUploadChange={this.handleFileUploadChange}
                fileUploadRef={this.fileUploadRef}
                prefillMessage={this.prefillMessage}
                disableSend={!this.isValidPersistentNotifications()}
                priorityLabel={this.hasPrioritySet() ? (
                    <PriorityLabels
                        canRemove={!this.props.shouldShowPreview}
                        hasError={!this.isValidPersistentNotifications()}
                        specialMentions={this.getSpecialMentions()}
                        onRemove={this.handleRemovePriority}
                        persistentNotifications={draft!.metadata!.priority?.persistent_notifications}
                        priority={draft!.metadata!.priority?.priority}
                        requestedAck={draft!.metadata!.priority?.requested_ack}
                    />
                ) : undefined}
                priorityControls={this.props.isPostPriorityEnabled ? (
                    <PostPriorityPickerOverlay
                        key='post-priority-picker-key'
                        settings={draft?.metadata?.priority}
                        onApply={this.handlePostPriorityApply}
                        onClose={this.handlePostPriorityHide}
                        disabled={this.props.shouldShowPreview}
                    />
                ) : undefined}
                formId={'create_post'}
                formClass={centerClass}
                onEditLatestPost={this.editLastPost}
                ctrlSend={this.props.ctrlSend}
                codeBlockOnCtrlEnter={this.props.codeBlockOnCtrlEnter}
                onMessageChange={this.onMessageChange}
                replyToLastPost={this.replyToLastPost}
                loadNextMessage={this.loadNextMessage}
                loadPrevMessage={this.loadPrevMessage}
                caretPosition={this.state.caretPosition}
                saveDraft={this.saveDraft}
            />
        );
    }
}

export default AdvancedCreatePost;
