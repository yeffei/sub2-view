# SUB2API 魔改开发计划与全站设计系统白皮书 (Zen-Tech Minimalism)

> 状态说明：本文件已被 [SST Unified Development Guide](./SST_UNIFIED_DEVELOPMENT_GUIDE.md) 作为默认统一规范接管；本文件保留为历史白皮书与创意来源参考。

## 一、 项目核心定位
* **目标受众**：专业开发者、重度 AI 团队、企业级客户及具备高审美的品质用户。
* **设计哲学**：**Tech-Minimalism（科技极简）** 与 **Oriental Zen（新中式禅意）** 的终极融合。通过宣纸肌理、泼墨留白、极细线条与等宽秩序感，彻底抹除传统中转站的“草根割韭菜感”，打造如同艺术品般的生产力工具。

---

## 二、 功能性调整蓝图（加减法原则）

AI 在开发时必须严格遵守以下业务逻辑的精简与重构：

### 1. 核心减法（消灭低端痕迹）
* **取消全站主动弹窗**：禁止任何进站强行弹出的公告。所有通知降噪为顶部 1 像素高的通栏 Banner（可手动关闭）或右上角的“未读消息小红点”。
* **隐藏后台渠道逻辑**：前台 UI 彻底对用户隐藏“渠道、节点、倍率、分组”等底层内部概念。所有模型计价在前端透明转化为：**“$ / 百万 Tokens”** 或 **“¥ / 百万 Tokens”**。
* **下架无用营销模块**：移除签到送额度、幸运抽奖、推广排行榜等破坏专业信任感的功能。将“邀请返利”更名为语义更克制的“团队引荐（Referrals）”，并隐藏在个人中心深处。

### 2. 核心加法（提升生产力溢价）
* **卡片化密钥管理 (Advanced Keys)**：支持为 API Key 分配项目标签（如 `Prod`, `Staging`, `Test`），并允许为单个 Key 设置每日/每周消耗上限。
* **无感代码沙盒 (Embedded Playground)**：前台集成极简测试面板，用户选择模型时，右侧实时生成基于该 Key 的 `cURL`、`Python` 和 `Node.js` 请求示例代码。
* **消费图表降噪 (Clean Analytics)**：拒绝密密麻麻的五彩表格。只提供一条单色（如墨黑色或翡翠绿）的无网格线消费趋势折线图，以及克制的模型消耗占比条形图。

---

## 三、 全站样式系统 (Design System)

### 1. Tailwind CSS 核心配置文件 (`tailwind.config.js`)
请将以下配置直接覆盖或融合进项目 Tailwind 配置的 `theme.extend` 中：

```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class', // 支持暗黑模式切换
  theme: {
    extend: {
      colors: {
        zen: {
          // 明亮模式（宣纸暖米）
          bg: '#f6f5f1',        // 宣纸底色，带微暖调
          card: '#ffffff',      // 玉石悬浮卡片背景
          text: '#09090b',      // 炭黑主文字（避免纯黑刺眼）
          sub: '#71717a',       // 远山灰（次要说明文字）
          line: '#e4e4e7',      // 极细微分割线
          
          // 暗黑模式（玄武泼墨）
          darkBg: '#09090b',    // 玄武岩深黑底色
          darkCard: '#121214',  // 泼墨深灰卡片背景
          darkText: '#f4f4f5',  // 月白主文字
          darkSub: '#a1a1aa',   // 烟雨灰（暗模式次要文字）
          darkLine: '#27272a',  // 暗极细分割线

          // 点缀色（克制使用，面积 < 1%）
          seal: '#b91c1c',      // 朱砂红（用于品牌印章或核心提示）
          emerald: '#10b981',   // 翡翠绿（用于 200 OK、Success 状态数字）
        }
      },
      fontFamily: {
        // 全站 UI 与大标题使用细腻的 Geist 或 Inter
        sans: ['Geist Sans', 'Inter', 'Noto Sans SC', 'sans-serif'],
        // 涉及代码、API Key、Token 数量、账单金额，强制使用等宽字体确保秩序感
        mono: ['Geist Mono', 'JetBrains Mono', 'monospace'],
        // 用于极少数大标题或新中式诗意文案的衬线体
        serif: ['Playfair Display', 'Noto Serif SC', 'serif'],
      },
      boxShadow: {
        // 复刻实体卡片悬浮感：极大模糊半径，极低透明度
        'stone': '0 30px 60px -15px rgba(0, 0, 0, 0.04), 0 10px 30px -10px rgba(0, 0, 0, 0.02)',
        'stone-dark': '0 30px 60px -15px rgba(0, 0, 0, 0.5), 0 10px 30px -10px rgba(0, 0, 0, 0.3)',
      },
      borderRadius: {
        'zen': '4px', // 全站统一采用 4px 微圆角或直角，彰显硬朗与严谨
      }
    },
  },
  plugins: [],
}
