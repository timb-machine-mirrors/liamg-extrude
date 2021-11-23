package elf

import "github.com/liamg/extrude/pkg/report"

func (m *Metadata) CreateReport() (report.Report, error) {
	rep := report.New()

	overview := report.NewSection("Overview")

	overview.AddKeyValue("File", m.File.Path)
	overview.AddKeyValue("Format", m.File.Format.String())

	if m.ELF != nil {
		overview.AddKeyValue("Platform", m.ELF.Machine.String())
		overview.AddKeyValue("OS/ABI", m.ELF.OSABI.String())
		overview.AddKeyValue("Type", m.ELF.Type.String())
		overview.AddKeyValue("Byte Order", m.ELF.ByteOrder.String())
	}

	overview.AddKeyValue("Compiler Name", m.Compiler.Name)
	overview.AddKeyValue("Compiler Version", m.Compiler.Version)
	overview.AddKeyValue("Source Language", m.Compiler.Language)

	rep.AddSection(overview)

	security := report.NewSection("Security")

	security.AddTest(
		"Position Independent Executable (PIE)",
		boolToResult(m.Hardening.PositionIndependent),
		`A PIE binary and all of its dependencies are loaded into random locations within virtual memory each time the application is executed. This makes Return Oriented Programming (ROP) attacks much more difficult to execute reliably.`,
	)

	security.AddTest(
		"Read-Only Relocations (RELRO)",
		boolToResult(m.Hardening.ReadOnlyRelocations),
		`Hardens ELF programs against loader memory area overwrites by having the loader mark any areas of the relocation table as read-only for any symbols resolved at load-time ("read-only relocations"). This reduces the area of possible GOT-overwrite-style memory corruption attacks. `,
	)

	security.AddTest(
		"Fortified Source Functions",
		boolToResult(m.Hardening.FortifySourceFunctions),
		`This is a security feature which applies to GLIBC functions vulnerable to buffer overflow attacks. It overrides the use of such functions with a safe variation and is enabled by default on most Linux platforms. If GLIBC functions are used within the binary, this test will fail if none are fortified.`,
	)

	security.AddTest(
		"Stack Protection",
		boolToResult(m.Hardening.StackProtected),
		`The basic idea behind stack protection is to push a "canary" (a randomly chosen integer) on the stack just after the function return pointer has been pushed. The canary value is then checked before the function returns; if it has changed, the program will abort. Generally, stack buffer overflow (aka "stack smashing") attacks will have to change the value of the canary as they write beyond the end of the buffer before they can get to the return pointer. Since the value of the canary is unknown to the attacker, it cannot be replaced by the attack. Thus, the stack protection allows the program to abort when that happens rather than return to wherever the attacker wanted it to go.`,
	)

	rep.AddSection(security)

	if len(m.Notes) > 0 {
		notes := report.NewSection("Other Findings")
		for _, note := range m.Notes {
			notes.AddTest(note.Heading, report.Warning, note.Content)
		}
		rep.AddSection(notes)
	}

	return rep, nil
}

func boolToResult(in bool) report.Result {
	if in {
		return report.Pass
	}
	return report.Fail
}
