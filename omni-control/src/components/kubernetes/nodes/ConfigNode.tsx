import { Position, Handle } from '@xyflow/react';
import type { Node, NodeProps } from '@xyflow/react';
import type { FC } from 'react';
import { KubernetesIcon } from '../KubernetesIcon'

export type ConfigData = {
  // custom data for current node
  name: string
  namespace: string
  configType: string;
}


export type ConfigNode = Node<ConfigData, 'config'>

export type ConfigNodeProps = NodeProps<ConfigNode> 


export const ConfigNode: FC<ConfigNodeProps> = (props: ConfigNodeProps) => {
  const { name, namespace, configType } = props.data || {}
  const type = (configType || 'configmap').toLowerCase()
  const kind = type === 'secret' ? 'secret' : 'configmap'
  const label = kind === 'secret' ? 'Secret' : 'ConfigMap'
  const color = kind === 'secret'
    ? {
        badge: 'px-2 py-0.5 rounded-full text-xs bg-rose-50 text-rose-600 ring-1 ring-rose-200 dark:bg-rose-950/40 dark:text-rose-300',
        iconWrap: 'h-6 w-6 rounded-md bg-rose-100 dark:bg-rose-900/40 flex items-center justify-center',
        icon: 'text-rose-600 dark:text-rose-300',
      }
    : {
        badge: 'px-2 py-0.5 rounded-full text-xs bg-violet-50 text-violet-600 ring-1 ring-violet-200 dark:bg-violet-950/40 dark:text-violet-300',
        iconWrap: 'h-6 w-6 rounded-md bg-violet-100 dark:bg-violet-900/40 flex items-center justify-center',
        icon: 'text-violet-600 dark:text-violet-300',
      }
  return (
    <div className="rounded-xl border ring-1 ring-slate-200 bg-white/90 dark:bg-slate-900 shadow-sm px-3 py-2 min-w-[240px] transition-all hover:shadow-md hover:ring-slate-300">
      <div className="flex items-center gap-2 text-sm font-medium text-slate-900 dark:text-slate-100">
        <div className={color.iconWrap}>
          <KubernetesIcon kind={kind} className={color.icon} size={16} />
        </div>
        <span className="truncate">{namespace}/{name}</span>
        <span className={color.badge}>{label}</span>
      </div>
      <Handle type="target" position={Position.Top} className="w-3 h-3 rounded-full border-2 border-white bg-violet-500 shadow ring-2 ring-white/50" />
      <Handle type="source" position={Position.Bottom} className="w-3 h-3 rounded-full border-2 border-white bg-cyan-500 shadow ring-2 ring-white/50" />
    </div>
  );
}
