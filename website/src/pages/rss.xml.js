import rss from "@astrojs/rss";
import { siteConfig } from "@/config/site";

export async function GET(context) {
  return rss({
    title: siteConfig.name,
    description: siteConfig.description,
    site: context.site,
    items: [],
    customData: `<language>en</language>`,
  });
}
