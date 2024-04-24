// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {
    useFloating,
    autoUpdate,
    autoPlacement,
    useTransitionStyles,
    useClick,
    useDismiss,
    useInteractions,
    useRole,
    FloatingFocusManager,
    FloatingOverlay,
    FloatingPortal,
} from '@floating-ui/react';
import classNames from 'classnames';
import type {HtmlHTMLAttributes, ReactNode} from 'react';
import React, {useState} from 'react';

import type {Channel} from '@mattermost/types/channels';
import type {UserProfile} from '@mattermost/types/users';

import ProfilePopover from 'components/profile_popover';

import {A11yClassNames} from 'utils/constants';

const PROFILE_POPOVER_OPENING_DELAY = 300;
const PROFILE_POPOVER_CLOSING_DELAY = 500;

interface Props<TriggerComponentType> {
    triggerComponentRootProps?: {
        as?: React.ElementType;
        id?: HtmlHTMLAttributes<TriggerComponentType>['id'];
        className?: HtmlHTMLAttributes<TriggerComponentType>['className'];
        style?: HtmlHTMLAttributes<TriggerComponentType>['style'];
    };

    /**
     * Source URL from the image to display in the popover
     */
    src: string;

    /**
     * This should be the trigger button for the popover, Do note that the root element of the trigger component should be passed in triggerComponentRoot
     */
    children: ReactNode;
    userId: UserProfile['id'];
    channelId?: Channel['id'];

    /**
     * The overwritten username that should be shown at the top of the popover
     */
    overwriteName?: string;

    /**
     * Source URL from the image that should override default image
     */
    overwriteIcon?: string;

    /**
     * Set to true of the popover was opened from a webhook post
     */
    fromWebhook?: boolean;
    hideStatus?: boolean;

    /**
     * Function to call to return focus to the previously focused element when the popover closes.
     * If not provided, the popover will automatically determine the previously focused element
     * and focus that on close. However, if the previously focused element is not correctly detected
     * by the popover, or the previously focused element will disappear after the popover opens,
     * it is necessary to provide this function to focus the correct element.
     */
    returnFocus?: () => void;
}

function ProfilePopoverController<TriggerComponentType = HTMLSpanElement>(props: Props<TriggerComponentType>) {
    const [isOpen, setOpen] = useState(false);

    const {refs, floatingStyles, context: floatingContext} = useFloating({
        open: isOpen,
        onOpenChange: setOpen,
        whileElementsMounted: autoUpdate,
        middleware: [autoPlacement()],
    });

    const {isMounted, styles: transitionStyles} = useTransitionStyles(floatingContext, {
        duration: {
            open: PROFILE_POPOVER_OPENING_DELAY,
            close: PROFILE_POPOVER_CLOSING_DELAY,
        },
    });

    const combinedFloatingStyles = Object.assign({}, floatingStyles, transitionStyles);

    const clickInteractions = useClick(floatingContext);

    const dismissInteraction = useDismiss(floatingContext);

    const role = useRole(floatingContext);

    const {getReferenceProps, getFloatingProps} = useInteractions([
        clickInteractions,
        dismissInteraction,
        role,
    ]);

    function handleHide() {
        setOpen(false);
    }

    const TriggerComponent = props.triggerComponentRootProps?.as ?? 'span';

    return (
        <>
            <TriggerComponent
                id={props.triggerComponentRootProps?.id}
                ref={refs.setReference}
                className={props.triggerComponentRootProps?.className}
                style={props.triggerComponentRootProps?.style}
                {...getReferenceProps()}
            >
                {props.children}
            </TriggerComponent>

            {isMounted && (
                <FloatingPortal id='user-profile-popover-portal'>
                    <FloatingOverlay
                        id='user-profile-popover-floating-overlay'
                        className='user-profile-popover-floating-overlay'
                        lockScroll={true}
                    >
                        <FloatingFocusManager context={floatingContext}>
                            <div
                                ref={refs.setFloating}
                                style={combinedFloatingStyles}
                                className={classNames('user-profile-popover', A11yClassNames.POPUP)}
                                {...getFloatingProps()}
                                id='user-profile-popover' // This is placed here as id is being overwritten by getFloatingProps
                            >
                                <ProfilePopover
                                    userId={props.userId}
                                    src={props.src}
                                    channelId={props.channelId}
                                    hideStatus={props.hideStatus}
                                    fromWebhook={props.fromWebhook}
                                    hide={handleHide}
                                    returnFocus={props.returnFocus}
                                    overwriteIcon={props.overwriteIcon}
                                    overwriteName={props.overwriteName}
                                />
                            </div>
                        </FloatingFocusManager>
                    </FloatingOverlay>
                </FloatingPortal>
            )}
        </>
    );
}

export default ProfilePopoverController;
