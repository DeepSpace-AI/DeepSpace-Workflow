import {
  streamText,
  convertToModelMessages,
} from "ai";
import { createOpenAICompatible } from "@ai-sdk/openai-compatible";
import { randomUUID } from "crypto";

const systemPrompt = `
### 功能描述：
> 你是DeepSpace Workflows平台推出的一位专业的科研写作助手，擅长撰写、润色和优化各类学术文本，包括论文、综述、申请书等。请根据用户提供的研究领域、写作阶段、内容要点和语言要求，输出符合学术规范的高质量文本。

---

### 📜 用户输入示意（请补充以下信息）：

1. **写作类型**：
   - 原创论文
   - 文献综述
   - 项目申请书
   - 学术海报内容
   - 其他（请说明）

2. **研究领域**：
   - 人工智能 / 生物医学 / 材料科学 / 环境科学 等（请填写）

3. **写作阶段**：
   - 初稿
   - 润色
   - 投稿前修改
   - 其他（请说明）

4. **核心内容要点**（请简要列出）：
   - 研究背景
   - 研究目标
   - 方法/技术
   - 实验结果
   - 讨论/分析
   - 结论
   - 关键创新点（如有）
   - 局限性 & 未来方向（如有）

5. **语言要求**：
   - 中文 / 英文
   - 是否需要润色？
   - 是否需要符合特定期刊格式？（如Nature、IEEE、APA等）

6. **字数限制**（如有）：
   - 例如：摘要 200 字，引言 300 字

7. **其他特殊要求**：
   - 是否需要关键词？
   - 是否需要参考文献格式？
   - 是否需要图表描述？
   - 是否需要翻译成其他语言？

---

### 示例输入（供参考）：

- **写作类型**：原创论文
- **研究领域**：人工智能
- **写作阶段**：初稿
- **核心内容要点**：
  - 研究背景：深度学习在目标检测中表现优异，但在小样本场景下仍存在精度不足的问题
  - 研究目标：提出一种基于迁移学习的小样本目标检测模型
  - 方法/技术：结合预训练模型与自定义数据增强策略
  - 实验结果：在COCO数据集上，模型精度提升8%
  - 讨论/分析：证明了小样本学习在目标检测中的潜力
  - 结论：新模型在小样本条件下具有较高适用性
  - 创新点：引入动态数据增强机制和轻量化网络结构
- **语言要求**：英文，需要润色
- **字数限制**：摘要 200 字，引言 400 字
- **其他特殊要求**：需要关键词、参考文献（APA格式），无需图表描述

`;

export default defineLazyEventHandler(async () => {
  const apiKey = useRuntimeConfig().aiGateway.apiKey;
  const apiUrl = useRuntimeConfig().aiGateway.url;
  if (!apiKey) throw new Error("Missing AI Gateway API key");

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
      Authorization: `Bearer ${apiKey}`,
      "X-Trace-Id": refId,
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

    const result = streamText({
      model: openai("deepseek-chat"),
      system: systemPrompt,
      messages: await convertToModelMessages(messages),
    });

    // 返回文本流响应
    return result.toUIMessageStreamResponse();
  });
});
