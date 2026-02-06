import {
  streamText,
  convertToModelMessages,
} from "ai";
import { createOpenAICompatible } from "@ai-sdk/openai-compatible";
import { randomUUID } from "crypto";
import { getGatewayBase } from "../../utils/gateway";

const defaultSystemPrompt = "你是一个可靠、简洁的通用 AI 助手，尽量给出清晰、可执行的回答。";

export default defineLazyEventHandler(async () => {
  const apiUrl = useRuntimeConfig().aiGateway.url;
  if (!apiUrl) throw new Error("Missing AI Gateway URL");

  return defineEventHandler(async (event: any) => {
    let body;
    try {
      body = await readBody(event);
    } catch (e) {
      throw createError({ statusCode: 400, statusMessage: "Invalid body" });
    }

    const { messages } = body || {};

    if (!messages || !Array.isArray(messages)) {
      throw createError({ statusCode: 400, statusMessage: "Messages required" });
    }

    // Preflight billing check to ensure proper 402 propagation when balance is insufficient.
    try {
      const base = getGatewayBase(apiUrl);
      const walletResp = await $fetch<{ wallet?: any }>(`${base}/api/billing/wallet`, {
        headers: { cookie: event.node.req.headers.cookie || "" },
      });
      const wallet = walletResp?.wallet ?? {};
      const balance =
        typeof wallet?.balance === "number"
          ? wallet.balance
          : typeof wallet?.Balance === "number"
            ? wallet.Balance
            : 0;
      if (balance <= 0) {
        throw createError({ statusCode: 402, statusMessage: "Payment Required" });
      }
    } catch (err: any) {
      if (err?.statusCode === 402) {
        throw err;
      }
      // If the check fails for any other reason, fall back to gateway enforcement.
    }

    const headers = event.node?.req?.headers ?? {};
    const headerAmount = headers["x-billing-amount"];
    const headerRef = headers["x-billing-ref-id"];
    const headerTrace = headers["x-trace-id"];

    const rawAmount = body?.billingAmount ?? headerAmount;
    const amount = rawAmount !== undefined ? Number(rawAmount) : 0;
    const hasAmount = Number.isFinite(amount) && amount > 0;

    const refId =
      (typeof body?.billingRefId === "string" && body.billingRefId) ||
      (typeof headerRef === "string" && headerRef) ||
      (typeof headerTrace === "string" && headerTrace) ||
      randomUUID();

    const extraHeaders: Record<string, string> = {
      "X-Trace-Id": refId,
      cookie: event.node.req.headers.cookie || "",
    };

    if (hasAmount) {
      extraHeaders["X-Billing-Amount"] = String(amount);
      extraHeaders["X-Billing-Ref-Id"] = refId;
    }

    const openai = createOpenAICompatible({
      name: "newapi",
      baseURL: apiUrl.endsWith("/v1") ? apiUrl : `${apiUrl}/v1`,
      headers: extraHeaders,
    });

    const modelName =
      typeof body?.model === "string" && body.model.trim()
        ? body.model.trim()
        : "deepseek-chat";
    const system =
      typeof body?.system === "string" && body.system.trim()
        ? body.system.trim()
        : defaultSystemPrompt;

    const maxTokens =
      typeof body?.max_tokens === "number" ? body.max_tokens : undefined;
    const temperature =
      typeof body?.temperature === "number" ? body.temperature : undefined;
    const topP = typeof body?.top_p === "number" ? body.top_p : undefined;
    const reasoningEffort =
      body?.reasoning_effort === "low" ||
      body?.reasoning_effort === "medium" ||
      body?.reasoning_effort === "high"
        ? body.reasoning_effort
        : undefined;

    const result = streamText({
      model: openai(modelName),
      system,
      maxTokens,
      temperature,
      topP,
      providerOptions: reasoningEffort
        ? {
            newapi: {
              reasoning_effort: reasoningEffort,
            },
          }
        : undefined,
      messages: await convertToModelMessages(messages),
    });

    // 返回文本流响应
    return result.toUIMessageStreamResponse();
  });
});
