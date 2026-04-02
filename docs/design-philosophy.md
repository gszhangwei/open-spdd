# AI Isn't Lacking Intelligence — It Has Too Many "Ideas"

> The stronger the model, the more you need to spell out your intent and boundaries

## A Counterintuitive Observation

AI coding tools are getting increasingly powerful. Claude Code, Cursor, Codex, various Coding Agents… model capabilities are advancing rapidly.

But if you're working on business development of any real complexity and you carefully review AI-generated code, you may have encountered these situations:

- You ask it to implement a feature, and it adds things you never asked for — you just wanted a simple registration, but it throws in OAuth2 and risk control
- You ask it to make a change, and it "optimizes" code it shouldn't have touched, introducing side effects you didn't anticipate
- You thought it "understood" the business context, but the generated code works technically while the design direction diverges from your intent

The code runs, but it needs rework. The rework isn't because it "wrote it wrong," but because it "did too much" or "went off course."

These issues don't happen every time. On simple tasks, AI is already very efficient. But once architectural decisions, business boundaries, and design trade-offs are involved — scenarios where "multiple reasonable choices exist" — the probability of AI going off course rises significantly.

This led me to think about a question: is the root cause of this rework AI's insufficient capability, or our failure to convey design intent precisely enough?

Based on this thinking, I built [OpenSPDD](https://github.com/gszhangwei/open-spdd) (Structured Prompt-Driven Development tool), attempting to address this problem through prompt structure. Through practice, there are some thoughts worth sharing.

## Capability vs. Control: A Key Distinction

Before diving into the specifics, one distinction needs to be established, as it is the foundation for understanding all subsequent arguments.

Improvements in AI coding can happen along two dimensions:

**Capability dimension** — Can AI understand requirements, generate correct code, and infer project conventions? This dimension is continuously improving with model evolution, without question.

**Control dimension** — How to ensure AI's understanding aligns with your intent, how to pick the one you want among multiple "correct" solutions, how to define what must not be done.

These two dimensions overlap but are not identical. Capability improvements can indeed alleviate some control problems — larger context windows, persistent memory, and better conversational clarification are all improving AI's understanding of design intent. But there is a class of control needs that capability improvements alone cannot eliminate: when multiple "correct" solutions exist, choosing which one is a human decision, not an AI inference; and this decision needs to be recorded, shared, and continuously maintained, otherwise the next conversation, the next team member, the next Agent all need to make it again from scratch.

A senior architect's capability is beyond question. Why do enterprises still require them to write design documents? Not because their capability is insufficient, but because teams need shared understanding, decisions need traceability, and future maintainers need context. Similarly, no matter how intelligent aircraft autopilot becomes, flight plans won't disappear — not because autopilot isn't good enough, but because air traffic control, crew, and dispatch need a shared, verifiable coordination baseline.

The control problem is not new to the AI coding era. When Fred Brooks distinguished "essential complexity" from "accidental complexity" in 1986, he pointed out that the hardest part of software is not coding but specification and design; practices like ADR (Architecture Decision Record) and Design by Contract also address "how to externalize and communicate design decisions."

But AI coding makes this old problem more urgent — not because the absolute cost of rework is higher (AI can regenerate quickly), but because deviations are more subtle. AI-generated code typically has decent quality, consistent style, and passes compilation and basic tests, making it easy to pass a cursory review. Deviations in design direction — using the wrong architectural pattern, doing things that shouldn't have been done — are often not discovered until later iterations or even after deployment, when the cost of correction is truly high.

This article does not attempt to offer new theoretical insights. It focuses on a practical question: how can these known control needs be operationalized into actionable workflows within the AI coding toolchain?

## Why Control Is So Hard

### The Ambiguity Amplification Effect of LLMs

LLMs are not deterministic compilers. When input admits multiple reasonable interpretations, they pick one to execute — and you often cannot predict which one. The issue is not that "AI gives different results each time," but that a single natural language statement inherently carries multiple possible intents, and AI amplifies this input-level ambiguity into output-level divergence.

A stronger AI can indeed reduce some ambiguity. It can scan your codebase and infer the project's prevailing patterns — for example, most Services use constructor injection, and the API style is RESTful. This kind of ambiguity at the coding style and convention level can increasingly be narrowed as AI gets stronger (provided the codebase itself is sufficiently consistent).

But there is a class of uncertainty that even the strongest AI cannot reduce through inference:

- Should this feature use event-driven or synchronous calls?
- Should it support multiple currencies or just a single currency?
- What things must not be done, must not be touched?

These are not "inference" problems but "decision" problems. There is no objectively "more correct" answer — it depends on the specific trade-offs of your project right now. If the codebase has no precedent for similar decisions, AI has no basis for inference; even if precedent exists, AI cannot distinguish "this was a deliberate choice" from "this was a historical compromise."

This is the core challenge of the control dimension — not helping AI infer coding style, but precisely conveying human design decisions, architectural trade-offs, and constraint boundaries to AI.

### "Negative Space" Is Easily Overlooked

A common cause of rework in AI coding is not "missing features" but what AI did extra or went off course on:

```
You: "Implement this Service"
AI:  [Added a caching layer you didn't ask for]
     [Decided on an exception retry strategy on its own]
     [Placed the transaction boundary at the caller instead of the service]
```

Each item individually seems "reasonable," but together they diverge from your design intent.

Human developers have implicit professional judgment for self-restraint — "this isn't my call to make," "this is outside the current scope." But AI lacks this sense of boundaries. If the prompt only defines the "positive space" (what to do) without defining the "negative space" (what not to do), AI tends to "overperform" within reasonable bounds.

Defining negative space doesn't rely solely on natural language declarations. A more reliable approach is semantic declarations combined with machine-verifiable constraints — the former being soft constraints, the latter hard constraints. But explicitly declaring negative space, even as just soft constraints, is better than not declaring it at all.

### Design Intent Easily Gets Lost Over Time

Consider this scenario:

```
v1.0: Plan defines "order service, synchronous processing, single currency USD"
      AI generates code accordingly
      ↓
Architecture review: "For scalability, switch to async event-driven"
Business decision: "European market requires EUR support"
      ↓
v1.1: Team refactors to async + multi-currency
      ↓
⚠️ Plan still describes "synchronous processing, single currency"
      ↓
v2.0: Need to generate "refund" functionality
```

After scanning the codebase, the Agent would likely see the async pattern and multi-currency support, and the generated refund code should be technically compatible. But what the Agent doesn't know is:

- Why was it made asynchronous? Was it to solve a performance bottleneck, or to prepare for a future message queue migration? Must refunds also go through the async pipeline?
- Does multi-currency mean "we just added EUR" or "we're building a universal multi-currency architecture"? Does the refund need to support arbitrary currencies?
- Which v1.0 code is "confirmed to be kept," and which is "not yet cleaned up"?

These are all design intent — invisible in the code, and rarely addressed in typical execution plans. Codebase scanning can prevent technical-level conflicts but cannot prevent design-direction drift.

It's worth noting that mainstream tool vendors are also exploring in the same direction: Claude officially positions CLAUDE.md as providing "persistent context it cannot infer from code," and Codex recommends writing corrections of agent's wrong assumptions into AGENTS.md. Vendor practices are not gospel truth, but when multiple independent teams converge on adding explicit context mechanisms beyond code scanning, it at least suggests that relying solely on code scanning is insufficient in practice — an empirical judgment that has been repeatedly validated.

### The Interpretation Space Is Hard to Narrow

```
Requirement: "Create a user registration service with email validation and password encryption"

AI-A: Simple Service, validate + encode + save, 50 lines
AI-B: Adds event publishing, audit logs, async email queue, 200 lines
AI-C: CQRS pattern, Command/Query separation + Event Sourcing, 500 lines
```

All three satisfy the literal requirement, but the complexity, maintenance cost, and architectural implications are entirely different.

Ambiguity in human development can be resolved through conversation. AI tools are also improving in this regard — Cursor's Plan Mode proactively asks questions, Claude Code clarifies before acting. But in complex scenarios, some ambiguities are hard to fully resolve through conversation, because developers themselves may not realize certain decisions need to be made explicit. A structured framework can serve as a checklist, helping discover those "you don't know you need to decide" decision points.

## What Current AI Coding Tools Are Doing

Mainstream tools have already done substantial work in the direction of "conveying design intent to AI," and the results are significant:

- **Cursor**'s Agent mode can deeply scan codebases, infer project conventions, and proactively clarify requirements. Plan Mode can produce structured plans with file paths and code references. For coding style consistency within a project, it already does a great job.
- **Kiro** defines specs as structured artifacts, split into requirements, design, and tasks files, going further in spec structuring.
- **Claude Code** emphasizes exploring the codebase and clarifying requirements before execution. CLAUDE.md is used to record coding standards and workflow rules, and its hooks mechanism supports machine-verifiable automated constraints.
- **Codex**'s AGENTS.md allows solidifying project conventions and recurring error corrections, supporting directory-level context configuration.
- **GitHub Copilot** supports project-level specification through `copilot-instructions.md`, and the `copilot-prompts/` directory can store reusable prompt templates.

These tools have far surpassed "simple code completion" and are all evolving toward spec-first, solving a large portion of real-world problems — especially aspects that can be inferred from the codebase like coding conventions, naming styles, and project structure.

However, as discussed in the previous section, there is a class of problems — design decisions, architectural trade-offs, negative-space constraints — that codebase scanning and conversational interaction alone still cannot fully cover. This is not a matter of tools being insufficient, but a structural challenge inherent in "precisely conveying design intent to AI."

## Codebase Scanning: Valuable, but with Boundaries

Codebase scanning is one of the most important capabilities of current AI coding tools. Understanding its boundaries helps us work with it more effectively.

**Code records "what is," not "what should be."**

```
After scanning the codebase, Agent finds:
- ServiceA uses field injection
- ServiceB uses constructor injection
- ServiceC uses setter injection

Agent's confusion: "What is this project's injection standard?"
Truth: ServiceA is legacy code from three years ago, ServiceB is the new standard after refactoring, ServiceC was written by an intern
```

A codebase is the accumulation of historical decisions, mixing best practices, legacy artifacts, and technical debt. One important use of CLAUDE.md and AGENTS.md is to help Agents distinguish "standards to follow" from "legacy to clean up."

This also responds to a premise mentioned earlier: AI's ability to infer coding conventions depends on the consistency of the codebase itself. But this premise doesn't hold automatically — three injection methods coexisting is precisely because the project lacked an explicit convention about "which one to use." This creates a cycle: AI needs a consistent codebase to infer standards, but the codebase's consistency itself requires some form of explicit constraints to maintain. Without control, there's no consistency; without consistency, AI's inference capability is diminished.

**Code doesn't record "why."**

```
Agent scans and finds: PaymentService uses eventual consistency

Agent cannot know:
- Was this a deliberate architectural decision? Or a temporary compromise?
- Why wasn't strong consistency used? Performance considerations? Or technical limitations?
- Is there a plan to migrate to strong consistency in the future?
```

Code is "the result of decisions," not "the process and reasoning behind decisions."

**The codebase looks backward; design intent looks forward.**

```
Scenario: Adding "multi-tenancy" support to the existing system

Agent scans codebase: Currently single-tenant architecture
Agent's inference: Continue generating code in single-tenant mode

But the design intent is: This time we're switching to multi-tenant
```

When your goal is to change the existing architecture, scanning the existing code may precisely lead the Agent in the wrong direction.

**Reverse-engineering design intent from code is lossy.**

A piece of code can correspond to multiple design intents:

```java
if (user.getRole() == Role.ADMIN) {
    return allRecords();
} else {
    return userRecords(user.getId());
}
```

This could mean: simple role-based access control (good enough for now); a temporary solution, with future migration to RBAC planned; or a product requirement to only distinguish between admin and regular users. The code is the same, but under three different intents, the next development direction is completely different.

Code is a projection of design. Recovering the original intent from a projection necessarily loses information.

Codebase scanning, explicit design intent declarations, and machine-verifiable constraints — the three are complementary, not mutually exclusive. Industry practice is also converging in this direction: all major tools are adding repo-level instruction files on top of codebase scanning to supplement context that cannot be inferred from code.

## An Attempt at Structured Prompts

The problem is clear: information like design decisions, architectural trade-offs, and negative-space constraints needs to be explicitly expressed. But "explicit expression" itself has varying degrees of quality.

Taking "create a user registration service" as an example, compare two approaches:

**Natural language requirement description:**

```markdown
Create a user registration service. Handle user registration business logic, including email deduplication, password encryption, and user saving.
Use constructor injection, injecting UserRepository, EmailValidator, and PasswordEncoder.
Throw an exception if the email already exists. New user status is PENDING_VERIFICATION.
Add @Service and @Transactional annotations.
```

This is already much better than a single "create UserRegistrationService." But notice what it omits: no explanation of business context (why only email registration?), and no boundary definitions (don't introduce OAuth? don't add caching?). This isn't because the writer wasn't careful enough — without structural guidance, people naturally tend to describe "what to do" and "how to do it" while omitting "why" and "what not to do."

**Expression guided by a structured framework (excerpt):**

```markdown
## Requirements (Why)

- Business context: MVP phase, only need email registration, no OAuth required

## Operations (How)

### Create Service Implementation - UserRegistrationServiceImpl

1. Responsibility: Handle user registration business logic
2. Location: `com.example.user.service.impl.UserRegistrationServiceImpl`
3. Dependencies (constructor injection):
   - `UserRepository userRepository`
   - `EmailValidator emailValidator`
   - `PasswordEncoder passwordEncoder`
4. Annotations: `@Service`, `@Transactional`
5. Method `register(UserRegistrationRequest request)`: UserRegistrationResponse
   - Logic:
     1. Call `emailValidator.validate(request.getEmail())`
     2. Call `userRepository.existsByEmail(request.getEmail())` → if true, throw `EmailAlreadyExistsException` with message "Email already registered"
     3. Call `passwordEncoder.encode(request.getPassword())` → get encodedPassword
     4. Create User entity with status `PENDING_VERIFICATION`
     5. Call `userRepository.save(user)` → get savedUser
     6. Return UserRegistrationResponse with savedUser.getId()

## Safeguards (What Not to Do)

- Do not introduce OAuth, social login, or other third-party authentication
- Do not add a custom caching layer
- Do not modify existing User entity field definitions
```

The difference is not just "level of detail." The latter structurally guides you to think about and declare dimensions that are easily overlooked — especially Safeguards (what not to do) and Requirements (why). It acts as a checklist effect: not helping you write more words, but helping you miss fewer critical dimensions.

But the checklist effect is just the starting point. A question checklist ("have you defined what not to do?") can serve as a reminder, but its output is a confirmation in the respondent's mind, not a consumable artifact. The additional value of a structured framework is: it not only reminds you to think about these dimensions but also provides a structure for filling in each dimension, producing a document that can be directly fed to AI for execution, version-controlled, and circulated across the team. From "reminding you to think of it" to "helping you write it down and letting the toolchain consume it" — this is the difference between a checklist and a framework.

Some might think: I don't need this much upfront structure; just let AI do it and fix the deviations. Iterative correction is indeed efficient within a single conversation — "remove the caching," "switch to synchronous," and a few rounds can get it right. But the problem is: the decisions implicit in these corrections are scattered throughout conversation history. Next time you start a new conversation, the same deviations may recur; Agents launched by other team members don't know about these conventions; when you switch models or tools, conversation history isn't portable. Iterative correction solves "this time's" problem; structured prompts solve "every time's" and "everyone's" problem.

"Structured prompts" is a philosophy, not a fixed format. The industry already has multiple practices: Kiro's three-tier structure of requirements/design/tasks, Claude's CLAUDE.md, Codex's AGENTS.md, traditional ADR. These practices have different emphases but all address the same class of problem.

What OpenSPDD's REASONS Canvas attempts to do is integrate these scattered concerns into a unified framework embeddable in the AI coding toolchain — a 7-dimension checklist:

```
┌─────────────────────────────────────────────────────────────────────┐
│                        REASONS Canvas                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  R - Requirements    Why to do it                ┐                  │
│  E - Entities        What concepts are involved  ├─ Strategic Layer  │
│  A - Approach        What approach to use        ┘   (Why/What)     │
│                                                                     │
│  S - Structure       How components are organized   ┐               │
│  O - Operations      Exactly how to do it           ├─ Implementation│
│                      (down to method signatures)    ┘   Layer (How)  │
│                                                                     │
│  N - Norms           What standards to follow       ┐               │
│  S - Safeguards      What must not be done          ├─ Constraint    │
│                                                     ┘   Layer       │
│                                                        (Boundary)   │
└─────────────────────────────────────────────────────────────────────┘
```

These three layers correspond to the three types of control problems discussed earlier:

- The strategic layer tells AI "why" and "what," addressing context gaps and design intent loss
- The implementation layer tells AI "precisely how," narrowing the interpretation space
- The constraint layer tells AI "what must not be done," addressing AI's "overperformance"

It is not "the only correct approach" but a way of thinking that reduces omissions through structural guidance. The question has never been "whether to express design intent" — the industry has long had consensus on this. The real difficulty is: in the fast-paced workflow of AI coding, how to do it in a way that is lightweight enough, embeddable in existing toolchains, and can evolve alongside the code. REASONS Canvas is one answer to this practical problem. If you're already using CLAUDE.md or AGENTS.md and it's working well, that's also an effective approach — the two operate at different granularities and are complementary rather than substitutive. CLAUDE.md / AGENTS.md are typically **repo-level** conventions: they record "how this project does things" (use constructor injection, API style is RESTful, error handling uses a unified format). REASONS Canvas targets **feature-level** decisions: it records "how this specific feature is designed" (this payment module uses synchronous calls, only does single currency, no caching layer). The former is like a team's coding standards handbook; the latter is like the conclusions from a particular architecture review. Both layers have value; which one you're missing depends on where your pain points are.

One thing to be clear about: structured prompts don't make AI stop "interpreting" entirely; they compress AI's interpretation space to a smaller, more controllable range. What they reduce is unnecessary interpretation, not the elimination of interpretation.

## AI Is Getting Stronger — Are Structured Prompts Still Needed?

Some may ask: once AI evolves to be "smart enough," can structured prompts be dropped?

For simple tasks like writing scripts, changing configurations, and adding logs — indeed they can. AI getting stronger makes these things increasingly effortless.

But in scenarios involving architectural decisions, business rules, and multi-person collaboration, my view is: you need them even more.

Suppose AI reaches "perfection" — infinite context, zero hallucination. Faced with "give me a payment system," a truly "perfect" AI's most likely response wouldn't be generating code directly, but asking follow-up questions — just as Cursor's Plan Mode and Claude Code are already doing today:

```
Perfect AI: "Do you want event-driven or synchronous calls? Multi-currency or single-currency?
             Eventual consistency or strong consistency for reconciliation?"
```

But this precisely proves the point: **Even if AI is smart enough to know what to ask, these decisions still need to be made by humans, and they need to be recorded.**

AI's questions don't produce decisions; human answers do. And if these answers only exist in one-off conversations:

- Next time you start a new conversation, the same questions need to be answered again
- Agents started by other team members don't know about these decisions
- Three months later, when you switch models or tools, conversation history isn't portable

A payment system can have 100 "perfect" designs, each one "correct," but only one is "what your project needs right now." These choices need to be solidified — not in one-off conversations, but in sustainably consumable artifacts.

Returning to the "capability vs. control" framework: AI getting stronger solves capability problems — understanding more accurately, generating better output. But control problems — how to pick the one you want among multiple "correct" options, how to define what cannot be done, how to let different Agents at different points in time share the same set of decisions — don't disappear because AI gets stronger.

## Shared Artifacts for the Collaboration Era

Software development is fundamentally distributed collaboration — whether the participants are humans or Agents. And distributed collaboration has a basic prerequisite: participants need to share some kind of artifact to coordinate their actions.

This shared artifact can take many forms:

- Structured design specifications (CLAUDE.md, AGENTS.md, REASONS Canvas, etc.)
- Machine-verifiable constraints (tests, type systems, interface contracts, Schemas, etc.)
- Decision records (ADR / Architecture Decision Record)
- CI strategies (pre-commit hooks, linters, automated checks)

In real projects, effective coordination is usually a combination of these forms rather than relying on just one.

Multi-Agent collaboration is already a reality: Cursor can launch sub-tasks for parallel processing, and Claude Code supports multiple agent instances collaborating. In these scenarios, shared artifacts become even more important — each Agent needs to know the global architectural decisions and constraint boundaries, otherwise the results of parallel execution will likely contradict each other.

Even with a single Agent, across time it effectively faces a "different entities" problem:

```
T1: Agent-v1 (GPT-5) developed system v1
T2: Agent-v2 (GPT-6) needs to develop v2 on top of v1
T3: Agent-v3 (Claude-5) needs to refactor the system
```

Without externalized decision records, Agent-v2 doesn't know why Agent-v1 chose event-driven over synchronous calls, and Agent-v3 doesn't know which designs are "deliberate trade-offs" and which are "historical compromises."

Shared artifacts are not just a collaboration baseline in the spatial dimension but also a knowledge carrier in the temporal dimension.

There's another easily overlooked angle: when humans step back from "executor" to "supervisor," humans' understanding of the system actually becomes harder. The Agent wrote all the code; humans need some kind of "telescope" to understand the big picture. Structured decision records can serve as that telescope.

## Not Just the Expression Layer — The Verification Layer Is Equally Important

At this point in the discussion, one thing must be emphasized: what truly makes AI output controllable is not just "more structured expression" but also a machine-verifiable validation loop.

Structured prompts address the "expression layer" — how to convey design intent to AI. But they are fundamentally soft constraints: AI can reference them and may also deviate from them.

The verification layer provides hard constraints — tests, type checking, Schema validation, Lint rules, pre-commit hooks. These don't depend on AI's "understanding" but mechanically check whether the output conforms to the rules.

The more robust combination is:

> Codebase scanning (understanding the current state) + Structured prompts (expressing intent and constraints) + Automated verification (machine-checking critical rules)

No single layer is a silver bullet. The three layers working together achieve a balance between "AI autonomy" and "human controllability."

## Applicable Scenarios and Limitations

We need an honest discussion of the applicable boundaries of structured prompts (including REASONS Canvas).

**Suitable scenarios:**

- Core business logic, systems requiring long-term maintenance and multi-person collaboration
- Architectural decision points, design choices involving significant trade-offs
- Regulated domains, scenarios requiring auditing and traceability
- Cross-team collaboration, projects where multiple Agents or humans need shared understanding

**Unsuitable scenarios:**

- Rapid prototyping, exploratory development, expected to be largely discarded
- One-off scripts, use-and-discard, no maintenance needed
- Simple tasks, fixing a typo, adding a log — not worth the investment in additional structured work

Structured prompts are externalization tools for important decisions, not a universal method for all AI coding tasks. Overuse actually increases cognitive burden, the same problem as "over-engineering" in traditional software engineering.

There is also a question that must be addressed head-on: **Structured prompts themselves can also become outdated.** Earlier, we used "the v1.0 Plan is already outdated by v1.1" to argue that design intent is easily lost, but the same criticism applies entirely to REASONS Canvas — if the code changes but the Canvas doesn't get updated, it transforms from a "design guide" into "misleading outdated documentation," which is more dangerous than no documentation at all.

This is not a problem to be glossed over. Frankly, bidirectional synchronization between design specifications and code is one of the biggest engineering challenges today, and no solution has fully solved it. Our current approach is: after each feature iteration, manually run `/spdd-sync` to have AI detect code changes and update the corresponding Canvas. This isn't yet a mandatory process, more of a team convention — "after changing the code, sync the spec while you're at it." It's not perfect, but it works in practice: changes to key decisions are basically captured, and the Canvas doesn't become seriously outdated. The future direction is to automate this step — for example, automatically triggering sync detection after each code update, and if changes are found, automatically updating the corresponding Canvas.

## Looking Ahead

As AI capabilities improve, a clear trend is: structured artifacts are increasingly generated with AI assistance or even AI-led. Taking OpenSPDD as an example, the current workflow is already: AI analyzes business requirements and codebase (`/spdd-analysis`) → AI automatically generates structured prompts (`/spdd-reasons-canvas`) → Humans review key decisions → AI executes (`/spdd-generate`) → Changes sync back to prompts (`/spdd-sync`).

If this trend continues, a possible evolution path is:

```
Current:    AI-assisted analysis → AI generates structured prompts → Human review → AI executes
Near-term:  AI autonomously makes routine decisions → AI autonomously handles part of the review → Humans only review key nodes → Multi-Agent collaborative execution
Long-term:  Agents autonomously generate and maintain most structured artifacts → Humans spot-check and govern
```

There is no fixed timeline for this evolution; it is speculation based on current trends.

There's a seemingly contradictory point here: if AI can ultimately generate and maintain these structured artifacts autonomously, do humans still need to care about this?

My understanding is: what changes is the "author" — from humans writing, to AI writing with human review, to AI maintaining autonomously. What doesn't change is the "need" — complex systems always need some form of structured, shareable, traceable intermediate artifacts to coordinate the actions of multiple participants. Just as compilers automatically generating machine code doesn't mean instruction sets can disappear, AI automatically generating design specifications doesn't mean the specifications themselves become unimportant. Only the mode of production changes; the necessity of existence doesn't.

Perhaps in the future it won't be called a prompt anymore. Perhaps it will evolve into a DSL, decision records, executable specs, or auto-syncing contract files. But the need behind it won't disappear.

## Conclusion

Returning to the core question: as AI gets stronger, can our approach to collaborating with AI become more casual?

For simple tasks, yes, it can be lighter. But for complex systems, it's precisely the opposite.

Understanding this requires three cognitive shifts:

**In content, from "prompts" to "constraints."** The question isn't just "how to state requirements clearly" but "how to state boundaries clearly." Much rework isn't because AI didn't deliver, but because AI did too much, went off course, or made decisions on your behalf that it shouldn't have.

**In time, from "one-time input" to "continuously evolving artifacts."** What's truly valuable isn't the prompt at that moment, but the intermediate artifact that evolves alongside the code and can be consumed by those who come after.

**In audience, from "human instructions to AI" to "a carrier for externalized decisions."** The value of structured prompts isn't just helping AI generate better code. More importantly, they externalize design intent — previously scattered across people's minds, conversation history, and code details — into a carrier that is easier to share, reuse, trace, and verify — for AI to consume, and for the team to consume.

The question has never been just "is AI smart enough." The more fundamental question is: as models grow stronger and collaboration grows more complex, have we written the truly important things — intent, boundaries, trade-offs — into a form that Agents can reliably consume and teams can sustainably reuse?

## If You Want to Try

### With Your Existing Tools

You don't need OpenSPDD; you can start right now:

- **Cursor**: Make good use of Agent mode and Plan Mode, letting AI fully clarify requirements before generating code; use the `.cursor/rules/` directory to solidify project conventions
- **Claude Code**: Maintain a good CLAUDE.md, recording project conventions and design decisions; leverage hooks for automated verification
- **Codex**: Make good use of AGENTS.md, solidifying recurring issues and project conventions
- **GitHub Copilot**: Use `copilot-instructions.md` to define project specifications, and `copilot-prompts/` to store reusable prompt templates

### Experience OpenSPDD

If you want to try a more structured approach with REASONS Canvas:

```bash
brew install gszhangwei/tools/openspdd
cd your-project
openspdd generate --all    # Generate SPDD commands

# Use in AI coding tools:
/spdd-analysis @requirements/feature.md      # Strategic analysis
/spdd-reasons-canvas @spdd/analysis/xxx.md   # Generate structured prompts
/spdd-generate @spdd/prompt/xxx.md           # Generate code from structured prompts
/spdd-sync @spdd/prompt/xxx.md               # Sync changes back to prompts
```

OpenSPDD is already being used in real projects for backend business logic. In our practice, for feature development with clear business rules, the reduction in rework and improvement in directional accuracy brought by structured prompts is significant — some features that previously required days of iterative adjustment can be completed in hours once specifications are clear. Of course, this efficiency difference is highly dependent on project type and task complexity and should not be generalized. The tool is still under continuous iteration; feel free to try it, and submit issues and PRs.

### Examine Your Current AI Workflow

Regardless of which tool you use, ask yourself a few questions:

- Does AI frequently "overdo" or "go off course"? → You may need clearer boundary definitions
- After code changes, is the previous design documentation still accurate? → You may need a better sync mechanism
- Is the code written by team members using AI consistent? → You may need shared convention files
- Have important design decisions been recorded? → You may need ADR or a similar mechanism

### Open Questions

This article presents a line of thinking, not a conclusion. Many questions are still worth exploring:

- What is the "optimal granularity" for structured prompts? When is it worth writing them, and when is it overkill?
- How can bidirectional synchronization between design specs and code be made reliable in engineering?
- What kind of structured framework do different teams and project types need?
- When Agents autonomously generate specs, how can humans effectively review them?

These questions have no standard answers and need to be explored through practice.

---

*This article explores the value and limitations of structured design intent expression in the AI coding era. [OpenSPDD](https://github.com/gszhangwei/open-spdd) is one such attempt — feel free to try it and share your feedback.*
