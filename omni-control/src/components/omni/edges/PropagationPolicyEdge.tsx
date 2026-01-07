import { getStraightPath, getBezierPath, EdgeLabelRenderer } from '@xyflow/react';
import type { Edge, EdgeProps } from '@xyflow/react';

export type PropagationPolicyEdgeData = {
    label?: string;
    style?: 'straight' | 'bezier';
    color?: string;
    width?: number;
    dash?: boolean;
};

export type PropagationPolicyEdgeType = Edge<PropagationPolicyEdgeData, 'propagation-policy'>;

export function PropagationPolicyEdge(props: EdgeProps<PropagationPolicyEdgeType>) {
    const { id, sourceX, sourceY, targetX, targetY } = props;
    const { label, style, color, width, dash } = props.data || {};
    const useBezier = style === 'bezier';
    const [edgePath, labelX, labelY] = (useBezier ? getBezierPath : getStraightPath)({
        sourceX,
        sourceY,
        targetX,
        targetY,
    });

    const stroke = color || '#3b82f6';
    const strokeWidth = width || 2;

    return <>
        <defs>
            <marker id={`${id}-arrow`} markerWidth="12" markerHeight="8" refX="8" refY="4" orient="auto">
                <path d="M0,0 L8,4 L0,8 Z" fill={stroke} />
            </marker>
        </defs>
        <path
            d={edgePath}
            style={{ stroke, strokeWidth, fill: 'none', strokeDasharray: dash ? '6 6' : undefined }}
            markerEnd={`url(#${id}-arrow)`}
        />
        <EdgeLabelRenderer>
            <div
                className="absolute pointer-events-auto"
                style={{
                    transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
                }}
            >
                <div className="flex items-center gap-2 px-2 py-1 rounded-md bg-white/90 dark:bg-slate-900 ring-1 ring-slate-200 shadow-sm text-xs">
                    <span className="text-slate-900 dark:text-slate-100 font-medium">{label || '-'}</span>
                    <span className="px-2 py-0.5 rounded-full bg-blue-50 text-blue-600 ring-1 ring-blue-200 dark:bg-blue-950/40 dark:text-blue-300">propagation policy</span>
                </div>
            </div>
        </EdgeLabelRenderer>
    </>
}
