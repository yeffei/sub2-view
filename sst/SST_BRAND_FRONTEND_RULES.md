# SST / 山枢庭品牌与前台改造规则

> 状态说明：本文件已被 [SST Unified Development Guide](./SST_UNIFIED_DEVELOPMENT_GUIDE.md) 作为默认统一规范接管；本文件保留为品牌与视觉细节历史参考。

创建日期：2026-06-15  
适用项目：`d:\sub2api-main`  
适用范围：Sub2API 前台首页、登录注册页、用户控制台入口层、导航信息架构、视觉资产方向。后端网关能力、计费、调度、账号池、风控、支付、监控等核心能力默认保留。

## 1. 品牌定位

品牌名：山枢庭  
副标：SST  
英文辅助：SST Gateway  
一句话：统一入口，安静流转。

山枢庭不是常规低价中转站，也不是堆满功能截图和接口示例的 SaaS 首页。它的产品定位是：以更高倍率换取稳定、秩序、品质客户与长期可维护服务。

市场对照：许多中转站主打 `0.2` 左右倍率，常见表达是便宜、快、接口多、注册即用。山枢庭可以接受 `0.4`、`0.5` 甚至更高倍率，因此页面必须解释一种不同价值：稳定供给、清晰规则、审慎准入、低噪声运营、对重度开发者和团队更可靠。

设计结论：不要塑造“折扣 API 超市”；要塑造“有门槛、有秩序、有长期服务能力的统一入口”。

## 2. 产品叙事

首页只负责品牌入口，不负责功能演示。

首页应该传达：

- 这是一个统一入口，多个上游与能力在背后安静流转。
- 服务可靠性、规则感、运维克制比低价噱头更重要。
- 用户进入后可以看到完整功能，但首页不需要把所有功能摊开。
- 山枢庭适合重度用户、品质客户、小团队、开发者，不主动迎合薅羊毛型流量。

首页不应该传达：

- “全网最低价”。
- “注册送大量试用额度”。
- “复制这个 `/v1/chat/completions` 示例立刻开干”。
- 大量模型 logo、功能卡片、价格标签、接口参数堆叠。

推荐文案语气：短句、留白、克制、可靠。少用感叹号，少用夸张最高级。

## 3. 当前前台结构结论

技术栈：`Vue 3 + Vite + TypeScript + Pinia + vue-router + TailwindCSS`。

关键入口：

- `frontend/src/main.ts`：初始化主题、Pinia、公共配置注入、i18n、router，然后挂载应用。
- `frontend/src/App.vue`：承载 `RouterView`、导航进度、Toast、公告、管理员合规弹窗，并拉取 public settings。
- `frontend/src/router/index.ts`：路由、鉴权、后台模式、功能开关路由限制。
- `frontend/src/stores/app.ts`：公共设置缓存，包括 `site_name`、`site_logo`、`site_subtitle`、`doc_url`、`home_content`、功能开关。
- `frontend/src/views/HomeView.vue`：默认首页实现。若 `home_content` 存在，则 iframe 或 HTML 接管首页。
- `frontend/src/components/layout/AppSidebar.vue`：控制台侧边栏与功能隐藏/显示逻辑。
- `frontend/src/components/layout/AuthLayout.vue`：登录注册页外壳，目前仍是普通渐变 SaaS 风格。
- `frontend/src/utils/featureFlags.ts`：前台功能开关统一解释层。

现有首页问题：

- 视觉是常规 SaaS：渐变、圆角卡片、功能卡、终端窗口、供应商标签。
- 结构超过品牌入口需要，信息堆叠偏多。
- 当前终端示例虽是 `/v1/messages`，但仍然强化“接口演示页”感；后续不应默认展示 `/v1/chat/completions` 或任何具体请求示例。
- 供应商和功能标签容易把山枢庭拉回“低价 API 聚合站”的竞争语境。

## 4. 改造边界

必须保留：

- `home_content` 后台接管逻辑。管理员配置了自定义首页时，仍优先使用自定义内容。
- 登录态判断：未登录进入 `/login`，已登录进入用户或管理员 dashboard。
- 语言切换、主题切换、文档链接、站点 logo 设置。
- 路由、鉴权、后台模式、设置注入、公告、Toast、合规弹窗。
- Sub2API 原有核心功能：账号管理、分组、密钥、用量、计费、订单、订阅、兑换、监控、代理、风控、支付、公告、自定义页面等。

可以重做或弱化：

- `HomeView.vue` 默认首页的完整结构和样式。
- 首页上的功能卡、供应商标签、接口演示、GitHub 强曝光。
- 登录注册页的外壳视觉，使其与山枢庭品牌统一。
- 用户侧导航顺序和信息密度，但不删除实际路由。
- 控制台 dashboard 的首屏信息优先级。

谨慎处理：

- 功能隐藏应优先走已有 `public settings` / `featureFlags` / `simple mode` / 自定义菜单机制。
- 不应为了美观删除路由、删除 store、删除 API 调用或破坏后台配置能力。
- 业务功能弱化优先表现为：不在首页宣传、不在默认首屏突出、折叠到二级、根据配置开关隐藏入口。

## 5. 首页信息架构

首页目标：`1.5` 屏品牌入口。

建议结构：

第一屏：

- 顶部：极简品牌标记、`SST Gateway`、语言/主题/文档/登录入口。
- 主视觉：山枢庭、SST、一句话“统一入口，安静流转。”
- 辅助短句：围绕稳定、秩序、长期服务，不出现低价或接口示例。
- 视觉主体：纸本、水墨、细线、朱砂点睛，不复刻参考图布局。
- CTA：一个主入口“入庭”或“进入控制台/登录”，一个弱入口“查看文档”。

半屏延伸：

- 只放 3 到 4 个秩序型价值点，不做卡片堆叠。
- 示例方向：`稳态调度`、`清晰计量`、`审慎准入`、`长期维护`。
- 可以用细线分栏、印章点、短句，不展示代码块。

不做：

- 不做多区块落地页。
- 不做大面积产品截图堆叠。
- 不做模型 logo 墙。
- 不做价格表。
- 不做接口 payload 示例。

## 6. 视觉方向

关键词：东方纸本感、水墨留白、细线秩序、朱砂点睛、安静、克制、可靠。

整体气质：

- 参考图片只取风格和气质，不取原始布局。
- 页面应像一张可交互的纸本卷轴或庭院入口图，而不是仪表盘预览。
- 留白是主要资产，控件只是入口。
- 线条表达秩序，水墨表达流转，朱砂表达品牌识别。

颜色建议：

- 纸色：`#f4efe4`、`#ede5d4`、`#d8cdb9`。
- 墨色：`#1f2320`、`#38413a`、`#59645a`。
- 山石青灰：`#4e5d55`、`#6f7a70`。
- 朱砂：`#a73a2a`、`#c24b35`，只做点睛，不做大面积背景。
- 暗色模式如保留，应像夜间宣纸/墨色，不走蓝紫霓虹。

字体建议：

- 中文标题优先走宋体/明朝/衬线气质：`Noto Serif SC`、`Source Han Serif SC`、`Songti SC`、`SimSun` 等可用栈。
- 英文辅助使用经典 serif 或低调 humanist sans，不用夸张科技字体。
- 控制台内部可继续使用现有系统字体，先保证可用性。

构图建议：

- 非对称留白，主标题不一定居中。
- 用细线形成隐约网格、轴线、分隔，不要厚重边框。
- 用一个小朱砂印或圆点建立记忆点。
- 动效轻：淡入、墨迹轻微扩散、线条显现即可，不做强动效。

禁用视觉：

- 紫色/蓝色 SaaS 渐变。
- 发光圆球、bokeh、科技网格背景、玻璃拟态堆叠。
- 大量圆角卡片、卡片套卡片。
- 表情化插画、3D 小物件、廉价 logo 墙。
- 直接复制参考图的左右卡片布局。

## 7. 背景图片与视觉资产风格

首页可使用 CSS 绘制或通过 `$imagegen` 生成位图背景。当前项目只有 `frontend/src/assets/icons/` 和 `frontend/public/logo.png`，尚无品牌资产目录；后续项目内生成资产建议新建 `frontend/src/assets/brand/`。

重新分析后的判断：位图只承担“纸本质感”和“水墨气氛”，结构线、按钮、入口、文字、状态信息全部由代码实现。这样可以保证响应式、暗色模式、可维护性和加载性能，不把 UI 固化进图片。

### 7.1 推荐资产拆分

第一优先级：生成一张纸本/水墨背景底图。

- 文件建议：`frontend/src/assets/brand/sst-paper-ink-bg.webp`。
- 用途：`HomeView.vue` 首屏背景，可叠加 CSS 线条和渐变遮罩。
- 画面：宣纸纹理、淡墨远山、低对比、大片留白、无文字、无 UI 面板。
- 尺寸建议：`2400x1600` 或接近 3:2，便于桌面横向裁切；移动端通过 `background-position` 保留留白。

第二优先级：生成一张轻量水墨山形或墨迹层。

- 文件建议：`frontend/src/assets/brand/sst-ink-mountain.webp`。
- 用途：作为首页下方半屏的淡墨延伸，或 CSS mask 的背景源。
- 画面：横向山影、透明感弱也可以，低对比，不承载文字。

第三优先级：朱砂印记优先用代码或 SVG，不优先用位图。

- 文件建议：`frontend/src/assets/brand/sst-seal.svg`。
- 用途：品牌记忆点、标题旁小印、按钮 hover 点睛。
- 原因：印记需要清晰缩放和颜色可控，SVG 比生成位图更可维护。

不建议生成：

- 带文字的 hero 图。图片生成文字容易出错，且无法国际化。
- 带 UI 面板或代码窗口的背景图。会把首页拉回 SaaS 模板。
- 大型 mockup 组合图。会削弱“品牌入口”定位。
- 复杂透明 PNG。除非确实需要，否则不走透明图流程。

### 7.2 `$imagegen` 使用规则

使用 `$imagegen` 时遵循内置工具优先：

- 普通背景、纹理、插画：使用内置 image generation。
- 生成结果若用于项目，必须移动或复制到 `frontend/src/assets/brand/`，不能只留在 `$CODEX_HOME/generated_images/`。
- 不覆盖已有资产，除非明确要求；新版本使用 `-v2`、`-draft`、`-final` 等命名。
- 生成后必须人工/工具检查：是否有乱码文字、是否出现 UI 面板、是否太像 stock、是否色彩过重、是否违背“安静克制”。
- 透明背景需求默认先走 chroma-key + 本地去背；只有复杂透明需求才询问是否切 CLI fallback。

推荐生成 prompt 草案：

```text
Use case: stylized-concept
Asset type: landing page background texture for SST Gateway / 山枢庭
Primary request: an elegant Eastern paper-and-ink background for a premium AI gateway brand
Scene/backdrop: warm handmade xuan paper texture with very faint ink-wash distant mountains and subtle archival fibers
Subject: no foreground object, no people, no interface, only atmospheric paper texture and ink landscape traces
Style/medium: restrained Chinese ink wash, editorial brand background, premium and quiet
Composition/framing: wide landscape composition with large empty negative space for web typography, visual weight slightly toward the lower corners
Lighting/mood: calm, reliable, understated, museum-paper softness
Color palette: warm rice paper, muted ink green-gray, tiny cinnabar accent allowed but no large red areas
Materials/textures: handmade paper grain, faint brush diffusion, fine line order implied but not dominant
Text (verbatim): none
Constraints: no readable text, no UI panels, no code, no logos, no people, no buildings as main subject, no stock-photo realism
Avoid: SaaS gradients, purple/blue glow, neon, glassmorphism, bokeh orbs, cyberpunk, busy composition, high contrast
```

如生成第二张山形层，可用：

```text
Use case: stylized-concept
Asset type: subtle lower-section ink mountain layer for a web landing page
Primary request: a quiet horizontal ink-wash mountain silhouette for SST Gateway / 山枢庭
Scene/backdrop: pale transparent-feeling paper, distant mountain ridge suggested by soft ink diffusion
Subject: low horizontal mountain forms only, no temple, no people, no UI
Style/medium: Chinese ink wash, restrained, low contrast, premium editorial
Composition/framing: extra-wide horizontal band, visual weight along the bottom edge, large blank top area
Color palette: muted ink gray, green-gray, warm paper, no saturated color except an optional tiny cinnabar fleck
Text (verbatim): none
Constraints: no text, no logo, no interface, no sharp illustration outlines, no photo realism
Avoid: decorative fantasy landscape, stock wallpaper, dramatic sunset, neon, blue-purple gradient
```

### 7.3 代码实现策略

首页背景建议分层：

1. 页面底色：CSS 纸色变量。
2. 生成位图：`sst-paper-ink-bg.webp`，低透明或 normal blend。
3. CSS 细线：`linear-gradient` 网格、轴线、圆弧、分割线。
4. CSS 墨迹：局部 `radial-gradient` 或 mask，不做发光 orb。
5. SVG/CSS 朱砂印：小而清晰。
6. 实际文字和按钮：HTML/CSS，保证可访问性和国际化。

性能建议：

- 背景图优先使用 `webp`，控制在约 400KB 以内；如过大，另做压缩。
- 首屏只加载一张主背景，第二张山形层可按需延后或内联轻量 CSS。
- 移动端可用同一张图裁切，不必先做多端图片，除非截图验证显示构图失衡。

验证建议：

- 生成图必须在桌面与移动端各看一次裁切效果。
- 首屏文字必须压在留白区域，不能压在重墨区域。
- 暗色模式若启用，优先 CSS 反相/覆盖，不直接复用过亮纸图。

### 7.4 背景图片方向

若生成图片，方向如下：

背景主图关键词：

- 东方宣纸纹理、淡墨山影、远山留白、细线经纬、微弱拓印感、朱砂小印。
- 画面低对比、大片留白、无明显人物、无现代办公场景。
- 不要写实摄影，不要黑暗赛博，不要 stock 山水图。
- 不要把 UI 面板画进背景图片，UI 应由代码实现。

推荐图片构成：

- 一张全屏纸本纹理底图，可作为 `public` 或 `src/assets` 静态资源。
- 一层 CSS 细线网格和圆弧秩序线。
- 一层水墨山形或墨迹，可用 CSS/SVG/PNG，但必须轻。
- 一个朱砂印章式品牌记号，可用 CSS 或 SVG 实现。

资产命名建议：

- `frontend/src/assets/brand/sst-paper-texture.png`
- `frontend/src/assets/brand/sst-ink-mountain.png`
- `frontend/src/assets/brand/sst-seal.svg`

如果不生成图片，优先用 CSS：多层 radial/linear gradients、noise-like texture、mask、细线伪元素。注意避免紫蓝渐变和 orb。

## 8. 功能策略

原则：功能不删减，入口可分层。

用户侧建议：

- 一级保留：Dashboard、API Keys、Usage、Profile。
- 按配置显示：Available Channels、Channel Status、Purchase、Orders、Affiliate。
- 可弱化：Redeem、Subscriptions、过多营销型入口。
- Dashboard 首屏应强调余额、密钥、用量、服务状态，减少“低价充值”心智。

管理员侧建议：

- 一级保留：Dashboard、Accounts、Groups、Channels、Users、Usage、Settings。
- 按配置显示：Ops、Risk Control、Payment、Affiliate、Available Channels。
- 复杂功能可以折叠到分组，不删除。

首页建议：

- 不展示完整功能矩阵。
- 不展示倍率宣传。
- 不展示“官方几折”对比。
- 可以暗示“稳定、清晰、长期、审慎”。

付费定位表达：

- 不主动说“更贵”。
- 可以说“稳定优先”、“面向长期使用”、“规则清晰”、“低噪声服务”。
- 价格、倍率、套餐应进入登录后或文档/后台配置，不放首页首屏。

## 9. 实施顺序建议

第一阶段：首页品牌入口。

- 改 `frontend/src/views/HomeView.vue` 默认首页。
- 保留 script 的公共设置、主题、登录态、`home_content` 接管逻辑。
- 删除默认首页的接口示例、功能卡堆叠、供应商标签。
- 首页高度控制在约 `150vh`，移动端自然分层。

第二阶段：登录注册外壳统一。

- 改 `frontend/src/components/layout/AuthLayout.vue`。
- 去掉 SaaS 渐变和 orb，改为纸本/细线/朱砂点睛。
- 表单卡片保持清晰，不牺牲可用性。

第三阶段：控制台入口层弱化。

- 优先调整 `AppSidebar.vue` 的排序、分组和显示策略。
- 使用既有 `featureFlags.ts` 和后台开关。
- 不删除路由，避免破坏深链和历史用户习惯。

第四阶段：用户 dashboard 品质化。

- 调整首屏信息优先级。
- 保留关键操作，但减少营销和低价导向。

## 10. 验证规则

每次改动后至少执行：

- `pnpm --dir frontend typecheck`
- 如改首页/登录页，启动 Vite 后做桌面与移动端截图检查。
- 检查 `/home`、`/login`、已登录 dashboard 路径、主题切换、语言切换、`home_content` 接管逻辑。

视觉验收：

- 首屏是否一眼看到“山枢庭 / SST / SST Gateway / 统一入口，安静流转。”
- 是否没有默认接口示例，尤其没有 `/v1/chat/completions`。
- 是否不像普通 SaaS 模板。
- 是否传达稳定、克制、品质，而不是低价、热闹、堆功能。
- 移动端文字不重叠，按钮不挤压，首页不变成长篇落地页。

## 11. 后续改造硬规则

- 修改前先确认当前文件职责和依赖，不做无关重构。
- 首页可以完全重做，但核心业务逻辑不动。
- 功能默认保留；隐藏和弱化必须可逆。
- 不用网络素材，除非明确要求联网搜图或生成图片。
- 不引入新的 UI 框架。
- 不把参考图布局照搬进项目。
- 不把山枢庭做成普通 AI API Gateway 模板站。
