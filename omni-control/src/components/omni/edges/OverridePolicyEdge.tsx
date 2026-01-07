import { getStraightPath, getBezierPath, EdgeLabelRenderer } from '@xyflow/react';
import type { Edge, EdgeProps } from '@xyflow/react';

export type OverridePolicyEdgeData = {
  label?: string;
  style?: 'straight' | 'bezier';
  color?: string;
  width?: number;
  dash?: boolean;
};

export type OverridePolicyEdgeType = Edge<OverridePolicyEdgeData, 'override-policy'>;

export function OverridePolicyEdge(props: EdgeProps<OverridePolicyEdgeType>) {
  const { id, sourceX, sourceY, targetX, targetY, data } = props;
  const { label, style } = props.data || {};
  const useBezier = style === 'bezier';
  const [edgePath, labelX, labelY] = (useBezier ? getBezierPath : getStraightPath)({
    sourceX,
    sourceY,
    targetX,
    targetY,
  });

  const stroke = data?.color || '#f59e0b';
  const strokeWidth = data?.width || 2;

  return <>
    <defs>
      <marker id={`${id}-arrow`} markerWidth="12" markerHeight="8" refX="8" refY="4" orient="auto">
        <path d="M0,0 L8,4 L0,8 Z" fill={stroke} />
      </marker>
    </defs>
    <path
      d={edgePath}
      style={{ stroke: stroke, strokeWidth: strokeWidth, fill: 'none', strokeDasharray: data?.dash ? '6 6' : undefined }}
      markerEnd={`url(#${id}-arrow)`}
    />
    {
      label ?
        <EdgeLabelRenderer>
          <div
            className="absolute pointer-events-auto"
            style={{
              transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            }}
          >
            <div className="flex items-center gap-2 px-2 py-1 rounded-md bg-white/90 dark:bg-slate-900 ring-1 ring-slate-200 shadow-sm text-xs">
              <span className="text-slate-900 dark:text-slate-100 font-medium">{label}</span>
              <span className="px-2 py-0.5 rounded-full bg-orange-50 text-orange-600 ring-1 ring-orange-200 dark:bg-orange-950/40 dark:text-orange-300">override policy</span>
            </div>
          </div>
        </EdgeLabelRenderer> :
        null
    }

  </>
}
