import type { PageServerLoad } from './$types';

interface Auto {
	id: string;
	categoria_id: string;
	numero: number;
	nombre: string;
	creador: string;
	edad: number;
	foto_url: string;
}

export const load: PageServerLoad = async ({ fetch, params }) => {
	const { id } = params;

	const [autosRes, catRes] = await Promise.all([
		fetch(`/api/categorias/${id}/autos`),
		fetch(`/api/categorias/${id}`)
	]);

	const autos: Auto[] = autosRes.ok ? await autosRes.json() : [];
	const categoria = catRes.ok ? await catRes.json() : null;

	return { autos, categoria, categoriaId: id };
};
