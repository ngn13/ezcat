<script>
  import { PUT, getjob, deljob } from "../../../lib/api.js";
  import Cardbar from "../../../lib/cardbar.svelte";
  import Input from "../../../lib/input.svelte";
  import { goto } from "$app/navigation";

  export let data;

  let success = "";
  let waiting = "";
  let error = "";
  let host = data["address"];

  async function checkjob(id) {
    const res = await getjob(fetch, id);

    if (res["waiting"]) {
      waiting = "Waiting for response";
    } else if (!res["waiting"] && res["success"]) {
      success = `Got a successful response: ${res["message"]}`;
    } else if (!res["waiting"] && !res["success"]) {
      error = `Got a failure response: ${res["message"]}`;
    }

    if (!res["waiting"]) {
      setTimeout(async () => {
        await goto("/");
      }, 2000);
      return await deljob(fetch, id);
    }

    setTimeout(async () => {
      await checkjob(id);
    }, 2000);
  }

  async function run(e) {
    e.preventDefault();
    const res = await PUT(
      fetch,
      "user/agent/run",
      {
        address: host,
        session: data.session,
      },
      true
    );

    if (res["error"] != undefined) {
      error = `Request failed: ${res["error"]}`;
      return;
    }

    await checkjob(res["job"]);
    setTimeout(async () => {
      await checkjob(res["job"]);
    }, 2000);
  }
</script>

<svelte:head>
  <title>ezcat | run</title>
  <meta content="ezcat | run" property="og:title" />
</svelte:head>

<main>
  <form
    on:submit={async (e) => {
      await run(e);
    }}
  >
    <Input bind:value={host} name="host" holder="IP:port">Host address</Input>
    <button type="submit">🚀 run</button>
  </form>
  <Cardbar bind:success bind:error bind:waiting></Cardbar>
</main>

<style>
  main {
    position: absolute;
    transform: translate(-50%, -50%);
    left: 50%;
    top: 50%;

    background: var(--black-1);
    box-shadow: var(--box-shadow);
    border-radius: var(--radius);
    border: solid 1px var(--white-5);
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 18px;
  }
</style>
