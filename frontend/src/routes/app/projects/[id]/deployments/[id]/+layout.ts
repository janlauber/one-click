import { UpdateFilterEnum, updateDataStores } from "$lib/stores/data";
import selectedProjectId from "$lib/stores/project";
import { get } from "svelte/store";
import type { PageLoad } from "../../../../$types";
import selectedDeploymentId from "$lib/stores/deployment";

export const load: PageLoad = async ({ params }: any) => {
    const { id } = params;

    // find the selected project id in $deployments with the id
    const projectId = get(selectedProjectId) || "";

    selectedDeploymentId.set(id);

    console.log(id);

    await updateDataStores({
        filter: UpdateFilterEnum.ALL,
        projectId: projectId,
        deploymentId: id
    }).catch((error) => {
        console.error(error);
    });
};
