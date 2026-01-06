import {
  Component,
  InjectionKey,
  inject,
  reactive,
} from 'vue';
import { Page, PageProps } from './types';

export const PagePropsKey: InjectionKey<PageProps> = Symbol('PageProps');
export const PageKey: InjectionKey<Page> = Symbol('Page');

export interface InertiaAppState {
  component: Component | null;
  page: Page;
  key: number;
}

export function createAppState(initialPage: Page): InertiaAppState {
  return reactive({
    component: null,
    page: initialPage,
    key: 0,
  });
}

export function usePageProps<T = PageProps>(): T {
  const props = inject(PagePropsKey);
  if (!props) {
    throw new Error('Page props not found. Make sure you are using Inertia app.');
  }
  return props as T;
}

export function usePage<T = PageProps>(): Page<T> {
  const page = inject(PageKey);
  if (!page) {
    throw new Error('Page not found. Make sure you are using Inertia app.');
  }
  return page as Page<T>;
}
