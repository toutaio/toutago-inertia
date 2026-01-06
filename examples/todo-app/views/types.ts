// Auto-generated TypeScript types from Go structs
// Do not edit manually

export interface Todo {
  id: number
  title: string
  description: string
  completed: boolean
  createdAt: string
  updatedAt: string
}

export interface TodosFilter {
  status: string
  search: string
}

export interface User {
  id: number
  name: string
  email: string
}

export interface HomePageProps {
  greeting: string
  user?: User
}

export interface TodosListPageProps {
  todos: Todo[]
  filter: TodosFilter
  flash?: Record<string, string>
}

export interface TodosEditPageProps {
  todo: Todo
  flash?: Record<string, string>
}

export interface LoginPageProps {
  flash?: Record<string, string>
}
