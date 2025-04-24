import { atom, computed } from "nanostores";
import { compile as compileProgram } from "../application/compile.js";
import { interpret as interpretProgram } from "../application/interpret.js";

export const source = atom(
  `
(import
  (scheme base)
  (scheme read)
  (scheme write))

(define (fibonacci x)
  (if (< x 2)
     x
     (+
        (fibonacci (- x 1))
        (fibonacci (- x 2)))))

(write \`(answer ,(fibonacci (read))))
  `.trim(),
);

const bytecodes = atom<Uint8Array | null>(new Uint8Array());

export const bytecodesReady = computed(
  bytecodes,
  (bytecodes) => bytecodes?.length,
);

export const compiling = computed(bytecodes, (output) => output === null);

const output = atom<Uint8Array | null>(new Uint8Array());

export const input = atom("10");
export const textOutput = computed(output, (output) =>
  output === null ? null : new TextDecoder().decode(output),
);

export const outputUrlStore = computed(output, (output) =>
  output?.length ? URL.createObjectURL(new Blob([output])) : null,
);

export const interpretingStore = computed(output, (output) => output === null);

export const compilerError = atom("");

export const interpreterError = atom("");

export const compile = async (): Promise<void> => {
  bytecodes.set(null);
  compilerError.set("");

  let value = new Uint8Array();

  try {
    value = await compileProgram(source.get());
  } catch (error) {
    compilerError.set((error as Error).message);
  }

  bytecodes.set(value);
};

export const interpret = async (): Promise<void> => {
  const bytecodeValue = bytecodes.get();

  if (!bytecodeValue) {
    return;
  }

  output.set(null);
  interpreterError.set("");

  let outputValue = new Uint8Array();

  try {
    outputValue = await interpretProgram(
      bytecodeValue,
      new TextEncoder().encode(input.get()),
    );
  } catch (error) {
    interpreterError.set((error as Error).message);
  }

  output.set(outputValue);
};
