async function fetchFromAPI(endpoint: string) {
    let data: ClusterInfoResponse | undefined;

    const token = localStorage.getItem("pocketbase_auth");
    if (!token) {
        return data;
    }
    const authHeader = { Authorization: `Bearer ${JSON.parse(token).token}` };

    const baseUrl = window.location.hostname === "localhost" ? `http://localhost:8090` : "";
    const url = `${baseUrl}/${endpoint}`;

    try {
        const response = await fetch(url, {
            headers: authHeader
        });
        data = await response.json();
    } catch (error) {
        // Handle error
    }

    return data;
}

export async function getClusterInfo() {
    return fetchFromAPI("pb/cluster-info");
}

export type ClusterInfoResponse = {
    ingressClasses: string[];
    storageClasses: string[];
};
