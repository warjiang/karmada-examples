import { BaseEdge, getStraightPath } from '@xyflow/react';
import type { Edge, EdgeProps } from '@xyflow/react';

export type OverridePolicyEdgeData = {
  label?: string;
};

export type OverridePolicyEdgeType = Edge<OverridePolicyEdgeData, 'override-policy'>;

export function OverridePolicyEdge(props: EdgeProps<OverridePolicyEdgeType>) {
  const { id, sourceX, sourceY, targetX, targetY } = props;

  const [edgePath] = getStraightPath({
    sourceX,
    sourceY,
    targetX,
    targetY,
  });

  return <BaseEdge id={id} path={edgePath} />;
}