import configurations from "@raviqqe/eslint-config";
import solid from "eslint-plugin-solid/configs/typescript";

export default [
  ...configurations,
  {
    rules: {
      "@typescript-eslint/triple-slash-reference": "off",
      "import-x/order": "off",
      "perfectionist/sort-named-imports": "off",
      "react/jsx-no-useless-fragment": "off",
      "react/no-unknown-property": "off",
    },
  },
  solid,
];
