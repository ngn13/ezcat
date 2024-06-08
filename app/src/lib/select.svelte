<script>
  export let title;
  export let value;
  export let name;
  export let on_change;
</script>

<div>
  <span>{title}</span>
  {#if on_change}
    <select
      on:change={(e) => {
        on_change(e);
      }}
      bind:value
      {name}
    >
      <slot></slot>
    </select>
  {:else}
    <select bind:value {name}>
      <slot></slot>
    </select>
  {/if}
</div>

<style>
  div {
    width: 100%;
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
