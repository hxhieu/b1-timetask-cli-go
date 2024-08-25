<script lang="ts">
	import { FetchTaskInputs } from '../../wailsjs/go/gui/App'
	import { errorToast } from '../toast'

	FetchTaskInputs()
		.then((tasks) => {
			console.log(tasks)
		})
		.catch(errorToast)
</script>

<div>
	<h1>Home</h1>
</div>
