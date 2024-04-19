/**
 * This file was @generated using pocketbase-typegen
 */

import type PocketBase from "pocketbase";
import type { RecordService } from "pocketbase";

export enum Collections {
    AutoUpdates = "autoUpdates",
    Blueprints = "blueprints",
    Deployments = "deployments",
    Projects = "projects",
    Rollouts = "rollouts",
    Users = "users"
}

// Alias types for improved usability
export type IsoDateString = string;
export type RecordIdString = string;
export type HTMLString = string;

// System fields
export type BaseSystemFields<T = never> = {
    id: RecordIdString;
    created: IsoDateString;
    updated: IsoDateString;
    collectionId: string;
    collectionName: Collections;
    expand?: T;
};

export type AuthSystemFields<T = never> = {
    email: string;
    emailVisibility: boolean;
    username: string;
    verified: boolean;
} & BaseSystemFields<T>;

// Record types for each collection

export type AutoUpdatesRecord = {
    deployment: RecordIdString;
    interval: string;
    pattern?: string;
    policy?: string;
    project: RecordIdString;
    user: RecordIdString;
};

export type BlueprintsRecord<Tmanifest = unknown> = {
    avatar?: string;
    description?: string;
    manifest: null | Tmanifest;
    name: string;
    owner: RecordIdString;
    private?: boolean;
    users?: RecordIdString[];
};

export type DeploymentsRecord = {
    avatar?: string;
    blueprint?: RecordIdString;
    name?: string;
    project: RecordIdString;
    user: RecordIdString;
};

export type ProjectsRecord = {
    avatar?: string;
    description?: string;
    name: string;
    tags?: string;
    user: RecordIdString;
};

export type RolloutsRecord<Tmanifest = unknown> = {
    deployment: RecordIdString;
    endDate?: IsoDateString;
    hide?: boolean;
    manifest: null | Tmanifest;
    project: RecordIdString;
    startDate?: IsoDateString;
    user: RecordIdString;
};

export type UsersRecord = {
    avatar?: string;
    name?: string;
};

// Response types include system fields and match responses from the PocketBase API
export type AutoUpdatesResponse<Texpand = unknown> = Required<AutoUpdatesRecord> &
    BaseSystemFields<Texpand>;
export type BlueprintsResponse<Tmanifest = unknown, Texpand = unknown> = Required<
    BlueprintsRecord<Tmanifest>
> &
    BaseSystemFields<Texpand>;
export type DeploymentsResponse<Texpand = unknown> = Required<DeploymentsRecord> &
    BaseSystemFields<Texpand>;
export type ProjectsResponse<Texpand = unknown> = Required<ProjectsRecord> &
    BaseSystemFields<Texpand>;
export type RolloutsResponse<Tmanifest = unknown, Texpand = unknown> = Required<
    RolloutsRecord<Tmanifest>
> &
    BaseSystemFields<Texpand>;
export type UsersResponse<Texpand = unknown> = Required<UsersRecord> & AuthSystemFields<Texpand>;

// Types containing all Records and Responses, useful for creating typing helper functions

export type CollectionRecords = {
    autoUpdates: AutoUpdatesRecord;
    blueprints: BlueprintsRecord;
    deployments: DeploymentsRecord;
    projects: ProjectsRecord;
    rollouts: RolloutsRecord;
    users: UsersRecord;
};

export type CollectionResponses = {
    autoUpdates: AutoUpdatesResponse;
    blueprints: BlueprintsResponse;
    deployments: DeploymentsResponse;
    projects: ProjectsResponse;
    rollouts: RolloutsResponse;
    users: UsersResponse;
};

// Type for usage with type asserted PocketBase instance
// https://github.com/pocketbase/js-sdk#specify-typescript-definitions

export type TypedPocketBase = PocketBase & {
    collection(idOrName: "autoUpdates"): RecordService<AutoUpdatesResponse>;
    collection(idOrName: "blueprints"): RecordService<BlueprintsResponse>;
    collection(idOrName: "deployments"): RecordService<DeploymentsResponse>;
    collection(idOrName: "projects"): RecordService<ProjectsResponse>;
    collection(idOrName: "rollouts"): RecordService<RolloutsResponse>;
    collection(idOrName: "users"): RecordService<UsersResponse>;
};
