export enum UpdateFilterEnum {
    ALL = "all",
}

export interface UpdateFilter {
    filter: UpdateFilterEnum;
    teamId?: string;
}

export async function updateDataStores(filter: UpdateFilter = { filter: UpdateFilterEnum.ALL }) {
}
