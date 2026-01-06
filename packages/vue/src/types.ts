import { App, Component, DefineComponent } from 'vue';

export type Method = 'get' | 'post' | 'put' | 'patch' | 'delete';

export interface PageProps {
  [key: string]: any;
}

export interface Page<TPageProps = PageProps> {
  component: string;
  props: TPageProps;
  url: string;
  version: string;
}

export interface VisitOptions {
  method?: Method;
  data?: any;
  replace?: boolean;
  preserveScroll?: boolean;
  preserveState?: boolean;
  only?: string[];
  headers?: Record<string, string>;
  onBefore?: () => void;
  onStart?: () => void;
  onProgress?: (progress: ProgressEvent) => void;
  onSuccess?: (page: Page) => void;
  onError?: (errors: Record<string, string>) => void;
  onFinish?: () => void;
}

export type PageResolver = (name: string) => Component | DefineComponent | Promise<Component | DefineComponent>;

export interface InertiaAppOptions {
  page?: Page;
  resolve: PageResolver;
  setup: {
    component: Component | DefineComponent;
    plugin?: (args: { app: App; props: { initialPage: Page } }) => void;
  };
}

export interface Errors {
  [key: string]: string;
}

export interface ErrorBag {
  [key: string]: string | string[];
}
