import { BaseEdge, getStraightPath, EdgeLabelRenderer } from '@xyflow/react';
import type { Edge, EdgeProps } from '@xyflow/react';

export type PropagationPolicyEdgeData = {
    label?: string;
};

export type PropagationPolicyEdgeType = Edge<PropagationPolicyEdgeData, 'propagation-policy'>;

export function PropagationPolicyEdge(props: EdgeProps<PropagationPolicyEdgeType>) {
    const { id, sourceX, sourceY, targetX, targetY } = props;
    const { label } = props.data || {};
    const [edgePath, labelX, labelY] = getStraightPath({
        sourceX,
        sourceY,
        targetX,
        targetY,
    });

    return <>
        <BaseEdge id={id} path={edgePath} />
        <EdgeLabelRenderer>
            <div style={{
                position: 'absolute',
                transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
                pointerEvents: 'all',
            }}>
                {label || '-'}
            </div>
        </EdgeLabelRenderer>
    </>
}