export async function getClusterInfo() {
    const response = await fetchClusterInfo();
    return response;
}

export type ClusterInfoResponse = {
    ingressClasses: string[];
    storageClasses: string[];
}

async function fetchClusterInfo() {
    let clusterInfo: ClusterInfoResponse | undefined;

    const token = localStorage.getItem("pocketbase_auth");
    if (!token) {
        return clusterInfo;
    }
    const authHeader = { Authorization: `Bearer ${JSON.parse(token).token}` };

    // if localhost, use localhost:8090 as base url
    if (window.location.hostname === "localhost") {
        try {
            const response = await fetch(
                `http://localhost:8090/cluster-info`,
                {
                    headers: authHeader
                }
            );
            clusterInfo = await response.json();
        } catch (error) {
        }

        return clusterInfo;
    }

    try {
        const response = await fetch(`/cluster-info`, {
            headers: authHeader
        });
        clusterInfo = await response.json();
    } catch (error) {
        // Handle error
    }

    return clusterInfo;
}
