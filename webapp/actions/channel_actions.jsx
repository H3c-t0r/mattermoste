// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import {browserHistory} from 'react-router/es6';
import * as Utils from 'utils/utils.jsx';
import TeamStore from 'stores/team_store.jsx';
import UserStore from 'stores/user_store.jsx';
import ChannelStore from 'stores/channel_store.jsx';
import * as AsyncClient from 'utils/async_client.jsx';
import * as UserAgent from 'utils/user_agent.jsx';
import Client from 'client/web_client.jsx';

export function goToChannel(channel) {
    if (channel.fake) {
        Utils.openDirectChannelToUser(
            UserStore.getProfileByUsername(channel.display_name),
            () => {
                browserHistory.push(TeamStore.getCurrentTeamRelativeUrl() + '/channels/' + channel.name);
            },
            null
        );
    } else {
        browserHistory.push(TeamStore.getCurrentTeamRelativeUrl() + '/channels/' + channel.name);
    }
}

export function executeCommand(channelId, message, suggest, success, error) {
    let msg = message;

    msg = msg.substring(0, msg.indexOf(' ')).toLowerCase() + msg.substring(msg.indexOf(' '), msg.length);

    if (message.indexOf('/shortcuts') !== -1) {
        if (UserAgent.isMobileApp()) {
            const err = {message: Utils.localizeMessage('create_post.shortcutsNotSupported', 'Keyboard shortcuts are not supported on your device')};
            error(err);
            return;
        } else if (Utils.isMac()) {
            msg += ' mac';
        } else if (message.indexOf('mac') !== -1) {
            msg = '/shortcuts';
        }
    }

    Client.executeCommand(channelId, msg, suggest, success, error);
}

export function setChannelAsRead(channelIdParam) {
    const channelId = channelIdParam || ChannelStore.getCurrentId();
    AsyncClient.updateLastViewedAt();
    ChannelStore.resetCounts(channelId);
    ChannelStore.emitChange();
    if (channelId === ChannelStore.getCurrentId()) {
        ChannelStore.emitLastViewed(Number.MAX_VALUE, false);
    }
}

export function addUserToChannel(channelId, userId, success, error) {
    Client.addChannelMember(
        channelId,
        userId,
        (data) => {
            UserStore.removeProfileNotInChannel(channelId, userId);
            const profile = UserStore.getProfile(userId);
            if (profile) {
                UserStore.saveProfileInChannel(channelId, profile);
                UserStore.emitInChannelChange();
            }
            UserStore.emitNotInChannelChange();

            if (success) {
                success(data);
            }
        },
        (err) => {
            AsyncClient.dispatchError(err, 'addChannelMember');

            if (error) {
                error(err);
            }
        }
    );
}

export function removeUserFromChannel(channelId, userId, success, error) {
    Client.removeChannelMember(
        channelId,
        userId,
        (data) => {
            UserStore.removeProfileInChannel(channelId, userId);
            const profile = UserStore.getProfile(userId);
            if (profile) {
                UserStore.saveProfileNotInChannel(channelId, profile);
                UserStore.emitNotInChannelChange();
            }
            UserStore.emitInChannelChange();

            if (success) {
                success(data);
            }
        },
        (err) => {
            AsyncClient.dispatchError(err, 'removeChannelMember');

            if (error) {
                error(err);
            }
        }
    );
}
