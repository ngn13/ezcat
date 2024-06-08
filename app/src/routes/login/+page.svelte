<script>
  import { tokenGet, tokenSet } from "../../lib/token.js";
  import Cardbar from "../../lib/cardbar.svelte";
  import Input from "../../lib/input.svelte";
  import Crash from "../../lib/crash.svelte";
  import { goto } from "$app/navigation";
  import { PUT } from "../../lib/api.js";

  if (tokenGet() != null) {
    goto("/");
  }

  export let data;
  let password = "";
  let error = "";

  async function login(e) {
    e.preventDefault();

    if (password == "") {
      error = "Password is required";
      return;
    }

    const res = await PUT(fetch, "login", {
      password: password,
    });

    if (res["error"] != undefined) {
      error = res["error"];
      return;
    }

    error = "";

    tokenSet(res["token"]);
    await goto("/");
  }
</script>

<svelte:head>
  <title>ezcat | login</title>
  <meta content="ezcat | login" property="og:title" />
</svelte:head>

{#if data.version}
  <main>
    <div class="header">
      <div class="title">
        <a class="logo" href="https://github.com/ngn13/ezcat"><h1>ezcat</h1></a>
        <p>v{data.version}</p>
      </div>
      <h1>🐱</h1>
    </div>
    <form
      on:submit={async (e) => {
        await login(e);
      }}
    >
      <Input
        bind:value={password}
        name="password"
        password={true}
        holder="**********"
        required={true}
      >
        Password
      </Input>
      <button type="submit">🔓 login</button>
    </form>
    <Cardbar bind:error></Cardbar>
  </main>
{:else}
  <Crash>{data.error}</Crash>
{/if}

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

    text-align: center;
  }

  .header {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    padding: 0 2px 6px 2px;
    border-bottom: solid 1px var(--white-4);
    padding: 18px;
  }

  .header h1 {
    font-size: 30px;
  }

  .title {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    gap: 8px;
  }

  .logo {
    color: var(--white-1);
    text-decoration: none;
    font-size: 16px;
  }

  .title p {
    color: var(--white-2);
    margin-top: 15px;
    font-size: 14px;
  }

  form {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 18px;
    gap: 7px;
  }

  button {
    width: 100%;
  }
</style>
