import { Position, Handle } from '@xyflow/react';
import type { Node, NodeProps } from '@xyflow/react';
import type { FC } from 'react';
import { KubernetesIcon } from '../KubernetesIcon'

export type ServiceData = {
  // custom data for current node
  name: string
  namespace: string
  serviceType: string;
}


export type ServiceNode = Node<ServiceData, 'service'>

export type ServiceNodeProps = NodeProps<ServiceNode>


export const ServiceNode: FC<ServiceNodeProps> = (props: ServiceNodeProps) => {
  const { name, namespace, serviceType } = props.data || {}
  const type = (serviceType || 'service').toLowerCase()
  const kind = type === 'ingress' ? 'ingress' : 'service'
  const label = kind === 'ingress' ? 'Ingress' : 'Service'
  const color = kind === 'ingress'
    ? {
        badge: 'px-2 py-0.5 rounded-full text-xs bg-emerald-50 text-emerald-600 ring-1 ring-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-300',
        iconWrap: 'h-6 w-6 rounded-md bg-emerald-100 dark:bg-emerald-900/40 flex items-center justify-center',
        icon: 'text-emerald-600 dark:text-emerald-300',
      }
    : {
        badge: 'px-2 py-0.5 rounded-full text-xs bg-sky-50 text-sky-600 ring-1 ring-sky-200 dark:bg-sky-950/40 dark:text-sky-300',
        iconWrap: 'h-6 w-6 rounded-md bg-sky-100 dark:bg-sky-900/40 flex items-center justify-center',
        icon: 'text-sky-600 dark:text-sky-300',
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
