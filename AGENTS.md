1. 在编写任何代码之前，请先描述你的方案并等待批准。如果需求不明确，在编写任何代码之前务必提出澄清问题。
2. 如果一项任务需要修改超过3个文件，请先停下来，将其分解成更小的任务。
2. 每次我纠正你之后，就在AGENTS.md文件中添加一条新规则，这样就不会再发生这种情况了。
3. 开发前请根据所需要的库和组件通过 context7 来读取相关文档，确保你理解了它们的工作原理和使用方法。
4. 项目首选为中文，所有的展示文案和代码注释都必须是中文。
5. 变量命名必须使用英文，禁止使用中文变量名。
6. 侧边栏用户信息按钮必须提供弹出菜单，包含个人信息与退出登录入口。
7. Nuxt项目中不用再重复引用 Vue 相关的库，如 vue, vue-router, vuex 等。
8. Gateway端接口文档位于 /docs/gateway-openapi.yaml 中。
9. 用户确认方案后直接实施修改并完成验证，不重复请求确认。
10. 用户提供官方示例纠正时，优先按示例对齐实现并说明原因。
11. 新增或调整样式时必须兼容暗色模式。
12. Gateway 接口文档默认以 /services/gateway/cmd/gateway/docs/swagger.json 为准，除非用户明确指定其他文档路径。
13. 对话模型下拉的 menuItems 不要改为 computed 响应式，保持 ref 同步，避免 Nuxt UI SelectMenu 渲染异常。
14. 对话模型下拉的 menuItems 必须是纯数组值，禁止使用 ref/computed 或其他 Vue 包装。
15. /api/users/me 的字段解析必须兼容 PascalCase（如 DisplayName/FullName/AvatarURL）与 snake_case 两种格式。
15. 新增风控策略必须提供类型选择并创建对应子规则（IP 规则/速率限制/预算上限）。
