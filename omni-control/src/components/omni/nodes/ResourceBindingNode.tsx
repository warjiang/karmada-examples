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
    <div>
      <div>
        {namespace}/{name}(rb)
      </div>
      <Handle type="target" position={Position.Top} />
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
}