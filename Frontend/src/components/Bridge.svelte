<script>
    import Web3 from 'web3';
    import { ProviderInfo } from "../stores/Provider.js"
    import { User } from "../stores/User.js"
    import { EthereumAccountInfo } from "../stores/AccountInfo.js"
    import { StableContractAbi } from '../utils/abis.js';
    import { StableMumbaiAddress, MumbaiChainId, StableSepoliaAddress, SepoliaChainId } from '../utils/contracts.js';
    import { backendUrl } from "../utils/backend.js"
    import { Buffer } from "buffer/"
    import { onMount, onDestroy } from 'svelte';
    import { signERC2612Permit } from 'eth-permit';

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

    let waitingSignature = false;
    let errorMessage = "";

    let StableContract = null;

    let amount = 0;
    let pixKey = "";

    let submitHandler = async () => {
        
        if (!Web3.utils.isAddress(ethAccInfo.account)) {
            errorMessage = "Invalid wallet address";
            return
        }
        await getBalanceInfoEth();
        if (amount * 10**16 > Number(ethAccInfo.balance)) {
            errorMessage = "Insufficient balance";
            return
        }
        errorMessage = "";
        waitingSignature = true;
        const msgToSign = Date.now().toString()
        const msg = `0x${Buffer.from(msgToSign, 'utf8').toString('hex')}`;
        try {
            const spender = await StableContract.methods.owner().call();
            const result = await signERC2612Permit(window.ethereum, StableContract._address, ethAccInfo.account, spender, Web3.utils.toWei((amount/100).toString(), 'ether'), Date.now()+60*60);

            const resDeposit = await fetch(backendUrl+'/user/bridge', {
                method: 'POST',
                body: JSON.stringify({
                    pixKey: pixKey,
                    walletAddress: ethAccInfo.account,
                    inputChain: ethAccInfo.chainId == MumbaiChainId ? "Polygon" : "Ethereum",
                    outputChain: ethAccInfo.chainId == MumbaiChainId ? "Ethereum" : "Polygon",
                    amount: amount,
                    permit: {
                        deadline: result.deadline,
                        nonce: Web3.utils.hexToNumber(result.nonce),
                        r: result.r,
                        s: result.s,
                        v: result.v
                    }
                }),
                credentials: 'include'
            })
            if (resDeposit.status == 201) {
                errorMessage = "Success. Your request is being processed."
            } else if (resDeposit.status == 429) {
                errorMessage = await resDeposit.text();
            } else {
                console.log("Error generating deposit");
            }

        } catch (err) {
            console.log(err)
        }
        waitingSignature = false;

    };

    const connectMetamask = async () => {
        try {
            connectingToMetamask = true;
            const accounts = await providerInfo.w3.eth.requestAccounts();
            ethAccInfo.account = accounts[0];
            ethAccInfo.chainId = await providerInfo.w3.eth.getChainId();
            EthereumAccountInfo.update(d => {
                d.account = ethAccInfo.account;
                d.chainId = ethAccInfo.chainId;
                return d;
            });
            await fetchContract();
            await getBalanceInfoEth();
        } catch (err) {
            console.log(err.message);
        }
        connectingToMetamask = false;
    }

    const getBalanceInfoEth = async () => {
        if (StableContract == null || ethAccInfo.account == null) {
            console.log("Contract is null can't fetch balance")
            EthereumAccountInfo.update(d => {
                d.balance = "0";
                return d
            })
            return
        }
        const balance = await StableContract.methods.balanceOf(ethAccInfo.account).call();
        EthereumAccountInfo.update(d => {
            d.balance = balance;
            return d
        })
        console.log("Balance fetched")
    }

    onMount(async () => {
        console.log("Brdige mounted");
        if (providerInfo.w3) {
            console.log("Here")
            await fetchContract();
            await getBalanceInfoEth();
            window.ethereum.on('accountsChanged', async (accounts) => {
                EthereumAccountInfo.update(d => {
                    d.account = accounts[0];
                    return d;
                });
                await getBalanceInfoEth();
            })
            window.ethereum.on('chainChanged', async (chainId) => {
                ethAccInfo = chainId;
                EthereumAccountInfo.update(d => {
                    d.chainId = chainId;
                    return d;
                });
                await fetchContract();
                await getBalanceInfoEth();
            })
        }
    })

    const fetchContract = async () => {
        if (!(ethAccInfo.chainId == SepoliaChainId || ethAccInfo.chainId == MumbaiChainId)) {
            StableContract = null;
            return;
        }
        let StableAddress = '';
        if (ethAccInfo.chainId == MumbaiChainId) {
            StableAddress = StableMumbaiAddress
        } else if (ethAccInfo.chainId == SepoliaChainId) {
            StableAddress = StableSepoliaAddress
        }
        StableContract = new providerInfo.w3.eth.Contract(StableContractAbi, StableAddress)
        console.log("Contract fetched")
    }

    onDestroy(() => {
        for (const prop in subscriptions) {
            subscriptions[prop]();
        }
    });

</script>

<div class="withdraw">
    <h2>Bridge</h2>
    
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
                    <p>Input Chain: {ethAccInfo.chainId == MumbaiChainId ? "Polygon" : "Ethereum"}</p>
                    <p>Output Chain: {ethAccInfo.chainId == MumbaiChainId ? "Ethereum" : "Polygon"}</p>
                    <p>Wallet Address: {ethAccInfo.account}</p>
                    <p>StableCoin Balance: {(Number(ethAccInfo.balance)/10**18).toFixed(2)}</p>
                    Amount :
                    <input type="number" step="1" name="amount" id="amount" required bind:value={amount}>
                    <br>
                    <p hidden={!errorMessage} style="color:red">{errorMessage}</p>
                    <button disabled={waitingSignature}>Bridge</button>
                </div>
            </form>
        {/if}
    {/if}
    
    
</div>

<style>
    .withdraw {
        padding: 8px;
        border-radius: 10px;
        max-width: 600px;
        margin: 2% auto;
        text-align: center;
        background: white;
        color: black;
    }
    .allowance {
        padding: 8px;
        border-radius: 10px;
        max-width: 600px;
        margin: 2% auto;
        text-align: center;
        background: rgb(184, 184, 184);
        color: black;
    }
</style>