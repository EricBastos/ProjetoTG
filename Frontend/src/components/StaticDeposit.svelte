<script>
    import Web3 from 'web3';
    import { ProviderInfo } from "../stores/Provider.js"
    import { EthereumAccountInfo } from "../stores/AccountInfo.js"
    import { User } from "../stores/User.js"
    import { backendUrl } from "../utils/backend.js"
    import { MumbaiChainId, SepoliaChainId } from '../utils/contracts.js';
    import { onMount, onDestroy } from 'svelte';
    import { pixKey } from "../utils/pix.js";
    import { Buffer } from "buffer/"
    import copy from 'copy-to-clipboard';

    let subscriptions = [];

    let providerInfo = {
        w3: null,
        error: '',
    };
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

    let ethAccInfo = {
        balance: '', 
        account: null,
        chainId: '',
    };
    subscriptions.push(
        EthereumAccountInfo.subscribe(data => {
            ethAccInfo = data;
        })
    );

    let connectingToMetamask = false;

    let depositCreated = false;
    let creatingDeposit = false;
    let requestError = '';

    let fields = {
        amount: 0,
        walletAddress: '',
        taxId: '',
    };

    const submitHandler = async () => {
        
        if (providerInfo == null || ethAccInfo.account == '' || fields.amount <= 0) {
            return;
        }

        requestError = '';
        if (!Web3.utils.isAddress(ethAccInfo.account)) {
            requestError = "Invalid wallet address";
            return
        }
        depositCreated = false;
        creatingDeposit = true;
        const msgToSign = Date.now().toString()
        const msg = `0x${Buffer.from(msgToSign, 'utf8').toString('hex')}`;
        try {
            fields.walletAddress = ethAccInfo.account

            const resDeposit = await fetch(backendUrl+'/user/mint', {
                method: 'POST',
                body: JSON.stringify({
                    chain: ethAccInfo.chainId == MumbaiChainId ? "Polygon" : "Ethereum",
                    amount: fields.amount,
                    walletAddress: fields.walletAddress,
                }),
                // headers: {
                //     'Authorization': 'Bearer ' + user.jwtToken
                // },
                credentials: 'include'
            })
            if (resDeposit.status == 201) {
                depositCreated = true
            } else if (resDeposit.status == 429) {
                requestError = await resDeposit.text();
            } else {
                console.log("Error generating deposit");
            }
        } catch (err) {
            console.log(err)
        }
        creatingDeposit = false;
        
    };

    const connectMetamask = async () => {
        try {
            connectingToMetamask = true;
            const accounts = await providerInfo.w3.eth.requestAccounts();
            ethAccInfo.account = accounts[0];
            ethAccInfo.chainId = await providerInfo.w3.eth.getChainId();
            EthereumAccountInfo.update(d => {
                d.account = ethAccInfo.account;
                d.chainId = ethAccInfo.chainId
                return d;
            });
        } catch (err) {
            console.log(err.message);
        }
        connectingToMetamask = false;
    }

    onMount(async () => {
        if (providerInfo.w3) {
        window.ethereum.on('accountsChanged', (accounts) => {
            EthereumAccountInfo.update(d => {
                d.account = accounts[0];
                return d;
            });
        })
        window.ethereum.on('chainChanged', (chainId) => {
            EthereumAccountInfo.update(d => {
                console.log(ethAccInfo.chainId)
                d.chainId = chainId;
                return d;
            });
        })
    }
    });

    onDestroy(() => {
        for (const prop in subscriptions) {
            subscriptions[prop]();
        }
    });

</script>

<div class="deposit">
    <h2>Static Deposit</h2>
    
    {#if providerInfo.w3 == null}
            <h3>Please install Metamask!</h3>
        {:else if providerInfo.error != ''}
            <h3>{providerInfo.error}</h3>
        {:else}
            {#if ethAccInfo.account == null}
                <button disabled={connectingToMetamask} on:click={connectMetamask}>Connect Metamask</button>
            {:else if !(ethAccInfo.chainId == MumbaiChainId || ethAccInfo.chainId == SepoliaChainId)}
                <h3>Please change to Sepolia or Mumbai Network!</h3>
            {:else}
                <form on:submit|preventDefault={submitHandler}>
                    <div class="form-field">
                        Amount:
                        <input type="number" min="1" name="amount" id="amount" required bind:value={fields.amount}>
                        <br>
                        <p>Wallet Address: {ethAccInfo.account}</p>
                        <p>Chain: {ethAccInfo.chainId == MumbaiChainId ? "Mumbai" : "Sepolia"}</p>
                        <button disabled={creatingDeposit}>Generate Invoice</button>
                        <p hidden={!requestError} style="color:red">{requestError}</p>
                        <br hidden={depositCreated == false}>
                        <!-- svelte-ignore a11y-missing-attribute -->
                        <img src="./img/static_qrcode.png" class="center" style={depositCreated ? "inline-block" : "display: none"}/>
                        <button type='button' hidden={depositCreated == false} on:click={() => {copy(pixKey)}}>Copy Pix Key</button>
                        <p hidden={depositCreated == false} style="color:red">Make sure to send this exact amount from this exact CPF/CNPJ to our account!</p>
                    </div>
                </form>
            {/if}
        {/if}
    
</div>

<style>
    .center {
        display: block;
        margin-left: auto;
        margin-right: auto;
        width: 50%;
    }
    .deposit {
        padding: 8px;
        border-radius: 10px;
        max-width: 600px;
        margin: 2% auto;
        text-align: center;
        background: white;
        color: black;
    }
</style>