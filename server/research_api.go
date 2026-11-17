package server

// research_api.go — the optional "research mode" toggle. When on, the four-phase
// vulnerability-research discipline (surface → novel → exploit → verify,
// source-first + verify-live, RCE objective) is appended to the planner/worker
// prompts via agent.researchModeBlock. Off by default, so existing orchestration
// behaves exactly as before. The toggle is a plain setting (flip via the settings
// API / DB), read per round by the resolvers below — the same pattern as the
// constraint-injection switches.

// settingResearchMode gates the research overlay for both planner and worker.
// Default false (GetBool's second arg) — the mode is opt-in.
const settingResearchMode = "research_mode_v1"

// researchModePlanner / researchModeWorker report whether the four-phase research
// discipline is appended to the respective agent's prompt. Wired as resolvers to
// the planner/worker, read each round so flipping the setting takes effect without
// rebuilding the agents.
func (s *Server) researchModePlanner() bool {
	return s.m.pg.GetBool(settingResearchMode, false)
}

func (s *Server) researchModeWorker() bool {
	return s.m.pg.GetBool(settingResearchMode, false)
}
