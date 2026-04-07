import ogImage from "@/assets/og-image.png";

export const siteConfig = {
  name: "Metric Neo",
  description:
    "Offline-first chronograph data platform for shooting sports. Manage profiles, projectiles, and sessions — fully air-gap ready.",
  url: "https://metric-neo.wellermann.de",
  lang: "en",
  locale: "en_US",
  author: "Daniel Wellermann",
  ogImage: ogImage,
  socialLinks: {
    github: "https://github.com/dwellermann/metric-neo",
    githubProfile: "https://github.com/dwellermann",
    releases: "https://github.com/dwellermann/metric-neo/releases/latest",
  },
  navLinks: [
    { text: "Home", href: "/" },
    { text: "Manual", href: "/manual" },
    { text: "Changelog", href: "/changelog" },
    { text: "About", href: "/about" },
  ],
  navLinksDE: [
    { text: "Start", href: "/de" },
    { text: "Handbuch", href: "/de/manual" },
    { text: "Änderungen", href: "/de/changelog" },
    { text: "Über", href: "/de/about" },
  ],
};
