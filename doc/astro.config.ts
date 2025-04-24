import { readFile, readdir, stat } from "node:fs/promises";
import { join, parse } from "node:path";
import sitemap from "@astrojs/sitemap";
import starlight from "@astrojs/starlight";
import { defineConfig } from "astro/config";
import { capitalize, sortBy } from "es-toolkit";

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
  integrations: [
    sitemap(),
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
          label: "Install",
          link: "/install",
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
