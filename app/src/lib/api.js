import { tokenClear, tokenGet } from "./token.js";
import { goto } from "$app/navigation";

function pathjoin(p1, p2) {
  if (p1.endsWith("/") && p2.startsWith("/")) return p1.substring(0, p1.length - 1) + p2;
  else if (p1.endsWith("/") && !p2.startsWith("/")) return p1 + p2;
  else if (!p1.endsWith("/") && p2.startsWith("/")) return p1 + p2;
  else return p1 + "/" + p2;
}

function geturl(path) {
  return new URL(pathjoin("api", path), import.meta.env.VITE_API_URL_DEV).href;
}

async function getjob(fetch, id) {
  return GET(fetch, `/user/job/get?id=${id}`, true);
}

async function deljob(fetch, id) {
  return DELETE(fetch, `/user/job/del?id=${id}`, true, false);
}

async function GET(fetch, path, authreq = false, isjson = true) {
  const token = tokenGet();

  if (authreq && token == null) {
    await goto("/login");
    return;
  } else if (token == null) {
    const res = await fetch(geturl(path));
    return await res.json();
  }

  const res = await fetch(geturl(path), {
    headers: {
      Authorization: token,
    },
  });

  if (res.status == 401) {
    tokenClear();
    await goto("/login");
    return;
  }

  if (isjson) return await res.json();
  else return res;
}

async function DELETE(fetch, path, authreq = false, isjson = true) {
  const token = tokenGet();

  if (authreq && token == null) {
    await goto("/login");
    return;
  } else if (token == null) {
    const res = await fetch(geturl(path));
    return await res.json();
  }

  const res = await fetch(geturl(path), {
    method: "DELETE",
    headers: {
      Authorization: token,
    },
  });

  if (res.status == 401) {
    tokenClear();
    await goto("/login");
    return;
  }

  if (isjson) return await res.json();
  else return res;
}

async function PUT(fetch, path, json, authreq = false, isjson = true) {
  const token = tokenGet();

  if (authreq && token == null) {
    await goto("/login");
    return;
  } else if (token == null) {
    const res = await fetch(geturl(path), {
      method: "PUT",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(json),
    });
    return await res.json();
  }

  const res = await fetch(geturl(path), {
    method: "PUT",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
      Authorization: token,
    },
    body: JSON.stringify(json),
  });

  if (res.status == 401) {
    tokenClear();
    await goto("/login");
    return;
  }

  if (isjson) return await res.json();
  else return res;
}

export { geturl, GET, PUT, getjob, deljob };
