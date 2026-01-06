import { App, createSSRApp, h } from 'vue';
import { InertiaAppOptions, Page, PageResolver } from './types';

export interface CreateInertiaAppSetupResult {
  app: App;
  props: {
    initialPage: Page;
    resolveComponent: PageResolver;
  };
}

export async function createInertiaApp(
  options: InertiaAppOptions
): Promise<CreateInertiaAppSetupResult> {
  const page = options.page || (window as any).INERTIA_PAGE;

  if (!page) {
    throw new Error(
      'Inertia page data is missing. Make sure the server is sending page data.'
    );
  }

  await options.resolve(page.component);

  const app = createSSRApp({
    render() {
      return h(options.setup.component, {
        initialPage: page,
        resolveComponent: options.resolve,
      });
    },
  });

  if (options.setup.plugin) {
    options.setup.plugin({ app, props: { initialPage: page } });
  }

  return {
    app,
    props: {
      initialPage: page,
      resolveComponent: options.resolve,
    },
  };
}
