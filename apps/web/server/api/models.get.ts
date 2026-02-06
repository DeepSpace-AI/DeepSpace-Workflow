import { getGatewayBase } from "../utils/gateway";

export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const base = getGatewayBase(aiGateway.url);
  try {
    return await $fetch(`${base}/v1/models`, {
      headers: {
        cookie: event.node.req.headers.cookie || "",
      },
    });
  } catch (error: any) {
    const status =
      error?.statusCode || error?.status || error?.response?.status;
    if (status === 402) {
      return {
        data: [
          { id: "deepseek-chat", object: "model", owned_by: "fallback" },
          { id: "gpt-4.1", object: "model", owned_by: "fallback" },
          { id: "claude-3.5", object: "model", owned_by: "fallback" },
        ],
      };
    }
    throw error;
  }
});
