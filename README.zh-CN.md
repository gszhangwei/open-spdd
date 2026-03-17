# OpenSPDD

> **结构化提示词驱动开发** — 让 AI 编码的 Prompt 成为可执行的设计契约

[English](README.md)

OpenSPDD 是一套面向 AI 编码时代的结构化提示词驱动开发方法论及跨平台 CLI 工具。它将 AI 编码的 Prompt 从"一次性输入"升级为"可执行的设计契约"，实现设计与实现的双向同步。

## 为什么需要 OpenSPDD？

现有的 AI 编码工具虽然也会生成 plan 文档或执行计划，但这些文档存在根本性的局限：

| 问题 | 典型 Plan 文档 | REASONS Canvas |
|------|----------------|----------------|
| **本质定位** | 任务清单（Task List） | 设计契约（Design Contract） |
| **约束力** | 无 — AI 可自由发挥 | 有 — Norms 定义"如何做"，Safeguards 定义"不能做什么" |
| **详细程度** | 高层描述：*"创建 BillingService"* | 精确规格：*方法签名、参数类型、错误处理、依赖注入方式* |
| **可追溯性** | 无 — 代码改了文档不更新 | 有 — `/spdd-sync` 支持反向同步 |
| **验证标准** | 模糊 — *"完成即可"* | 明确 — Safeguards 中定义精确的错误消息、HTTP状态码 |
| **依赖管理** | 隐式 — AI 自行推断 | 显式 — Operations 定义严格的执行顺序和依赖关系 |

**核心洞察**：Plan 是"建议"，REASONS Canvas 是"契约"。

## REASONS Canvas 框架

REASONS Canvas 是一个 7 维度的结构化设计框架：

```
┌─────────────────────────────────────────────────────────────────────┐
│                        REASONS Canvas                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  R - Requirements    需求本质，回答"为什么做"                         │
│  E - Entities        领域模型（Mermaid 类图），回答"涉及什么概念"      │
│  A - Approach        方案策略与权衡，回答"用什么方式"                  │
│  S - Structure       架构层次/继承/依赖，回答"组件如何组织"           │
│  O - Operations      精确的实现任务序列，回答"具体怎么做"             │
│  N - Norms           编码规范与模式，回答"按什么标准"                 │
│  S - Safeguards      约束与护栏，回答"什么不能做"                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**为什么是 7 个维度？**
- **R+E+A** = 设计决策的"Why"和"What"
- **S+O** = 实现路径的"How"
- **N+S** = 质量保障的"Guardrails"

三者缺一不可：缺少 N+S，AI 会自由发挥；缺少 S+O，AI 会随意架构；缺少 R+E+A，AI 不理解上下文。

## 核心工作流

```
┌─────────────────────────────────────────────────────────────────────┐
│                         SPDD 完整工作流                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  业务需求                                                            │
│      │                                                               │
│      ▼                                                               │
│  /spdd-analysis ──────→ 战略级分析（概念识别、方案方向、风险评估）     │
│      │                                                               │
│      ▼                                                               │
│  /spdd-reasons-canvas ─→ REASONS Canvas 结构化设计文档               │
│      │                                                               │
│      ▼                                                               │
│  /spdd-generate ───────→ AI 按契约生成代码（不自由发挥）             │
│      │                                                               │
│      ▼                                                               │
│  代码审查/重构                                                       │
│      │                                                               │
│      ▼                                                               │
│  /spdd-sync ───────────→ 代码变更反向同步回设计文档                  │
│      │                                                               │
│      ▼                                                               │
│  设计文档与代码保持一致 ──────→ 下一轮开发基于准确的设计             │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**核心原则**：*"When reality diverges, fix the prompt first — then update the code."*（当现实与设计分歧时，先修改 Prompt，再更新代码）

## 核心特性

- **跨平台支持**：适配 Cursor、Claude Code、GitHub Copilot、Antigravity
- **自动检测**：自动识别当前 AI 编码环境
- **单一二进制**：所有模板通过 Go embed 嵌入，无外部依赖
- **双向同步**：设计文档与代码保持同步
- **交互式 UI**：现代化终端界面进行命令选择

## 安装

### Homebrew (macOS/Linux)

```bash
brew install gszhangwei/tools/openspdd
```

或：

```bash
brew tap gszhangwei/tools
brew install openspdd
```

升级到最新版本：

```bash
brew upgrade openspdd
```

### Go Install

```bash
go install github.com/gszhangwei/open-spdd@latest
```

### 下载二进制

从 [GitHub Releases](https://github.com/gszhangwei/open-spdd/releases) 下载。

## 快速开始

```bash
# 进入项目目录
cd your-project

# 初始化（自动检测 AI 工具）
openspdd init

# 生成 SPDD 命令
openspdd generate --all
```

然后在 AI 编码工具中，按照完整的 SPDD 工作流操作：

```bash
# 第一步：战略级分析（复杂功能推荐）
/spdd-analysis @requirements/user-registration.md

# 第二步：根据分析生成 REASONS Canvas
/spdd-reasons-canvas @spdd/analysis/xxx.md

# 第三步：根据 REASONS Canvas 生成代码
/spdd-generate @spdd/prompt/xxx.md

# 第四步：代码审查/重构后，同步变更回设计文档
/spdd-sync @spdd/prompt/xxx.md
```

对于简单功能，可以跳过第一步，直接提供需求描述：

```bash
/spdd-reasons-canvas 实现用户注册功能，支持邮箱验证
```

## 使用方法

### 初始化环境

```bash
# 自动检测并初始化
openspdd init

# 手动指定工具
openspdd --tool cursor init
```

### 列出命令

```bash
# 列出可用命令（核心 + 工具特定）
openspdd list

# 列出可选命令
openspdd list --optional

# 列出所有命令
openspdd list --all

# 按类别筛选
openspdd list -c Development
```

### 生成命令

```bash
# 生成所有默认命令
openspdd generate --all

# 交互式选择
openspdd generate

# 生成特定命令
openspdd generate spdd-generate

# 强制覆盖
openspdd generate --force spdd-generate
```

### 全局标志

```bash
openspdd --tool cursor <command>
openspdd --tool claude-code <command>
openspdd --tool antigravity <command>
openspdd --tool github-copilot <command>
```

## 支持的环境

| 工具           | 检测方式                                               | 配置目录                   |
| -------------- | ------------------------------------------------------ | -------------------------- |
| Cursor         | `.cursor/`, `.cursorrules`                             | `.cursor/commands/`        |
| Claude Code    | `.claude/`, `CLAUDE.md`                                | `.claude/commands/`        |
| Antigravity    | `.antigravity/`                                        | `.antigravity/commands/`   |
| GitHub Copilot | `.github/copilot-instructions.md`, `.github/copilot-prompts/` | `.github/copilot-prompts/` |

### GitHub Copilot 文件结构

```
.github/
├── copilot-instructions.md     # 主指令文件（支持标记合并）
└── copilot-prompts/
    ├── spdd-reasons-canvas.md
    ├── spdd-generate.md
    └── spdd-sync.md
```

## 可用命令

### 核心命令

| 命令                  | 描述                                     |
| --------------------- | ---------------------------------------- |
| `spdd-generate`       | 从结构化 SPDD Prompt 文件生成代码        |
| `spdd-sync`           | 将代码变更同步回 SPDD Prompt 文件        |
| `spdd-reasons-canvas` | 生成 REASONS-Canvas 结构化 Prompt        |
| `spdd-analysis`       | 需求的战略级分析                         |

### 工具特定命令

| 工具           | 命令                   | 描述                        |
| -------------- | ---------------------- | --------------------------- |
| GitHub Copilot | `copilot-instructions` | Copilot 主指令文件          |

### 可选命令

```bash
# 列出可选命令
openspdd list --optional

# 安装特定可选命令
openspdd generate <optional-command-name>
```

## Plan vs REASONS Canvas：示例对比

**场景**：实现用户注册功能

**典型 Plan 文档**：
```
1. 创建 UserRegistrationController
2. 创建 UserRegistrationService
3. 创建 UserRegistrationRequest DTO
4. 实现邮箱验证
5. 保存用户到数据库
```

**REASONS Canvas（Operations 节选）**：
```markdown
### 创建 UserRegistrationService - `UserRegistrationServiceImpl`

1. **职责**: 处理用户注册业务逻辑
2. **包路径**: `com.example.user.service.impl`
3. **实现接口**: `UserRegistrationService`
4. **依赖注入** (构造器注入):
   - `UserRepository userRepository`
   - `EmailValidator emailValidator`
   - `PasswordEncoder passwordEncoder`
5. **方法**:
   - `register(UserRegistrationRequest request): UserRegistrationResponse`
     - **输入校验**: 调用 `emailValidator.validate(request.getEmail())`
     - **业务逻辑**:
       1. 通过 `userRepository.existsByEmail()` 检查邮箱是否已存在
       2. 如果存在，抛出 `EmailAlreadyExistsException`，消息为 "Email already registered"
       3. 通过 `passwordEncoder.encode()` 加密密码
       4. 创建 User 实体，状态为 `PENDING_VERIFICATION`
       5. 通过 `userRepository.save()` 保存
     - **异常处理**: 让异常传播到 GlobalExceptionHandler
6. **注解**: `@Service`, `@Transactional`
```

**差距一目了然**：Plan 说"做什么"，REASONS Canvas 规定"精确怎么做"。

## 适用场景

| 场景 | 推荐程度 | 理由 |
|------|----------|------|
| 企业级功能开发 | ⭐⭐⭐⭐⭐ | 需要设计-实现可追溯性，长期可维护 |
| 团队协作项目 | ⭐⭐⭐⭐⭐ | 统一的 AI 编码规范，减少风格冲突 |
| 复杂重构任务 | ⭐⭐⭐⭐ | Operations 的严格顺序防止依赖混乱 |
| 跨 AI 工具协作 | ⭐⭐⭐⭐ | 同一份 REASONS Canvas 在不同工具间通用 |
| 快速原型/MVP | ⭐⭐ | 可能过重，但如果后续需要维护仍值得 |
| 一次性脚本 | ⭐ | 投入产出比不高 |

## 从源码构建

```bash
git clone https://github.com/gszhangwei/open-spdd.git
cd open-spdd
go build -o openspdd .
go install .
```

## 测试

```bash
# 运行所有测试
go test ./tests/...

# 详细输出
go test ./tests/... -v

# 运行特定模块测试
go test ./tests/detector/...
go test ./tests/templates/...
```

## 许可证

[MIT License](LICENSE)
