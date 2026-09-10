# Chatbot TO BE — Afirmativo, Tiempo Pasado
## Autómatas, Gramáticas y Lenguajes

Chatbot en Go que valida frases en inglés con el verbo TO BE
**exclusivamente en su forma AFIRMATIVA y tiempo PASADO** (was/were),
usando expresiones regulares.

Frases en presente, negativas o interrogativas se reconocen y se
rechazan con un mensaje explicando por qué no son válidas para este caso.

## Requisitos
- Go 1.21+

## Ejecutar
```
go run .
```

## Compilar
```
go build -o tobebot .
./tobebot
```

## Estructura
- `main.go`       — bucle conversacional del chatbot (saludo, entrada/salida, control del ciclo).
- `validator.go`  — expresión regular principal (afirmativa + pasado) y lógica de
  concordancia sujeto-verbo (was/were), más ER de apoyo para detectar y explicar
  otros casos inválidos (presente, negativa, pregunta).

## Ejemplos válidos
- `I was tired.`
- `He was happy.`
- `You were a good student.`
- `The dog was furious.`
- `Maria was sick last week.`
- `These pencils were black.`

## Ejemplos inválidos (con motivo)
- `I am a teacher.` → presente, no pasado
- `The cat was not Brown.` → negativa, no afirmativa
- `Were you happy?` → interrogativa, no afirmativa
- `They was here.` → mala concordancia (debería ser "were")
