<script>
    import { createEventDispatcher } from "svelte";
    
    const dispatch = createEventDispatcher();
    export let tabs;
    export let currentTab;
</script>

<div class="tabs">
    <ul>
        {#each tabs as tab}
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <li on:click={() => dispatch('tabChange', tab)}>
                <div class:active={tab === currentTab}>{tab}</div>
            </li>
        {/each}
    </ul>
</div>

<style>
    .tabs {
        margin-bottom: 40px;
        line-height: 1.5;
    }
    ul {
        display: flex;
        justify-content: center;
        padding: 0;
        list-style-type: none;
    }
    li {
        margin: 0 16px;
        font-size: 18px;
        color: #555;
        cursor: pointer;
    }
    .active {
        color: #d91b42;
        border-bottom: 2px solid;
        padding-bottom: 8px;
    }
</style>