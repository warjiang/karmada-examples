import {useEffect, useState} from 'react';
import {useXTerm} from "@/xterm.tsx";
import {FitAddon} from '@xterm/addon-fit'
import {LogDetails, request} from "@/utils/request.ts";
import {useRequest} from 'ahooks';

const fitAddon = new FitAddon();
const addons = [
    fitAddon
]

const LineSize = 10;
export default function App() {
    const {instance, ref} = useXTerm({
        // options: {},
        addons
    })
    useEffect(() => {
        if (!instance) return;
        fitAddon.fit();
    }, [instance]);
    const [logOptions, setLogOptions] = useState({
        referenceTimestamp: 'oldest',
        logFilePosition: 'beginning',
        offsetFrom: 0,
        offsetTo: LineSize,
    });
    useRequest(async () => {
        if (!instance?.writeln) return
        const namespace = import.meta.env.VITE_NAMESPACE_NAME;
        const pod = import.meta.env.VITE_POD_NAME;
        const container = import.meta.env.VITE_CONTAINER_NAME;
        const resp = await request.get<{ data: LogDetails }>(`log/${namespace}/${pod}/${container}`, {
            params: logOptions
        })
        const data = resp.data.data;
        data.logs.forEach(log => {
            instance?.writeln(log.content)
        })
        return resp.data
    }, {
        refreshDeps: [instance, logOptions],
    });
    return (
        <div style={{height: '100vh', width: '100%'}}>
            <div ref={ref} style={{width: '1000px', height: '220px', margin: '20px auto 0 auto'}}/>
            <div style={{width: '1000px', display: 'flex', flexDirection: 'row', margin: '5px auto 0 auto'}}>
                {/*<button>prev</button>*/}
                <button onClick={() => {
                    const {offsetTo} = logOptions;
                    const _offsetFrom = offsetTo;
                    const _offsetTo = _offsetFrom + LineSize ;
                    setLogOptions({
                        ...logOptions,
                        offsetFrom: _offsetFrom,
                        offsetTo: _offsetTo,
                    })
                }}>next
                </button>
            </div>
        </div>
    )
}
