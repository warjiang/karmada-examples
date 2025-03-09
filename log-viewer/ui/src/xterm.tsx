import {type ITerminalAddon, type ITerminalInitOnlyOptions, type ITerminalOptions, ITheme, Terminal} from '@xterm/xterm'
import '@xterm/xterm/css/xterm.css'
import { type ComponentPropsWithoutRef, useEffect, useRef, useState } from 'react'

export interface UseXTermProps {
    addons?: ITerminalAddon[]
    options?: ITerminalOptions & ITerminalInitOnlyOptions
    listeners?: {
        onBinary?(data: string): void
        onCursorMove?(): void
        onData?(data: string): void
        onKey?: (event: { key: string; domEvent: KeyboardEvent }) => void
        onLineFeed?(): void
        onScroll?(newPosition: number): void
        onSelectionChange?(): void
        onRender?(event: { start: number; end: number }): void
        onResize?(event: { cols: number; rows: number }): void
        onTitleChange?(newTitle: string): void
        customKeyEventHandler?(event: KeyboardEvent): boolean
    }
}

export function useXTerm({ options, addons, listeners }: UseXTermProps = {}) {
    const terminalRef = useRef<HTMLDivElement>(null)
    const [terminalInstance, setTerminalInstance] = useState<Terminal | null>(null)

    useEffect(() => {
        const instance = new Terminal({
            fontFamily: 'operator mono,SFMono-Regular,Consolas,Liberation Mono,Menlo,monospace',
            fontSize: 14,
            // theme: {
            //     background: '#101420',
            // },
            // cursorStyle: 'underline',
            // cursorBlink: false,
            convertEol: true,
            theme: {
                foreground: '#d2d2d2',
                background: '#2b2b2b',
                cursor: '#adadad',
                black: '#000000',
                red: '#d81e00',
                green: '#5ea702',
                yellow: '#cfae00',
                blue: '#427ab3',
                magenta: '#89658e',
                cyan: '#00a7aa',
                white: '#dbded8',
                brightBlack: '#686a66',
                brightRed: '#f54235',
                brightGreen: '#99e343',
                brightYellow: '#fdeb61',
                brightBlue: '#84b0d8',
                brightMagenta: '#bc94b7',
                brightCyan: '#37e6e8',
                brightWhite: '#f1f1f0',
            } as ITheme,
            ...options,
        })

        // Load addons if the prop exists
        if (addons) {
            addons.forEach((addon) => {
                instance.loadAddon(addon)
            })
        }

        // Listeners
        if (listeners) {
            if (listeners.onBinary) instance.onBinary(listeners.onBinary)
            if (listeners.onCursorMove) instance.onCursorMove(listeners.onCursorMove)
            if (listeners.onLineFeed) instance.onLineFeed(listeners.onLineFeed)
            if (listeners.onScroll) instance.onScroll(listeners.onScroll)
            if (listeners.onSelectionChange) instance.onSelectionChange(listeners.onSelectionChange)
            if (listeners.onRender) instance.onRender(listeners.onRender)
            if (listeners.onResize) instance.onResize(listeners.onResize)
            if (listeners.onTitleChange) instance.onTitleChange(listeners.onTitleChange)
            if (listeners.onKey) instance.onKey(listeners.onKey)
            if (listeners.onData) instance.onData(listeners.onData)

            // Add Custom Key Event Handler
            if (listeners.customKeyEventHandler) {
                instance.attachCustomKeyEventHandler(listeners.customKeyEventHandler)
            }
        }

        if (terminalRef.current) {
            // Mount terminal
            instance.open(terminalRef.current)
            instance.focus()
        }

        setTerminalInstance(instance)

        return () => {
            instance.dispose()
            setTerminalInstance(null)
        }
    }, [
        terminalRef,
        options,
        addons,
        listeners,
        listeners?.onBinary,
        listeners?.onCursorMove,
        listeners?.onData,
        listeners?.onKey,
        listeners?.onLineFeed,
        listeners?.onScroll,
        listeners?.onSelectionChange,
        listeners?.onRender,
        listeners?.onResize,
        listeners?.onTitleChange,
        listeners?.customKeyEventHandler,
    ])

    return {
        ref: terminalRef,
        instance: terminalInstance,
    }
}

export interface XTermProps extends Omit<ComponentPropsWithoutRef<'div'>, 'onResize' | 'onScroll'>, UseXTermProps {}

export function XTerm({ className = '', options, addons, listeners, ...props }: XTermProps) {
    const { ref } = useXTerm({
        options,
        addons,
        listeners,
    })

    return <div className={className} ref={ref} {...props} />
}

export default XTerm