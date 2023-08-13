import { writable } from "svelte/store";

const User = writable({
    name: null, 
});

export {User};