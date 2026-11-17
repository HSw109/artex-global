package agent

import "strings"

import "testing"

func TestResearchModeBlockContent(t *testing.T) {
	b := researchModeBlock("")
	for _, want := range []string{"surface", "novel", "exploit", "verify", "source-first", "RCE"} {
		if !strings.Contains(b, want) {
			t.Fatalf("research block missing %q:\n%s", want, b)
		}
	}
}

func TestResearchModeBlockRoleTailoring(t *testing.T) {
	if !strings.Contains(researchModeBlock("planner"), "作为规划者") {
		t.Fatal("planner block missing planner-specific guidance")
	}
	if !strings.Contains(researchModeBlock("worker"), "作为执行者") {
		t.Fatal("worker block missing worker-specific guidance")
	}
	// The shared core is present in both.
	if !strings.Contains(researchModeBlock("planner"), "研究模式") ||
		!strings.Contains(researchModeBlock("worker"), "研究模式") {
		t.Fatal("shared research-mode core missing from a role block")
	}
}

func TestResearchModeGatingDefaultOff(t *testing.T) {
	// nil resolver ⇒ off, so existing orchestration is unchanged until enabled.
	if (&Planner{}).wantResearch() {
		t.Fatal("planner research must default OFF with a nil resolver")
	}
	if (&Worker{}).wantResearch() {
		t.Fatal("worker research must default OFF with a nil resolver")
	}
	// Explicit resolvers flip it.
	pOn := &Planner{}
	pOn.SetResearchMode(func() bool { return true })
	if !pOn.wantResearch() {
		t.Fatal("planner research should be ON when the resolver returns true")
	}
	wOff := &Worker{}
	wOff.SetResearchMode(func() bool { return false })
	if wOff.wantResearch() {
		t.Fatal("worker research should be OFF when the resolver returns false")
	}
}
