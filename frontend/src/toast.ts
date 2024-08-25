import { getToastStore, type ToastStore } from '@skeletonlabs/skeleton'

let toast: ToastStore

const initToast = () => {
	toast = getToastStore()
}

const errorToast = (err: string) => {
	toast.trigger({
		message: err,
		background: 'variant-filled-error',
		autohide: false,
	})
}

export { initToast, errorToast }
