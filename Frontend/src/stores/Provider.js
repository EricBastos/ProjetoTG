import { writable } from "svelte/store";

const ProviderInfo = writable({
    w3: null,
    error: '',
});

export {ProviderInfo};