import { getGatewayBase } from "#server/utils/gateway";

export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const body = await readBody(event);
  const base = getGatewayBase(aiGateway.url);

  const res = await fetch(`${base}/api/users/me/password`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      cookie: event.node.req.headers.cookie || "",
    },
    body: JSON.stringify(body ?? {}),
  });

  const data = await res.json();
  if (!res.ok) {
    const msg =
      typeof data?.error === "string"
        ? data.error
        : data?.error?.message || "Password change failed";
    throw createError({ statusCode: res.status, statusMessage: msg });
  }

  return data;
});
