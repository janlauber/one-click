import { client } from "$lib/pocketbase";
import type {
    RolloutsResponse,
    ProjectsResponse,
    AutoUpdatesResponse,
    BlueprintsResponse,
    UsersResponse,
    DeploymentsResponse
} from "$lib/pocketbase/generated-types";
import { getClusterInfo, type ClusterInfoResponse } from "$lib/api/cluster-info";
import { get, writable, type Writable } from "svelte/store";
import type { RolloutStatusResponse } from "$lib/types/status";
import { getRolloutStatus } from "$lib/api/rollouts";
import selectedProjectId from "./project";
import selectedDeploymentId from "./deployment";

// Generic type for expandable responses
export type ExpandableResponse<T, U> = T & { expand?: U };

// Generic writable store
export function createWritableStore<T>(initialValue: T) {
    return writable<T>(initialValue);
}

// Loading
export const loading = createWritableStore<boolean>(false);

// Blueprints
export type Bexpand = {
    owner: UsersResponse;
    users?: UsersResponse[];
};
export const blueprints = createWritableStore<ExpandableResponse<BlueprintsResponse[], Bexpand>>(
    []
);

// Rollouts
export type Rexpand = {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    spec: any;
    deployment: DeploymentsResponse;
    project: ProjectsResponse;
};
export const rollouts = createWritableStore<
    ExpandableResponse<RolloutsResponse[], Rexpand | undefined>
>([]);
export const currentRollout = createWritableStore<
    ExpandableResponse<RolloutsResponse, Rexpand> | undefined
>(undefined);

// Rollout Status
export const currentRolloutStatus = createWritableStore<RolloutStatusResponse | undefined>(
    undefined
);

// Deployments
export type Dexpand = {
    project: ProjectsResponse;
    blueprint?: BlueprintsResponse;
};
export const deployments = createWritableStore<ExpandableResponse<DeploymentsResponse[], Dexpand>>(
    []
);
export const selectedDeployment = createWritableStore<
    ExpandableResponse<DeploymentsResponse, Dexpand> | undefined
>(undefined);

// Projects
export const projects = createWritableStore<ProjectsResponse[]>([]);
export const selectedProject = createWritableStore<ProjectsResponse | undefined>(undefined);

// Auto Updates
export type Aexpand = {
    project: ProjectsResponse;
};
export const autoUpdates = createWritableStore<ExpandableResponse<AutoUpdatesResponse[], Aexpand>>(
    []
);

// Cluster Info
export const clusterInfo = createWritableStore<ClusterInfoResponse | undefined>(undefined);

export enum UpdateFilterEnum {
    ALL = "all"
}

export interface UpdateFilter {
    filter: UpdateFilterEnum;
    projectId?: string;
    deploymentId?: string;
    blueprintId?: string;
}

async function updateDataStore<T, U>(
    collectionName: string,
    store: Writable<T[]>,
    filterFunc?: (item: T) => boolean,
    filter?: UpdateFilter,
    expand?: string
) {
    try {
        const queryOptions = {
            sort: "-created",
            expand: expand
        };

        const response = await client.collection(collectionName).getFullList<U>(queryOptions);

        // if filter.projectId is defined, filter the response
        if (filterFunc) {
            // set the currentRollout if the filter is deploymentId and collectionName is rollouts
            if (filter?.deploymentId && collectionName === "rollouts") {
                // currentRollout should be the one that has no endDate
                const cRollout = response.find(
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    (rollout: any) =>
                        rollout.deployment === filter.deploymentId && rollout.endDate == ""
                );
                if (cRollout) {
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    currentRollout.set(cRollout as any);
                }
                // @ts-expect-error filterFunc is defined
                store.set(response.filter(filterFunc) as T[]);
            } else {
                // @ts-expect-error filterFunc is defined
                store.set(response.filter(filterFunc) as T[]);
            }
        } else {
            store.set(response as unknown as T[]);
        }
    } catch (error) {
        // Handle error
    }
}

export async function updateClusterInfo() {
    const response = await getClusterInfo();
    clusterInfo.set(response);
}

export async function updateDataStores(filter: UpdateFilter = { filter: UpdateFilterEnum.ALL }) {
    loading.set(true);
    if (filter.filter === UpdateFilterEnum.ALL) {
        await updateDataStore(
            "projects",
            projects,
            (project) => !filter.projectId || project.id === filter.projectId,
            filter,
            "blueprint"
        );
        await updateDataStore(
            "deployments",
            deployments,
            (deployment) => !filter.projectId || deployment.project === filter.projectId,
            filter,
            "project,blueprint"
        );
        await updateDataStore(
            "blueprints",
            blueprints,
            (blueprint) => !filter.blueprintId || blueprint.id === filter.blueprintId,
            filter,
            "owner,users"
        );
        await updateDataStore(
            "autoUpdates",
            autoUpdates,
            (update) => update.deployment === filter.deploymentId,
            filter,
            "project"
        );
        await updateDataStore(
            "rollouts",
            rollouts,
            (rollout) =>
                !filter.deploymentId ||
                (rollout.project === filter.projectId &&
                    rollout.deployment === filter.deploymentId),
            filter,
            "project,deployment"
        );
        await updateClusterInfo();
    }

    // wait for all stores to be updated
    loading.set(false);
}

export async function updateCurrentRolloutStatus() {
    const tempProjectId = get(selectedProjectId);
    const tempDeploymentId = get(selectedDeploymentId);

    const response: RolloutStatusResponse | undefined = await getRolloutStatus(
        tempProjectId,
        tempDeploymentId
    );

    currentRolloutStatus.set(response);
}
