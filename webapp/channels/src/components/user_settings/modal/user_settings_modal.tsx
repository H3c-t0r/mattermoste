// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {Modal} from 'react-bootstrap';
import ReactDOM from 'react-dom';
import {injectIntl} from 'react-intl';
import type {IntlShape} from 'react-intl';
import Constants from 'utils/constants';
import {cmdOrCtrlPressed, isKeyPressed} from 'utils/keyboard';
import {stopTryNotificationRing} from 'utils/notification_sounds';

import type {PreferencesType} from '@mattermost/types/lib/preferences';
import type {UserProfile} from '@mattermost/types/users';

import type {ActionResult} from 'mattermost-redux/types/actions';

import ConfirmModal from 'components/confirm_modal';
import SettingsSidebar from 'components/settings_sidebar';
import UserSettings from 'components/user_settings';

import type {PluginConfiguration} from 'types/plugins/user_settings';

export type OwnProps = {
    userID: string;
    adminMode?: boolean;
    currentUser?: UserProfile;
    isContentProductSettings: boolean;
    userPreferences?: PreferencesType;
}

export type Props = OwnProps & {
    onExited: () => void;
    intl: IntlShape;
    actions: {
        sendVerificationEmail: (email: string) => Promise<ActionResult>;
    };
    pluginSettings: {[pluginId: string]: PluginConfiguration};
}

type State = {
    active_tab?: string;
    active_section: string;
    showConfirmModal: boolean;
    enforceFocus?: boolean;
    show: boolean;
    resendStatus: string;
}

class UserSettingsModal extends React.PureComponent<Props, State> {
    private requireConfirm: boolean;
    private customConfirmAction: ((handleConfirm: () => void) => void) | null;
    private modalBodyRef: React.RefObject<Modal>;
    private afterConfirm: (() => void) | null;

    constructor(props: Props) {
        super(props);

        console.log({
            adminMode: this.props.adminMode,
            currentUser: this.props.currentUser,
        });

        this.state = {
            active_tab: props.isContentProductSettings ? 'notifications' : 'profile',
            active_section: '',
            showConfirmModal: false,
            enforceFocus: true,
            show: true,
            resendStatus: '',
        };

        this.requireConfirm = false;

        // Used when settings want to override the default confirm modal with their own
        // If set by a child, it will be called in place of showing the regular confirm
        // modal. It will be passed a function to call on modal confirm
        this.customConfirmAction = null;
        this.afterConfirm = null;

        this.modalBodyRef = React.createRef();
    }

    handleResend = (email: string) => {
        this.setState({resendStatus: 'sending'});

        this.props.actions.sendVerificationEmail(email).then(({data, error: err}) => {
            if (data) {
                this.setState({resendStatus: 'success'});
            } else if (err) {
                this.setState({resendStatus: 'failure'});
            }
        });
    };

    componentDidMount() {
        document.addEventListener('keydown', this.handleKeyDown);
    }

    componentWillUnmount() {
        document.removeEventListener('keydown', this.handleKeyDown);
    }

    componentDidUpdate(prevProps: Props, prevState: State) {
        if (this.state.active_tab !== prevState.active_tab) {
            const el = ReactDOM.findDOMNode(this.modalBodyRef.current) as any;
            el.scrollTop = 0;
        }
    }

    handleKeyDown = (e: KeyboardEvent) => {
        if (cmdOrCtrlPressed(e) && e.shiftKey && isKeyPressed(e, Constants.KeyCodes.A)) {
            e.preventDefault();
            this.handleHide();
        }
    };

    // Called when the close button is pressed on the main modal
    handleHide = () => {
        if (this.requireConfirm) {
            this.showConfirmModal(() => this.handleHide());
            return;
        }

        // Cancel any ongoing notification sound, if any (from DesktopNotificationSettings)
        stopTryNotificationRing();

        this.setState({
            show: false,
        });
    };

    // called after the dialog is fully hidden and faded out
    handleHidden = () => {
        this.setState({
            active_tab: this.props.isContentProductSettings ? 'notifications' : 'profile',
            active_section: '',
        });
        this.props.onExited();
    };

    // Called to hide the settings pane when on mobile
    handleCollapse = () => {
        const el = ReactDOM.findDOMNode(this.modalBodyRef.current) as HTMLDivElement;
        el.closest('.modal-dialog')!.classList.remove('display--content');

        this.setState({
            active_tab: '',
            active_section: '',
        });
    };

    handleConfirm = () => {
        this.setState({
            showConfirmModal: false,
            enforceFocus: true,
        });

        this.requireConfirm = false;
        this.customConfirmAction = null;

        if (this.afterConfirm) {
            this.afterConfirm();
            this.afterConfirm = null;
        }
    };

    handleCancelConfirmation = () => {
        this.setState({
            showConfirmModal: false,
            enforceFocus: true,
        });

        this.afterConfirm = null;
    };

    showConfirmModal = (afterConfirm: () => void) => {
        if (afterConfirm) {
            this.afterConfirm = afterConfirm;
        }

        if (this.customConfirmAction) {
            this.customConfirmAction(this.handleConfirm);
            return;
        }

        this.setState({
            showConfirmModal: true,
            enforceFocus: false,
        });
    };

    // Called by settings tabs when their close button is pressed
    closeModal = () => {
        if (this.requireConfirm) {
            this.showConfirmModal(this.closeModal);
        } else {
            this.handleHide();
        }
    };

    // Called by settings tabs when their back button is pressed
    collapseModal = () => {
        if (this.requireConfirm) {
            this.showConfirmModal(this.collapseModal);
        } else {
            this.handleCollapse();
        }
    };

    updateTab = (tab?: string, skipConfirm?: boolean) => {
        if (!skipConfirm && this.requireConfirm) {
            this.showConfirmModal(() => this.updateTab(tab, true));
        } else {
            this.setState({
                active_tab: tab,
                active_section: '',
            });
        }
    };

    updateSection = (section?: string, skipConfirm?: boolean) => {
        if (!skipConfirm && this.requireConfirm) {
            this.showConfirmModal(() => this.updateSection(section, true));
        } else {
            this.setState({
                active_section: section ?? '',
            });
        }
    };

    getUserSettingsTabs = () => {
        return [
            {
                name: 'notifications',
                uiName: this.props.intl.formatMessage({id: 'user.settings.modal.notifications', defaultMessage: 'Notifications'}),
                icon: 'icon icon-bell-outline',
                iconTitle: this.props.intl.formatMessage({id: 'user.settings.notifications.icon', defaultMessage: 'Notification Settings Icon'}),
            },
            {
                name: 'display',
                uiName: this.props.intl.formatMessage({id: 'user.settings.modal.display', defaultMessage: 'Display'}),
                icon: 'icon icon-eye-outline',
                iconTitle: this.props.intl.formatMessage({id: 'user.settings.display.icon', defaultMessage: 'Display Settings Icon'}),
            },
            {
                name: 'sidebar',
                uiName: this.props.intl.formatMessage({id: 'user.settings.modal.sidebar', defaultMessage: 'Sidebar'}),
                icon: 'icon icon-dock-left',
                iconTitle: this.props.intl.formatMessage({id: 'user.settings.sidebar.icon', defaultMessage: 'Sidebar Settings Icon'}),
            },
            {
                name: 'advanced',
                uiName: this.props.intl.formatMessage({id: 'user.settings.modal.advanced', defaultMessage: 'Advanced'}),
                icon: 'icon icon-tune',
                iconTitle: this.props.intl.formatMessage({id: 'user.settings.advance.icon', defaultMessage: 'Advanced Settings Icon'}),
            },
        ];
    };

    getProfileSettingsTab = () => {
        return [
            {
                name: 'profile',
                uiName: this.props.intl.formatMessage({id: 'user.settings.modal.profile', defaultMessage: 'Profile'}),
                icon: 'icon icon-settings-outline',
                iconTitle: this.props.intl.formatMessage({id: 'user.settings.profile.icon', defaultMessage: 'Profile Settings Icon'}),
            },
            {
                name: 'security',
                uiName: this.props.intl.formatMessage({id: 'user.settings.modal.security', defaultMessage: 'Security'}),
                icon: 'icon icon-lock-outline',
                iconTitle: this.props.intl.formatMessage({id: 'user.settings.security.icon', defaultMessage: 'Security Settings Icon'}),
            },
        ];
    };

    getPluginsSettingsTab = () => {
        return Object.values(this.props.pluginSettings).map((v) => ({
            name: v.id,
            uiName: v.uiName,
            icon: v.icon ? {url: v.icon} : 'icon-power-plug-outline',
            iconTitle: v.uiName,
        }));
    };

    render() {
        const {formatMessage} = this.props.intl;
        if (this.props.currentUser == null) {
            return (<div/>);
        }

        const modalTitle = this.props.isContentProductSettings ? formatMessage({
            id: 'global_header.productSettings',
            defaultMessage: 'Settings',
        }) : formatMessage({
            id: 'user.settings.modal.title',
            defaultMessage: 'Profile',
        });

        return (
            <Modal
                id='accountSettingsModal'
                dialogClassName='a11y__modal settings-modal'
                show={this.state.show}
                onHide={this.handleHide}
                onExited={this.handleHidden}
                enforceFocus={this.state.enforceFocus}
                role='dialog'
                aria-label={modalTitle}
            >
                <Modal.Header
                    id='accountSettingsHeader'
                    closeButton={true}
                >
                    <Modal.Title
                        componentClass='h1'
                        id='accountSettingsModalLabel'
                    >
                        {modalTitle}
                    </Modal.Title>
                </Modal.Header>
                <Modal.Body ref={this.modalBodyRef}>
                    <div className='settings-table'>
                        <div className='settings-links'>
                            <SettingsSidebar
                                tabs={this.props.isContentProductSettings ? this.getUserSettingsTabs() : this.getProfileSettingsTab()}
                                pluginTabs={this.props.isContentProductSettings ? this.getPluginsSettingsTab() : []}
                                activeTab={this.state.active_tab}
                                updateTab={this.updateTab}
                            />
                        </div>
                        <div className='settings-content minimize-settings'>
                            <UserSettings
                                activeTab={this.state.active_tab}
                                activeSection={this.state.active_section}
                                updateSection={this.updateSection}
                                updateTab={this.updateTab}
                                closeModal={this.closeModal}
                                collapseModal={this.collapseModal}
                                setEnforceFocus={(enforceFocus?: boolean) => this.setState({enforceFocus})}
                                setRequireConfirm={
                                    (requireConfirm?: boolean, customConfirmAction?: () => () => void) => {
                                        this.requireConfirm = requireConfirm!;
                                        this.customConfirmAction = customConfirmAction!;
                                    }
                                }
                                pluginSettings={this.props.pluginSettings}
                                user={this.props.currentUser}
                                adminMode={this.props.adminMode}
                                userPreferences={this.props.userPreferences}
                            />
                        </div>
                    </div>
                </Modal.Body>
                <ConfirmModal
                    title={formatMessage({id: 'user.settings.modal.confirmTitle', defaultMessage: 'Discard Changes?'})}
                    message={formatMessage({
                        id: 'user.settings.modal.confirmMsg',
                        defaultMessage: 'You have unsaved changes, are you sure you want to discard them?',
                    })}
                    confirmButtonText={formatMessage({id: 'user.settings.modal.confirmBtns', defaultMessage: 'Yes, Discard'})}
                    show={this.state.showConfirmModal}
                    onConfirm={this.handleConfirm}
                    onCancel={this.handleCancelConfirmation}
                />
            </Modal>
        );
    }
}

export default injectIntl(UserSettingsModal);
