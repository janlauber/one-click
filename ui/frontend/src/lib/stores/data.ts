import { client } from "$lib/pocketbase";
import type {
    DeploymentsResponse,
    ProjectsResponse,
    TechnologiesResponse
} from "$lib/pocketbase/generated-types";
import { writable, type Writable } from "svelte/store";
import selectedProjectId from "./project";

export const technologies: Writable<TechnologiesResponse[]> = writable<TechnologiesResponse[]>([]);
export type Pexpand = {
    technology: TechnologiesResponse;
    deployments: DeploymentsResponse[];
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
        await updateTechnologies();
        await updateProjects(filter.projectId);
    }
}

export async function updateTechnologies() {
    await client
        .collection("technologies")
        .getFullList({
            sort: "name"
        })
        .then((response: unknown) => {
            technologies.set(response as TechnologiesResponse[]);
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

        sortProjectsDeployments();
    } catch (error) {
        // Handle error
    }
}

async function fetchProjects(): Promise<ProjectsResponse<Pexpand>[]> {
    const queryOptions = {
        sort: "-created",
        expand: "technology,deployments"
    };

    return await client.collection("projects").getFullList<ProjectsResponse<Pexpand>>(queryOptions);
}

function sortProjectsDeployments() {
    projects.update((projects) => {
        projects.forEach((project) => {
            project.expand?.deployments.sort((a, b) => {
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
