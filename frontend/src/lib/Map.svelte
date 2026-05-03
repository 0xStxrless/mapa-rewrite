<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { getPins, getCategories } from '$lib/api';
	import { createPinMarker } from '$lib/mapUtils';

	let mapEl: HTMLDivElement;
	let map: any;

	onMount(async () => {
		const token = localStorage.getItem('token');
		if (!token) return;

		const L = (await import('leaflet')).default;
		await import('leaflet/dist/leaflet.css');

		map = L.map(mapEl).setView([54.35, 18.65], 14);
		L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
			attribution: '© OpenStreetMap contributors'
		}).addTo(map);

		try {
			const [pins, categories] = await Promise.all([getPins(), getCategories()]);
			const colorMap = Object.fromEntries(categories.map((c) => [c.name, c.color]));

			for (const pin of pins) {
				const color = colorMap[pin.category] ?? '#888';
				const marker = await createPinMarker(pin, color);
				marker.addTo(map).on('click', (e: any) => {
					e.originalEvent.stopPropagation();
					goto(`/pin/${pin.id}`);
				});
				marker.addTo(map).on('click', (e: any) => {
					e.originalEvent.stopPropagation();
					goto(`/pin/${pin.id}`);
				});

				// keyboard access
				marker.getElement()?.addEventListener('keydown', (e: KeyboardEvent) => {
					if (e.key === 'Enter' || e.key === ' ') {
						goto(`/pin/${pin.id}`);
					}
				});
			}
		} catch (e) {
			console.error('Failed to load map data:', e);
		}
	});

	onDestroy(() => map?.remove());
</script>

<div bind:this={mapEl} style="width: 100%; height: calc(100vh - 41px);"></div>
