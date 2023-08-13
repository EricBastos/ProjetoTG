<script>
// @ts-nocheck

    import { Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell} from 'flowbite-svelte';
    import { backendUrl } from "../utils/backend.js"
    import { User } from "../stores/User.js"
    import { onMount, onDestroy } from 'svelte';

    let subscriptions = [];

    let user = {
		name: null, 
	};
	subscriptions.push(
		User.subscribe(data => {
			user = data;
		})
	);

    let historyData = [];

    onMount(async () => {
        const resHistory = await fetch(backendUrl+'/user/mint/static-pix/history?page=1&pageSize=100', {
            method: 'GET',
            // headers: {
            //     'Authorization': 'Bearer ' + user.jwtToken
            // },
            credentials: 'include'
        })
        if (resHistory.status == 200) {
            const bodyJson = await resHistory.json();
            historyData = bodyJson['depositsLogs'];
            if (!historyData) historyData = [];
        }
    });

    onDestroy(() => {
        for (const prop in subscriptions) {
            subscriptions[prop]();
        }
    });

</script>

<div class="history">
    <Table shadow>
        <TableHead>
        <TableHeadCell>Id</TableHeadCell>
        <TableHeadCell>Status</TableHeadCell>
        <TableHeadCell>Amount</TableHeadCell>
        <TableHeadCell>Tax Id</TableHeadCell>
        <TableHeadCell>Wallet Address</TableHeadCell>
        <TableHeadCell>Created At</TableHeadCell>
        <TableHeadCell>Chain</TableHeadCell>
        <TableHeadCell>Blockchain Hash</TableHeadCell>
        </TableHead>
        <TableBody class="divide-y">
            {#each historyData as row}
            <TableBodyRow>
                <TableBodyCell>{row.id}</TableBodyCell>
                <TableBodyCell>{row.status}</TableBodyCell>
                <TableBodyCell>{row.amount/100}</TableBodyCell>
                <TableBodyCell>{row.taxId}</TableBodyCell>
                <TableBodyCell>{row.walletAddress}</TableBodyCell>
                <TableBodyCell>{row.createdAt}</TableBodyCell>
                <TableBodyCell>{row.chain}</TableBodyCell>
                <TableBodyCell>
                    {#if row.mintOps?.[0]?.smartContractOps?.[0]?.tx == undefined}
                        Not available
                    {:else}
                        {#if row.chain == "Polygon"}
                            <a href="https://mumbai.polygonscan.com/tx/{row.mintOps?.[0]?.smartContractOps?.[0]?.tx}" target="_blank">Click to check</a>
                        {:else}
                            <a href="https://sepolia.etherscan.io/tx/{row.mintOps?.[0]?.smartContractOps?.[0]?.tx}" target="_blank">Click to check</a>
                        {/if}
                    {/if}
                </TableBodyCell>
            </TableBodyRow>
            {/each}
        </TableBody>
    </Table>
</div>

<style>
    .history {
        padding: 8px;
        border-radius: 10px;
        max-width: 1200px;
        margin: 2% auto;
        text-align: center;
        background: white;
        color: black;
    }
</style>