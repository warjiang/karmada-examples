import { Position, Handle } from '@xyflow/react';
import type { Node, NodeProps } from '@xyflow/react';
import type { FC } from 'react';

export type WorkData = {
  // custom data for current node
  name: string
  namespace: string
}


export type WorkNode = Node<WorkData, 'work'>

export type WorkNodeProps = NodeProps<WorkNode>



export const WorkNode: FC<WorkNodeProps> = (props: WorkNodeProps) => {
  const { name, namespace } = props.data || {}
  return (
    <div className="rounded-lg border ring-1 ring-slate-200 bg-white/90 dark:bg-slate-900 shadow-sm px-3 py-2 min-w-[220px]">
      <div className="flex items-center gap-2 text-sm font-medium text-slate-900 dark:text-slate-100">
        <span className="truncate">{name}</span>
        <span className="px-2 py-0.5 rounded-full text-xs bg-emerald-50 text-emerald-600 ring-1 ring-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-300">work</span>
      </div>
      <Handle type="target" position={Position.Top} className="w-3 h-3 rounded-full border-2 border-white bg-emerald-500 shadow ring-2 ring-white/50" />
    </div>
  );
}
