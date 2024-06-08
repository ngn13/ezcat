import { GET } from "../../lib/api.js";

export async function load({ fetch, params }) {
  try {
    return await GET(fetch, "info");
  } catch (error) {
    return {
      error: `Failed to connect to the API (${error})`,
    };
  }
}
