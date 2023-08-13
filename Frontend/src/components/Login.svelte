<script>

    import jwt_decode from 'jwt-decode';
    import { User } from "../stores/User.js"
    import { backendUrl } from "../utils/backend.js"
    import { onDestroy } from 'svelte';

    let subscriptions = [];

    let fields = {
        email: '',
        password: '',
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

    const submitHandler = async () => {
        loginError = '';
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
    };

    onDestroy(() => {
        for (const prop in subscriptions) {
            subscriptions[prop]();
        }
    });

</script>

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
</div>

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