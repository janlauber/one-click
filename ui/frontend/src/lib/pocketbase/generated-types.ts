/**
* This file was @generated using pocketbase-typegen
*/

import type PocketBase from 'pocketbase'
import type { RecordService } from 'pocketbase'

export enum Collections {
	Deployments = "deployments",
	Projects = "projects",
	Tags = "tags",
	Technologies = "technologies",
	TechnologyVersionings = "technologyVersionings",
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

export type DeploymentsRecord = {
	endDate?: IsoDateString
	startDate: IsoDateString
	tags?: RecordIdString[]
	user: RecordIdString
	values: string
	version: string
}

export type ProjectsRecord = {
	deployments?: RecordIdString[]
	endPoint?: string
	name: string
	statusEndPoint: string
	tags?: RecordIdString[]
	technology: RecordIdString
	user: RecordIdString
}

export enum TagsColorOptions {
	"default" = "default",
	"dark" = "dark",
	"red" = "red",
	"green" = "green",
	"yellow" = "yellow",
	"indigo" = "indigo",
	"purple" = "purple",
	"pink" = "pink",
}
export type TagsRecord = {
	color: TagsColorOptions
	name: string
	user: RecordIdString
}

export type TechnologiesRecord = {
	logo: string
	name: string
	url: string
}

export type TechnologyVersioningsRecord = {
	technology?: RecordIdString
	version?: string
}

export type UsersRecord = {
	avatar?: string
	name?: string
}

// Response types include system fields and match responses from the PocketBase API
export type DeploymentsResponse<Texpand = unknown> = Required<DeploymentsRecord> & BaseSystemFields<Texpand>
export type ProjectsResponse<Texpand = unknown> = Required<ProjectsRecord> & BaseSystemFields<Texpand>
export type TagsResponse<Texpand = unknown> = Required<TagsRecord> & BaseSystemFields<Texpand>
export type TechnologiesResponse<Texpand = unknown> = Required<TechnologiesRecord> & BaseSystemFields<Texpand>
export type TechnologyVersioningsResponse<Texpand = unknown> = Required<TechnologyVersioningsRecord> & BaseSystemFields<Texpand>
export type UsersResponse<Texpand = unknown> = Required<UsersRecord> & AuthSystemFields<Texpand>

// Types containing all Records and Responses, useful for creating typing helper functions

export type CollectionRecords = {
	deployments: DeploymentsRecord
	projects: ProjectsRecord
	tags: TagsRecord
	technologies: TechnologiesRecord
	technologyVersionings: TechnologyVersioningsRecord
	users: UsersRecord
}

export type CollectionResponses = {
	deployments: DeploymentsResponse
	projects: ProjectsResponse
	tags: TagsResponse
	technologies: TechnologiesResponse
	technologyVersionings: TechnologyVersioningsResponse
	users: UsersResponse
}

// Type for usage with type asserted PocketBase instance
// https://github.com/pocketbase/js-sdk#specify-typescript-definitions

export type TypedPocketBase = PocketBase & {
	collection(idOrName: 'deployments'): RecordService<DeploymentsResponse>
	collection(idOrName: 'projects'): RecordService<ProjectsResponse>
	collection(idOrName: 'tags'): RecordService<TagsResponse>
	collection(idOrName: 'technologies'): RecordService<TechnologiesResponse>
	collection(idOrName: 'technologyVersionings'): RecordService<TechnologyVersioningsResponse>
	collection(idOrName: 'users'): RecordService<UsersResponse>
}
