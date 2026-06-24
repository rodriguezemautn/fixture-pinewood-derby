const API_BASE = '';

function getToken(): string | null {
	if (typeof localStorage === 'undefined') return null;
	return localStorage.getItem('auth_token');
}

export async function apiFetch(url: string, options: RequestInit = {}): Promise<Response> {
	const token = getToken();
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...(options.headers as Record<string, string>)
	};

	if (token && options.method && options.method !== 'GET') {
		headers['Authorization'] = `Bearer ${token}`;
	}

	return fetch(`${API_BASE}${url}`, {
		...options,
		headers
	});
}
