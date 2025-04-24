import sitemap from "@astrojs/sitemap";
import starlight from "@astrojs/starlight";
import { defineConfig } from "astro/config";

export default defineConfig({
  base: "/muffet",
  integrations: [
    sitemap(),
    starlight({
      title: "Muffet",
      favicon: "/icon.svg",
      head: [
        {
          tag: "link",
          attrs: {
            rel: "manifest",
            href: "/muffet/manifest.json",
          },
        },
        {
          tag: "meta",
          attrs: {
            property: "og:image",
            content: "/muffet/icon.svg",
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
          href: "https://github.com/raviqqe/muffet",
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
        {
          label: "Usage",
          link: "/usage",
        },
      ],
    }),
  ],
  prefetch: { prefetchAll: true },
  site: "https://raviqqe.github.io/muffet",
});
