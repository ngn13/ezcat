<script>
  import { GET, getjob, deljob } from "./api.js";
  import Cardbar from "../lib/cardbar.svelte";
  import { goto } from "$app/navigation";

  export let data;
  let waiting = "";
  let success = "";
  let error = "";

  async function checkjob(id) {
    const res = await getjob(fetch, id);

    if (res["active"]) {
      waiting = res["message"];
    } else if (!res["active"] && res["success"]) {
      success = res["message"];
      return await deljob(fetch, id);
    } else if (!res["active"] && !res["success"]) {
      error = res["message"];
      return await deljob(fetch, id);
    }

    setTimeout(async () => {
      await checkjob(id);
    }, 2000);
  }

  async function run() {
    await goto(`/run/${data.id}`);
  }

  async function kill() {
    const res = await GET(fetch, `user/agent/kill?id=${data.id}`, true);

    if (res["error"] != undefined) {
      error = res["error"];
      return;
    }

    await checkjob(res["job"]);
    setTimeout(async () => {
      await checkjob(res["job"]);
    }, 2000);
  }
</script>

<main>
  <div class="content">
    <div class="details">
      <h1>[{data.username}@{data.hostname}]</h1>
      <div class="detail">
        <h3>System</h3>
        <p>{data.kernel}</p>
      </div>
      <div class="detail">
        <h3>Address</h3>
        <p>{data.ip}</p>
      </div>
      <div class="detail">
        <h3>PID</h3>
        <p>{data.pid}</p>
      </div>
    </div>
    <div class="buttons">
      <button
        on:click={async () => {
          await run();
        }}>🚀 run</button
      >
      {#if error != ""}
        <button
          class="button-error"
          on:click={async () => {
            await kill();
          }}>☠️ kill</button
        >
      {:else if success != ""}
        <button
          class="button-success"
          on:click={async () => {
            await kill();
          }}>☠️ kill</button
        >
      {:else if waiting != ""}
        <button
          class="button-waiting"
          on:click={async () => {
            await kill();
          }}>☠️ kill</button
        >
      {:else}
        <button
          on:click={async () => {
            await kill();
          }}>☠️ kill</button
        >
      {/if}
    </div>
  </div>
  <Cardbar bind:success bind:error bind:waiting></Cardbar>
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    gap: 0;

    border: solid 1px var(--white-5);
    border-radius: var(--radius);
  }

  .content {
    background: var(--black-1);
    color: var(--white-2);
    padding: 18px;

    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }

  main h1 {
    color: var(--white-1);
    font-weight: 900;
    margin-bottom: 8px;
  }

  .details {
    display: flex;
    flex-direction: column;
    gap: 5px;
  }

  .detail {
    display: flex;
    flex-direction: row;
    gap: 7px;
  }

  .detail h3 {
    font-size: 15px;
    color: var(--white-2);
    font-weight: 900;
  }

  .detail h3::after {
    content: ":";
  }

  .detail p {
    font-size: 15px;
    color: var(--white-1);
  }

  .buttons {
    display: flex;
    flex-direction: column;
    gap: 7px;

    align-items: end;
    justify-content: center;
  }

  .buttons button {
    width: 120px;
  }

  .button-waiting {
    color: var(--yellow-2);
    background: var(--yellow-1);
    border: solid 1px var(--white-4);
  }

  .button-success {
    color: var(--green-2);
    background: var(--green-1);
    border: solid 1px var(--white-4);
  }

  .button-error {
    color: var(--red-2);
    background: var(--red-1);
    border: solid 1px var(--white-4);
  }
</style>
