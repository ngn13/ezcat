<svelte:head>
  <title>ezcat | CRASH!!</title>
  <meta content="ezcat | CRASH!!" property="og:title" />
</svelte:head>

<main>
  <div>
    <h1>😿 Something went wrong!</h1>
    <p>
      This may be a bug, or an issue with your setup. Either way, if you don't know how to fix it, <a
        href="https://github.com/ngn13/ezcat/issues/new">create an issue</a
      >. Here are some extra details:
    </p>
  </div>
  <pre><slot></slot></pre>
</main>

<style>
  main {
    position: absolute;
    transform: translate(-50%, -50%);
    left: 50%;
    top: 50%;

    box-shadow: var(--box-shadow);
    text-align: center;
  }

  div {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    text-align: left;
    color: var(--white-2);
    padding: 18px;
    gap: 15px;

    background: var(--black-1);
    border-radius: var(--radius) var(--radius) 0 0;
    border-bottom: solid 1px var(--red-2);
  }

  h1 {
    font-weight: 900;
    font-size: 25px;
    color: var(--white-1);
  }

  p {
    font-size: 18px;
  }

  a {
    color: var(--white-1);
  }

  pre {
    font-size: 16px;
    font-style: italic;
    color: var(--red-2);
    padding: 18px;
    background: var(--red-1);
    border-radius: 0 0 var(--radius) var(--radius);
  }
</style>
