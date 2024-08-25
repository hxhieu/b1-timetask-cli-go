<script lang="ts">
	import './app.css'
	import Router, { replace } from 'svelte-spa-router'
	import { wrap } from 'svelte-spa-router/wrap'
	import { LightSwitch } from '@skeletonlabs/skeleton'
	import { initializeStores, Toast } from '@skeletonlabs/skeleton'
	import { InitUser } from '../wailsjs/go/gui/App'
	import { initToast } from './toast'
	import UserPanel from './components/UserPanel.svelte'

	initializeStores()
	initToast()

	let userEmail: string
	let fetchingUser = true

	// Default unauthenticated route
	replace('/401')
	InitUser()
		.then((email) => {
			userEmail = email
			if (userEmail) {
				replace('/')
			}
		})
		.catch(console.error)
		.finally(() => {
			fetchingUser = false
		})
</script>

<Toast />
<div class="grid h-screen grid-rows-[auto_1fr_auto]">
	<!-- Header -->
	<header class="sticky top-0 bg-gradient-to-br from-purple-950 to-blue-950 flex justify-between z-10 p-3">
		<section></section>
		<section class="flex items-center">
			<UserPanel email={userEmail} loading={fetchingUser} />
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
				'/401': wrap({
					asyncComponent: () => import('./routes/Unauthenticated.svelte'),
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
