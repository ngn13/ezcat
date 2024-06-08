<script>
  import { onDestroy, onMount } from "svelte";
  import Navbar from "../lib/navbar.svelte";
  import Agent from "../lib/agent.svelte";
  import Crash from "../lib/crash.svelte";
  import { GET } from "../lib/api.js";

  let megamind = true;
  let agents = [];
  let error = "";
  let interval;

  onMount(async () => {
    try {
      const info_data = await GET(fetch, "info");
      megamind = info_data["megamind"];
    } catch (err) {
      error = `Failed to fetch info: ${err}`;
      return;
    }

    await update_agents();
    interval = setInterval(async () => {
      try {
        await update_agents();
      } catch (err) {
        error = `Failed to update agent list: ${err}`;
        return;
      }
    }, 5000);
  });

  onDestroy(() => {
    clearInterval(interval);
  });

  async function update_agents() {
    const agents_data = await GET(fetch, "user/agent/list", true);
    agents = agents_data["list"];
  }
</script>

<svelte:head>
  <title>ezcat | dashboard</title>
  <meta content="ezcat | dashboard" property="og:title" />
</svelte:head>

{#if error == ""}
  <Navbar></Navbar>
  <main>
    {#if agents != null}
      <div class="agents">
        {#each agents as agent}
          <Agent class="agent" data={agent}></Agent>
        {/each}
      </div>
    {:else if agents == null && megamind}
      <div class="megamind">
        <p>no shells?</p>
        <img src="megamind.png" alt="no shells?" />
      </div>
    {/if}
  </main>
{:else}
  <Crash>{error}</Crash>
{/if}

<style>
  main {
    padding: 30px 10%;
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .agents {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .megamind {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
  }

  .megamind img {
    width: 30%;
    opacity: 50%;
    border-radius: var(--radius);
  }

  .megamind p {
    font-size: 30px;
    color: var(--white-1);
    opacity: 50%;
  }
</style>
