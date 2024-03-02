import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params }: any) => {
    const { id } = params;

    // return the id as blueprintId
    return {
        props: {
            blueprintId: id
        }
    };
};
