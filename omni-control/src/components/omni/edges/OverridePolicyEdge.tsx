import { BaseEdge, getStraightPath, getBezierPath } from '@xyflow/react';
import type { Edge, EdgeProps } from '@xyflow/react';

export type OverridePolicyEdgeData = {
  label?: string;
  style?: 'straight' | 'bezier';
};

export type OverridePolicyEdgeType = Edge<OverridePolicyEdgeData, 'override-policy'>;

export function OverridePolicyEdge(props: EdgeProps<OverridePolicyEdgeType>) {
  const { id, sourceX, sourceY, targetX, targetY, data } = props;
  const useBezier = data?.style === 'bezier';

  const [edgePath] = (useBezier ? getBezierPath : getStraightPath)({
    sourceX,
    sourceY,
    targetX,
    targetY,
  });

  return <BaseEdge id={id} path={edgePath} style={{ stroke: '#f59e0b', strokeWidth: 2 }} />;
}
