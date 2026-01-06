import { Component, DefineComponent, h } from 'vue';

/**
 * Wraps a page component with a layout component
 * @param page - The page component to wrap
 * @param layout - The layout component to use
 * @returns A new component with the layout applied
 */
export function withLayout(
  page: Component | DefineComponent,
  layout: Component | DefineComponent
): DefineComponent {
  return {
    ...page,
    setup(props, ctx) {
      // Call the original setup if it exists
      const pageSetup = typeof page.setup === 'function' 
        ? page.setup(props, ctx)
        : {};

      return () => h(layout, {}, {
        default: () => h(page, props)
      });
    }
  } as DefineComponent;
}

/**
 * Type for components that have a layout property
 */
export interface PageWithLayout extends Component {
  layout?: Component | DefineComponent | (() => Component | DefineComponent);
}

/**
 * Resolves the layout for a page component
 * @param page - The page component
 * @returns The layout component or undefined
 */
export function resolvePageLayout(page: PageWithLayout): Component | DefineComponent | undefined {
  if (!page.layout) {
    return undefined;
  }

  return typeof page.layout === 'function' ? page.layout() : page.layout;
}
