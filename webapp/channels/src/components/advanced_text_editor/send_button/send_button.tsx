// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {memo, useCallback, useMemo} from 'react';
import {defineMessage, useIntl} from 'react-intl';
import {useSelector} from 'react-redux';
import styled from 'styled-components';

import {SendIcon} from '@mattermost/compass-icons/components';
import type {SchedulingInfo} from '@mattermost/types/schedule_post';

import {isSendOnCtrlEnter} from 'selectors/preferences';

import {SendPostOptions} from 'components/advanced_text_editor/send_button/send_post_options';

import './style.scss';
import classNames from 'classnames';

import WithTooltip from 'components/with_tooltip';
import type {ShortcutDefinition} from 'components/with_tooltip/shortcut';
import {ShortcutKeys} from 'components/with_tooltip/shortcut';

type SendButtonProps = {
    handleSubmit: (e: React.FormEvent, schedulingInfo?: SchedulingInfo) => void;
    disabled: boolean;
}

const SendButtonContainer = styled.button`
    cursor: pointer;
    place-content: center;
    place-items: center;
    transition: color 150ms;

    .android &,
    .ios & {
        display: flex;
    }
`;

const SendButton = ({disabled, handleSubmit}: SendButtonProps) => {
    const {formatMessage} = useIntl();

    const sendMessage = useCallback((e: React.FormEvent, schedulingInfo?: SchedulingInfo) => {
        e.stopPropagation();
        e.preventDefault();
        handleSubmit(e, schedulingInfo);
    }, [handleSubmit]);

    const sendOnCtrlEnter = useSelector(isSendOnCtrlEnter);

    const sendNowKeyboardShortcutDescriptor = useMemo<ShortcutDefinition>(() => {
        const shortcutDefinition: ShortcutDefinition = {
            default: [
                defineMessage({
                    id: 'shortcuts.generic.enter',
                    defaultMessage: 'Enter',
                }),
            ],
            mac: [
                defineMessage({
                    id: 'shortcuts.generic.enter',
                    defaultMessage: 'Enter',
                }),
            ],
        };

        if (sendOnCtrlEnter) {
            shortcutDefinition.default.unshift(ShortcutKeys.ctrl);
            shortcutDefinition.mac?.unshift(ShortcutKeys.cmd);
        }

        return shortcutDefinition;
    }, [sendOnCtrlEnter]);

    return (
        <div className={classNames('splitSendButton', {disabled})}>
            <WithTooltip
                placement='top'
                id='send_post_now_tooltip'
                title={formatMessage({id: 'create_post_button.option.send_now', defaultMessage: 'Send Now'})}
                shortcut={sendNowKeyboardShortcutDescriptor}
                disabled={disabled}
            >
                <SendButtonContainer
                    className={classNames('SendMessageButton', {disabled})}
                    data-testid='SendMessageButton'
                    tabIndex={0}
                    aria-label={formatMessage({
                        id: 'create_post.send_message',
                        defaultMessage: 'Send a message',
                    })}
                    disabled={disabled}
                    onClick={sendMessage}
                >
                    <SendIcon
                        size={18}
                        color='currentColor'
                        aria-label={formatMessage({
                            id: 'create_post.icon',
                            defaultMessage: 'Create a post',
                        })}
                    />
                </SendButtonContainer>
            </WithTooltip>

            <SendPostOptions
                disabled={disabled}
                onSelect={sendMessage}
            />
        </div>
    );
};

export default memo(SendButton);
