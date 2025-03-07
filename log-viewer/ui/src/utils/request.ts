import axios from 'axios';

export const request = axios.create({
    baseURL: "/",
});

export interface IResponse<Data = {}> {
    code: number;
    data: Data;
}

export interface LogInfo {
    podName: string;
    containerName: string;
    initContainerName: string;
    fromDate: string;
    toDate: string;
    truncated: boolean;
}

export interface LogLineId {
    timestamp: string;
    lineNum: number;
}
export interface Selection {
    referencePoint: LogLineId;
    offsetFrom: number;
    offsetTo: number;
    logFilePosition: string;
}

export interface LogLine {
    timestamp: string;
    content: string;
}

export interface LogDetails {
    info: LogInfo;
    selection: Selection;
    logs: LogLine[];
}