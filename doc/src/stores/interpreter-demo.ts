import { atom, computed, task } from "nanostores";
import { run as runProgram } from "../application/run.js";

export const source = atom(
  `
(import (scheme base) (scheme write))

(define (fibonacci x)
  (if (< x 2)
     x
     (+
        (fibonacci (- x 1))
        (fibonacci (- x 2)))))

(display "Answer: ")
(write (fibonacci 10))
(newline)
  `.trim(),
);

const run = computed(source, (source) =>
  task(async () => {
    try {
      return await runProgram(source);
    } catch (error) {
      return error as Error;
    }
  }),
);

export const output = computed(run, (output) =>
  output instanceof Error ? null : new TextDecoder().decode(output),
);

export const error = computed(run, (error) =>
  error instanceof Error ? error : null,
);
