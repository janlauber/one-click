import { client } from "$lib/pocketbase";

export function avatarUrl(): string {
    if (client.authStore) {
        return (
            "/api/files/" +
            client.authStore.model?.collectionId +
            "/" +
            client.authStore.model?.id +
            "/" +
            client.authStore.model?.avatar
        );
    }
    return "";
}

export function avatarUrlById(collectionId: string, userId: string, avatar: string): string {
    return "/api/files/" + collectionId + "/" + userId + "/" + avatar;
}
