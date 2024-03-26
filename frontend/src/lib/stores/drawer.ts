import { writable } from "svelte/store";

export type NodeObject = {
    kind: "pod" | "service" | "ingress" | "secret" | "pvc";
    name: string;
    namespace: string;
    labels: Map<string, string>;
    status?: "ADDED" | "MODIFIED" | "DELETED" | "ERROR";
    containerStatusClass?: string;
    icon?: string;
    object: any;
};

export const currentNodeObject = writable<NodeObject | null>(null);

export const drawerHidden = writable<boolean>(true);

export const selectedNode = writable<NodeObject | null>(null);
