import { GET } from "../../../lib/api.js";

export async function load({ fetch, params }) {
  const address = await GET(fetch, "user/payload/addr?port=4000", true);

  if (address == undefined) {
    return {};
  }

  return {
    address: address["address"],
    id: params.id,
  };
}
