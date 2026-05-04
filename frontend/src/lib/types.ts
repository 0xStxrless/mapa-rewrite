export interface AppUpdate {
	id: number;
	version: string;
	title: string;
	description: string;
	features: string;
	released_at: string;
}

export interface Category {
	name: string;
	color: string;
}

export interface User {
	id: number;
	email: string;
	must_change_password: boolean;
	created_at: string;
	last_login: string | null;
}

export interface Pin {
	id: number;
	title: string;
	description: string | null;
	lat: number;
	lng: number;
	category: string;
	image_url: string | null;
	created_at: string;
	updated_at: string;
	version: number;
	visits_count: number;
}

export interface Visit {
	id: number;
	pin_id: number;
	name: string;
	note: string | null;
	image_url: string | null;
	visited_at: string;
}

export interface VisitWithPin {
	id: number;
	pin_id: number;
	name: string;
	note: { String: string; Valid: boolean } | null;
	image_url: { String: string; Valid: boolean } | null;
	visited_at: string;
	pin_title: string;
	pin_category: string;
}

export interface PatrolPlan {
	id: number;
	name: string;
	date: string;
	created_at: string;
	updated_at: string;
}

export interface PatrolPlanPin {
	id: number;
	patrol_plan_id: number;
	pin_id: number;
	sort_order: number;
	created_at: string;
}

export interface PatrolPlanWithPins extends PatrolPlan {
	pins: Pin[];
}

export interface StreetworkStat {
	id: number;
	worker_name: string;
	month: string;
	interactions: number;
	new_contacts: number;
	interventions: number;
	created_at: string;
	updated_at: string;
	avatar: string | null;
	bg_color: string | null;
}

export interface UserUpdateViewed {
	id: number;
	user_id: number;
	update_id: number;
	viewed_at: string;
}

// input for creating/updating entities - excludes read-only fields like id, created_at, etc.

export interface CreatePinInput {
	title: string;
	description?: string;
	lat: number;
	lng: number;
	category: string;
	image_url?: string;
}

export interface CreateVisitInput {
	pin_id: number;
	name: string;
	note?: string;
	image_url?: string;
	visited_at?: string;
}

export interface CreatePatrolPlanInput {
	name: string;
	date: string;
}

export interface CreateCategoryInput {
	name: string;
	color: string;
}

export interface UpsertStatInput {
	worker_name: string;
	month: string;
	interactions: number;
	new_contacts: number;
	interventions: number;
	avatar?: string;
	bg_color?: string;
}
