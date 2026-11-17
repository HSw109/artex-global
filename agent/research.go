package agent

// research.go — the vulnerability-research discipline overlay ("研究模式"). It is
// an OPTIONAL, additive block appended to the planner/worker system prompt in the
// exact same way as constraintBlock: it never rewrites the tuned base prompt, it
// only layers a four-phase research cycle (surface → novel → exploit → verify),
// source-first with live verification, and an objective bias toward RCE, on top.
//
// Gated by a resolver (SetResearchMode); nil/false ⇒ the block is empty and the
// orchestration behaves exactly as before. This mirrors the constraint-injection
// toggle so the mode can be flipped per deployment without rebuilding the agents.

// researchModeBlock renders the four-phase research discipline for the given role
// ("planner" or "worker"). Empty string disables it. The framing sits BELOW the
// operation constraints (which stay highest priority) but ABOVE the base
// exploration heuristics, so it shapes HOW the loop explores without overriding a
// declared boundary.
func researchModeBlock(role string) string {
	// Shared discipline — the four phases, source-first + verify-live, and the
	// RCE objective. Written in the base prompt's language to preserve the tuned
	// agents' benchmarked reasoning.
	const core = `

【研究模式（漏洞研究纪律，叠加在下方常规启发式之上；操作约束仍最高优先）】：
本任务按"漏洞研究"而非泛化渗透推进，遵循四阶段循环（是纪律，不是僵硬流程）：
1. **攻击面(surface)**：优先在【源码/反编译代码】上铺开攻击面——用 Grep/Read/Bash 静态定位危险 sink（反序列化、命令执行、模板/表达式求值、路径与上传、认证前可达面、类型混淆等）与可疑数据流，而非一上来就打 live。把候选面连同 file:line 锚点记成事实/资产。
2. **新颖推理(novel)**：对每个候选面推理【具体、可触发】的漏洞假设——攻击者输入如何到达 sink、需要哪些前置原语。只挑能形成可验证假设的方向，给出明确的触发路径，而非含糊"疑似"。
3. **利用/PoC(exploit)**：为前置已满足的假设构造最小可触发 PoC；串行利用链一步步来，后一步依赖前一步的真实产出。能把多个原语在本会话内串成一条可触发链正是研究的价值所在。
4. **验证(verify)**：在授权 live 目标上独立触发该 PoC 取得可复现证据，或在无 live 条件时以源码级数据流证明可达——用与首次【不同的路径/命令】复现，证伪与证实同等有价值；仅凭版本/CVE/"看起来可注入"绝不算已确认。

贯穿纪律：
- **source-first, verify-live**：源码上找 bug，live 目标上验证/触发 PoC；两者都在授权范围内。
- **以距 RCE 的距离排序**：优先能【链成 RCE/提权/数据外泄】的原语与路径，把高价值入口往深打穿，而不是为覆盖度铺广浅测。
- **不重复、不空转**：每个方向只认真走一次；确认封锁的方向，仅出现材料性新机理才重开（呼应下方反重复规则）。`

	switch role {
	case "planner":
		return core + `
作为规划者：把意图按上述四阶段表述并排期——先派 surface（源码攻击面）意图，据其产出的 fact 再派 novel→exploit→verify；串行链用 TodoWrite 分步派，别并行front-load。`
	case "worker":
		return core + `
作为执行者：先判断本意图属于哪个阶段，只做该阶段的活并把结果写回知识图谱；验证阶段严格执行对抗式自检，未独立复现不得上报 finding。`
	default:
		return core
	}
}
