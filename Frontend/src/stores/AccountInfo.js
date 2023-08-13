import { writable } from "svelte/store";

const EthereumAccountInfo = writable({
    balance: '', 
    account: null,
    chainId: '',
});

export {EthereumAccountInfo};