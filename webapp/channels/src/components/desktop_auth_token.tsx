// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useState} from 'react';
import {FormattedMessage} from 'react-intl';
import {useDispatch} from 'react-redux';
import {useHistory, useLocation} from 'react-router-dom';

import crypto from 'crypto';
import classNames from 'classnames';

import {UserProfile} from '@mattermost/types/users';
import {DispatchFunc} from 'mattermost-redux/types/actions';

import {loginWithDesktopToken} from 'actions/views/login';

import './desktop_auth_token.scss';

const BOTTOM_MESSAGE_TIMEOUT = 10000;
const DESKTOP_AUTH_PREFIX = 'desktop_auth_client_token';

enum DesktopAuthStatus {
    None,
    WaitingForBrowser,
    LoggedIn,
    Authenticating,
    Error,
}

type Props = {
    href: string;
    onLogin: (userProfile: UserProfile) => void;
}

const DesktopAuthToken: React.FC<Props> = ({href, onLogin}: Props) => {
    const dispatch = useDispatch<DispatchFunc>();
    const history = useHistory();
    const {search} = useLocation();
    const query = new URLSearchParams(search);

    const serverToken = query.get('server_token');
    const clientToken = sessionStorage.getItem(DESKTOP_AUTH_PREFIX);
    const [status, setStatus] = useState(serverToken ? DesktopAuthStatus.LoggedIn : DesktopAuthStatus.None);
    const [showBottomMessage, setShowBottomMessage] = useState<React.ReactNode>();

    const tryDesktopLogin = async () => {
        if (!(clientToken && serverToken)) {
            setStatus(DesktopAuthStatus.Error);
            return;
        }

        sessionStorage.removeItem(DESKTOP_AUTH_PREFIX);
        const {data: userProfile, error: loginError} = await dispatch(loginWithDesktopToken(clientToken, serverToken));

        if (loginError && loginError.server_error_id && loginError.server_error_id.length !== 0) {
            if (loginError.server_error_id === 'app.desktop_token.validate.expired') {
                setStatus(DesktopAuthStatus.Error);
            }
            return;
        }

        setStatus(DesktopAuthStatus.LoggedIn);
        await onLogin(userProfile as UserProfile);
    };

    const openExternalLoginURL = () => {
        const desktopToken = crypto.randomBytes(32).toString('hex');
        sessionStorage.setItem(DESKTOP_AUTH_PREFIX, desktopToken);
        const parsedURL = new URL(href);

        const params = new URLSearchParams(parsedURL.searchParams);
        params.set('desktop_token', desktopToken);

        window.open(`${parsedURL.origin}${parsedURL.pathname}?${params.toString()}`);
        setStatus(DesktopAuthStatus.WaitingForBrowser);
    };

    const forwardToDesktopApp = () => {
        const url = new URL(window.location.href);
        url.protocol = 'mattermost';
        window.location.href = url.toString();
    };

    useEffect(() => {
        setShowBottomMessage(false);

        const timeout = setTimeout(() => {
            setShowBottomMessage(true);
        }, BOTTOM_MESSAGE_TIMEOUT) as unknown as number;

        return () => {
            clearTimeout(timeout);
        };
    }, [status]);

    useEffect(() => {
        if (status === DesktopAuthStatus.LoggedIn) {
            if (clientToken && serverToken) {
                tryDesktopLogin();
            } else {
                forwardToDesktopApp();
            }
            return;
        }

        openExternalLoginURL();
    }, []);

    let mainMessage;
    let subMessage;
    let bottomMessage;

    if (status === DesktopAuthStatus.WaitingForBrowser) {
        mainMessage = (
            <FormattedMessage
                id='desktop_auth_token.polling.redirectingToBrowser'
                defaultMessage='Redirecting to browser...'
            />
        );
        subMessage = (
            <FormattedMessage
                id='desktop_auth_token.polling.awaitingToken'
                defaultMessage='Authenticating in the browser, awaiting valid token.'
            />
        );

        bottomMessage = null;
    }

    if (status === DesktopAuthStatus.LoggedIn) {
        mainMessage = (
            <FormattedMessage
                id='desktop_auth_token.complete.youAreNowLoggedIn'
                defaultMessage='You are now logged in'
            />
        );
        subMessage = (
            <FormattedMessage
                id='desktop_auth_token.complete.openMattermost'
                defaultMessage='Click on <b>Open Mattermost</b> in the browser prompt to <a>launch the desktop app</a>'
                values={{
                    a: (chunks: React.ReactNode) => {
                        return (
                            <a onClick={forwardToDesktopApp}>
                                {chunks}
                            </a>
                        );
                    },
                    b: (chunks: React.ReactNode) => (<b>{chunks}</b>),
                }}
            />
        );

        bottomMessage = (
            <FormattedMessage
                id='desktop_auth_token.complete.havingTrouble'
                defaultMessage='Having trouble logging in? <a>Open Mattermost in your browser</a>'
                values={{
                    a: (chunks: React.ReactNode) => {
                        return (
                            <a onClick={() => history.push('/')}>
                                {chunks}
                            </a>
                        );
                    },
                }}
            />
        );
    }

    if (status === DesktopAuthStatus.Error) {
        mainMessage = (
            <FormattedMessage
                id='desktop_auth_token.error.somethingWentWrong'
                defaultMessage='Something went wrong'
            />
        );
        subMessage = (
            <FormattedMessage
                id='desktop_auth_token.error.restartFlow'
                defaultMessage={'Click <a>here</a> to try again.'}
                values={{
                    a: (chunks: React.ReactNode) => {
                        return (
                            <a onClick={() => history.push('/')}>
                                {chunks}
                            </a>
                        );
                    },
                }}
            />
        );
        bottomMessage = null;
    }

    return (
        <div className='DesktopAuthToken'>
            <h1 className='DesktopAuthToken__main'>
                {mainMessage}
            </h1>
            <p className={classNames('DesktopAuthToken__sub', {complete: status === DesktopAuthStatus.LoggedIn})}>
                {subMessage}
            </p>
            <div className='DesktopAuthToken__bottom'>
                {showBottomMessage ? bottomMessage : null}
            </div>
        </div>
    );
};

export default DesktopAuthToken;
