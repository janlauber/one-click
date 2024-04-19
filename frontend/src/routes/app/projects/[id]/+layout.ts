import { UpdateFilterEnum, updateDataStores } from "$lib/stores/data";
import type { PageLoad } from "../../$types";

export const load: PageLoad = async ({ params }: any) => {
    const { id } = params;

    // only update the projects if the current path is /projects/:id and not /projects/:id/deployments/:id
    if (
        window.location.pathname === `/app/projects/${id}` ||
        window.location.pathname === `/app/projects/${id}/`
    ) {
        console.log("projects");
        await updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: id
        }).catch((error) => {
            console.error(error);
        });
    }
};
