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

	const res = await fetch(`${API_BASE}${url}`, {
		...options,
		headers
	});

	// Si el token expiró o es inválido, redirigir al login
	if (res.status === 401 && token) {
		localStorage.removeItem('auth_token');
		localStorage.removeItem('auth_role');
		if (typeof window !== 'undefined') {
			window.location.href = '/login';
		}
	}

	return res;
}
