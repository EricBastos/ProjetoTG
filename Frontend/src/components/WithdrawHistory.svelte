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
        const resHistory = await fetch(backendUrl+'/user/burn/history?page=1&pageSIze=100', {
            method: 'GET',
            credentials: 'include',
        })
        if (resHistory.status == 200) {
            const bodyJson = await resHistory.json();
            historyData = bodyJson['transfersLogs'];
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
        <TableHeadCell>Chain</TableHeadCell>
        <TableHeadCell>Wallet Address</TableHeadCell>
        <TableHeadCell>Amount</TableHeadCell>
        <TableHeadCell>Blockchain Hash</TableHeadCell>
        <TableHeadCell>Bank Code</TableHeadCell>
        <TableHeadCell>Branch Code</TableHeadCell>
        <TableHeadCell>Account Number</TableHeadCell>
        <TableHeadCell>Created At</TableHeadCell>
        </TableHead>
        <TableBody class="divide-y">
            {#each historyData as row}
            <TableBodyRow>
                <TableBodyCell>{row.id}</TableBodyCell>
                <TableBodyCell>{row.chain}</TableBodyCell>
                <TableBodyCell>{row.walletAddress}</TableBodyCell>
                <TableBodyCell>{row.amount/100}</TableBodyCell>
                <TableBodyCell>
                {#if row.smartContractOps?.[0]?.tx == undefined}
                    Not available
                {:else}
                    {#if row.chain == "Polygon"}
                        <a href="https://mumbai.polygonscan.com/tx/{row.smartContractOps?.[0]?.tx}" target="_blank">Click to check</a>
                    {:else}
                        <a href="https://sepolia.etherscan.io/tx/{row.smartContractOps?.[0]?.tx}" target="_blank">Click to check</a>
                    {/if}
                {/if}
                </TableBodyCell>
                <TableBodyCell>{row.transfers?.[0]?.bankCode == undefined ? "Processing" : row.transfers?.[0]?.bankCode}</TableBodyCell>
                <TableBodyCell>{row.transfers?.[0]?.branchCode == undefined ? "Processing" : row.transfers?.[0]?.branchCode}</TableBodyCell>
                <TableBodyCell>{row.transfers?.[0]?.accountNumber == undefined ? "Processing" : row.transfers?.[0]?.accountNumber}</TableBodyCell>
                <TableBodyCell>{row.createdAt}</TableBodyCell>
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