import { writable } from 'svelte/store';
import { login as apiLogin } from './api';

interface AuthState {
	token: string;
	must_change_password: boolean;
}

function loadAuth(): AuthState | null {
	if (typeof localStorage === 'undefined') return null;
	const token = localStorage.getItem('token');
	if (!token) return null;
	return { token, must_change_password: false };
}

export const auth = writable<AuthState | null>(loadAuth());

export async function login(email: string, password: string) {
	const data = await apiLogin({ email, password });
	localStorage.setItem('token', data.token);
	auth.set({ token: data.token, must_change_password: data.must_change_password });
	return data;
}

export function logout() {
	localStorage.removeItem('token');
	auth.set(null);
}
