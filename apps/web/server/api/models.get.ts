import { getGatewayBase } from "../utils/gateway";

export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const base = getGatewayBase(aiGateway.url);
  try {
    const response = await $fetch<any>(`${base}/api/models`, {
      headers: {
        cookie: event.node.req.headers.cookie || "",
      },
    });
    const items = Array.isArray(response?.items)
      ? response.items
      : Array.isArray(response?.data)
        ? response.data
        : [];
    return {
      ...response,
      items,
      data: items,
    };
  } catch (error: any) {
    const status =
      error?.statusCode || error?.status || error?.response?.status;
    if (status === 402) {
      const fallbackItems = [
        { id: "deepseek-chat", name: "deepseek-chat", object: "model", owned_by: "fallback" },
        { id: "gpt-4.1", name: "gpt-4.1", object: "model", owned_by: "fallback" },
        { id: "claude-3.5", name: "claude-3.5", object: "model", owned_by: "fallback" },
      ];
      return {
        items: fallbackItems,
        data: fallbackItems,
      };
    }
    throw error;
  }
});
