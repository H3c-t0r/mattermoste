// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useState} from 'react';
import {useSelector, useDispatch} from 'react-redux';
import {FormattedMessage} from 'react-intl';
import styled from 'styled-components';

import Constants, {RHSStates} from 'utils/constants';
import * as Keyboard from 'utils/keyboard';
import {isDesktopApp, getDesktopVersion, isMacApp} from 'utils/user_agent';
import {isServerVersionGreaterThanOrEqualTo} from 'utils/server_version';
import {getCurrentChannelNameForSearchShortcut} from 'mattermost-redux/selectors/entities/channels';
import Popover from 'components/widgets/popover';
import {
    updateSearchTerms,
    updateSearchTermsForShortcut,
    showSearchResults,
    showChannelFiles,
    showMentions,
    showFlaggedPosts,
    closeRightHandSide,
    updateRhsState,
    setRhsExpanded,
    openRHSSearch,
    filterFilesSearchByExt,
    updateSearchType,
} from 'actions/views/rhs';

import SearchBox from './search_box';

type Props = {
    enableFindShortcut: boolean;
}

const PopoverStyled = styled(Popover)`
    min-width: 600px;
    left: -90px;
    top: -10px;

    .popover-content {
        padding: 0px;
    }
`

const NewSearchContainer = styled.div`
    display: flex;
    position: relative;
    align-items: center;
    height: 28px;
    width: 100%;
    background-color: rgba(var(--sidebar-text-rgb), 0.08);
    color: rgba(var(--sidebar-text-rgb), 0.56);
    font-size: 10px;
    font-weight: 600;
    border-radius: 8px;
    padding: 4px;
    cursor: text;
`

const NewSearch = ({enableFindShortcut}: Props): JSX.Element => {
    const currentChannelName = useSelector(getCurrentChannelNameForSearchShortcut);
    const dispatch = useDispatch();
    const [focused, setFocused] = useState<boolean>(false);

    useEffect(() => {
        if (!enableFindShortcut) {
            return undefined;
        }

        const isDesktop = isDesktopApp() && isServerVersionGreaterThanOrEqualTo(getDesktopVersion(), '4.7.0');

        const handleKeyDown = (e: KeyboardEvent) => {
            if (Keyboard.cmdOrCtrlPressed(e) && Keyboard.isKeyPressed(e, Constants.KeyCodes.F)) {
                if (!isDesktop && !e.shiftKey) {
                    return;
                }

                // Special case for Mac Desktop xApp where Ctrl+Cmd+F triggers full screen view
                if (isMacApp() && e.ctrlKey) {
                    return;
                }

                e.preventDefault();
                if (currentChannelName) {
                    dispatch(updateSearchTermsForShortcut());
                }
                setFocused(true)
            }
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, [currentChannelName]);

    return (
        <NewSearchContainer onClick={() => setFocused(true)}>
            <i className='icon icon-magnify'/>
            <FormattedMessage
                id='search_bar.search'
                defaultMessage='Search'
            />
            {focused && (
                <PopoverStyled placement='bottom'>
                    <SearchBox
                        onClose={() => setFocused(false)}
                        onSearch={(searchType: string, searchTerms: string) => {
                            dispatch(updateSearchType(searchType))
                            dispatch(updateSearchTerms(searchTerms))
                            dispatch(showSearchResults(false))
                            setFocused(false)
                        }}
                    />
                </PopoverStyled>
            )}
        </NewSearchContainer>
    );
};

export default NewSearch;
