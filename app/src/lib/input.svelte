<script>
  export let holder;
  export let value;
  export let name;

  export let password;
  export let required;
</script>

<div>
  <span><slot></slot></span>
  {#if required && password}
    <input bind:value {name} type="password" placeholder={holder} required />
  {:else if !required && password}
    <input bind:value {name} type="password" placeholder={holder} />
  {:else if required && !password}
    <input bind:value {name} type="text" placeholder={holder} required />
  {:else}
    <input bind:value {name} type="text" placeholder={holder} />
  {/if}
</div>

<style>
  div {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  div span {
    color: var(--white-2);
    text-align: left;
    margin-left: 3px;
    font-size: 15px;
  }
</style>
