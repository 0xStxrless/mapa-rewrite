import type { Pin } from '$lib/types';

function timeSinceVisit(updatedAt: string): string {
	const diff = Date.now() - new Date(updatedAt).getTime();
	const days = Math.floor(diff / (1000 * 60 * 60 * 24));
	if (days === 0) return 'dziś';
	if (days === 1) return '1 dzień';
	return `${days} dni`;
}

export async function createPinMarker(pin: Pin, categoryColor: string) {
	const L = (await import('leaflet')).default;
	const time = timeSinceVisit(pin.updated_at);

	const icon = L.divIcon({
		className: 'pin-marker-icon',
		html: `
    <div class="pin-marker" role="button" tabindex="0" 
         aria-label="${pin.title}, kategoria ${pin.category}, ${time}">
      <div class="pin-bubble" style="background: ${categoryColor}">${time}</div>
      <div class="pin-label">${pin.title}</div>
    </div>
  `,
		iconAnchor: [40, 12]
	});

	return L.marker([pin.lat, pin.lng], { icon });
}
