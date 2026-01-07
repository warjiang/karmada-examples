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
  ReactFlowProvider,
  Panel,
  useReactFlow,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { useState, useCallback } from 'react';
import type { Edge, Node, OnConnect } from '@xyflow/react';
import {
  ResourceTemplateNode, ResourceBindingNode, WorkNode, PlaneGroupNode,
  PropagationPolicyEdge, OverridePolicyEdge,
} from './components/omni';
import Dagre from '@dagrejs/dagre';


// NodeType extends Node = Node, EdgeType extends Edge = Edge

const initialNodes: OmniNode[] = [
  // {
  //   id: 'A',
  //   type: 'planeGroupNode',
  //   position: { x: 0, y: 0 },
  //   style: {
  //     width:800,
  //     height: 460,
  //   },
  //   data: {
  //     label: 'control-plane'
  //   },
  // },

  {
    id: 'node1',
    position: { x: 80, y: 40 },
    data: {
      name: 'nginx-deployment',
      namespace: 'example',
    },
    type: 'resourceTemplateNode',
    // parentId: 'A',
    // extent: 'parent',
  },
  {
    id: 'node2',
    position: { x: 60, y: 200 },
    data: {
      name: 'nginx-deployment-deployment',
      namespace: 'example',
    },
    type: 'resourceBindingNode',
    // parentId: 'A',
    // extent: 'parent',
  },
  {
    id: 'node3',
    position: { x: 40, y: 400 },
    data: {
      name: 'nginx-deployment-697bcd7948',
      namespace: 'karmada-es-member1',
    },
    type: 'workNode',
    // parentId: 'A',
    // extent: 'parent',
  },
  {
    id: 'node4',
    position: { x: 500, y: 400 },
    data: {
      name: 'nginx-deployment-697bcd7950',
      namespace: 'karmada-es-member2',
    },
    type: 'workNode',
    // parentId: 'A',
    // extent: 'parent',
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
      label: 'op1-deployment',
    },
  },
  {
    id: 'node2-node4',
    source: 'node2',
    target: 'node4',
    type: 'overridePolicy',
    data: {
      style: 'bezier',
      label: 'op2-deployment',
    },
  },
];

type OmniNodeData = {
}
type OmniNode = Node<OmniNodeData, string>;
type OmniEdge = Edge

const nodeColor = (node: OmniNode) => {
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
  planeGroupNode: PlaneGroupNode,
};
const edgeTypes = {
  propagationPolicy: PropagationPolicyEdge,
  overridePolicy: OverridePolicyEdge,
};

const getLayoutedElements = (nodes: OmniNode[], edges: OmniEdge[], options: { direction: 'TB' | 'LR' }) => {
  const g = new Dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));
  g.setGraph({ rankdir: options.direction });

  edges.forEach((edge) => g.setEdge(edge.source, edge.target));
  nodes.forEach((node) =>
    g.setNode(node.id, {
      ...node,
      width: node.measured?.width ?? 0,
      height: node.measured?.height ?? 0,
    }),
  );

  Dagre.layout(g);

  return {
    nodes: nodes.map((node: OmniNode) => {
      const position = g.node(node.id);
      // We are shifting the dagre node position (anchor=center center) to the top left
      // so it matches the React Flow node anchor point (top left).
      const x = position.x - (node.measured?.width ?? 0) / 2;
      const y = position.y - (node.measured?.height ?? 0) / 2;

      return { ...node, position: { x, y } };
    }),
    edges,
  };
};


export function LayoutFlow() {
  const { fitView } = useReactFlow();

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
  
  const onLayout = useCallback((direction: 'TB' | 'LR') => {
    console.log(nodes);
    const layouted = getLayoutedElements(nodes, edges, { direction });

    setNodes([...layouted.nodes]);
    setEdges([...layouted.edges]);

    fitView();
  }, [nodes, edges]);
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
            <Panel position="top-right">
              <button onClick={() => onLayout('TB')}>vertical layout</button>
              <button onClick={() => onLayout('LR')}>horizontal layout</button>
            </Panel>
          </ReactFlow>
      </div>
    </div>
  );
}


export default function App() {
  return (
    <ReactFlowProvider>
      <LayoutFlow />
    </ReactFlowProvider>
  );
}