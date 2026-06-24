import type { PageServerLoad } from './$types';

interface Categoria {
	id: string;
	nombre: string;
	edad_min: number;
	edad_max: number;
}

export const load: PageServerLoad = async ({ fetch }) => {
	try {
		const res = await fetch('/api/categorias');
		if (!res.ok) throw new Error('Failed to fetch');
		const categorias: Categoria[] = await res.json();
		return { categorias };
	} catch {
		return { categorias: [] };
	}
};
