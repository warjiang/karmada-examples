import { BaseEdge, getStraightPath, getBezierPath, EdgeLabelRenderer } from '@xyflow/react';
import type { Edge, EdgeProps } from '@xyflow/react';

export type PropagationPolicyEdgeData = {
    label?: string;
    style?: 'straight' | 'bezier';
};

export type PropagationPolicyEdgeType = Edge<PropagationPolicyEdgeData, 'propagation-policy'>;

export function PropagationPolicyEdge(props: EdgeProps<PropagationPolicyEdgeType>) {
    const { id, sourceX, sourceY, targetX, targetY } = props;
    const { label, style } = props.data || {};
    const useBezier = style === 'bezier';
    const [edgePath, labelX, labelY] = (useBezier ? getBezierPath : getStraightPath)({
        sourceX,
        sourceY,
        targetX,
        targetY,
    });

    return <>
        <BaseEdge id={id} path={edgePath} style={{ stroke: '#3b82f6', strokeWidth: 2 }} />
        <EdgeLabelRenderer>
            <div
                className="absolute pointer-events-auto select-none px-2 py-1 rounded-full bg-blue-600/80 text-white text-xs shadow ring-1 ring-white/40"
                style={{
                    transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
                }}
            >
                {label || '-'}
            </div>
        </EdgeLabelRenderer>
    </>
}
