// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import $ from 'jquery';

import SearchResults from './search_results.jsx';
import RhsThread from './rhs_thread.jsx';
import Webrtc from './webrtc/webrtc_controller.jsx';
import SearchStore from 'stores/search_store.jsx';
import PostStore from 'stores/post_store.jsx';
import UserStore from 'stores/user_store.jsx';
import PreferenceStore from 'stores/preference_store.jsx';
import WebrtcStore from 'stores/webrtc_store.jsx';
import Constants from 'utils/constants.jsx';
import * as Utils from 'utils/utils.jsx';

import React from 'react';

export default class SidebarRight extends React.Component {
    constructor(props) {
        super(props);

        this.plScrolledToBottom = true;

        this.onPreferenceChange = this.onPreferenceChange.bind(this);
        this.onSelectedChange = this.onSelectedChange.bind(this);
        this.onSearchChange = this.onSearchChange.bind(this);
        this.onUserChange = this.onUserChange.bind(this);
        this.onShowSearch = this.onShowSearch.bind(this);
        this.onInitializeVideoCall = this.onInitializeVideoCall.bind(this);

        this.doStrangeThings = this.doStrangeThings.bind(this);

        this.state = {
            searchVisible: SearchStore.getSearchResults() !== null,
            isMentionSearch: SearchStore.getIsMentionSearch(),
            postRightVisible: !!PostStore.getSelectedPost(),
            fromSearch: false,
            currentUser: UserStore.getCurrentUser(),
            useMilitaryTime: PreferenceStore.getBool(Constants.Preferences.CATEGORY_DISPLAY_SETTINGS, Constants.Preferences.USE_MILITARY_TIME, false),
            videoCallVisible: false,
            isCaller: false,
            videoCallWithUserId: null
        };
    }

    componentDidMount() {
        SearchStore.addSearchChangeListener(this.onSearchChange);
        PostStore.addSelectedPostChangeListener(this.onSelectedChange);
        SearchStore.addShowSearchListener(this.onShowSearch);
        UserStore.addChangeListener(this.onUserChange);
        PreferenceStore.addChangeListener(this.onPreferenceChange);
        WebrtcStore.addCallListener(this.onInitializeVideoCall);
        this.doStrangeThings();
    }

    componentWillUnmount() {
        SearchStore.removeSearchChangeListener(this.onSearchChange);
        PostStore.removeSelectedPostChangeListener(this.onSelectedChange);
        SearchStore.removeShowSearchListener(this.onShowSearch);
        UserStore.removeChangeListener(this.onUserChange);
        PreferenceStore.removeChangeListener(this.onPreferenceChange);
        WebrtcStore.removeCallListener(this.onInitializeVideoCall);
    }

    shouldComponentUpdate(nextProps, nextState) {
        return !Utils.areObjectsEqual(nextState, this.state);
    }

    componentWillUpdate(nextProps, nextState) {
        const isOpen = this.state.searchVisible || this.state.postRightVisible || this.state.videoCallVisible;
        const willOpen = nextState.searchVisible || nextState.postRightVisible || nextState.videoCallVisible;

        if (isOpen !== willOpen) {
            PostStore.jumpPostsViewSidebarOpen();
        }
    }

    doStrangeThings() {
        // We should have a better way to do this stuff
        // Hence the function name.
        $('.app__body .inner-wrap').removeClass('.move--right');
        $('.app__body .inner-wrap').addClass('move--left');
        $('.app__body .sidebar--left').removeClass('move--right');
        $('.app__body .sidebar--right').addClass('move--left');

        //$('.sidebar--right').prepend('<div class="sidebar__overlay"></div>');
        if (!this.state.searchVisible && !this.state.postRightVisible && !this.state.videoCallVisible) {
            $('.app__body .inner-wrap').removeClass('move--left').removeClass('move--right');
            $('.app__body .sidebar--right').removeClass('move--left');
            return (
                <div></div>
            );
        }

        /*setTimeout(() => {
            $('.sidebar__overlay').fadeOut('200', () => {
                $('.sidebar__overlay').remove();
            });
            }, 500);*/
        return null;
    }

    componentDidUpdate() {
        this.doStrangeThings();
    }

    onPreferenceChange() {
        this.setState({
            useMilitaryTime: PreferenceStore.getBool(Constants.Preferences.CATEGORY_DISPLAY_SETTINGS, Constants.Preferences.USE_MILITARY_TIME, false)
        });
    }

    onSelectedChange(fromSearch) {
        this.setState({
            videoCallVisible: false,
            postRightVisible: !!PostStore.getSelectedPost(),
            fromSearch
        });
    }

    onSearchChange() {
        this.setState({
            videoCallVisible: false,
            searchVisible: SearchStore.getSearchResults() !== null,
            isMentionSearch: SearchStore.getIsMentionSearch()
        });
    }

    onUserChange() {
        this.setState({
            currentUser: UserStore.getCurrentUser()
        });
    }

    onShowSearch() {
        if (!this.state.searchVisible) {
            this.setState({
                videoCallVisible: false,
                searchVisible: true
            });
        }
    }

    onInitializeVideoCall(userId, isCaller) {
        this.setState({
            searchVisible: false,
            isMentionSearch: false,
            postRightVisible: false,
            videoCallVisible: true,
            isCaller,
            videoCallWithUserId: userId
        });
    }

    render() {
        let content = null;

        if (this.state.searchVisible) {
            content = (
                <SearchResults
                    isMentionSearch={this.state.isMentionSearch}
                    useMilitaryTime={this.state.useMilitaryTime}
                />
            );
        } else if (this.state.postRightVisible) {
            content = (
                <RhsThread
                    fromSearch={this.state.fromSearch}
                    isMentionSearch={this.state.isMentionSearch}
                    currentUser={this.state.currentUser}
                    useMilitaryTime={this.state.useMilitaryTime}
                />
            );
        } else if (this.state.videoCallVisible) {
            content = (
                <Webrtc
                    isMentionSearch={this.state.isMentionSearch}
                    currentUser={this.state.currentUser}
                    userId={this.state.videoCallWithUserId}
                    isCaller={this.state.isCaller}
                />
            );
        }

        return (
            <div
                className='sidebar--right'
                id='sidebar-right'
            >
                <div className='sidebar-right-container'>
                    {content}
                </div>
            </div>
        );
    }
}
