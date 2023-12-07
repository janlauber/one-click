/**
* This file was @generated using pocketbase-typegen
*/

import type PocketBase from 'pocketbase'
import type { RecordService } from 'pocketbase'

export enum Collections {
	Builds = "builds",
	Frameworks = "frameworks",
	Plans = "plans",
	Projects = "projects",
	Repositories = "repositories",
	RepositoryCredentials = "repositoryCredentials",
	Rollouts = "rollouts",
	Users = "users",
}

// Alias types for improved usability
export type IsoDateString = string
export type RecordIdString = string
export type HTMLString = string

// System fields
export type BaseSystemFields<T = never> = {
	id: RecordIdString
	created: IsoDateString
	updated: IsoDateString
	collectionId: string
	collectionName: Collections
	expand?: T
}

export type AuthSystemFields<T = never> = {
	email: string
	emailVisibility: boolean
	username: string
	verified: boolean
} & BaseSystemFields<T>

// Record types for each collection

export type BuildsRecord = {
	project: RecordIdString
	repository: RecordIdString
	user: RecordIdString
}

export type FrameworksRecord<Tsettings = unknown> = {
	application?: boolean
	logo: string
	name: string
	settings: null | Tsettings
	url: string
}

export type PlansRecord<Tmanifest = unknown> = {
	description: string
	framework?: RecordIdString
	manifest: null | Tmanifest
	name: string
	price?: number
}

export type ProjectsRecord = {
	avatar?: string
	description?: string
	framework: RecordIdString
	name: string
	repository?: RecordIdString
	selectedPlan?: RecordIdString
	tags?: string
	user: RecordIdString
}

export type RepositoriesRecord = {
	private?: boolean
	repositoryCredential?: RecordIdString
	url: string
	user: RecordIdString
}

export type RepositoryCredentialsRecord = {
	password?: string
	user?: RecordIdString
	username?: string
}

export type RolloutsRecord<Tmanifest = unknown> = {
	endDate?: IsoDateString
	manifest: null | Tmanifest
	project: RecordIdString
	startDate?: IsoDateString
	user: RecordIdString
}

export type UsersRecord = {
	avatar?: string
	name?: string
}

// Response types include system fields and match responses from the PocketBase API
export type BuildsResponse<Texpand = unknown> = Required<BuildsRecord> & BaseSystemFields<Texpand>
export type FrameworksResponse<Tsettings = unknown, Texpand = unknown> = Required<FrameworksRecord<Tsettings>> & BaseSystemFields<Texpand>
export type PlansResponse<Tmanifest = unknown, Texpand = unknown> = Required<PlansRecord<Tmanifest>> & BaseSystemFields<Texpand>
export type ProjectsResponse<Texpand = unknown> = Required<ProjectsRecord> & BaseSystemFields<Texpand>
export type RepositoriesResponse<Texpand = unknown> = Required<RepositoriesRecord> & BaseSystemFields<Texpand>
export type RepositoryCredentialsResponse<Texpand = unknown> = Required<RepositoryCredentialsRecord> & BaseSystemFields<Texpand>
export type RolloutsResponse<Tmanifest = unknown, Texpand = unknown> = Required<RolloutsRecord<Tmanifest>> & BaseSystemFields<Texpand>
export type UsersResponse<Texpand = unknown> = Required<UsersRecord> & AuthSystemFields<Texpand>

// Types containing all Records and Responses, useful for creating typing helper functions

export type CollectionRecords = {
	builds: BuildsRecord
	frameworks: FrameworksRecord
	plans: PlansRecord
	projects: ProjectsRecord
	repositories: RepositoriesRecord
	repositoryCredentials: RepositoryCredentialsRecord
	rollouts: RolloutsRecord
	users: UsersRecord
}

export type CollectionResponses = {
	builds: BuildsResponse
	frameworks: FrameworksResponse
	plans: PlansResponse
	projects: ProjectsResponse
	repositories: RepositoriesResponse
	repositoryCredentials: RepositoryCredentialsResponse
	rollouts: RolloutsResponse
	users: UsersResponse
}

// Type for usage with type asserted PocketBase instance
// https://github.com/pocketbase/js-sdk#specify-typescript-definitions

export type TypedPocketBase = PocketBase & {
	collection(idOrName: 'builds'): RecordService<BuildsResponse>
	collection(idOrName: 'frameworks'): RecordService<FrameworksResponse>
	collection(idOrName: 'plans'): RecordService<PlansResponse>
	collection(idOrName: 'projects'): RecordService<ProjectsResponse>
	collection(idOrName: 'repositories'): RecordService<RepositoriesResponse>
	collection(idOrName: 'repositoryCredentials'): RecordService<RepositoryCredentialsResponse>
	collection(idOrName: 'rollouts'): RecordService<RolloutsResponse>
	collection(idOrName: 'users'): RecordService<UsersResponse>
}
