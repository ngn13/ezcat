<script>
  export let waiting;
  export let success;
  export let error;

  function hide() {
    waiting = "";
    success = "";
    error = "";
  }
</script>

{#if error != "" && error != undefined}
  <button class="error" on:click={hide}>
    {error}
  </button>
{:else if success != "" && success != undefined}
  <button class="success" on:click={hide}>
    {success}
  </button>
{:else if waiting != "" && waiting != undefined}
  <button class="waiting" on:click={hide}>
    {waiting}
  </button>
{/if}

<style>
  button {
    margin: 0;
    border: 0;
    outline: none;
    background: none;
    width: 100%;

    padding: 18px;
    font-size: 15px;
    text-align: center;
    border-radius: 0 0 var(--radius) var(--radius);

    font-family: monospace;
  }

  .error {
    color: var(--red-2);
    background: var(--red-1);
    border-top: solid 1px var(--red-2);
  }

  .success {
    color: var(--green-2);
    background: var(--green-1);
    border-top: solid 1px var(--green-2);
  }

  .waiting {
    color: var(--yellow-2);
    background: var(--yellow-1);
    border-top: solid 1px var(--yellow-2);
  }
</style>
