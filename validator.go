package main

import (
	"regexp"
	"strings"
)

// ---------------------------------------------------------------------
// EXPRESIONES REGULARES DEL PROYECTO
// Caso asignado a este grupo: TO BE, forma AFIRMATIVA, tiempo PASADO.
//   SUJETO   = I | (He|She|It) | (You|We|They) | (The N) | ((This|That|These|Those) N) | PN
//   VERBO_PA = was | were
//   COMPL    = (letra (letra|espacio)*)
// ---------------------------------------------------------------------

// SUJETO = I | (He|She|It|You|We|They) | ((The|This|That|These|Those) N) | PN
const subj = `(I|(?i:He|She|It|You|We|They)|(?:(?i:The|This|That|These|Those) [a-z]+)|[A-Z][a-z]+)`

// COMPLEMENTO = letra (letra | espacio)*
const compl = `([A-Za-z][A-Za-z ]*)`

var (
	// (SUJETO)(was|were)(COMPLEMENTO).  -> ÚNICA forma válida: afirmativa pasado
	reAffPast = regexp.MustCompile(`^` + subj + ` (was|were) ` + compl + `\.$`)

	// ER de apoyo para clasificar el tipo de sujeto (insensibles a
	// mayúsculas porque el artículo/demostrativo puede ir en minúscula).
	reDemonstrativeSg = regexp.MustCompile(`(?i)^(this|that) ([a-z]+)$`)
	reDemonstrativePl = regexp.MustCompile(`(?i)^(these|those) ([a-z]+)$`)
	reCommonNoun      = regexp.MustCompile(`(?i)^the ([a-z]+)$`)
	reProperNoun      = regexp.MustCompile(`^[A-Z][a-z]+$`)

	// ER "laxa" para diagnosticar qué escribió el usuario aunque no
	// cumpla la forma exigida (detecta presente, negativa o pregunta
	// para dar feedback específico en vez de un "inválida" genérico).
	reLoosePresent  = regexp.MustCompile(`^` + subj + ` (am|is|are) (not )?` + compl + `\.$`)
	reLooseNegative = regexp.MustCompile(`^` + subj + ` (was|were) not ` + compl + `\.$`)
	reLooseQuestion = regexp.MustCompile(`^(Was|Were|Am|Is|Are) ` + subj + ` ` + compl + `\?$`)
)

var irregularPlurals = map[string]bool{
	"children": true, "men": true, "women": true, "people": true,
}

// isPlural decide si un sustantivo común es plural. Regla: termina en
// "s" (y no es una excepción de una sola sílaba tipo "bus", "gas") o
// pertenece a la lista de plurales irregulares.
func isPlural(noun string) bool {
	if irregularPlurals[noun] {
		return true
	}
	exceptions := map[string]bool{"bus": true, "gas": true, "glass": true}
	if exceptions[noun] {
		return false
	}
	return len(noun) > 1 && noun[len(noun)-1] == 's'
}

// expectedPastVerb calcula, para un sujeto dado, cuál es la forma
// correcta del verbo TO BE en pasado: was o were.
func expectedPastVerb(subject string) (string, bool) {
	switch strings.ToLower(subject) {
	case "i", "he", "she", "it":
		return "was", true
	case "you", "we", "they":
		return "were", true
	}

	if reDemonstrativeSg.MatchString(subject) {
		return "was", true
	}
	if reDemonstrativePl.MatchString(subject) {
		return "were", true
	}
	if m := reCommonNoun.FindStringSubmatch(subject); m != nil {
		if isPlural(m[1]) {
			return "were", true
		}
		return "was", true
	}
	if reProperNoun.MatchString(subject) {
		return "was", true
	}
	return "", false
}

// ValidationResult resume el veredicto de una frase.
type ValidationResult struct {
	Valid   bool
	Message string
}

// Validate solo acepta frases TO BE afirmativas en tiempo pasado
// (was/were), con la concordancia sujeto-verbo correcta. Cualquier
// otro caso (presente, negativa, interrogativa o mala concordancia)
// se marca inválido con un mensaje que explica el motivo.
func Validate(sentence string) ValidationResult {
	// La forma negativa se revisa ANTES que la afirmativa: una frase
	// negativa ("The cat was not Brown.") también calza de forma laxa
	// con la ER afirmativa (compl acepta "not Brown" como complemento),
	// así que la más específica debe evaluarse primero.
	if m := reLooseNegative.FindStringSubmatch(sentence); m != nil { //FindStringSubmatch FindStringSubmatch
		return ValidationResult{false, "Invalid: this is a negative sentence ('" + m[1] +
			" " + m[2] + " not ...'). This project only accepts the affirmative form."}
	}

	if m := reAffPast.FindStringSubmatch(sentence); m != nil {
		subject, verb := m[1], m[2]
		if exp, ok := expectedPastVerb(subject); ok && exp != verb {
			return ValidationResult{false, "Invalid sentence: '" + verb +
				"' does not agree with '" + subject + "'. Expected '" + exp + "'."}
		}
		return ValidationResult{true, "Correct sentence: affirmative, past tense."}
	}

	if m := reLoosePresent.FindStringSubmatch(sentence); m != nil {
		return ValidationResult{false, "Invalid: this is in present tense ('" + m[2] +
			"'). This project only accepts the past tense (was/were)."}
	}
	if reLooseQuestion.MatchString(sentence) { //aqui usamos MatchString para ver si una cadena de texto coincide con la expresión regular esto devuelve true o false
		return ValidationResult{false, "Invalid: this is a question. This project only accepts affirmative statements."}
	}

	return ValidationResult{false, "Invalid sentence: it does not match the affirmative past tense structure (Subject + was/were + complement)."}
}
