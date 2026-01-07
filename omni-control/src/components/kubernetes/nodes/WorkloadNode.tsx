import { Position, Handle } from '@xyflow/react';
import type { Node, NodeProps } from '@xyflow/react';
import type { FC } from 'react';
import { KubernetesIcon } from '../KubernetesIcon'

export type WorkloadData = {
  // custom data for current node
  name: string
  namespace: string
  workloadType: string;
}


export type WorkloadNode = Node<WorkloadData, 'workload'>

export type WorkloadNodeProps = NodeProps<WorkloadNode>


export const WorkloadNode: FC<WorkloadNodeProps> = (props: WorkloadNodeProps) => {
  const { name, namespace, workloadType } = props.data || {}
  const type = (workloadType || 'deployment').toLowerCase()
  const kind = ['statefulset', 'daemonset', 'cronjob', 'job'].includes(type) ? (type as any) : 'deployment'
  const label = kind === 'statefulset'
    ? 'StatefulSet'
    : kind === 'daemonset'
    ? 'DaemonSet'
    : kind === 'cronjob'
    ? 'CronJob'
    : kind === 'job'
    ? 'Job'
    : 'Deployment'
  console.log(kind, label, type)
  const color = kind === 'statefulset'
    ? {
        badge: 'px-2 py-0.5 rounded-full text-xs bg-amber-50 text-amber-700 ring-1 ring-amber-200 dark:bg-amber-950/40 dark:text-amber-300',
        iconWrap: 'h-6 w-6 rounded-md bg-amber-100 dark:bg-amber-900/40 flex items-center justify-center',
        icon: 'text-amber-700 dark:text-amber-300',
      }
    : kind === 'daemonset'
    ? {
        badge: 'px-2 py-0.5 rounded-full text-xs bg-cyan-50 text-cyan-700 ring-1 ring-cyan-200 dark:bg-cyan-950/40 dark:text-cyan-300',
        iconWrap: 'h-6 w-6 rounded-md bg-cyan-100 dark:bg-cyan-900/40 flex items-center justify-center',
        icon: 'text-cyan-700 dark:text-cyan-300',
      }
    : kind === 'cronjob'
    ? {
        badge: 'px-2 py-0.5 rounded-full text-xs bg-slate-100 text-slate-700 ring-1 ring-slate-300 dark:bg-slate-900/40 dark:text-slate-300',
        iconWrap: 'h-6 w-6 rounded-md bg-slate-200 dark:bg-slate-800 flex items-center justify-center',
        icon: 'text-slate-700 dark:text-slate-300',
      }
    : kind === 'job'
    ? {
        badge: 'px-2 py-0.5 rounded-full text-xs bg-blue-50 text-blue-700 ring-1 ring-blue-200 dark:bg-blue-950/40 dark:text-blue-300',
        iconWrap: 'h-6 w-6 rounded-md bg-blue-100 dark:bg-blue-900/40 flex items-center justify-center',
        icon: 'text-blue-700 dark:text-blue-300',
      }
    : {
        badge: 'px-2 py-0.5 rounded-full text-xs bg-indigo-50 text-indigo-700 ring-1 ring-indigo-200 dark:bg-indigo-950/40 dark:text-indigo-300',
        iconWrap: 'h-6 w-6 rounded-md bg-indigo-100 dark:bg-indigo-900/40 flex items-center justify-center',
        icon: 'text-indigo-700 dark:text-indigo-300',
      }
  return (
    <div className="rounded-xl border ring-1 ring-slate-200 bg-white/90 dark:bg-slate-900 shadow-sm px-3 py-2 min-w-[240px] transition-all hover:shadow-md hover:ring-slate-300">
      <div className="flex items-center gap-2 text-sm font-medium text-slate-900 dark:text-slate-100">
        <div className={color.iconWrap}>
          <KubernetesIcon kind={kind as any} className={color.icon} size={16} />
        </div>
        <span className="truncate">{namespace}/{name}</span>
        <span className={color.badge}>{label}</span>
      </div>
      <Handle type="target" position={Position.Top} className="w-3 h-3 rounded-full border-2 border-white bg-violet-500 shadow ring-2 ring-white/50" />
      <Handle type="source" position={Position.Bottom} className="w-3 h-3 rounded-full border-2 border-white bg-cyan-500 shadow ring-2 ring-white/50" />
    </div>
  );
}
