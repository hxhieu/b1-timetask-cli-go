<script lang="ts">
	import Router from 'svelte-spa-router'
	import { wrap } from 'svelte-spa-router/wrap'
	import { LightSwitch } from '@skeletonlabs/skeleton'
	import SettingsIcon from '~icons/cil/settings'

	import './app.css'
</script>

<div class="grid h-screen grid-rows-[auto_1fr_auto]">
	<!-- Header -->
	<header class="sticky top-0 bg-gradient-to-br from-purple-950 to-blue-950 flex justify-between z-10 p-3">
		<section></section>
		<section class="flex items-center">
			<button type="button" class="btn btn-icon btn-icon-sm variant-ghost text-primary-100 mx-4">
				<SettingsIcon class="w-5 h-5" />
			</button>
			<LightSwitch />
		</section>
	</header>
	<!-- Main -->
	<main class="p-4 space-y-4">
		<Router
			routes={{
				'/': wrap({
					asyncComponent: () => import('./routes/Home.svelte'),
				}),
				'*': wrap({
					asyncComponent: () => import('./routes/NotFound.svelte'),
				}),
			}}
		/>
	</main>
	<!-- Footer -->
	<footer class="bg-gradient-to-br from-purple-950 to-blue-950 flex justify-between p-3">(footer)</footer>
</div>

<style>
</style>
