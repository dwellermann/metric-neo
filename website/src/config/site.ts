import ogImage from "@/assets/og-image.png";

export const siteConfig = {
  name: "Metric Neo",
  description:
    "Offline-first chronograph data platform for shooting sports. Manage profiles, projectiles, and sessions — fully air-gap ready.",
  url: "https://metric-neo.wellermann.de",
  lang: "de",
  locale: "de_DE",
  author: "Daniel Wellermann",
  ogImage: ogImage,
  socialLinks: {
    github: "https://github.com/dwellermann/metric-neo",
    githubProfile: "https://github.com/dwellermann",
    releases: "https://github.com/dwellermann/metric-neo/releases/latest",
  },
  navLinks: [
    { text: "Start", href: "/" },
    { text: "Handbuch", href: "/manual" },
    { text: "Änderungen", href: "/changelog" },
    { text: "Über", href: "/about" },
  ],
  navLinksEN: [
    { text: "Home", href: "/en" },
    { text: "Manual", href: "/en/manual" },
    { text: "Changelog", href: "/en/changelog" },
    { text: "About", href: "/en/about" },
  ],
};
