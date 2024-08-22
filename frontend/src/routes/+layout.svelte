<script lang="ts">
	import { autoModeWatcher } from '@skeletonlabs/skeleton'
</script>

<svelte:head>{@html '<script>(' + autoModeWatcher.toString() + ')();</script>'}</svelte:head>
<h1>Hello world</h1>
<slot />
