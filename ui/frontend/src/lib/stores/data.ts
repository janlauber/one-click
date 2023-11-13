import { client } from "$lib/pocketbase";
import type {
    DeploymentsResponse,
    ProjectsResponse,
    TechnologiesResponse
} from "$lib/pocketbase/generated-types";
import { writable, type Writable } from "svelte/store";

export const technologies: Writable<TechnologiesResponse[]> = writable<TechnologiesResponse[]>([]);
export type Pexpand = {
    technology: TechnologiesResponse;
    deployments: DeploymentsResponse[];
};
export const projects: Writable<ProjectsResponse<Pexpand>[]> = writable<
    ProjectsResponse<Pexpand>[]
>([]);
export enum UpdateFilterEnum {
    ALL = "all"
}

export interface UpdateFilter {
    filter: UpdateFilterEnum;
}

export async function updateDataStores(filter: UpdateFilter = { filter: UpdateFilterEnum.ALL }) {
    if (filter.filter === UpdateFilterEnum.ALL) {
        await updateTechnologies();
        await updateProjects();
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

export async function updateProjects() {
    await client
        .collection("projects")
        .getFullList<ProjectsResponse<Pexpand>>({
            sort: "-created",
            expand: "technology,deployments"
        })
        .then((response: unknown) => {
            projects.set(response as ProjectsResponse<Pexpand>[]);
            // sort projects.expand?.deployments by startDate descending
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
        })
        .catch((error) => {
        });
}
