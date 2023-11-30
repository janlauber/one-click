import { UpdateFilterEnum, updateDataStores } from "$lib/stores/data";
import type { PageLoad } from "../../$types";

export const load: PageLoad = async ({ params }: any) => {
    const { id } = params;
    await updateDataStores({
        filter: UpdateFilterEnum.ALL,
        projectId: id
    }).catch((error) => {
        console.error(error);
    });
};
