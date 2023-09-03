<script>

    import jwt_decode from 'jwt-decode';
    import { User } from "../stores/User.js"
    import { backendUrl } from "../utils/backend.js"
    import { onDestroy } from 'svelte';

    let subscriptions = [];

    let loginScreen = true;

    let fields = {
        email: '',
        password: '',
    };

    let createAccountFields = {
        name: '',
        email: '',
        password: '',
        confirmPassword: '',
        taxId: '',
    };

    let user = {
		name: null, 
	};
    subscriptions.push(
        User.subscribe(data => {
            user = data;
        })
    );

    let loginError = '';
    let createAccountMessage = '';

    const submitHandler = async () => {
        loginError = '';
        try {
            const resTok = await fetch(backendUrl+'/user/login', {
                method: 'POST',
                body: JSON.stringify({
                    email: fields.email,
                    password: fields.password,
                }),
                credentials: 'include'
            })
            if (resTok.status == 200) {
                // Authenticated, jwt cookie stored
                const bodyJson = await resTok.json();
                const token = bodyJson.accessToken;
                const tokenPayload = jwt_decode(token);

                User.update(u => {
                    u.name = tokenPayload.name
                    return u;
                });
            } else if (resTok.status == 429) {
                loginError = 'Rate limited. Try again later'
            } else {
                loginError = 'Wrong credentials';
            }
        } catch(error) {
            loginError = 'Connection error'
        }
    };

    const createAccountHandler = async () => {
        createAccountMessage = '';
        try {
            const resTok = await fetch(backendUrl+'/user/create', {
                method: 'POST',
                body: JSON.stringify({
                    name: createAccountFields.name,
                    email: createAccountFields.email,
                    password: createAccountFields.password,
                    confirmPassword: createAccountFields.confirmPassword,
                    taxId: createAccountFields.taxId
                }),
                credentials: 'include'
            })
            if (resTok.status == 201) {
                // Authenticated, jwt cookie stored
                createAccountMessage = "Account created with success. Redirecting to login page in 3 seconds."
                setTimeout(() => {
                    loginScreen = true;
                    createAccountMessage = '';
                }, 3000)
            } else if (resTok.status == 429) {
                createAccountMessage = 'Rate limited. Try again later'
            } else {
                const bodyJson = await resTok.json();
                createAccountMessage = bodyJson.error;
            }
        } catch(error) {
            createAccountMessage = 'Connection error';
        }
    };

    onDestroy(() => {
        for (const prop in subscriptions) {
            subscriptions[prop]();
        }
    });

</script>

<main>
    {#if loginScreen}
        <div class="login">
            <h2>Sign in</h2>
            <form on:submit|preventDefault={submitHandler}>
                Email:
                <input type="text" placeholder="user@email.com" id="email" bind:value={fields.email}>
                <br>
                Password:
                <input type="password" id="password" bind:value={fields.password}>
                <br>
                <p hidden={!loginError} style="color:red">{loginError}</p>
                <button>Login</button>
            </form>
            <div style="cursor: pointer; color: #3333AA" on:click|preventDefault={() => {loginScreen = false}}>Or click here to create an account.</div>
        </div>
    {:else}
        <div class="login">
            <h2>Create an account</h2>
            <form on:submit|preventDefault={createAccountHandler}>
                Name:
                <input type="text" placeholder="Pedro" id="email" bind:value={createAccountFields.name}>
                <br>
                Email:
                <input type="text" placeholder="user@email.com" id="email" bind:value={createAccountFields.email}>
                <br>
                Password:
                <input type="password" id="password" bind:value={createAccountFields.password}>
                <br>
                Confirm Password:
                <input type="password" id="confirmPassword" bind:value={createAccountFields.confirmPassword}>
                <br>
                TaxId:
                <input type="text" placeholder="123.456.789-00" id="email" bind:value={createAccountFields.taxId}>
                <br>
                <p hidden={!createAccountMessage} style="color:red">{createAccountMessage}</p>
                <button>Create</button>
            </form>
            <div style="cursor: pointer; color: #3333AA" on:click|preventDefault={() => {loginScreen = true}}>Or click here to login with an existing account.</div>
        </div>
    {/if}
</main>

<style>
    .login {
        padding: 8px;
        border-radius: 10px;
        max-width: 600px;
        margin: 2% auto;
        text-align: center;
        background: white;
        color: black;
    }
</style>