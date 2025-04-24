import { readFile, readdir, stat } from "node:fs/promises";
import { join, parse } from "node:path";
import sitemap from "@astrojs/sitemap";
import solid from "@astrojs/solid-js";
import starlight from "@astrojs/starlight";
import { defineConfig } from "astro/config";
import { capitalize, sortBy } from "es-toolkit";
import wasm from "vite-plugin-wasm";

type Item = { label: string; link: string } | { label: string; items: Item[] };

const documentDirectory = "src/content/docs";

const listItems = async (directory: string): Promise<Item[]> =>
  sortBy(
    await Promise.all(
      (await readdir(join(documentDirectory, directory)))
        .filter((path) => !path.startsWith("."))
        .map(async (path) => {
          const fullPath = join(documentDirectory, directory, path);
          const { name } = parse(path);
          const linkPath = join(directory, name);

          return (await stat(fullPath)).isDirectory()
            ? {
                label: capitalize(name.replace("-", " ")),
                items: await listItems(linkPath),
              }
            : {
                label:
                  (await readFile(fullPath, "utf-8"))
                    .split("\n")
                    .find((line) => line.startsWith("title: "))
                    ?.replace("title: ", "")
                    .trim() ?? "",
                link: linkPath,
              };
        }),
    ),
    [({ label, link }) => [!link, label]],
  );

export default defineConfig({
  base: "/stak",
  vite: {
    plugins: [wasm()],
    worker: {
      format: "es",
    },
  },
  integrations: [
    sitemap(),
    solid(),
    starlight({
      title: "Stak Scheme",
      customCss: ["./src/index.css"],
      favicon: "/icon.svg",
      head: [
        {
          tag: "link",
          attrs: {
            rel: "manifest",
            href: "/stak/manifest.json",
          },
        },
        {
          tag: "meta",
          attrs: {
            property: "og:image",
            content: "/stak/icon.svg",
          },
        },
      ],
      logo: {
        src: "./public/icon.svg",
      },
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/raviqqe/stak",
        },
      ],
      sidebar: [
        {
          label: "Home",
          link: "/",
        },
        {
          label: "Guides",
          items: [
            {
              label: "Install",
              link: "/install",
            },
            {
              label: "Embedding Scheme in Rust",
              link: "/embedding-scripts",
            },
            {
              label: "Running in no-std and no-alloc environment",
              link: "/no-std-no-alloc",
            },
            {
              label: "Hot reloading",
              link: "/hot-reload",
            },
            {
              label: "Writing a Scheme subset",
              link: "/writing-scheme-subset",
            },
          ],
        },
        {
          label: "Limitations",
          link: "/limitations",
        },
        {
          label: "Demo",
          items: [
            {
              label: "Interpreter",
              link: "/demo/interpreter",
            },
            {
              label: "Compiler",
              link: "/demo/compiler",
            },
          ],
        },

        {
          label: "Examples",
          items: await listItems("examples"),
        },
      ],
    }),
  ],
  prefetch: { prefetchAll: true },
  redirects: {
    "/demo": "/interpreter-demo",
  },
  site: "https://raviqqe.github.io/stak",
});
