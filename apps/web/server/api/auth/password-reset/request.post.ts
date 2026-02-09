import { getGatewayBase } from "#server/utils/gateway";

export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const body = await readBody(event);
  const base = getGatewayBase(aiGateway.url);

  const res = await fetch(`${base}/api/auth/password-reset/request`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    const msg =
      typeof data?.error === "string"
        ? data.error
        : data?.error?.message || data?.message || "Password reset request failed";
    throw createError({ statusCode: res.status, statusMessage: msg });
  }

  return data;
});
