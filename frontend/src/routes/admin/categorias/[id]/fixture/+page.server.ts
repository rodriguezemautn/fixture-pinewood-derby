import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const { id } = params;

	const [fixtureRes, posicionesRes, catRes] = await Promise.all([
		fetch(`/api/categorias/${id}/fixture`),
		fetch(`/api/categorias/${id}/posiciones`),
		fetch(`/api/categorias/${id}`)
	]);

	const fixture = fixtureRes.ok ? await fixtureRes.json() : null;
	const posiciones = posicionesRes.ok ? await posicionesRes.json() : [];
	const categoria = catRes.ok ? await catRes.json() : null;

	return { fixture, posiciones, categoriaId: id, categoria };
};
