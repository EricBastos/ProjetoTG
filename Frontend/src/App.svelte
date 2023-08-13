<script>
	import StaticDeposit from "./components/StaticDeposit.svelte";
    import Login from "./components/Login.svelte";
    import Tabs from "./components/Tabs.svelte";
	import UserBar from "./components/UserBar.svelte";
    import Withdraw from "./components/Withdraw.svelte";
	import WithdrawHistory from "./components/WithdrawHistory.svelte";
    import StaticDepositHistory from "./components/StaticDepositHistory.svelte";
	import { ProviderInfo } from "./stores/Provider.js"
	import { User } from "./stores/User.js"

	import detectEthereumProvider from '@metamask/detect-provider';
	import Web3 from 'web3';
	import { Buffer } from "buffer/"
	import { onMount, onDestroy } from 'svelte';
	import { backendUrl } from "./utils/backend.js";

	if (!window.Buffer) window.Buffer = Buffer

	let subscriptions = [];

	const tabs = ['Deposit', 'Withdraw', 'Deposit History', 'Withdraw History'];
	let currentTab = 'Deposit';

	let providerInfo = {
		w3: null,
		error: '',
	}
	subscriptions.push(
		ProviderInfo.subscribe(data => {
			providerInfo = data;
		})
	);

	let user = {
		name: null, 
	};
	subscriptions.push(
		User.subscribe(data => {
			user = data;
		})
	);

	onMount(async() => {
        const provider = await detectEthereumProvider();
        if (provider) {
			const web3 = new Web3(provider);
            providerInfo.w3 = web3;
			
			if (provider != window.ethereum) {
				providerInfo.error = 'Do you have multiple wallets installed?';
			}
        }
		ProviderInfo.set(providerInfo);

		// Check if logged in with a requisition
		
		const resTok = await fetch(backendUrl+'/user/info', {
            method: 'GET',
            credentials: 'include'
        })
        if (resTok.status == 200) {
            // Authenticated, jwt cookie stored
            const bodyJson = await resTok.json();

            User.update(u => {
                u.name = bodyJson.name
                return u;
            });
        }

    });


	const handleTabChange = (e) => {
		currentTab = e.detail;
	}

	onDestroy(() => {
        for (const prop in subscriptions) {
            subscriptions[prop]();
        }
    });

</script>

<!-- <Header /> -->
<UserBar />
<main>
	{#if user.name === null}
		<Login /> 
	{:else}
		<Tabs {tabs} {currentTab} on:tabChange={handleTabChange}/>
		{#if currentTab === 'Deposit'}
			<StaticDeposit />
		{:else if currentTab === 'Withdraw'}
			<Withdraw />
		{:else if currentTab === 'Deposit History'}
			<StaticDepositHistory />
		{:else if currentTab === 'Withdraw History'}
			<WithdrawHistory />
		{/if}
	{/if}
</main>

<style>

</style>