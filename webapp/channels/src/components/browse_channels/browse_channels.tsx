// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {FormattedMessage} from 'react-intl';

import {ActionResult} from 'mattermost-redux/types/actions';
import {RelationOneToOne} from '@mattermost/types/utilities';
import {Channel, ChannelMembership, ChannelSearchOpts} from '@mattermost/types/channels';
import Permissions from 'mattermost-redux/constants/permissions';

import NewChannelModal from 'components/new_channel_modal/new_channel_modal';
import SearchableChannelList from 'components/searchable_channel_list';
import TeamPermissionGate from 'components/permissions_gates/team_permission_gate';

import {ModalData} from 'types/actions';
import {RhsState} from 'types/store/rhs';

import {getHistory} from 'utils/browser_history';
import {ModalIdentifiers, RHSStates, StoragePrefixes} from 'utils/constants';
import {getRelativeChannelURL} from 'utils/url';
import {GenericModal} from '@mattermost/components';
import classNames from 'classnames';
import {localizeMessage} from 'utils/utils';
import LoadingScreen from 'components/loading_screen';

import './browse_channels.scss';

const CHANNELS_CHUNK_SIZE = 50;
const CHANNELS_PER_PAGE = 50;
const SEARCH_TIMEOUT_MILLISECONDS = 100;
export const FILTER = {
    all: 'all' as const,
    public: 'public' as const,
    private: 'private' as const,
    archived: 'archived' as const,
};
export type FilterType = keyof typeof FILTER;

type Actions = {
    getChannels: (teamId: string, page: number, perPage: number) => Promise<ActionResult<Channel[], Error>>;
    getArchivedChannels: (teamId: string, page: number, channelsPerPage: number) => Promise<ActionResult<Channel[], Error>>;
    joinChannel: (currentUserId: string, teamId: string, channelId: string) => Promise<ActionResult>;
    searchAllChannels: (term: string, opts?: ChannelSearchOpts) => Promise<ActionResult<Channel[], Error>>;
    openModal: <P>(modalData: ModalData<P>) => void;
    closeModal: (modalId: string) => void;

    /*
     * Function to set a key-value pair in the local storage
     */
    setGlobalItem: (name: string, value: string) => void;
    closeRightHandSide: () => void;
    getChannelsMemberCount: (channelIds: string[]) => Promise<ActionResult>;
}

export type Props = {
    channels: Channel[];
    archivedChannels: Channel[];
    privateChannels: Channel[];
    currentUserId: string;
    teamId: string;
    teamName: string;
    channelsRequestStarted?: boolean;
    canShowArchivedChannels?: boolean;
    myChannelMemberships: RelationOneToOne<Channel, ChannelMembership>;
    shouldHideJoinedChannels: boolean;
    rhsState?: RhsState;
    rhsOpen?: boolean;
    channelsMemberCount?: Record<string, number>;
    actions: Actions;
}

type State = {
    loading: boolean;
    filter: FilterType;
    search: boolean;
    searchedChannels: Channel[];
    serverError: React.ReactNode | string;
    searching: boolean;
    searchTerm: string;
}

export default class BrowseChannels extends React.PureComponent<Props, State> {
    public searchTimeoutId: number;
    activeChannels: Channel[] = [];

    constructor(props: Props) {
        super(props);

        this.searchTimeoutId = 0;

        this.state = {
            loading: true,
            filter: FILTER.all,
            search: false,
            searchedChannels: [],
            serverError: null,
            searching: false,
            searchTerm: '',
        };
    }

    componentDidMount() {
        const promises = [
            this.props.actions.getChannels(this.props.teamId, 0, CHANNELS_CHUNK_SIZE * 2),
        ];

        if (this.props.canShowArchivedChannels) {
            promises.push(this.props.actions.getArchivedChannels(this.props.teamId, 0, CHANNELS_CHUNK_SIZE * 2));
        }

        Promise.all(promises).then((results) => {
            const channelIDsForMemberCount = results.flatMap((result) => {
                return result.data ? result.data.map((channel) => channel.id) : [];
            },
            );
            if (channelIDsForMemberCount.length > 0) {
                this.props.actions.getChannelsMemberCount(channelIDsForMemberCount);
            }
        });
        this.loadComplete();
    }

    loadComplete = () => {
        this.setState({loading: false});
    };

    handleNewChannel = () => {
        this.handleExit();
        this.closeEditRHS();
        this.props.actions.openModal({
            modalId: ModalIdentifiers.NEW_CHANNEL_MODAL,
            dialogType: NewChannelModal,
        });
    };

    handleExit = () => {
        this.props.actions.closeModal(ModalIdentifiers.MORE_CHANNELS);
    };

    closeEditRHS = () => {
        if (this.props.rhsOpen && this.props.rhsState === RHSStates.EDIT_HISTORY) {
            this.props.actions.closeRightHandSide();
        }
    };

    onChange = (force: boolean) => {
        if (this.state.search && !force) {
            return;
        }

        this.setState({
            searchedChannels: [],
            serverError: null,
        });
    };

    nextPage = (page: number) => {
        this.props.actions.getChannels(this.props.teamId, page + 1, CHANNELS_PER_PAGE).then((result) => {
            if (result.data && result.data.length > 0) {
                this.props.actions.getChannelsMemberCount(result.data.map((channel) => channel.id));
            }
        });
    };

    handleJoin = async (channel: Channel, done: () => void) => {
        const {actions, currentUserId, teamId, teamName} = this.props;
        let result;

        if (!this.isMemberOfChannel(channel.id)) {
            result = await actions.joinChannel(currentUserId, teamId, channel.id);
        }

        if (result?.error) {
            this.setState({serverError: result.error.message});
        } else {
            this.props.actions.getChannelsMemberCount([channel.id]);
            getHistory().push(getRelativeChannelURL(teamName, channel.name));
            this.closeEditRHS();
        }

        if (done) {
            done();
        }
    };

    search = (term: string) => {
        clearTimeout(this.searchTimeoutId);

        if (term === '') {
            this.onChange(true);
            this.setState({search: false, searchedChannels: [], searching: false, searchTerm: term});
            this.searchTimeoutId = 0;
            return;
        }
        this.setState({search: true, searching: true, searchTerm: term});

        const searchTimeoutId = window.setTimeout(
            async () => {
                try {
                    // todo sinan - integrate page and per page similar to similar to getChannels
                    const {data} = await this.props.actions.searchAllChannels(term, {team_ids: [this.props.teamId], nonAdminSearch: false});
                    if (searchTimeoutId !== this.searchTimeoutId) {
                        return;
                    }

                    if (data) {
                        this.setSearchResults(data);
                    } else {
                        this.setState({searchedChannels: [], searching: false});
                    }
                } catch (ignoredErr) {
                    this.setState({searchedChannels: [], searching: false});
                }
            },
            SEARCH_TIMEOUT_MILLISECONDS,
        );

        this.searchTimeoutId = searchTimeoutId;
    };

    //todo sinan fix. filter based on the filter
    setSearchResults = (channels: Channel[]) => {
        this.setState({searchedChannels: this.state.shouldShowArchivedChannels ? channels.filter((c) => c.delete_at !== 0) : channels.filter((c) => c.delete_at === 0), searching: false});
    };

    changeFilter = (filter: FilterType) => {
        // search again when switching channels to update search results
        this.search(this.state.searchTerm);
        this.setState({filter});
    };

    isMemberOfChannel(channelId: string) {
        return this.props.myChannelMemberships[channelId];
    }

    handleShowJoinedChannelsPreference = (shouldHideJoinedChannels: boolean) => {
        // search again when switching channels to update search results
        this.search(this.state.searchTerm);
        this.props.actions.setGlobalItem(StoragePrefixes.HIDE_JOINED_CHANNELS, shouldHideJoinedChannels.toString());
    };

    getActiveChannels = (options: Pick<Props, 'channels' | 'archivedChannels' | 'privateChannels' | 'shouldHideJoinedChannels'> & Pick<State, 'search' | 'filter' | 'searchedChannels'>) => {
        const {channels, archivedChannels, privateChannels, shouldHideJoinedChannels, search, filter, searchedChannels} = options;

        const allChannels = channels.concat(archivedChannels, privateChannels);
        const filterChannelsWithoutJoined = (channelList: Channel[]) => channelList.filter((channel) => !this.isMemberOfChannel(channel.id));
        const allChannelsWithoutJoined = filterChannelsWithoutJoined(allChannels);
        const otherChannelsWithoutJoined = filterChannelsWithoutJoined(channels);
        const archivedChannelsWithoutJoined = filterChannelsWithoutJoined(archivedChannels);
        const privateChannelsWithoutJoined = filterChannelsWithoutJoined(privateChannels);

        if (search) {
            return searchedChannels;
        }

        if (filter === FILTER.archived) {
            return shouldHideJoinedChannels ? archivedChannelsWithoutJoined : archivedChannels;
        }

        if (filter === FILTER.private) {
            return shouldHideJoinedChannels ? privateChannelsWithoutJoined : privateChannels;
        }

        if (filter === FILTER.public) {
            return shouldHideJoinedChannels ? otherChannelsWithoutJoined : channels;
        }

        return shouldHideJoinedChannels ? allChannelsWithoutJoined : allChannels;
    };

    render() {
        const {
            channels,
            archivedChannels,
            teamId,
            channelsRequestStarted,
            shouldHideJoinedChannels,
            privateChannels,
        } = this.props;

        const {
            search,
            searchedChannels,
            serverError: serverErrorState,
            searching,
            filter,
        } = this.state;

        this.activeChannels = this.getActiveChannels({channels, archivedChannels, privateChannels, shouldHideJoinedChannels, search, filter, searchedChannels});

        let serverError;
        if (serverErrorState) {
            serverError =
                <div className='form-group has-error'><label className='control-label'>{serverErrorState}</label></div>;
        }

        const createNewChannelButton = (className: string, icon?: JSX.Element) => {
            const buttonClassName = classNames('btn', className);
            return (
                <TeamPermissionGate
                    teamId={teamId}
                    permissions={[Permissions.CREATE_PUBLIC_CHANNEL]}
                >
                    <button
                        type='button'
                        id='createNewChannelButton'
                        className={buttonClassName}
                        onClick={this.handleNewChannel}
                        aria-label={localizeMessage('more_channels.create', 'Create New Channel')}
                    >
                        {icon}
                        <FormattedMessage
                            id='more_channels.create'
                            defaultMessage='Create New Channel'
                        />
                    </button>
                </TeamPermissionGate>
            );
        };

        const noResultsText = (
            <>
                <p className='secondary-message'>
                    <FormattedMessage
                        id='more_channels.searchError'
                        defaultMessage='Try searching different keywords, checking for typos or adjusting the filters.'
                    />
                </p>
                {createNewChannelButton('primaryButton', <i className='icon-plus'/>)}
            </>
        );

        const body = this.state.loading ? <LoadingScreen/> : (
            <React.Fragment>
                <SearchableChannelList
                    channels={this.activeChannels}
                    channelsPerPage={CHANNELS_PER_PAGE}
                    nextPage={this.nextPage}
                    isSearch={search}
                    search={this.search}
                    handleJoin={this.handleJoin}
                    noResultsText={noResultsText}
                    loading={search ? searching : channelsRequestStarted}
                    changeFilter={this.changeFilter}
                    filter={this.state.filter}
                    canShowArchivedChannels={this.props.canShowArchivedChannels}
                    myChannelMemberships={this.props.myChannelMemberships}
                    closeModal={this.props.actions.closeModal}
                    hideJoinedChannelsPreference={this.handleShowJoinedChannelsPreference}
                    rememberHideJoinedChannelsChecked={shouldHideJoinedChannels}
                    channelsMemberCount={this.props.channelsMemberCount}
                />
                {serverError}
            </React.Fragment>
        );

        const title = (
            <FormattedMessage
                id='more_channels.title'
                defaultMessage='Browse Channels'
            />
        );

        return (
            <GenericModal
                onExited={this.handleExit}
                id='browseChannelsModal'
                aria-labelledby='browseChannelsModalLabel'
                compassDesign={true}
                modalHeaderText={title}
                headerButton={createNewChannelButton('outlineButton')}
                autoCloseOnConfirmButton={false}
                aria-modal={true}
                enforceFocus={false}
                bodyPadding={false}
            >
                {body}
            </GenericModal>
        );
    }
}
