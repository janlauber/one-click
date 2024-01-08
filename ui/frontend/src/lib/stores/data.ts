import { client } from "$lib/pocketbase";
import type {
    RolloutsResponse,
    ProjectsResponse,
    AutoUpdatesResponse,
    BlueprintsResponse,
    UsersResponse
} from "$lib/pocketbase/generated-types";
import { writable, type Writable } from "svelte/store";
import selectedProjectId from "./project";

// Blueprints //
export type Bexpand = {
    owner: UsersResponse;
    users?: UsersResponse[];
};
export const blueprints: Writable<BlueprintsResponse[]> = writable<BlueprintsResponse[]>([]);

// Rollouts //
export type Rexpand = {
    spec: any;
    project: ProjectsResponse;
};
export const rollouts: Writable<RolloutsResponse<Rexpand>[]> = writable<
    RolloutsResponse<Rexpand>[]
>([]);
export const currentRollout: Writable<RolloutsResponse<Rexpand> | undefined> = writable<
    RolloutsResponse<Rexpand> | undefined
>(undefined);

// Projects //
export type Pexpand = {
    blueprint?: BlueprintsResponse;
};
export const projects: Writable<ProjectsResponse<Pexpand>[]> = writable<
    ProjectsResponse<Pexpand>[]
>([]);
export const selectedProject: Writable<ProjectsResponse<Pexpand> | undefined> = writable<
    ProjectsResponse<Pexpand> | undefined
>(undefined);

// Auto Updates //
export type Aexpand = {
    project: ProjectsResponse;
}
export const autoUpdates: Writable<AutoUpdatesResponse<Aexpand>[]> = writable<AutoUpdatesResponse<Aexpand>[]>([]);

export enum UpdateFilterEnum {
    ALL = "all"
}

export interface UpdateFilter {
    filter: UpdateFilterEnum;
    projectId?: string;
}

export async function updateDataStores(filter: UpdateFilter = { filter: UpdateFilterEnum.ALL }) {
    if (filter.filter === UpdateFilterEnum.ALL) {
        await updateBlueprints();
        await updateProjects(filter.projectId);
        await updateRollouts(filter.projectId);
        await updateAutoUpdates(filter.projectId);
    }
}

export async function updateBlueprints() {
    await client
        .collection("blueprints")
        .getFullList({
            sort: "-created",
            expand: "owner,users"
        })
        .then((response: unknown) => {
            blueprints.set(response as BlueprintsResponse[]);
        })
        .catch((error) => {
            console.error(error);
        });
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

            return;
        }
        rollouts.set(response);
    } catch (error) {
        // Handle error
    }
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

async function fetchProjects(): Promise<ProjectsResponse<Pexpand>[]> {
    const queryOptions = {
        sort: "-created",
        expand: "blueprint"
    };

    return await client.collection("projects").getFullList<ProjectsResponse<Pexpand>>(queryOptions);
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

    return await client.collection("autoUpdates").getFullList<AutoUpdatesResponse<Aexpand>>(queryOptions);
}
