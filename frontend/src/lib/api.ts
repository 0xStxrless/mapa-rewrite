import type {
	Pin,
	CreatePinInput,
	Category,
	CreateCategoryInput,
	Visit,
	CreateVisitInput,
	PatrolPlan,
	PatrolPlanWithPins,
	CreatePatrolPlanInput,
	StreetworkStat,
	UpsertStatInput,
	AppUpdate
} from './types';

const BASE_URL = 'http://localhost:8080';

async function request<T>(method: string, endpoint: string, body?: unknown): Promise<T> {
	const token = localStorage.getItem('token');
	const res = await fetch(`${BASE_URL}${endpoint}`, {
		method,
		headers: {
			'Content-Type': 'application/json',
			...(token ? { Authorization: `Bearer ${token}` } : {})
		},
		...(body ? { body: JSON.stringify(body) } : {})
	});

	if (!res.ok) throw new Error(`${method} ${endpoint} failed: ${res.status}`);

	const text = await res.text();
	if (!text) return undefined as T;
	return JSON.parse(text) as T;
}

// Auth
export const login = (credentials: { email: string; password: string }) =>
	request<{ token: string; must_change_password: boolean }>('POST', '/login', credentials);

export const changePassword = (data: { old_password: string; new_password: string }) =>
	request<void>('POST', '/auth/change-password', data);

// Pins
export const getPins = () => request<Pin[]>('GET', '/pins');
export const getPin = (id: number) => request<Pin>('GET', `/pin/${id}`);
export const createPin = (data: CreatePinInput) => request<Pin>('POST', '/pins', data);
export const updatePin = (id: number, data: Partial<CreatePinInput>) =>
	request<Pin>('PUT', `/pin/${id}`, data);
export const deletePin = (id: number) => request<void>('DELETE', `/pin/${id}`);
export const getPinsByCategory = (category: string) =>
	request<Pin[]>('GET', `/pins/category/${category}`);

// Categories
export const getCategories = () => request<Category[]>('GET', '/categories');
export const createCategory = (data: CreateCategoryInput) =>
	request<Category>('POST', '/categories', data);
export const deleteCategory = (category: string) =>
	request<void>('DELETE', `/category/${category}`);
export const updateCategory = (data: Category) => request<Category>('PUT', '/category', data);

// Visits
export const createVisit = (data: CreateVisitInput) => request<Visit>('POST', '/visits', data);
export const deleteVisit = (id: number) => request<void>('DELETE', `/visit/${id}`);
export const getVisitsByPin = (id: number) => request<Visit[]>('GET', `/visits/pin/${id}`);
export const updateVisit = (data: Visit) => request<Visit>('PUT', '/visit', data);

// Patrol Plans
export const getPatrolPlans = () => request<PatrolPlan[]>('GET', '/patrol-plans');
export const getPatrolPlan = (id: number) => request<PatrolPlan>('GET', `/patrol-plan/${id}`);
export const getPatrolPlanWithPins = (id: number) =>
	request<PatrolPlanWithPins>('GET', `/patrol-plan/${id}/pins`);
export const createPatrolPlan = (data: CreatePatrolPlanInput) =>
	request<PatrolPlan>('POST', '/patrol-plans', data);
export const deletePatrolPlan = (id: number) => request<void>('DELETE', `/patrol-plan/${id}`);
export const addPinToPatrolPlan = (id: number, data: { pin_id: number }) =>
	request<void>('POST', `/patrol-plan/${id}/pins`, data);
export const removePinFromPatrolPlan = (id: number, pinId: number) =>
	request<void>('DELETE', `/patrol-plan/${id}/pins/${pinId}`);

// Stats
export const upsertStat = (data: UpsertStatInput) =>
	request<StreetworkStat>('POST', '/stats', data);
export const getStatsByMonth = (month: string) =>
	request<StreetworkStat[]>('GET', `/stats/month/${month}`);
export const getStatsByWorker = (worker: string) =>
	request<StreetworkStat[]>('GET', `/stats/worker/${worker}`);
export const getAllStats = () => request<StreetworkStat[]>('GET', '/all-stats');

export const getWorkers = () => request<string[]>('GET', '/workers');
