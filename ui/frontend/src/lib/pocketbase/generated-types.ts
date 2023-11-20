/**
* This file was @generated using pocketbase-typegen
*/

import type PocketBase from 'pocketbase'
import type { RecordService } from 'pocketbase'

export enum Collections {
	Frameworks = "frameworks",
	Projects = "projects",
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

export type FrameworksRecord = {
	logo: string
	name: string
	url: string
}

export type ProjectsRecord = {
	framework: RecordIdString
	name: string
	tags?: string
	user: RecordIdString
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
export type FrameworksResponse<Texpand = unknown> = Required<FrameworksRecord> & BaseSystemFields<Texpand>
export type ProjectsResponse<Texpand = unknown> = Required<ProjectsRecord> & BaseSystemFields<Texpand>
export type RolloutsResponse<Tmanifest = unknown, Texpand = unknown> = Required<RolloutsRecord<Tmanifest>> & BaseSystemFields<Texpand>
export type UsersResponse<Texpand = unknown> = Required<UsersRecord> & AuthSystemFields<Texpand>

// Types containing all Records and Responses, useful for creating typing helper functions

export type CollectionRecords = {
	frameworks: FrameworksRecord
	projects: ProjectsRecord
	rollouts: RolloutsRecord
	users: UsersRecord
}

export type CollectionResponses = {
	frameworks: FrameworksResponse
	projects: ProjectsResponse
	rollouts: RolloutsResponse
	users: UsersResponse
}

// Type for usage with type asserted PocketBase instance
// https://github.com/pocketbase/js-sdk#specify-typescript-definitions

export type TypedPocketBase = PocketBase & {
	collection(idOrName: 'frameworks'): RecordService<FrameworksResponse>
	collection(idOrName: 'projects'): RecordService<ProjectsResponse>
	collection(idOrName: 'rollouts'): RecordService<RolloutsResponse>
	collection(idOrName: 'users'): RecordService<UsersResponse>
}
