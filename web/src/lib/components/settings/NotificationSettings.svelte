<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { getDeviceId } from '$lib/device';
	import { notifications } from '$lib/managers/notifications.svelte';
	import {
		apiErrorMessage,
		getPrinterNotificationSettings,
		testNotification,
		updatePrinterNotificationSettings,
		type Printer
	} from '$lib/sdk';
	import { onMount } from 'svelte';

	type Props = {
		printer: Printer;
	};

	let { printer }: Props = $props();

	let loading = $state(false);
	let saving = $state(false);
	let testing = $state(false);
	let error = $state('');
	let enabled = $state(false);
	let notifyComplete = $state(false);
	let notifyError = $state(false);

	onMount(loadSettings);

	async function loadSettings() {
		loading = true;
		error = '';
		try {
			const settings = await getPrinterNotificationSettings(getDeviceId(), printer.serial);
			enabled = settings.enabled ?? false;
			notifyComplete = settings.notifyComplete ?? false;
			notifyError = settings.notifyError ?? false;
		} catch (err) {
			error = apiErrorMessage(err, 'Failed to load notification settings');
		} finally {
			loading = false;
		}
	}

	async function persist() {
		saving = true;
		error = '';
		try {
			const settings = await updatePrinterNotificationSettings(getDeviceId(), printer.serial, {
				enabled,
				notifyComplete,
				notifyError
			});
			enabled = settings.enabled ?? false;
			notifyComplete = settings.notifyComplete ?? false;
			notifyError = settings.notifyError ?? false;
		} catch (err) {
			error = apiErrorMessage(err, 'Failed to update notification settings');
		} finally {
			saving = false;
		}
	}

	async function handleMasterToggle(next: boolean) {
		error = '';

		if (next) {
			try {
				const subscribed = await notifications.isSubscribed();
				if (!subscribed) {
					const ok = await notifications.enable();
					if (!ok) {
						enabled = false;
						error = 'Notifications were blocked. Enable them in your browser settings.';
						return;
					}
				}
			} catch (err) {
				enabled = false;
				error =
					err instanceof Error
						? `${err.name}: ${err.message || 'unknown error'}`
						: apiErrorMessage(err, 'Failed to enable notifications');
				return;
			}
		}

		enabled = next;
		await persist();
	}

	async function handleSubToggle(field: 'notifyComplete' | 'notifyError', next: boolean) {
		if (field === 'notifyComplete') notifyComplete = next;
		else notifyError = next;
		await persist();
	}

	async function handleTest() {
		testing = true;
		error = '';
		try {
			await testNotification({ deviceId: getDeviceId(), serial: printer.serial });
		} catch (err) {
			error = apiErrorMessage(err, 'Failed to send test notification');
		} finally {
			testing = false;
		}
	}
</script>

<div class="grid gap-4">
	{#if !notifications.supported}
		<p class="text-sm text-muted-foreground">
			Push notifications are not supported in this browser.
		</p>
	{:else if loading}
		<div class="flex items-center gap-2 text-sm text-muted-foreground">
			<Spinner class="size-4" />
			Loading settings...
		</div>
	{:else}
		<div class="flex items-center justify-between gap-4">
			<div class="flex flex-col gap-1">
				<Label>Enable notifications</Label>
				<p class="text-xs text-muted-foreground">
					Receive push notifications for this printer on this device.
				</p>
			</div>
			<Switch checked={enabled} disabled={saving} onCheckedChange={handleMasterToggle} />
		</div>

		{#if enabled}
			<div class="grid gap-4 border-t pt-4">
				<div class="flex items-center justify-between gap-4">
					<Label>Print complete</Label>
					<Switch
						checked={notifyComplete}
						disabled={saving}
						onCheckedChange={(v) => handleSubToggle('notifyComplete', v)}
					/>
				</div>
				<div class="flex items-center justify-between gap-4">
					<Label>Print error</Label>
					<Switch
						checked={notifyError}
						disabled={saving}
						onCheckedChange={(v) => handleSubToggle('notifyError', v)}
					/>
				</div>

				<div class="flex justify-end">
					<Button variant="outline" size="sm" disabled={testing} onclick={handleTest}>
						{#if testing}
							<Spinner class="me-2 size-4" />
						{/if}
						Test
					</Button>
				</div>
			</div>
		{/if}
	{/if}

	{#if error}
		<p class="text-sm text-destructive">{error}</p>
	{/if}
</div>
