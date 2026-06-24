import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const { id } = params;

	const [fixtureRes, posicionesRes, catRes, autosRes] = await Promise.all([
		fetch(`/api/categorias/${id}/fixture`),
		fetch(`/api/categorias/${id}/posiciones`),
		fetch(`/api/categorias/${id}`),
		fetch(`/api/categorias/${id}/autos`)
	]);

	const fixture = fixtureRes.ok ? await fixtureRes.json() : null;
	const posiciones = posicionesRes.ok ? await posicionesRes.json() : [];
	const categoria = catRes.ok ? await catRes.json() : null;
	const autos = autosRes.ok ? await autosRes.json() : [];

	// Crear mapa autoId → nombre
	const autoNombres: Record<string, string> = {};
	const autoNumeros: Record<string, number> = {};
	for (const a of autos) {
		autoNombres[a.id] = a.nombre;
		autoNumeros[a.id] = a.numero;
	}

	return { fixture, posiciones, categoria, autoNombres, autoNumeros };
};
