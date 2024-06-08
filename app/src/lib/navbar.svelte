<script>
  import { tokenClear } from "./token.js";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { GET } from "./api.js";

  let version;

  onMount(async () => {
    const info_data = await GET(fetch, "info");
    version = info_data["version"];
  });

  async function logout() {
    await GET(fetch, "user/logout");
    tokenClear();
    await goto("/login");
  }

  async function generate() {
    await goto("/generate");
  }
</script>

<navbar>
  <div class="header">
    <h1>🐱</h1>
    <a href="https://github.com/ngn13/ezcat"><h1>ezcat</h1></a>
    {#if version}
      <p>v{version}</p>
    {/if}
  </div>
  <div class="buttons">
    <button on:click={async () => generate()}>⚙️ generate</button>
    <button on:click={async () => logout()}>🔒 logout</button>
  </div>
</navbar>

<style>
  navbar {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    padding: 15px 20px;

    box-shadow: var(--box-shadow);
    background: var(--black-1);
    border-bottom: solid 1px var(--white-5);
  }

  navbar div {
    display: flex;
    flex-direction: row;
    align-items: end;
    justify-content: center;
    gap: 8px;

    font-size: 14px;
    color: var(--white-2);
  }

  .header h1 {
    text-shadow: var(--text-shadow);
  }

  .header a {
    color: var(--white-1);
    text-decoration: none;
  }

  .buttons {
    display: flex;
    flex-direction: row;
    gap: 10px;
  }
</style>
