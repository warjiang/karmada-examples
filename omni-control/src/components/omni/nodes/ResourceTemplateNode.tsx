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
        <div>
            <div>
                {namespace}/{name}(template)
            </div>
            <Handle type="source" position={Position.Bottom} />
        </div>
    );
}