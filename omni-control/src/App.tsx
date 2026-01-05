// App.jsx
import { ReactFlow, Background, Controls } from '@xyflow/react';
import '@xyflow/react/dist/style.css';

const initialNodes = [
  // 顶部节点
  {
    id: 'resource-template',
    type: 'default',
    data: { label: 'Resource Template\n(Deployment / Service / ConfigMap / Secret ... )\nExactly same with K8s APIs' },
    position: { x: 200, y: 50 },
    style: { background: '#cce5ff', borderRadius: '8px', padding: '10px', width: '300px', textAlign: 'center' },
  },
  {
    id: 'propagation-policy',
    type: 'default',
    data: { label: 'Propagation Policy' },
    position: { x: 500, y: 100 },
    style: { background: '#cce5ff', borderRadius: '8px', padding: '10px', width: '150px', textAlign: 'center' },
  },

  // Resource Binding（堆叠效果）
  {
    id: 'resource-binding-1',
    type: 'default',
    data: { label: 'Resource Binding' },
    position: { x: 250, y: 200 },
    style: { background: '#cce5ff', borderRadius: '8px', padding: '10px', width: '200px', textAlign: 'center', zIndex: 10 },
  },
  {
    id: 'resource-binding-2',
    type: 'default',
    data: { label: 'Resource Binding' },
    position: { x: 255, y: 205 },
    style: { background: '#cce5ff', borderRadius: '8px', padding: '10px', width: '200px', textAlign: 'center', opacity: 0.9, zIndex: 9 },
  },
  {
    id: 'resource-binding-3',
    type: 'default',
    data: { label: 'Resource Binding' },
    position: { x: 260, y: 210 },
    style: { background: '#cce5ff', borderRadius: '8px', padding: '10px', width: '200px', textAlign: 'center', opacity: 0.8, zIndex: 8 },
  },

  // Override Policy
  {
    id: 'override-policy',
    type: 'default',
    data: { label: 'Override Policy' },
    position: { x: 500, y: 250 },
    style: { background: '#cce5ff', borderRadius: '8px', padding: '10px', width: '150px', textAlign: 'center' },
  },

  // ExecutionSpace ClusterA
  {
    id: 'exec-space-a',
    type: 'group',
    data: { label: 'ExecutionSpace\nClusterA' },
    position: { x: 100, y: 350 },
    style: { background: '#f9f9f9', border: '2px solid #ccc', borderRadius: '8px', padding: '10px', width: '250px', height: '120px' },
  },
  {
    id: 'work-a1',
    type: 'default',
    data: { label: 'work' },
    position: { x: 120, y: 400 },
    style: { background: '#cce5ff', borderRadius: '8px', padding: '5px', width: '60px', textAlign: 'center' },
    parentNode: 'exec-space-a',
  },
  {
    id: 'work-a2',
    type: 'default',
    data: { label: 'work' },
    position: { x: 190, y: 400 },
    style: { background: '#cce5ff', borderRadius: '8px', padding: '5px', width: '60px', textAlign: 'center' },
    parentNode: 'exec-space-a',
  },
  {
    id: 'work-a3',
    type: 'default',
    data: { label: 'work' },
    position: { x: 260, y: 400 },
    style: { background: '#cce5ff', borderRadius: '8px', padding: '5px', width: '60px', textAlign: 'center' },
    parentNode: 'exec-space-a',
  },

  // ClusterA 内容
  {
    id: 'cluster-a',
    type: 'group',
    data: { label: 'ClusterA' },
    position: { x: 50, y: 550 },
    style: { background: '#f0f0f0', border: '2px solid #ddd', borderRadius: '8px', padding: '10px', width: '300px', height: '150px' },
  },
  {
    id: 'deployment-a',
    type: 'default',
    data: { label: 'Deployment' },
    position: { x: 80, y: 600 },
    style: { background: '#c6e2b5', borderRadius: '8px', padding: '5px', width: '100px', textAlign: 'center' },
    parentNode: 'cluster-a',
  },
  {
    id: 'secret-a',
    type: 'default',
    data: { label: 'Secret' },
    position: { x: 200, y: 600 },
    style: { background: '#ffd6a5', borderRadius: '8px', padding: '5px', width: '100px', textAlign: 'center' },
    parentNode: 'cluster-a',
  },
  {
    id: 'configmap-a',
    type: 'default',
    data: { label: 'Configmap' },
    position: { x: 140, y: 670 },
    style: { background: '#ffb3ba', borderRadius: '8px', padding: '5px', width: '100px', textAlign: 'center' },
    parentNode: 'cluster-a',
  },

  // ExecutionSpace ClusterB
  {
    id: 'exec-space-b',
    type: 'group',
    data: { label: 'ExecutionSpace\nClusterB' },
    position: { x: 400, y: 350 },
    style: { background: '#f9f9f9', border: '2px solid #ccc', borderRadius: '8px', padding: '10px', width: '250px', height: '120px' },
  },
  {
    id: 'work-b1',
    type: 'default',
    data: { label: 'work' },
    position: { x: 420, y: 400 },
    style: { background: '#cce5ff', borderRadius: '8px', padding: '5px', width: '60px', textAlign: 'center' },
    parentNode: 'exec-space-b',
  },
  {
    id: 'work-b2',
    type: 'default',
    data: { label: 'work' },
    position: { x: 490, y: 400 },
    style: { background: '#cce5ff', borderRadius: '8px', padding: '5px', width: '60px', textAlign: 'center' },
    parentNode: 'exec-space-b',
  },
  {
    id: 'work-b3',
    type: 'default',
    data: { label: 'work' },
    position: { x: 560, y: 400 },
    style: { background: '#cce5ff', borderRadius: '8px', padding: '5px', width: '60px', textAlign: 'center' },
    parentNode: 'exec-space-b',
  },

  // ClusterB 内容
  {
    id: 'cluster-b',
    type: 'group',
    data: { label: 'ClusterB' },
    position: { x: 350, y: 550 },
    style: { background: '#f0f0f0', border: '2px solid #ddd', borderRadius: '8px', padding: '10px', width: '300px', height: '150px' },
  },
  {
    id: 'deployment-b',
    type: 'default',
    data: { label: 'Deployment' },
    position: { x: 380, y: 600 },
    style: { background: '#c6e2b5', borderRadius: '8px', padding: '5px', width: '100px', textAlign: 'center' },
    parentNode: 'cluster-b',
  },
  {
    id: 'secret-b',
    type: 'default',
    data: { label: 'Secret' },
    position: { x: 500, y: 600 },
    style: { background: '#ffd6a5', borderRadius: '8px', padding: '5px', width: '100px', textAlign: 'center' },
    parentNode: 'cluster-b',
  },
  {
    id: 'configmap-b',
    type: 'default',
    data: { label: 'Configmap' },
    position: { x: 440, y: 670 },
    style: { background: '#ffb3ba', borderRadius: '8px', padding: '5px', width: '100px', textAlign: 'center' },
    parentNode: 'cluster-b',
  },
];

const initialEdges = [
  // Resource Template → Resource Binding
  { id: 'e1', source: 'resource-template', target: 'resource-binding-1', type: 'default', animated: false },
  // Propagation Policy → Resource Binding
  { id: 'e2', source: 'propagation-policy', target: 'resource-binding-1', type: 'default', animated: false },
  // Resource Binding → Override Policy
  { id: 'e3', source: 'resource-binding-1', target: 'override-policy', type: 'default', animated: false },
  // Override Policy → Work nodes (with labels)
  { id: 'e4', source: 'override-policy', target: 'work-a1', type: 'default', label: 'Apply override', style: { fontSize: '12px' }, animated: false },
  { id: 'e5', source: 'override-policy', target: 'work-b1', type: 'default', label: 'Apply override', style: { fontSize: '12px' }, animated: false },
  // Work → Resources
  { id: 'e6', source: 'work-a1', target: 'deployment-a', type: 'default', animated: false },
  { id: 'e7', source: 'work-a2', target: 'configmap-a', type: 'default', animated: false },
  { id: 'e8', source: 'work-a3', target: 'secret-a', type: 'default', animated: false },
  { id: 'e9', source: 'work-b1', target: 'deployment-b', type: 'default', animated: false },
  { id: 'e10', source: 'work-b2', target: 'configmap-b', type: 'default', animated: false },
  { id: 'e11', source: 'work-b3', target: 'secret-b', type: 'default', animated: false },
];

export default function App() {
  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <ReactFlow
        nodes={initialNodes}
        edges={initialEdges}
        fitView
        proOptions={{ hideAttribution: true }}
      >
        <Background />
        <Controls />
      </ReactFlow>
    </div>
  );
}