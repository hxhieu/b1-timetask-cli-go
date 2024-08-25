<script lang="ts">
	import SettingsIcon from '~icons/cil/settings'
	import LoadingIcon from '~icons/cil/loop-circular'
	import LoginIcon from '~icons/cil/power-standby'

	export let email: string
	export let loading: boolean
</script>

{#if loading}
	<button type="button" class="btn btn-icon btn-icon-sm variant-ghost text-primary-100 mx-4">
		<LoadingIcon class="w-4 h-4 animate-spin" />
	</button>
{:else if !!email}
	<button type="button" class="btn btn-sm variant-ghost text-primary-100 mx-4">
		<span>{email}</span>
		<SettingsIcon class="w-5 h-5" />
	</button>
{:else}
	<button type="button" class="btn btn-sm variant-ghost text-primary-100 mx-4">
		<span>Login</span>
		<LoginIcon class="w-5 h-5" />
	</button>
{/if}
