import './App.css'
import {containerTerminal} from '@/utils/terminal.ts';
import {useEffect, useRef} from 'react';

function App() {
    const terminal = containerTerminal;
    const terminalContainerRef = useRef<HTMLDivElement | null>(null);
    const initRef = useRef(false);
    useEffect(() => {
        if(initRef.current) return;


        if (!terminalContainerRef.current) return;
        console.log('init terminal');
        console.log(terminal)
        // eslint-disable-next-line @typescript-eslint/no-floating-promises
        terminal.getSessionId().then(() => {
            terminal.open(terminalContainerRef.current!);
            terminal.connect();
            initRef.current = true;
        });
    }, [terminalContainerRef.current]);
    return (
        <>
            <div ref={terminalContainerRef} style={{ height: '500px', width: "800px" }}></div>
        </>
    )
}

export default App
