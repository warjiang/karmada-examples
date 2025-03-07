import './App.css'
import {useXTerm} from "@/xterm.tsx";
import {useEffect} from "react";
import {LogDetails, request} from "@/utils/request.ts";

function App() {
    const {ref} = useXTerm()

    useEffect(() => {
        (async () => {
            const namespace = import.meta.env.VITE_NAMESPACE_NAME;
            const pod = import.meta.env.VITE_POD_NAME;
            const container = import.meta.env.VITE_CONTAINER_NAME;
            const resp = await request.get<LogDetails>(`log/${namespace}/${pod}/${container}`)
            console.log(resp)
        })()
    }, []);
    /*
    console.log({
        namespace: import.meta.env.VITE_NAMESPACE_NAME,
        pod: import.meta.env.VITE_POD_NAME,
        container: import.meta.env.VITE_CONTAINER_NAME,
    })
    */
    /*
    instance?.writeln('Hello from react-xtermjs!')
    instance?.onData((data) => instance?.write(data))
    */
    return <div ref={ref} style={{width: '800', height: '600'}}/>

}

export default App
