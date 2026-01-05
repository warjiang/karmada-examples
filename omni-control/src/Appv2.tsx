import { useCallback, useMemo, useRef, useState } from 'react';
import {
  ReactFlow,
  Background,
  Controls,
  addEdge,
  useEdgesState,
  useNodesState,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import type { Connection, Edge, Node, ReactFlowInstance } from '@xyflow/react';

type PaletteItem = {
  key: string;
  label: string;
  node: Pick<Node, 'type' | 'data' | 'style'>;
};

export default function App() {
  const initialNodes: Node[] = [
    {
      id: 'resource-template',
      type: 'default',
      data: {
        label:
          'Resource Template\n(Deployment / Service / ConfigMap / Secret ... )\nExactly same with K8s APIs',
      },
      position: { x: 200, y: 50 },
      style: {
        background: '#cce5ff',
        borderRadius: '8px',
        padding: '10px',
        width: '300px',
        textAlign: 'center',
      },
    },
    {
      id: 'resource-binding',
      type: 'default',
      data: { label: 'Resource Binding' },
      position: { x: 250, y: 200 },
      style: {
        background: '#cce5ff',
        borderRadius: '8px',
        padding: '10px',
        width: '200px',
        textAlign: 'center',
        zIndex: 10,
      },
    },
  ];

  const initialEdges: Edge[] = [
    {
      id: 'e1',
      source: 'resource-template',
      target: 'resource-binding',
      type: 'default',
      animated: true,
      markerEnd: {
        type: 'arrowclosed',
        width: 15,
        height: 15,
        color: '#000',
      },
    },
  ];

  const palette = useMemo<PaletteItem[]>(
    () => [
      {
        key: 'resource-template',
        label: 'Resource Template',
        node: {
          type: 'default',
          data: {
            label:
              'Resource Template\n(Deployment / Service / ConfigMap / Secret ... )\nExactly same with K8s APIs',
          },
          style: {
            background: '#cce5ff',
            borderRadius: '8px',
            padding: '10px',
            width: '300px',
            textAlign: 'center',
          },
        },
      },
      {
        key: 'resource-binding',
        label: 'Resource Binding',
        node: {
          type: 'default',
          data: { label: 'Resource Binding' },
          style: {
            background: '#cce5ff',
            borderRadius: '8px',
            padding: '10px',
            width: '200px',
            textAlign: 'center',
            zIndex: 10,
          },
        },
      },
      {
        key: 'propagation-policy',
        label: 'Propagation Policy',
        node: {
          type: 'default',
          data: { label: 'Propagation Policy' },
          style: {
            background: '#cce5ff',
            borderRadius: '8px',
            padding: '10px',
            width: '180px',
            textAlign: 'center',
          },
        },
      },
      {
        key: 'override-policy',
        label: 'Override Policy',
        node: {
          type: 'default',
          data: { label: 'Override Policy' },
          style: {
            background: '#cce5ff',
            borderRadius: '8px',
            padding: '10px',
            width: '180px',
            textAlign: 'center',
          },
        },
      },
    ],
    [],
  );

  const idRef = useRef(1);
  const nextId = useCallback(() => {
    const id = idRef.current;
    idRef.current += 1;
    return `dnd-${id}`;
  }, []);

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

  const [reactFlowInstance, setReactFlowInstance] = useState<ReactFlowInstance | null>(null);

  const onConnect = useCallback(
    (connection: Connection) => {
      setEdges((prev) => addEdge({ ...connection, type: 'default' }, prev));
    },
    [setEdges],
  );

  const onDragStart = useCallback((event: React.DragEvent<HTMLDivElement>, key: string) => {
    event.dataTransfer.setData('application/reactflow', key);
    event.dataTransfer.effectAllowed = 'move';
  }, []);

  const onDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
  }, []);

  const onDrop = useCallback(
    (event: React.DragEvent) => {
      event.preventDefault();
      const key = event.dataTransfer.getData('application/reactflow');
      const preset = palette.find((p) => p.key === key);
      if (!preset || !reactFlowInstance) return;

      const position = reactFlowInstance.screenToFlowPosition({
        x: event.clientX,
        y: event.clientY,
      });

      const newNode: Node = {
        id: nextId(),
        position,
        ...preset.node,
      };

      setNodes((prev) => prev.concat(newNode));
    },
    [nextId, palette, reactFlowInstance, setNodes],
  );

  return (
    <div style={{ width: '100vw', height: '100vh', display: 'flex' }}>
      <div
        style={{
          width: 260,
          padding: 12,
          borderRight: '1px solid rgba(0,0,0,0.12)',
          background: '#ffffff',
          color: '#213547',
          display: 'flex',
          flexDirection: 'column',
          gap: 10,
        }}
      >
        <div style={{ fontSize: 14, fontWeight: 600 }}>拖拽创建节点</div>
        {palette.map((item) => (
          <div
            key={item.key}
            draggable
            onDragStart={(event) => onDragStart(event, item.key)}
            style={{
              padding: '10px 12px',
              border: '1px solid rgba(0,0,0,0.15)',
              borderRadius: 8,
              background: '#f6f8fa',
              cursor: 'grab',
              userSelect: 'none',
            }}
          >
            {item.label}
          </div>
        ))}
      </div>

      <div style={{ flex: 1 }}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          onInit={setReactFlowInstance}
          onDragOver={onDragOver}
          onDrop={onDrop}
          fitView
          proOptions={{ hideAttribution: true }}
          nodesDraggable
        >
          <Background />
          <Controls />
        </ReactFlow>
      </div>
    </div>
  );
}