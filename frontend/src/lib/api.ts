const API_BASE = '';

function getToken(): string | null {
	if (typeof localStorage === 'undefined') return null;
	return localStorage.getItem('auth_token');
}

export async function apiFetch(url: string, options: RequestInit = {}): Promise<Response> {
	const token = getToken();
	const isFormData = options.body instanceof FormData;
	const headers: Record<string, string> = {};

	if (!isFormData) {
		headers['Content-Type'] = 'application/json';
	}

	if (token && options.method && options.method !== 'GET') {
		headers['Authorization'] = `Bearer ${token}`;
	}

	return fetch(`${API_BASE}${url}`, {
		...options,
		headers
	});
}
