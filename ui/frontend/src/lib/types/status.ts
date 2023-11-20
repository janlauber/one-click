// Define a type for the status values
type Status = "OK" | "Pending" | "Error";

// Define a type for the deployment object
interface Deployment {
    replicas: number;
    podNames: string[];
    status: Status;
}

// Define a type for the service objects
interface Service {
    name: string;
    ports: number[];
    status: Status;
}

// Define a type for the ingress objects
interface Ingress {
    name: string;
    hosts: string[];
    status: Status;
}

export interface RolloutStatusResponse {
    deployment: Deployment;
    services?: Service[];
    ingresses?: Ingress[];
}
