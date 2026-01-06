import { createSSRApp, Component, h } from 'vue';
import { renderToString } from 'vue/server-renderer';
import type { Page } from './types';

export interface SSRContext {
  page: Page;
  resolveComponent: (name: string) => Component;
}

export interface SSRResult {
  head: string[];
  body: string;
}

/**
 * Server-side rendering for Inertia.js with Vue
 * 
 * @param context - SSR context containing page data and component resolver
 * @returns Promise resolving to HTML head and body
 */
export async function createInertiaSSRApp(context: SSRContext): Promise<SSRResult> {
  const { page, resolveComponent } = context;
  
  // Resolve the component for this page
  const component = resolveComponent(page.component);
  
  if (!component) {
    throw new Error(`Component "${page.component}" not found`);
  }
  
  // Create SSR app instance
  const app = createSSRApp({
    render: () => h(component, {
      ...page.props,
      // Provide Inertia context to the component
      initialPage: page,
      initialComponent: page.component,
    }),
  });
  
  // Render to string
  const body = await renderToString(app);
  
  // Extract head tags (title, meta tags, etc.)
  // This is a simplified version - in production you'd use something like @vueuse/head
  const head: string[] = [];
  
  // Add page version as meta tag
  if (page.version) {
    head.push(`<meta name="inertia-version" content="${page.version}">`);
  }
  
  return {
    head,
    body,
  };
}

/**
 * Creates a minimal SSR page template
 * 
 * @param body - Rendered body HTML
 * @param page - Page data for hydration
 * @param head - Array of head tags
 * @returns Complete HTML page
 */
export function createSSRPage(body: string, page: Page, head: string[] = []): string {
  const pageJson = JSON.stringify(page).replace(/</g, '\\u003c');
  
  return `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  ${head.join('\n  ')}
</head>
<body>
  <div id="app" data-page='${pageJson}'>${body}</div>
</body>
</html>`;
}
