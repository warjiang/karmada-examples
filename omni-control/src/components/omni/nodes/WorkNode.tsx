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
    <div>
      <div>
        {namespace}/{name}(work)
      </div>
      <Handle type="target" position={Position.Top} />
    </div>
  );
}