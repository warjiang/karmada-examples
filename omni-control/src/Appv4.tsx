import {
    ReactFlow,
    Background,
    Controls,
    applyNodeChanges,
    applyEdgeChanges,
    type OnNodesChange,
    type OnEdgesChange,
    ReactFlowProvider,
    Handle,
    Position,
    type NodeProps,
    type Edge,
    type Node,
    MarkerType,
    BaseEdge,
    getBezierPath,
    type EdgeProps,
    useReactFlow,
    Panel,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { useState, useCallback, useEffect, useRef } from 'react';
import Dagre from '@dagrejs/dagre';
import { FaHeart, FaCheckCircle, FaEllipsisV, FaLayerGroup, FaNetworkWired, FaSyncAlt, FaFileAlt, FaLink, FaPlus, FaMinus } from 'react-icons/fa';

// ==================== Types ====================
type NodeStatus = 'healthy' | 'progressing' | 'degraded' | 'unknown';

interface ArgoNodeData {
    name: string;
    type: string;
    subType?: string;
    status: NodeStatus;
    syncStatus?: 'synced' | 'outOfSync' | 'unknown';
    time?: string;
    icon?: 'app' | 'svc' | 'deploy' | 'pp' | 'rb';
    // Expand/Collapse related
    collapsed?: boolean;
    childCount?: number; // Number of hidden children when collapsed
    onToggleCollapse?: () => void;
    [key: string]: unknown;
}

type ArgoNode = Node<ArgoNodeData, string>;
type ArgoEdge = Edge;

// ==================== Status Icons ====================
const StatusIcon = ({ status, syncStatus }: { status: NodeStatus; syncStatus?: 'synced' | 'outOfSync' | 'unknown' }) => {
    const heartColor = status === 'healthy' ? '#2FBDA5' : status === 'progressing' ? '#F5A623' : status === 'degraded' ? '#E96D76' : '#9CA3AF';
    const syncColor = syncStatus === 'synced' ? '#2FBDA5' : syncStatus === 'outOfSync' ? '#F5A623' : '#9CA3AF';

    return (
        <div className="flex items-center gap-0.5">
            <FaHeart size={12} color={heartColor} />
            <FaCheckCircle size={12} color={syncColor} />
        </div>
    );
};

// ==================== Node Icon ====================
const NodeIcon = ({ icon }: { icon?: string }) => {
    const iconStyle = "w-10 h-10 flex items-center justify-center rounded-lg text-white";

    switch (icon) {
        case 'app':
            return (
                <div className={`${iconStyle}`} style={{ background: 'linear-gradient(135deg, #5B6B7C 0%, #4A5568 100%)' }}>
                    <FaLayerGroup size={18} />
                </div>
            );
        case 'svc':
            return (
                <div className={`${iconStyle}`} style={{ background: 'linear-gradient(135deg, #6B7280 0%, #4B5563 100%)' }}>
                    <FaNetworkWired size={18} />
                </div>
            );
        case 'deploy':
            return (
                <div className={`${iconStyle}`} style={{ background: 'linear-gradient(135deg, #2FBDA5 0%, #1E9583 100%)' }}>
                    <FaSyncAlt size={18} />
                </div>
            );
        case 'pp':
            return (
                <div className={`${iconStyle} bg-gray-200 text-gray-600 font-bold text-sm`}>
                    PP
                </div>
            );
        case 'rb':
            return (
                <div className={`${iconStyle} bg-gray-200 text-gray-600 font-bold text-sm`}>
                    RB
                </div>
            );
        default:
            return (
                <div className={`${iconStyle} bg-gray-400`}>
                    <FaFileAlt size={18} />
                </div>
            );
    }
};

// ==================== Argo CD Style Node ====================
const ArgoNode = ({ data }: NodeProps<ArgoNode>) => {
    const hasChildren = (data.childCount ?? 0) > 0;

    return (
        <div className="bg-white rounded-lg shadow-md border border-gray-200 min-w-[200px] max-w-[280px] hover:shadow-lg transition-shadow duration-200">
            {/* Left Handle */}
            <Handle
                type="target"
                position={Position.Left}
                className="!w-3 !h-3 !bg-gray-400 !border-2 !border-white"
            />

            <div className="flex items-start p-3 gap-3">
                {/* Icon */}
                <NodeIcon icon={data.icon} />

                {/* Content */}
                <div className="flex-1 min-w-0">
                    {/* Name */}
                    <div className="font-medium text-gray-800 text-sm truncate" title={data.name}>
                        {data.name}
                    </div>

                    {/* Type & Status */}
                    <div className="flex items-center gap-2 mt-1">
                        <span className="text-xs text-gray-500">{data.type}</span>
                        {(data.subType && data.subType !== data.type) && (
                            <>
                                <span className="text-gray-300">|</span>
                                <span className="text-xs text-gray-400">{data.subType}</span>
                            </>
                        )}
                    </div>

                    {/* Status Icons */}
                    <div className="mt-1.5">
                        <StatusIcon status={data.status} syncStatus={data.syncStatus} />
                    </div>
                </div>

                {/* Menu Button */}
                <button className="p-1.5 rounded hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors">
                    <FaEllipsisV size={12} />
                </button>
            </div>

            {/* Time Badge */}
            {data.time && (
                <div className="absolute -bottom-6 left-1/2 -translate-x-1/2">
                    <span className="px-2 py-0.5 bg-gray-100 text-gray-500 text-xs rounded border border-gray-200">
                        {data.time}
                    </span>
                </div>
            )}

            {/* Right Handle with Expand/Collapse Button */}
            <Handle
                type="source"
                position={Position.Right}
                className="!w-3 !h-3 !bg-gray-400 !border-2 !border-white"
            />

            {/* Expand/Collapse Button (shown when node has children) */}
            {hasChildren && (
                <div
                    onClick={(e) => {
                        e.stopPropagation();
                        data.onToggleCollapse?.();
                    }}
                    className="absolute -right-5 top-1/2 -translate-y-1/2 w-7 h-7 bg-white rounded-md border border-gray-300 flex items-center justify-center shadow-sm hover:bg-gray-50 hover:border-gray-400 transition-all z-10"
                    title={data.collapsed ? `展开 ${data.childCount} 个子节点` : '收起子节点'}
                >
                    {data.collapsed ? (
                        <FaPlus size={10} className="text-gray-500" />
                    ) : (
                        <FaMinus size={10} className="text-gray-500" />
                    )}
                </div>
            )}
        </div>
    );
};

// ==================== Argo CD Style Root Node ====================
const ArgoRootNode = ({ data }: NodeProps<ArgoNode>) => {
    const hasChildren = (data.childCount ?? 0) > 0;

    return (
        <div className="bg-white rounded-lg shadow-md border border-gray-200 min-w-[220px] hover:shadow-lg transition-shadow duration-200">
            <div className="flex items-center p-3 gap-3">
                {/* App Icon */}
                <div className="w-12 h-12 flex items-center justify-center rounded-xl" style={{ background: 'linear-gradient(135deg, #5B6B7C 0%, #4A5568 100%)' }}>
                    <FaLayerGroup size={24} className="text-white" />
                </div>

                {/* Content */}
                <div className="flex-1 min-w-0">
                    {/* Name */}
                    <div className="font-semibold text-gray-800 text-sm truncate" title={data.name}>
                        {data.name}
                    </div>

                    {/* Status Icons */}
                    <div className="mt-1.5">
                        <StatusIcon status={data.status} syncStatus={data.syncStatus} />
                    </div>
                </div>

                {/* Menu Button */}
                <button className="p-1.5 rounded hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors">
                    <FaEllipsisV size={12} />
                </button>
            </div>

            {/* Time Badge */}
            {data.time && (
                <div className="absolute -bottom-6 left-1/2 -translate-x-1/2">
                    <span className="px-2 py-0.5 bg-gray-100 text-gray-500 text-xs rounded border border-gray-200">
                        {data.time}
                    </span>
                </div>
            )}

            {/* Right Handle */}
            <Handle
                type="source"
                position={Position.Right}
                className="!w-3 !h-3 !bg-gray-400 !border-2 !border-white"
            />

            {/* Expand/Collapse Button (shown when node has children) */}
            {hasChildren && (
                <div
                    onClick={(e) => {
                        e.stopPropagation();
                        data.onToggleCollapse?.();
                    }}
                    className="absolute -right-5 top-1/2 -translate-y-1/2 w-7 h-7 bg-white rounded-md border border-gray-300 flex items-center justify-center shadow-sm hover:bg-gray-50 hover:border-gray-400 transition-all z-10"
                    title={data.collapsed ? `展开 ${data.childCount} 个子节点` : '收起子节点'}
                >
                    {data.collapsed ? (
                        <FaPlus size={10} className="text-gray-500" />
                    ) : (
                        <FaMinus size={10} className="text-gray-500" />
                    )}
                </div>
            )}
        </div>
    );
};

// ==================== Custom Edge with Link Icon ====================
const ArgoEdge = ({ id, sourceX, sourceY, targetX, targetY, sourcePosition, targetPosition, style, markerEnd }: EdgeProps) => {
    const [edgePath] = getBezierPath({
        sourceX,
        sourceY,
        sourcePosition,
        targetX,
        targetY,
        targetPosition,
    });

    // Calculate midpoint for the link icon
    const midX = (sourceX + targetX) / 2;
    const midY = (sourceY + targetY) / 2;

    return (
        <>
            <BaseEdge id={id} path={edgePath} style={style} markerEnd={markerEnd} />
            {/* Link Icon at midpoint */}
            <foreignObject
                width={24}
                height={24}
                x={midX - 12}
                y={midY - 12}
                className="overflow-visible pointer-events-none"
            >
                <div className="w-6 h-6 bg-white rounded-full border border-gray-300 flex items-center justify-center shadow-sm">
                    <FaLink size={10} className="text-gray-400" />
                </div>
            </foreignObject>
        </>
    );
};

// ==================== Node Types ====================
const nodeTypes = {
    argoNode: ArgoNode,
    argoRootNode: ArgoRootNode,
};

const edgeTypes = {
    argoEdge: ArgoEdge,
};

// ==================== Graph Structure (parent-child relationships) ====================
interface GraphRelation {
    parentId: string;
    childIds: string[];
}

const graphRelations: GraphRelation[] = [
    { parentId: 'app', childIds: ['svc', 'deploy', 'pp'] },
    { parentId: 'svc', childIds: ['rb-svc'] },
    { parentId: 'deploy', childIds: ['rb-deploy'] },
];

// Helper function to get all descendant node IDs
const getAllDescendants = (nodeId: string, relations: GraphRelation[]): string[] => {
    const relation = relations.find(r => r.parentId === nodeId);
    if (!relation) return [];

    const descendants: string[] = [...relation.childIds];
    relation.childIds.forEach(childId => {
        descendants.push(...getAllDescendants(childId, relations));
    });
    return descendants;
};

// Helper function to get direct children
const getDirectChildren = (nodeId: string, relations: GraphRelation[]): string[] => {
    const relation = relations.find(r => r.parentId === nodeId);
    return relation?.childIds ?? [];
};

// ==================== Initial Data (ArgoCD Topology Example) ====================
const createInitialNodes = (): ArgoNode[] => [
    // Root Application Node
    {
        id: 'app',
        type: 'argoRootNode',
        position: { x: 0, y: 200 },
        data: {
            name: 'guestbook-multi-cluster',
            type: 'application',
            status: 'healthy',
            syncStatus: 'synced',
            time: '15 minutes',
            collapsed: false,
            childCount: 3,
        },
    },
    // Service Node
    {
        id: 'svc',
        type: 'argoNode',
        position: { x: 300, y: 50 },
        data: {
            name: 'guestbook-ui',
            type: 'svc',
            subType: 'Service',
            status: 'healthy',
            syncStatus: 'synced',
            time: '15 minutes',
            icon: 'svc',
            collapsed: false,
            childCount: 1,
        },
    },
    // Deployment Node
    {
        id: 'deploy',
        type: 'argoNode',
        position: { x: 300, y: 200 },
        data: {
            name: 'guestbook-ui',
            type: 'deploy',
            subType: 'Deployment',
            status: 'healthy',
            syncStatus: 'synced',
            time: '15 minutes',
            icon: 'deploy',
            collapsed: false,
            childCount: 1,
        },
    },
    // PropagationPolicy Node
    {
        id: 'pp',
        type: 'argoNode',
        position: { x: 300, y: 350 },
        data: {
            name: 'guestbook',
            type: 'propagationpolicy',
            status: 'healthy',
            syncStatus: 'synced',
            time: '37 minutes',
            icon: 'pp',
            collapsed: false,
            childCount: 0,
        },
    },
    // ResourceBinding for Service
    {
        id: 'rb-svc',
        type: 'argoNode',
        position: { x: 600, y: 50 },
        data: {
            name: 'guestbook-ui-service',
            type: 'resourcebinding',
            status: 'healthy',
            syncStatus: 'synced',
            time: '15 minutes',
            icon: 'rb',
            collapsed: false,
            childCount: 0,
        },
    },
    // ResourceBinding for Deployment
    {
        id: 'rb-deploy',
        type: 'argoNode',
        position: { x: 600, y: 200 },
        data: {
            name: 'guestbook-ui-deployment',
            type: 'resourcebinding',
            status: 'healthy',
            syncStatus: 'synced',
            time: '15 minutes',
            icon: 'rb',
            collapsed: false,
            childCount: 0,
        },
    },
];

const initialEdges: ArgoEdge[] = [
    {
        id: 'app-svc',
        source: 'app',
        target: 'svc',
        type: 'argoEdge',
        style: { stroke: '#9CA3AF', strokeWidth: 1.5 },
        markerEnd: { type: MarkerType.ArrowClosed, color: '#9CA3AF', width: 15, height: 15 },
    },
    {
        id: 'app-deploy',
        source: 'app',
        target: 'deploy',
        type: 'argoEdge',
        style: { stroke: '#9CA3AF', strokeWidth: 1.5 },
        markerEnd: { type: MarkerType.ArrowClosed, color: '#9CA3AF', width: 15, height: 15 },
    },
    {
        id: 'app-pp',
        source: 'app',
        target: 'pp',
        type: 'argoEdge',
        style: { stroke: '#9CA3AF', strokeWidth: 1.5 },
        markerEnd: { type: MarkerType.ArrowClosed, color: '#9CA3AF', width: 15, height: 15 },
    },
    {
        id: 'svc-rb-svc',
        source: 'svc',
        target: 'rb-svc',
        type: 'argoEdge',
        style: { stroke: '#9CA3AF', strokeWidth: 1.5 },
        markerEnd: { type: MarkerType.ArrowClosed, color: '#9CA3AF', width: 15, height: 15 },
    },
    {
        id: 'deploy-rb-deploy',
        source: 'deploy',
        target: 'rb-deploy',
        type: 'argoEdge',
        style: { stroke: '#9CA3AF', strokeWidth: 1.5 },
        markerEnd: { type: MarkerType.ArrowClosed, color: '#9CA3AF', width: 15, height: 15 },
    },
];

// ==================== Dagre Layout ====================
const getLayoutedElements = (nodes: ArgoNode[], edges: ArgoEdge[], options: { direction: 'TB' | 'LR' }) => {
    const g = new Dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));
    g.setGraph({
        rankdir: options.direction,
        nodesep: 80,
        ranksep: 150,
        marginx: 50,
        marginy: 50,
    });

    // Only layout visible nodes
    const visibleNodes = nodes.filter(n => !n.hidden);
    const visibleEdges = edges.filter(e => !e.hidden);

    visibleEdges.forEach((edge) => g.setEdge(edge.source, edge.target));
    visibleNodes.forEach((node) =>
        g.setNode(node.id, {
            ...node,
            width: node.measured?.width ?? 220,
            height: node.measured?.height ?? 80,
        }),
    );

    Dagre.layout(g);

    return {
        nodes: nodes.map((node: ArgoNode) => {
            if (node.hidden) return node;
            const position = g.node(node.id);
            if (!position) return node;
            const x = position.x - (node.measured?.width ?? 220) / 2;
            const y = position.y - (node.measured?.height ?? 80) / 2;

            return { ...node, position: { x, y } };
        }),
        edges,
    };
};

// ==================== Main Component ====================
function ArgoTopologyFlow() {
    const { fitView } = useReactFlow();
    const [nodes, setNodes] = useState<ArgoNode[]>(createInitialNodes());
    const [edges, setEdges] = useState(initialEdges);
    const [collapsedNodes, setCollapsedNodes] = useState<Set<string>>(new Set());
    const layoutApplied = useRef(false);

    // Toggle collapse state for a node
    const toggleCollapse = useCallback((nodeId: string) => {
        setCollapsedNodes(prev => {
            const newSet = new Set(prev);
            if (newSet.has(nodeId)) {
                newSet.delete(nodeId);
            } else {
                newSet.add(nodeId);
            }
            return newSet;
        });
        // Reset layout flag to trigger re-layout
        layoutApplied.current = false;
    }, []);

    // Update nodes and edges visibility based on collapsed state
    useEffect(() => {
        // Calculate all hidden node IDs
        const hiddenNodeIds = new Set<string>();
        collapsedNodes.forEach(collapsedNodeId => {
            const descendants = getAllDescendants(collapsedNodeId, graphRelations);
            descendants.forEach(id => hiddenNodeIds.add(id));
        });

        // Update nodes with visibility and collapse toggle callback
        setNodes(prevNodes => prevNodes.map(node => {
            const directChildCount = getDirectChildren(node.id, graphRelations).length;
            const isCollapsed = collapsedNodes.has(node.id);

            // Calculate visible child count (children that would be shown if expanded)
            const totalDescendants = getAllDescendants(node.id, graphRelations).length;

            return {
                ...node,
                hidden: hiddenNodeIds.has(node.id),
                data: {
                    ...node.data,
                    collapsed: isCollapsed,
                    childCount: directChildCount > 0 ? (isCollapsed ? totalDescendants : directChildCount) : 0,
                    onToggleCollapse: directChildCount > 0 ? () => toggleCollapse(node.id) : undefined,
                },
            };
        }));

        // Update edges visibility
        setEdges(prevEdges => prevEdges.map(edge => ({
            ...edge,
            hidden: hiddenNodeIds.has(edge.source) || hiddenNodeIds.has(edge.target),
        })));
    }, [collapsedNodes, toggleCollapse]);

    const onNodesChange: OnNodesChange<ArgoNode> = useCallback(
        (changes) => setNodes((nodesSnapshot) => applyNodeChanges(changes, nodesSnapshot)),
        [],
    );
    const onEdgesChange: OnEdgesChange<ArgoEdge> = useCallback(
        (changes) => setEdges((edgesSnapshot) => applyEdgeChanges(changes, edgesSnapshot)),
        [],
    );

    const onLayout = useCallback(() => {
        const layouted = getLayoutedElements(nodes, edges, { direction: 'LR' });
        setNodes([...layouted.nodes]);
        setEdges([...layouted.edges]);
        // Use requestAnimationFrame to ensure the DOM is updated before fitView
        window.requestAnimationFrame(() => {
            fitView({ padding: 0.3 });
        });
    }, [nodes, edges, fitView]);

    // Auto-layout when nodes are first measured or when visibility changes
    useEffect(() => {
        // Check if all visible nodes have been measured
        const visibleNodes = nodes.filter(n => !n.hidden);
        const allNodesMeasured = visibleNodes.every(node => node.measured?.width && node.measured?.height);

        if (allNodesMeasured && !layoutApplied.current && visibleNodes.length > 0) {
            layoutApplied.current = true;
            const layouted = getLayoutedElements(nodes, edges, { direction: 'LR' });
            setNodes([...layouted.nodes]);
            setEdges([...layouted.edges]);
            window.requestAnimationFrame(() => {
                fitView({ padding: 0.3 });
            });
        }
    }, [nodes, edges, fitView]);

    return (
        <div className="w-screen h-screen" style={{ background: '#E8EDEF' }}>
            <ReactFlow<ArgoNode, ArgoEdge>
                nodes={nodes}
                edges={edges}
                onNodesChange={onNodesChange}
                onEdgesChange={onEdgesChange}
                nodeTypes={nodeTypes}
                edgeTypes={edgeTypes}
                fitView
                fitViewOptions={{ padding: 0.3 }}
                defaultEdgeOptions={{
                    type: 'argoEdge',
                }}
            >
                <Background color="#ccc" gap={20} />
                <Controls />

                {/* Control Panel */}
                <Panel position="top-right">
                    <button
                        onClick={onLayout}
                        className="px-4 py-2 bg-white rounded-lg shadow-md border border-gray-200 text-gray-700 hover:bg-gray-50 transition-colors text-sm font-medium"
                    >
                        Auto Layout
                    </button>
                </Panel>
            </ReactFlow>
        </div>
    );
}

// ==================== Export ====================
export default function App() {
    return (
        <ReactFlowProvider>
            <ArgoTopologyFlow />
        </ReactFlowProvider>
    );
}
