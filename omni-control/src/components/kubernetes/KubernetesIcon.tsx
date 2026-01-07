import type { FC } from 'react'
import { MdSettings, MdLock, MdLan, MdAltRoute, MdViewModule, MdStorage, MdSync, MdSchedule, MdWork } from 'react-icons/md'

export type KubernetesKind =
  | 'configmap'
  | 'secret'
  | 'service'
  | 'ingress'
  | 'deployment'
  | 'statefulset'
  | 'daemonset'
  | 'cronjob'
  | 'job'

export type KubernetesIconProps = {
  kind: KubernetesKind
  className?: string
  size?: number
}

const kindIconMap: Record<KubernetesKind, FC<{ className?: string; size?: number }>> = {
  configmap: MdSettings,
  secret: MdLock,
  service: MdLan,
  ingress: MdAltRoute,
  deployment: MdViewModule,
  statefulset: MdStorage,
  daemonset: MdSync,
  cronjob: MdSchedule,
  job: MdWork,
}

export const KubernetesIcon: FC<KubernetesIconProps> = ({ kind, className, size = 16 }) => {
  const Icon = kindIconMap[kind]
  return <Icon className={className} size={size} />
}

