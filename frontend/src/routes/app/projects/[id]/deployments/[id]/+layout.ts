import { UpdateFilterEnum, updateDataStores } from "$lib/stores/data";
import selectedProjectId from "$lib/stores/project";
import { deployments } from "$lib/stores/data";
import { get } from "svelte/store";
import type { PageLoad } from "../../../../$types";

export const load: PageLoad = async ({ params }: any) => {
    const { id } = params;

    // find the selected project id in $deployments with the id
    let projectId = get(selectedProjectId) || "";

    await updateDataStores({
        filter: UpdateFilterEnum.ALL,
        projectId: projectId,
        deploymentId: id
    }).catch((error) => {
        console.error(error);
    });
};
