<script>
  import { PUT } from "../../lib/api.js";
  import Cardbar from "../../lib/cardbar.svelte";
  import Select from "../../lib/select.svelte";
  import Input from "../../lib/input.svelte";
  import { goto } from "$app/navigation";

  export let data;

  let selected_host = data["address"];
  let selected_type = "";
  let selected_os = "";

  let payloads = data["payloads"];
  let type_list = [];
  let os_list = [];

  let waiting = "";
  let success = "";
  let error = "";

  function os_code(os) {
    return os["name"].toLowerCase() + "_" + os["arch"].toLowerCase();
  }

  function has_os(code) {
    for (let i = 0; i < os_list.length; i++) {
      if (os_list[i]["code"] == code) {
        return true;
      }
    }
    return false;
  }

  function get_types() {
    type_list = [];
    payloads.forEach((p) => {
      p["os"].forEach((os) => {
        if (os_code(os) == selected_os) {
          type_list.push(p["name"]);
        }
      });
    });
    selected_type = type_list[0];
  }

  function change(e) {
    selected_os = e.target.value;
    get_types();
  }

  async function generate(e) {
    e.preventDefault();
    waiting = "Waiting for the build"

    const res = await PUT(
      fetch,
      "user/payload/build",
      {
        address: selected_host,
        type: selected_type,
        os: selected_os,
      },
      true
    );

    if (res["error"] != undefined) {
      error = res["error"];
    }

    try {
      navigator.clipboard.writeText(res["payload"]);
    } catch (error) {
      error = "Failed to copy payload to the clipboard";
    }
    success = "Copied payload to the clipboard";

    setTimeout(async () => {
      await goto("/");
    }, 2500);
  }

  payloads.forEach((p) => {
    p["os"].forEach((os) => {
      let code = os_code(os);
      if (has_os(code)) return;
      os_list.push({
        name: `${os["name"]} (${os["arch"]})`,
        code: code,
      });
    });
  });

  if (os_list.length > 0) selected_os = os_list[0]["code"];
  get_types();
</script>

<svelte:head>
  <title>ezcat | generate</title>
  <meta content="ezcat | generate" property="og:title" />
</svelte:head>

<main>
  <h1>⚙️ generate payload</h1>
  <form
    on:submit={async (e) => {
      await generate(e);
    }}
  >
    <Input bind:value={selected_host} name="host" holder="IP:port">Host address</Input>
    <div class="options">
      <Select on_change={change} bind:value={selected_os} name="os" title="Operating system">
        {#each os_list as os}
          <option value={os.code}>{os.name}</option>
        {/each}
      </Select>
      <Select bind:value={selected_type} name="type" title="Type">
        {#each type_list as type}
          <option value={type}>{type}</option>
        {/each}
      </Select>
    </div>
    <button type="submit">📋 generate & copy</button>
  </form>
  <Cardbar bind:success bind:error bind:waiting></Cardbar>
</main>

<style>
  main {
    position: absolute;
    transform: translate(-50%, -50%);
    left: 50%;
    top: 50%;

    width: 50%;

    background: var(--black-1);
    box-shadow: var(--box-shadow);
    border-radius: var(--radius);
    border: solid 1px var(--white-5);
  }

  main h1 {
    color: var(--white-1);
    padding: 15px;
    border-bottom: solid 1px var(--white-4);
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 18px;
  }

  .options {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    gap: 7px;
  }

  main button {
    width: 100%;
  }
</style>
