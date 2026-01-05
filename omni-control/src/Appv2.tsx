import { useCallback, useMemo, useRef, useState } from 'react';
import {
  ReactFlow,
  Background,
  Controls,
  Handle,
  Position,
  addEdge,
  useEdgesState,
  useNodesState,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import type { Connection, Edge, Node, NodeProps, ReactFlowInstance } from '@xyflow/react';

type NodeCategory = 'resource' | 'policy';

type NodeKind =
  | 'ResourceTemplate'
  | 'ResourceBinding'
  | 'Work'
  | 'PropagationPolicy'
  | 'OverridePolicy';

type OmniNodeData = {
  title: string;
  subtitle?: string;
  category: NodeCategory;
  kind: NodeKind;
};

type PaletteKey = 'resource-template' | 'resource-binding' | 'work' | 'propagation-policy' | 'override-policy';

type ResourceFlowNode = Node<OmniNodeData, 'resource'>;
type PolicyFlowNode = Node<OmniNodeData, 'policy'>;
type OmniFlowNode = ResourceFlowNode | PolicyFlowNode;

type PaletteItem = {
  key: PaletteKey;
  label: string;
  node: Pick<OmniFlowNode, 'type' | 'data'>;
};

function getNodeTheme(kind: NodeKind) {
  if (kind === 'Work') {
    return {
      bg: '#EAF7EE',
      border: '#22C55E',
      badgeBg: '#16A34A',
      width: 160,
    };
  }

  if (kind === 'PropagationPolicy' || kind === 'OverridePolicy') {
    return {
      bg: '#FFF7ED',
      border: '#F97316',
      badgeBg: '#EA580C',
      width: 190,
    };
  }

  if (kind === 'ResourceTemplate') {
    return {
      bg: '#E8F2FF',
      border: '#3B82F6',
      badgeBg: '#2563EB',
      width: 320,
    };
  }

  return {
    bg: '#E8F2FF',
    border: '#3B82F6',
    badgeBg: '#2563EB',
    width: 220,
  };
}

function ResourceNode({ data, selected }: NodeProps<ResourceFlowNode>) {
  const theme = getNodeTheme(data.kind);

  return (
    <div
      style={{
        width: theme.width,
        padding: 12,
        borderRadius: 12,
        background: theme.bg,
        border: `2px solid ${selected ? theme.badgeBg : theme.border}`,
        color: '#0f172a',
        boxShadow: selected ? '0 0 0 3px rgba(37, 99, 235, 0.2)' : 'none',
      }}
    >
      <Handle type="target" position={Position.Top} style={{ background: theme.border }} />
      <Handle type="source" position={Position.Bottom} style={{ background: theme.border }} />

      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 10 }}>
        <div style={{ fontSize: 14, fontWeight: 650, whiteSpace: 'pre-line' }}>{data.title}</div>
        <div
          style={{
            fontSize: 11,
            fontWeight: 650,
            padding: '2px 8px',
            borderRadius: 999,
            background: theme.badgeBg,
            color: '#fff',
            flexShrink: 0,
          }}
        >
          Resource
        </div>
      </div>

      {data.subtitle ? (
        <div style={{ marginTop: 6, fontSize: 12, opacity: 0.85, whiteSpace: 'pre-line' }}>{data.subtitle}</div>
      ) : null}
    </div>
  );
}

function PolicyNode({ data, selected }: NodeProps<PolicyFlowNode>) {
  const theme = getNodeTheme(data.kind);

  return (
    <div
      style={{
        width: theme.width,
        padding: 12,
        borderRadius: 12,
        background: theme.bg,
        border: `2px solid ${selected ? theme.badgeBg : theme.border}`,
        color: '#0f172a',
        boxShadow: selected ? '0 0 0 3px rgba(234, 88, 12, 0.2)' : 'none',
      }}
    >
      <Handle type="target" position={Position.Left} style={{ background: theme.border }} />
      <Handle type="source" position={Position.Right} style={{ background: theme.border }} />

      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 10 }}>
        <div style={{ fontSize: 14, fontWeight: 650, whiteSpace: 'pre-line' }}>{data.title}</div>
        <div
          style={{
            fontSize: 11,
            fontWeight: 650,
            padding: '2px 8px',
            borderRadius: 999,
            background: theme.badgeBg,
            color: '#fff',
            flexShrink: 0,
          }}
        >
          Policy
        </div>
      </div>

      {data.subtitle ? (
        <div style={{ marginTop: 6, fontSize: 12, opacity: 0.85, whiteSpace: 'pre-line' }}>{data.subtitle}</div>
      ) : null}
    </div>
  );
}

export default function App() {
  const nodeTypes = useMemo(() => ({ resource: ResourceNode, policy: PolicyNode }), []);

  const initialNodes: OmniFlowNode[] = [
    {
      id: 'resource-template',
      type: 'resource',
      data: {
        title: 'Resource Template',
        subtitle:
          '(Deployment / Service / ConfigMap / Secret ... )\nExactly same with K8s APIs',
        category: 'resource',
        kind: 'ResourceTemplate',
      },
      position: { x: 200, y: 50 },
    },
    {
      id: 'resource-binding',
      type: 'resource',
      data: {
        title: 'Resource Binding',
        category: 'resource',
        kind: 'ResourceBinding',
      },
      position: { x: 250, y: 200 },
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
          type: 'resource',
          data: {
            title: 'Resource Template',
            subtitle:
              '(Deployment / Service / ConfigMap / Secret ... )\nExactly same with K8s APIs',
            category: 'resource',
            kind: 'ResourceTemplate',
          },
        },
      },
      {
        key: 'resource-binding',
        label: 'Resource Binding',
        node: {
          type: 'resource',
          data: {
            title: 'Resource Binding',
            category: 'resource',
            kind: 'ResourceBinding',
          },
        },
      },
      {
        key: 'work',
        label: 'Work',
        node: {
          type: 'resource',
          data: {
            title: 'Work',
            category: 'resource',
            kind: 'Work',
          },
        },
      },
      {
        key: 'propagation-policy',
        label: 'Propagation Policy',
        node: {
          type: 'policy',
          data: {
            title: 'Propagation Policy',
            category: 'policy',
            kind: 'PropagationPolicy',
          },
        },
      },
      {
        key: 'override-policy',
        label: 'Override Policy',
        node: {
          type: 'policy',
          data: {
            title: 'Override Policy',
            category: 'policy',
            kind: 'OverridePolicy',
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

  const [nodes, setNodes, onNodesChange] = useNodesState<OmniFlowNode>(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

  const [reactFlowInstance, setReactFlowInstance] = useState<ReactFlowInstance<OmniFlowNode, Edge> | null>(null);

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

      const newNode: OmniFlowNode = {
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
        <ReactFlow<OmniFlowNode, Edge>
          nodes={nodes}
          edges={edges}
          nodeTypes={nodeTypes}
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