import { GET } from "../../lib/api.js";

export async function load({ fetch, params }) {
  const payload = await GET(fetch, "user/payload/list", true);
  const address = await GET(fetch, "user/payload/addr", true);

  if (payload == undefined || address == undefined) {
    return {};
  }

  return {
    payloads: payload["list"],
    address: address["address"],
  };
}
