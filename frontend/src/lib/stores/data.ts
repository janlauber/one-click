import { client } from "$lib/pocketbase";
import type {
    RolloutsResponse,
    ProjectsResponse,
    AutoUpdatesResponse,
    BlueprintsResponse,
    UsersResponse,
    DeploymentsResponse
} from "$lib/pocketbase/generated-types";
import { getClusterInfo, type ClusterInfoResponse } from "$lib/utils/cluster-info";
import { get, writable, type Writable } from "svelte/store";
import selectedProjectId from "./project";
import type { RolloutStatusResponse } from "$lib/types/status";
import { getRolloutStatus } from "$lib/utils/rollouts";
import { c } from "svelte-highlight/languages/index";

// Blueprints //
export type Bexpand = {
    owner: UsersResponse;
    users?: UsersResponse[];
};
export const blueprints: Writable<BlueprintsResponse[]> = writable<BlueprintsResponse[]>([]);

// Rollouts //
export type Rexpand = {
    spec: any;
    deployment: DeploymentsResponse;
    project: ProjectsResponse;
};
export const rollouts: Writable<RolloutsResponse<Rexpand>[]> = writable<
    RolloutsResponse<Rexpand>[]
>([]);
export const currentRollout: Writable<RolloutsResponse<Rexpand> | undefined> = writable<
    RolloutsResponse<Rexpand> | undefined
>(undefined);

// Rollout Status //
export const currentRolloutStatus: Writable<RolloutStatusResponse | undefined> = writable<
    RolloutStatusResponse | undefined
>(undefined);

// Deployments //
export type Dexpand = {
    project: ProjectsResponse;
    blueprint?: BlueprintsResponse;
};
export const deployments: Writable<DeploymentsResponse<Dexpand>[]> = writable<DeploymentsResponse<Dexpand>[]>([]);

// Projects //
export const projects: Writable<ProjectsResponse[]> = writable<ProjectsResponse[]>([]);
export const selectedProject: Writable<ProjectsResponse | undefined> = writable<
    ProjectsResponse | undefined
>(undefined);

// Auto Updates //
export type Aexpand = {
    project: ProjectsResponse;
};
export const autoUpdates: Writable<AutoUpdatesResponse<Aexpand>[]> = writable<
    AutoUpdatesResponse<Aexpand>[]
>([]);

// Cluster Info //
export const clusterInfo: Writable<ClusterInfoResponse | undefined> = writable<
    ClusterInfoResponse | undefined
>(undefined);

export enum UpdateFilterEnum {
    ALL = "all"
}

export interface UpdateFilter {
    filter: UpdateFilterEnum;
    projectId?: string;
    deploymentId?: string;
    blueprintId?: string;
}

export async function updateDataStores(filter: UpdateFilter = { filter: UpdateFilterEnum.ALL }) {
    if (filter.filter === UpdateFilterEnum.ALL) {
        await updateProjects(filter.projectId);
        await updateDeployments(filter.projectId);
        await updateBlueprints(filter.deploymentId);
        await updateAutoUpdates(filter.deploymentId);
        await updateRollouts(filter.deploymentId);
        await updateClusterInfo();
    }
}

export async function updateRollouts(projectId?: string) {
    try {
        const response = await fetchRollouts();
        if (projectId) {
            // set selected project
            selectedProjectId.set(projectId);
            // @ts-ignore
            rollouts.set(response.filter((rollout) => rollout.project === projectId));

            // set the current rollout to the one without an endDate and the project id
            // @ts-ignore
            currentRollout.set(
                response.find((rollout) => rollout.project === projectId && !rollout.endDate)
            );

            if (currentRollout && projectId && get(currentRollout) !== undefined) {
                const response = await getRolloutStatus(projectId, get(currentRollout)!.id);
                currentRolloutStatus.set(response);
            } else {
                currentRolloutStatus.set(undefined);
            }

            return;
        }
        rollouts.set(response);
    } catch (error) {
        // Handle error
    }
}

export async function updateDeployments(projectId?: string) {
    try {
        const response = await fetchDeployments();
        if (projectId) {
            // set selected project
            selectedProjectId.set(projectId);
            // @ts-ignore
            deployments.set(response.filter((deployment) => deployment.project === projectId));
            return;
        }
        deployments.set(response);
    } catch (error) {
        // Handle error
    }
}

async function fetchDeployments(): Promise<DeploymentsResponse<Dexpand>[]> {
    const queryOptions = {
        sort: "-created",
        expand: "project,blueprint"
    };

    return await client.collection("deployments").getFullList<DeploymentsResponse<Dexpand>>(queryOptions);
}

export async function updateBlueprints(blueprintId?: string) {
    await client
        .collection("blueprints")
        .getFullList({
            sort: "-created",
            expand: "owner,users"
        })
        .then((response: unknown) => {
            // if blueprintId is set, set the selected blueprint filtered by the id
            if (blueprintId) {
                // @ts-ignore
                blueprints.set(response.filter((blueprint) => blueprint.id === blueprintId));
                return;
            }
            blueprints.set(response as BlueprintsResponse[]);
        })
        .catch((error) => {
            console.error(error);
        });
}

export async function updateCurrentRolloutStatus(projectId: string) {
    const rolloutId = get(currentRollout)!.id;
    const response = await getRolloutStatus(projectId, rolloutId);
    currentRolloutStatus.set(response);
}

async function fetchRollouts(): Promise<RolloutsResponse<Rexpand>[]> {
    const queryOptions = {
        sort: "-created",
        expand: "project"
    };

    return await client.collection("rollouts").getFullList<RolloutsResponse<Rexpand>>(queryOptions);
}

export async function updateProjects(projectId?: string) {
    try {
        const response = await fetchProjects();
        if (projectId) {
            // set selected project
            selectedProjectId.set(projectId);
            selectedProject.set(response.find((project) => project.id === projectId));
        }
        projects.set(response);
    } catch (error) {
        // Handle error
    }
}

async function fetchProjects(): Promise<ProjectsResponse[]> {
    const queryOptions = {
        sort: "-created",
        expand: "blueprint"
    };

    return await client.collection("projects").getFullList<ProjectsResponse>(queryOptions);
}

export async function updateAutoUpdates(projectId?: string) {
    try {
        const response = await fetchAutoUpdates();
        if (projectId) {
            // set selected project
            selectedProjectId.set(projectId);
            // @ts-ignore
            autoUpdates.set(response.filter((update) => update.project === projectId));
            return;
        }
        autoUpdates.set(response);
    } catch (error) {
        // Handle error
    }
}

async function fetchAutoUpdates(): Promise<AutoUpdatesResponse<Aexpand>[]> {
    const queryOptions = {
        sort: "-created",
        expand: "project"
    };

    return await client
        .collection("autoUpdates")
        .getFullList<AutoUpdatesResponse<Aexpand>>(queryOptions);
}

export async function updateClusterInfo() {
    const response = await getClusterInfo();
    clusterInfo.set(response);
}
