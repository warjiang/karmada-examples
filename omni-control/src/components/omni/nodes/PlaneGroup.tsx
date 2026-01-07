import { NodeResizer } from '@xyflow/react';
import type { Node, NodeProps } from '@xyflow/react';
import type { FC } from 'react';

export type PlaneGroupData = {
  label?: string;
};

export type PlaneGroupNode = Node<PlaneGroupData, 'planeGroupNode'>;
export type PlaneGroupNodeProps = NodeProps<PlaneGroupNode>;

export const PlaneGroupNode: FC<PlaneGroupNodeProps> = ({ data, selected }) => {
  return (
    <div
      className={
        'relative w-full h-full rounded-lg border bg-slate-50 shadow-sm ' +
        (selected ? 'ring-2 ring-blue-300 border-blue-300' : 'ring-1 ring-slate-200')
      }
    >
      <NodeResizer isVisible={selected} minWidth={260} minHeight={160} />

      <div className="absolute left-0 right-0 top-0 px-3 py-2">
        <div className="font-semibold text-sm text-slate-900">{data?.label ?? 'Group'}</div>
      </div>
    </div>
  );
};