package base

import "strings"

type Tense string

const (
	PresentSimple            Tense = "present simple"
	PresentContinuous        Tense = "present continuous"
	PresentPerfect           Tense = "present perfect"
	PresentPerfectContinuous Tense = "present perfect continuous"

	PastSimple            Tense = "past simple"
	PastContinuous        Tense = "past continuous"
	PastPerfect           Tense = "past perfect"
	PastPerfectContinuous Tense = "past perfect continuous"

	FutureSimple            Tense = "future simple"
	FutureContinuous        Tense = "future continuous"
	FuturePerfect           Tense = "future perfect"
	FuturePerfectContinuous Tense = "future perfect continuous"

	ConditionalSimple  Tense = "conditional simple"
	ConditionalPerfect Tense = "conditional perfect"

	Imperative         Tense = "imperative"
	SubjunctivePresent Tense = "subjunctive present"
	SubjunctivePast    Tense = "subjunctive past"
	Infinitive         Tense = "infinitive"
	Gerund             Tense = "gerund"
	Participle         Tense = "participle"
)

func GetTense(rawTense string) Tense {
	rawTense = strings.ToLower(rawTense)
	rawTense = strings.ReplaceAll(rawTense, "-", " ")

	switch rawTense {
	case "present simple":
		return PresentSimple
	case "present continuous":
		return PresentContinuous
	case "present perfect":
		return PresentPerfect
	case "present perfect continuous":
		return PresentPerfectContinuous

	case "past simple":
		return PastSimple
	case "past continuous":
		return PastContinuous
	case "past perfect":
		return PastPerfect
	case "past perfect continuous":
		return PastPerfectContinuous

	case "future simple":
		return FutureSimple
	case "future continuous":
		return FutureContinuous
	case "future perfect":
		return FuturePerfect
	case "future perfect continuous":
		return FuturePerfectContinuous

	case "conditional simple":
		return ConditionalSimple
	case "conditional perfect":
		return ConditionalPerfect

	case "imperative":
		return Imperative
	case "subjunctive present":
		return SubjunctivePresent
	case "subjunctive past":
		return SubjunctivePast
	case "infinitive":
		return Infinitive
	case "gerund":
		return Gerund
	case "participle":
		return Participle
	}

	return PresentSimple
}
