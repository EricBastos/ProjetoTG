<script>
    import { User } from "../stores/User.js"
    import { onDestroy } from 'svelte';
    import { backendUrl } from "../utils/backend.js"
    
    let subscriptions = [];
    
    let user = {
		name: null, 
	};
    subscriptions.push(
        User.subscribe(data => {
            user = data;
        })
    );

    const logout = async () => {
        const resTok = await fetch(backendUrl+'/user/logout', {
            method: 'POST',
            credentials: 'include'
        })
        if (resTok.status == 200) {
            User.update(u => {
                u.name = null
                return u;
            });
        }
    };

	onDestroy(() => {
        for (const prop in subscriptions) {
            subscriptions[prop]();
        }
        //clearInterval(interval);
    })

</script>

<main>
    {#if user.name === null}
        <ul>
            <div class="alignright">Login to start minting.</div>
            <div style="clear: both;"></div>
        </ul>
    {:else}
        <ul>
            <div class="alignleft" style="cursor: pointer;" on:click|preventDefault={logout}>Click to Logout.</div>
            
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div class="alignright">Hello, {user.name}</div>
            <div style="clear: both;"></div>
        </ul>
    {/if}
</main>

<style>

    main {
        line-height: 3;
    }
    .alignleft {
	    float: left;
        margin: 0 16px;
        font-size: 18px;
    }
    .alignright {
	    float: right;
        margin: 0 16px;
        font-size: 18px;
    }
</style>