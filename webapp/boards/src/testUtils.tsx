// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import {IntlProvider, createIntl} from 'react-intl'
import React, {} from 'react'
import {DndProvider} from 'react-dnd'
import {HTML5Backend} from 'react-dnd-html5-backend'
import configureStore, {MockStoreEnhanced} from 'redux-mock-store'
import {Middleware} from 'redux'

import {DragDropContext, Droppable} from 'react-beautiful-dnd'

import userEvent from '@testing-library/user-event'

import {render} from '@testing-library/react'

import defaultMessages from 'i18n/en.json'

import {Block} from './blocks/block'

type SetupOpts = {user?: Parameters<typeof userEvent.setup>[0], render?: Parameters<typeof render>[1]}
export function setup(element: Parameters<typeof render>[0], opts: SetupOpts = {}) {
    return {
        user: userEvent.setup(opts.user),
        ...render(element, opts.render),
    }
}

export const defaultIntl = createIntl({
    locale: 'en',
    defaultLocale: 'en',
    messages: defaultMessages,
})

export const wrapIntl = (children?: React.ReactNode): JSX.Element => <IntlProvider {...defaultIntl}>{children}</IntlProvider>
export const wrapDNDIntl = (children?: React.ReactNode): JSX.Element => {
    return (
        <DndProvider backend={HTML5Backend}>
            {wrapIntl(children)}
        </DndProvider>
    )
}

export const wrapRBDNDContext = (children?: React.ReactNode): JSX.Element => {
    return (
        <DragDropContext onDragEnd={() => {}}>
            {children}
        </DragDropContext>
    )
}

export const wrapRBDNDDroppable = (children?: React.ReactNode): JSX.Element => {
    const draggableComponent = (
        <Droppable droppableId='droppable_id'>
            {(provided) => (
                <div
                    ref={provided.innerRef}
                    {...provided.droppableProps}
                >
                    {children}
                </div>
            )}
        </Droppable>
    )

    return wrapRBDNDContext(draggableComponent)
}

export function mockDOM(): void {
    window.focus = jest.fn()
    document.createRange = () => {
        const range = new Range()
        range.getBoundingClientRect = jest.fn()
        range.getClientRects = () => {
            return {
                item: () => null,
                length: 0,
                [Symbol.iterator]: jest.fn(),
            }
        }

        return range
    }
}
export function mockMatchMedia(result: {matches: boolean}): void {
    // We check if system preference is dark or light theme.
    // This is required to provide it's definition since
    // window.matchMedia doesn't exist in Jest.
    Object.defineProperty(window, 'matchMedia', {
        writable: true,
        value: jest.fn().mockImplementation(() => {
            return result

            // return ({
            //     matches: true,
            // })
        }),
    })
}

export function mockStateStore(middleware: Middleware[], state: unknown): MockStoreEnhanced<unknown, unknown> {
    const mockStore = configureStore(middleware)

    return mockStore(state)
}

export type BlocksById<BlockType> = {[key: string]: BlockType}

export function blocksById<BlockType extends Block>(blocks: BlockType[]): BlocksById<BlockType> {
    return blocks.reduce((res, block) => {
        res[block.id] = block

        return res
    }, {} as BlocksById<BlockType>)
}
