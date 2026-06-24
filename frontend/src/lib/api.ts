const API_BASE = '';

function getToken(): string | null {
	if (typeof localStorage === 'undefined') return null;
	return localStorage.getItem('auth_token');
}

export async function apiFetch(url: string, options: RequestInit = {}): Promise<Response> {
	const token = getToken();

	// Si body es FormData, no poner Content-Type (fetch lo setea con boundary)
	const isFormData = options.body instanceof FormData;

	const headers: Record<string, string> = {};

	if (!isFormData) {
		headers['Content-Type'] = 'application/json';
	}

	// Merge headers from options (only if not FormData)
	if (options.headers && !isFormData) {
		Object.assign(headers, options.headers as Record<string, string>);
	}

	if (token && options.method && options.method !== 'GET') {
		headers['Authorization'] = `Bearer ${token}`;
	}

	return fetch(`${API_BASE}${url}`, {
		...options,
		headers
	});
}
