import { Position, Handle } from '@xyflow/react';
import type { Node, NodeProps } from '@xyflow/react';
import type { FC } from 'react';

export type ResourceBindingData = {
  // custom data for current node
  name: string
  namespace: string
}


export type ResourceBindingNode = Node<ResourceBindingData, 'resource-binding'>

export type ResourceBindingNodeProps = NodeProps<ResourceBindingNode>


export const ResourceBindingNode: FC<ResourceBindingNodeProps> = (props: ResourceBindingNodeProps) => {
  const { name, namespace } = props.data || {}
  return (
    <div className="rounded-lg border ring-1 ring-slate-200 bg-white/90 dark:bg-slate-900 shadow-sm px-3 py-2 min-w-[220px]">
      <div className="flex items-center gap-2 text-sm font-medium text-slate-900 dark:text-slate-100">
        <span className="truncate">{namespace}/{name}</span>
        <span className="px-2 py-0.5 rounded-full text-xs bg-violet-50 text-violet-600 ring-1 ring-violet-200 dark:bg-violet-950/40 dark:text-violet-300">resource-binding</span>
      </div>
      <Handle type="target" position={Position.Top} className="w-3 h-3 rounded-full border-2 border-white bg-violet-500 shadow ring-2 ring-white/50" />
      <Handle type="source" position={Position.Bottom} className="w-3 h-3 rounded-full border-2 border-white bg-cyan-500 shadow ring-2 ring-white/50" />
    </div>
  );
}
