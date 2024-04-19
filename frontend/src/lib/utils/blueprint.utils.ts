export function recordLogoUrl(record: any): string {
    if (!record) {
        return "";
    }
    if (!record.avatar) {
        return "";
    }

    return "/api/files/" + record?.collectionId + "/" + record?.id + "/" + record?.avatar;
}
