import type { Node, Edge } from "@xyflow/svelte";

export const initialNodes: Node[] = [
    {
        id: "1",
        position: { x: 250, y: 125 },
        data: { icon: "ingress", title: "ingress" },
        type: "turbo"
    },
    {
        id: "2",
        position: { x: 500, y: 125 },
        data: { icon: "service", title: "service" },
        type: "turbo"
    },
    {
        id: "3",
        position: { x: 750, y: 125 },
        data: { icon: "pod", title: "pod" },
        type: "turbo"
    }
];

export const initialEdges: Edge[] = [
    {
        id: "e1-2",
        source: "1",
        target: "2",
        animated: true
    },
    {
        id: "e3-4",
        source: "2",
        target: "3",
        animated: true
    }
];
