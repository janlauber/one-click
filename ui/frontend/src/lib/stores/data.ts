import { client } from "$lib/pocketbase";
import type {
    RolloutsResponse,
    ProjectsResponse,
    FrameworksResponse,
    PlansResponse
} from "$lib/pocketbase/generated-types";
import { get, writable, type Writable } from "svelte/store";
import selectedProjectId from "./project";

export const frameworks: Writable<FrameworksResponse[]> = writable<FrameworksResponse[]>([]);
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
export type Pexpand = {
    framework: FrameworksResponse;
};
export const projects: Writable<ProjectsResponse<Pexpand>[]> = writable<
    ProjectsResponse<Pexpand>[]
>([]);
export type Plexpand = {
    framework: FrameworksResponse;
};
export const plans: Writable<PlansResponse<Plexpand>[]> = writable<PlansResponse<Plexpand>[]>([]);

export const selectedProject: Writable<ProjectsResponse<Pexpand> | undefined> = writable<
    ProjectsResponse<Pexpand> | undefined
>(undefined);

export enum UpdateFilterEnum {
    ALL = "all"
}

export interface UpdateFilter {
    filter: UpdateFilterEnum;
    projectId?: string;
}

export async function updateDataStores(filter: UpdateFilter = { filter: UpdateFilterEnum.ALL }) {
    if (filter.filter === UpdateFilterEnum.ALL) {
        await updateFrameworks();
        await updateProjects(filter.projectId);
        await updateRollouts(filter.projectId);
        await updatePlans();
    }
}

export async function updateFrameworks() {
    await client
        .collection("frameworks")
        .getFullList({
            sort: "application,name"
        })
        .then((response: unknown) => {
            frameworks.set(response as FrameworksResponse[]);
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

            // set the current rollout to the one without an endDate
            // @ts-ignore
            currentRollout.set(response.find((rollout) => rollout.endDate === ""));
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
        expand: "framework"
    };

    return await client.collection("projects").getFullList<ProjectsResponse<Pexpand>>(queryOptions);
}

export async function updatePlans() {
    try {
        const response = await fetchPlans();
        plans.set(response);
    } catch (error) {
        // Handle error
    }
}

async function fetchPlans(): Promise<PlansResponse<Plexpand>[]> {
    const queryOptions = {
        sort: "-created",
        expand: "framework"
    };

    return await client.collection("plans").getFullList<PlansResponse<Plexpand>>(queryOptions);
}
