import {
  ReactFlow,
  Background,
  Controls,
  applyNodeChanges,
  applyEdgeChanges,
  type OnNodesChange,
  type OnEdgesChange,
  addEdge,
  MiniMap,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { useState, useCallback } from 'react';
import type { Edge, Node, OnConnect } from '@xyflow/react';
import { 
  ResourceTemplateNode, ResourceBindingNode, WorkNode ,
  PropagationPolicyEdge, OverridePolicyEdge,
} from './components/omni';


// NodeType extends Node = Node, EdgeType extends Edge = Edge

const initialNodes: OmniNode[] = [
  {
    id: 'node1',
    position: { x: 0, y: 0 },
    data: { 
      name: 'nginx-deployment',
      namespace: 'example',
    },
    type: 'resourceTemplateNode',
  },
  {
    id: 'node2',
    position: { x: 0, y: 200 },
    data: { 
      name: 'nginx-deployment-deployment',
      namespace: 'example',
    },
    type: 'resourceBindingNode',
  },
  {
    id: 'node3',
    position: { x: 0, y: 400 },
    data: { 
      name: 'nginx-deployment-697bcd7948',
      namespace: 'karmada-es-member1',
    },
    type: 'workNode',
  },
];
const initialEdges: OmniEdge[] = [
  {
    id: 'node1-node2',
    source: 'node1',
    target: 'node2',
    type: 'propagationPolicy',
    data: {
      label: 'pp-workload-deployment',
      style: 'bezier',
    },
  },
  {
    id: 'node2-node3',
    source: 'node2',
    target: 'node3',
    type: 'overridePolicy',
    data: {
      style: 'bezier',
    },
  },
];

type OmniNodeData = {
}
type OmniNode = Node<OmniNodeData, string>;
type OmniEdge = Edge

const nodeColor = (node:OmniNode) => {
  switch (node.type) {
    case 'input':
      return '#6ede87';
    case 'output':
      return '#6865A5';
    default:
      return '#ff0072';
  }
};


const nodeTypes = {
  resourceTemplateNode: ResourceTemplateNode,
  resourceBindingNode: ResourceBindingNode,
  workNode: WorkNode,
};
const edgeTypes = {
  propagationPolicy: PropagationPolicyEdge,
  overridePolicy: OverridePolicyEdge,
};

export default function App() {
  const [nodes, setNodes] = useState<OmniNode[]>(initialNodes);
  const [edges, setEdges] = useState(initialEdges);
  const onNodesChange: OnNodesChange<OmniNode> = useCallback(
    (changes) => setNodes((nodesSnapshot) => applyNodeChanges(changes, nodesSnapshot)),
    [],
  );
  const onEdgesChange: OnEdgesChange<OmniEdge> = useCallback(
    (changes) => {
      console.log("edge change", changes)
      setEdges((edgesSnapshot) => applyEdgeChanges(changes, edgesSnapshot));
    },
    [],
  );
  const onConnect: OnConnect = useCallback(
    (params) => setEdges((edgesSnapshot) => addEdge(params, edgesSnapshot)),
    [],
  );
  return (
    <div className="w-screen h-screen flex">
      <div className="flex-1">
        <ReactFlow<OmniNode, OmniEdge>
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          fitView
          nodeTypes={nodeTypes}
          edgeTypes={edgeTypes}
        >
          <Background />
          <Controls />
          <MiniMap nodeColor={nodeColor} nodeStrokeWidth={3} zoomable pannable />
        </ReactFlow>
      </div>
    </div>
  );
}
