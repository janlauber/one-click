import { client } from "$lib/pocketbase";
import type {
    RolloutsResponse,
    ProjectsResponse,
    FrameworksResponse
} from "$lib/pocketbase/generated-types";
import { writable, type Writable } from "svelte/store";
import selectedProjectId from "./project";

export const frameworks: Writable<FrameworksResponse[]> = writable<FrameworksResponse[]>([]);
export type Pexpand = {
    framework: FrameworksResponse;
    rollouts: RolloutsResponse[];
};
export const projects: Writable<ProjectsResponse<Pexpand>[]> = writable<
    ProjectsResponse<Pexpand>[]
>([]);

export const selectedProject: Writable<ProjectsResponse<Pexpand> | undefined> = writable<
    ProjectsResponse<Pexpand> | undefined
>(undefined);

export enum UpdateFilterEnum {
    ALL = "all",
}

export interface UpdateFilter {
    filter: UpdateFilterEnum;
    projectId?: string;
}

export async function updateDataStores(filter: UpdateFilter = { filter: UpdateFilterEnum.ALL }) {
    if (filter.filter === UpdateFilterEnum.ALL) {
        await updateFrameworks();
        await updateProjects(filter.projectId);
    }
}

export async function updateFrameworks() {
    await client
        .collection("frameworks")
        .getFullList({
            sort: "name"
        })
        .then((response: unknown) => {
            frameworks.set(response as FrameworksResponse[]);
        })
        .catch((error) => {
            console.error(error);
        });
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

        sortProjectsRollouts();
    } catch (error) {
        // Handle error
    }
}

async function fetchProjects(): Promise<ProjectsResponse<Pexpand>[]> {
    const queryOptions = {
        sort: "-created",
        expand: "framework,rollouts"
    };

    return await client.collection("projects").getFullList<ProjectsResponse<Pexpand>>(queryOptions);
}

function sortProjectsRollouts() {
    projects.update((projects) => {
        projects.forEach((project) => {
            project.expand?.rollouts.sort((a, b) => {
                if (a.startDate > b.startDate) {
                    return -1;
                }
                if (a.startDate < b.startDate) {
                    return 1;
                }
                return 0;
            });
        });
        return projects;
    });
}
