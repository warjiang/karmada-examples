import { Handle, Position } from '@xyflow/react';
import type { Node, NodeProps } from '@xyflow/react';
import { useCallback } from 'react';
import type {FC} from 'react'


export type ResourceTemplateData = {
    // custom data for current node
    name: string
    namespace: string
}

export type ResourceTemplateNode = Node<ResourceTemplateData, 'resource-template'>

export type ResourceTemplateNodeProps = NodeProps<ResourceTemplateNode>


export const ResourceTemplateNode: FC<ResourceTemplateNodeProps> = (props: ResourceTemplateNodeProps) => {
    const onChange = useCallback((evt: any) => {
        console.log(evt.target.value);
    }, []);
    const { name, namespace } = props.data || {}
    return (
        <div className="rounded-lg border ring-1 ring-slate-200 bg-white/90 dark:bg-slate-900 shadow-sm px-3 py-2 min-w-[220px]">
            <div className="flex items-center gap-2 text-sm font-medium text-slate-900 dark:text-slate-100">
                <span className="truncate">{namespace}/{name}</span>
                <span className="px-2 py-0.5 rounded-full text-xs bg-blue-50 text-blue-600 ring-1 ring-blue-200 dark:bg-blue-950/40 dark:text-blue-300">resource-template</span>
            </div>
            <Handle type="source" position={Position.Bottom} className="w-3 h-3 rounded-full border-2 border-white bg-blue-500 shadow ring-2 ring-white/50" />
        </div>
    );
}
