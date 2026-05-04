<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getVisitsByRange, getPins, getPin } from '$lib/api';
	import type { VisitWithPin, Pin } from '$lib/types';

	type Tab = 'kalendarz' | 'pins' | 'kartoteki';
	let activeTab: Tab = $state('kalendarz');

	// calendar state
	let currentDate = $state(new Date());
	let visits: VisitWithPin[] = $state([]);
	let pins: Pin[] = $state([]);
	let loading = $state(true);
	let selectedDay: string | null = $state(null);

	$effect(() => {
		loadMonth(currentDate);
	});

	async function loadMonth(date: Date) {
		loading = true;
		const year = date.getFullYear();
		const month = String(date.getMonth() + 1).padStart(2, '0');
		const start = `${year}-${month}-01`;
		const daysInMonth = new Date(year, date.getMonth() + 1, 0).getDate();
		const end = `${year}-${month}-${daysInMonth}`;
		try {
			const [v, p] = await Promise.all([getVisitsByRange(start, end), getPins()]);
			visits = v ?? [];
			pins = p ?? [];
		} catch (e) {
			console.error(e);
		} finally {
			loading = false;
		}
	}

	function prevMonth() {
		const d = new Date(currentDate);
		d.setMonth(d.getMonth() - 1);
		currentDate = d;
		selectedDay = null;
	}

	function nextMonth() {
		const d = new Date(currentDate);
		d.setMonth(d.getMonth() + 1);
		currentDate = d;
		selectedDay = null;
	}

	function monthLabel() {
		return currentDate.toLocaleDateString('pl-PL', { month: 'long', year: 'numeric' });
	}

	function getDaysInMonth() {
		const year = currentDate.getFullYear();
		const month = currentDate.getMonth();
		const days = new Date(year, month + 1, 0).getDate();
		const firstDay = new Date(year, month, 1).getDay();
		// Monday first: 0=Mon, ..., 6=Sun
		const offset = (firstDay + 6) % 7;
		return { days, offset };
	}

	function dateKey(day: number) {
		const year = currentDate.getFullYear();
		const month = String(currentDate.getMonth() + 1).padStart(2, '0');
		return `${year}-${month}-${String(day).padStart(2, '0')}`;
	}

	function visitsForDay(day: number) {
		const key = dateKey(day);
		return visits.filter((v) => v.visited_at.startsWith(key));
	}

	function newPinsForDay(day: number) {
		const key = dateKey(day);
		return pins.filter((p) => p.created_at.startsWith(key));
	}

	function visitIntensity(count: number): string {
		if (count === 0) return 'bg-gray-800';
		if (count <= 2) return 'bg-emerald-900';
		if (count <= 5) return 'bg-emerald-700';
		if (count <= 10) return 'bg-emerald-500';
		return 'bg-emerald-400';
	}

	function totalVisits() {
		return visits.length;
	}

	function totalNewPins() {
		return pins.filter((p) => {
			const year = currentDate.getFullYear();
			const month = String(currentDate.getMonth() + 1).padStart(2, '0');
			return p.created_at.startsWith(`${year}-${month}`);
		}).length;
	}

	function selectedVisits() {
		if (!selectedDay) return [];
		return visits.filter((v) => v.visited_at.startsWith(selectedDay!));
	}

	function formatTime(d: string) {
		return new Date(d).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
	}

	function formatDayLabel(key: string) {
		return new Date(key).toLocaleDateString('pl-PL', {
			weekday: 'long',
			day: 'numeric',
			month: 'long',
			year: 'numeric'
		});
	}

	async function flyToPin(pinId: number) {
		try {
			const pin = await getPin(pinId);
			goto(`/?lat=${pin.lat}&lng=${pin.lng}&zoom=16`);
		} catch (e) {
			console.error(e);
		}
	}

	const weekdays = ['Pn', 'Wt', 'Śr', 'Cz', 'Pt', 'Sb', 'Nd'];
</script>

<div class="flex h-full flex-col overflow-y-auto bg-gray-950 text-white">
	<!-- Tab bar -->
	<div class="grid shrink-0 grid-cols-3 border-b border-gray-800">
		{#each [['kalendarz', 'Kalendarz'], ['pins', 'Dane pinów'], ['kartoteki', 'Kartoteki']] as [key, label]}
			<button
				type="button"
				onclick={() => (activeTab = key as Tab)}
				class="border-b-2 py-3 text-sm font-medium transition-colors {activeTab === key
					? 'border-blue-500 text-blue-400'
					: 'border-transparent text-gray-500 hover:text-gray-300'}"
			>
				{label}
			</button>
		{/each}
	</div>

	<!-- KALENDARZ TAB -->
	{#if activeTab === 'kalendarz'}
		<div class="flex flex-1 overflow-hidden">
			<!-- Calendar column -->
			<div class="flex-1 overflow-y-auto p-4">
				<!-- Month nav -->
				<div class="mb-4 flex items-center justify-between">
					<button
						type="button"
						onclick={prevMonth}
						class="rounded-lg p-2 text-gray-400 transition-colors hover:bg-gray-800 hover:text-white"
						aria-label="Poprzedni miesiąc"
					>
						<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M15 19l-7-7 7-7"
							></path>
						</svg>
					</button>
					<h2 class="text-lg font-bold text-white capitalize">{monthLabel()}</h2>
					<button
						type="button"
						onclick={nextMonth}
						class="rounded-lg p-2 text-gray-400 transition-colors hover:bg-gray-800 hover:text-white"
						aria-label="Następny miesiąc"
					>
						<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"
							></path>
						</svg>
					</button>
				</div>

				<!-- Summary cards -->
				<div class="mb-4 grid grid-cols-2 gap-3">
					<div class="rounded-xl border border-gray-800 bg-gray-900 p-4">
						<p class="mb-1 text-xs text-gray-500">Nowe pinezki</p>
						<p class="text-2xl font-bold text-blue-400">+{totalNewPins()}</p>
					</div>
					<div class="rounded-xl border border-gray-800 bg-gray-900 p-4">
						<p class="mb-1 text-xs text-gray-500">Wizyty</p>
						<p class="text-2xl font-bold text-amber-400">↺{totalVisits()}</p>
					</div>
				</div>

				<!-- Calendar grid -->
				{#if loading}
					<div class="flex h-40 items-center justify-center">
						<p class="text-sm text-gray-500">Ładowanie...</p>
					</div>
				{:else}
					<div class="mb-1 grid grid-cols-7 gap-1">
						{#each weekdays as wd}
							<div class="py-1 text-center text-xs text-gray-600">{wd}</div>
						{/each}
					</div>

					{@const { days, offset } = getDaysInMonth()}
					<div class="grid grid-cols-7 gap-1" role="grid" aria-label="Kalendarz wizyt">
						{#each Array(offset) as _}
							<div class="aspect-square"></div>
						{/each}

						{#each Array(days) as _, i}
							{@const day = i + 1}
							{@const dayVisits = visitsForDay(day)}
							{@const dayPins = newPinsForDay(day)}
							{@const key = dateKey(day)}
							{@const isSelected = selectedDay === key}
							<button
								type="button"
								onclick={() => (selectedDay = isSelected ? null : key)}
								class="flex aspect-square flex-col items-center justify-center gap-0.5 rounded-lg border-2 transition-all {isSelected
									? 'border-blue-500'
									: 'border-transparent'} {visitIntensity(dayVisits.length)} hover:opacity-80"
								aria-label="Dzień {day}, {dayVisits.length} wizyt"
								aria-pressed={isSelected}
							>
								<span class="text-xs font-medium text-white">{day}</span>
								{#if dayVisits.length > 0}
									<span class="text-[10px] text-amber-300">↺{dayVisits.length}</span>
								{/if}
								{#if dayPins.length > 0}
									<span class="text-[10px] text-blue-300">+{dayPins.length}</span>
								{/if}
							</button>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Day detail panel -->
			{#if selectedDay}
				<div
					class="w-80 shrink-0 overflow-y-auto border-l border-gray-800 bg-gray-900"
					role="region"
					aria-label="Wizyty wybranego dnia"
				>
					<div class="sticky top-0 border-b border-gray-800 bg-gray-900 p-4">
						<p class="text-sm font-semibold text-white capitalize">{formatDayLabel(selectedDay)}</p>
						<p class="mt-0.5 text-xs text-amber-400">↺{selectedVisits().length} wizyt</p>
					</div>

					<div class="flex flex-col gap-2 p-3">
						{#each selectedVisits() as visit}
							<button
								type="button"
								onclick={() => flyToPin(visit.pin_id)}
								class="w-full rounded-xl bg-gray-800 p-3 text-left transition-colors hover:bg-gray-700 focus:ring-2 focus:ring-blue-500 focus:outline-none"
								aria-label="Przejdź do pinu {visit.pin_title}"
							>
								<p class="text-sm font-medium text-white">{visit.pin_title}</p>
								<p class="mt-0.5 text-xs text-gray-400">
									Odwiedził: {visit.name} • {formatTime(visit.visited_at)}
								</p>
								{#if visit.note?.Valid && visit.note.String}
									<p class="mt-1 text-xs text-gray-500 italic">"{visit.note.String}"</p>
								{/if}
							</button>
						{/each}
					</div>
				</div>
			{/if}
		</div>

		<!-- DANE PINÓW TAB -->
	{:else if activeTab === 'pins'}
		<div class="flex-1 overflow-y-auto p-4">
			{#if loading}
				<p class="text-sm text-gray-500">Ładowanie...</p>
			{:else}
				<div class="flex flex-col gap-2">
					{#each pins as pin}
						<a
							href="/pin/{pin.id}"
							class="flex items-center gap-3 rounded-xl border border-gray-800 bg-gray-900 p-4 transition-colors hover:bg-gray-800 focus:ring-2 focus:ring-blue-500 focus:outline-none"
							aria-label="{pin.title}, kategoria {pin.category}, {pin.visits_count} wizyt"
						>
							<div class="flex-1">
								<p class="text-sm font-medium text-white">{pin.title}</p>
								<p class="mt-0.5 text-xs text-gray-500">{pin.category}</p>
								{#if pin.description}
									<p class="mt-1 truncate text-xs text-gray-600">{pin.description}</p>
								{/if}
							</div>
							<div class="shrink-0 text-right">
								<p class="text-sm font-medium text-amber-400">↺{pin.visits_count}</p>
								<p class="text-xs text-gray-600">wizyt</p>
							</div>
						</a>
					{/each}
				</div>
			{/if}
		</div>

		<!-- KARTOTEKI TAB -->
	{:else if activeTab === 'kartoteki'}
		<div class="flex flex-1 items-center justify-center p-8">
			<div class="text-center">
				<p class="text-sm text-gray-500">Kartoteki</p>
				<p class="mt-1 text-xs text-gray-700">Funkcja w przygotowaniu</p>
			</div>
		</div>
	{/if}
</div>
