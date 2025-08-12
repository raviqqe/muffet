import sitemap from "@astrojs/sitemap";
import starlight from "@astrojs/starlight";
import { defineConfig } from "astro/config";

export default defineConfig({
  base: "/muffet",
  integrations: [
    sitemap(),
    starlight({
      customCss: ["./src/index.css"],
      favicon: "/icon.svg",
      head: [
        {
          attrs: {
            href: "/muffet/manifest.json",
            rel: "manifest",
          },
          tag: "link",
        },
        {
          attrs: {
            content: "/muffet/icon.svg",
            property: "og:image",
          },
          tag: "meta",
        },
        {
          attrs: {
            "data-domain": "raviqqe.com",
            defer: true,
            src: "https://plausible.io/js/plausible.js",
          },
          tag: "script",
        },
      ],
      logo: {
        src: "./src/icon.svg",
      },
      sidebar: [
        {
          label: "Home",
          link: "/",
        },
        {
          label: "Install",
          link: "/install",
        },
        {
          label: "Usage",
          link: "/usage",
        },
      ],
      social: [
        {
          href: "https://github.com/raviqqe/muffet",
          icon: "github",
          label: "GitHub",
        },
      ],
      title: "Muffet",
    }),
  ],
  prefetch: { prefetchAll: true },
  site: "https://raviqqe.github.io/muffet",
});
