<script lang="ts">
	import GeneralSettings from '$lib/components/settings/GeneralSettings.svelte';
	import NotificationSettings from '$lib/components/settings/NotificationSettings.svelte';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { IsMobile } from '$lib/hooks/is-mobile.svelte.js';
	import { type Printer } from '$lib/sdk';
	import { cn } from '$lib/utils.js';
	import { Bell, SlidersHorizontal } from '@lucide/svelte';

	type Props = {
		open: boolean;
		printer: Printer;
	};

	let { open = $bindable(), printer }: Props = $props();

	const tabs = [
		{ id: 'general', label: 'General', icon: SlidersHorizontal },
		{ id: 'notifications', label: 'Notifications', icon: Bell }
	] as const;

	type TabId = (typeof tabs)[number]['id'];

	let activeTab = $state<TabId>('general');

	const isMobile = new IsMobile();

	$effect(() => {
		if (open) activeTab = 'general';
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-2xl">
		<Dialog.Header>
			<Dialog.Title>Printer Settings</Dialog.Title>
			<Dialog.Description>Manage settings for {printer.name}.</Dialog.Description>
		</Dialog.Header>

		<div class="flex h-[60vh] flex-col gap-4 sm:h-112 md:flex-row">
			<nav
				class={cn(
					'flex gap-1',
					isMobile.current
						? 'flex-row overflow-x-auto pb-1'
						: 'w-40 shrink-0 flex-col border-e pe-2'
				)}
			>
				{#each tabs as tab (tab.id)}
					{@const Icon = tab.icon}
					<button
						type="button"
						onclick={() => (activeTab = tab.id)}
						class={cn(
							buttonVariants({ variant: activeTab === tab.id ? 'secondary' : 'ghost' }),
							'shrink-0 justify-start'
						)}
					>
						<Icon class="size-4" />
						{tab.label}
					</button>
				{/each}
			</nav>

			<div class="min-w-0 flex-1 overflow-y-auto">
				{#if activeTab === 'general'}
					<GeneralSettings {printer} onClose={() => (open = false)} />
				{:else if activeTab === 'notifications'}
					<NotificationSettings {printer} />
				{/if}
			</div>
		</div>
	</Dialog.Content>
</Dialog.Root>
