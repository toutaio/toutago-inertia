<template>
  <div class="todos-index">
    <div class="level">
      <div class="level-left">
        <div class="level-item">
          <h1 class="title">Todos</h1>
        </div>
      </div>
      <div class="level-right">
        <div class="level-item">
          <button @click="showCreateForm = true" class="button is-primary">
            + New Todo
          </button>
        </div>
      </div>
    </div>

    <!-- Filter tabs -->
    <div class="tabs">
      <ul>
        <li :class="{ 'is-active': filter.status === 'all' }">
          <Link href="/todos?status=all">All</Link>
        </li>
        <li :class="{ 'is-active': filter.status === 'active' }">
          <Link href="/todos?status=active">Active</Link>
        </li>
        <li :class="{ 'is-active': filter.status === 'completed' }">
          <Link href="/todos?status=completed">Completed</Link>
        </li>
      </ul>
    </div>

    <!-- Create form modal -->
    <div class="modal" :class="{ 'is-active': showCreateForm }">
      <div class="modal-background" @click="showCreateForm = false"></div>
      <div class="modal-card">
        <header class="modal-card-head">
          <p class="modal-card-title">Create Todo</p>
          <button class="delete" @click="showCreateForm = false"></button>
        </header>
        <section class="modal-card-body">
          <form @submit.prevent="createTodo">
            <div class="field">
              <label class="label">Title</label>
              <div class="control">
                <input
                  v-model="createForm.title"
                  type="text"
                  class="input"
                  :class="{ 'is-danger': createForm.errors.title }"
                  placeholder="Enter todo title"
                />
              </div>
              <p v-if="createForm.errors.title" class="help is-danger">
                {{ createForm.errors.title }}
              </p>
            </div>

            <div class="field">
              <label class="label">Description</label>
              <div class="control">
                <textarea
                  v-model="createForm.description"
                  class="textarea"
                  placeholder="Enter description (optional)"
                ></textarea>
              </div>
            </div>
          </form>
        </section>
        <footer class="modal-card-foot">
          <button
            @click="createTodo"
            class="button is-primary"
            :class="{ 'is-loading': createForm.processing }"
            :disabled="createForm.processing"
          >
            Create
          </button>
          <button @click="showCreateForm = false" class="button">Cancel</button>
        </footer>
      </div>
    </div>

    <!-- Todos list -->
    <div v-if="todos.length === 0" class="notification">
      <p>No todos found. Create one to get started!</p>
    </div>

    <div v-else class="todos-list">
      <div
        v-for="todo in todos"
        :key="todo.id"
        class="box todo-item"
        :class="{ 'is-completed': todo.completed }"
      >
        <div class="columns is-vcentered">
          <div class="column is-narrow">
            <input
              type="checkbox"
              :checked="todo.completed"
              @change="toggleComplete(todo)"
              class="checkbox"
            />
          </div>
          <div class="column">
            <h3 class="title is-5" :class="{ 'has-text-grey-light': todo.completed }">
              {{ todo.title }}
            </h3>
            <p v-if="todo.description" class="subtitle is-6">
              {{ todo.description }}
            </p>
          </div>
          <div class="column is-narrow">
            <div class="buttons">
              <Link :href="`/todos/${todo.id}/edit`" class="button is-small">
                Edit
              </Link>
              <button
                @click="deleteTodo(todo)"
                class="button is-small is-danger"
                :disabled="deleteForm.processing"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Link, useForm } from '@toutaio/inertia-vue'
import type { TodosListPageProps, Todo } from '../../types'

const props = defineProps<TodosListPageProps>()

const showCreateForm = ref(false)

const createForm = useForm({
  title: '',
  description: '',
})

const deleteForm = useForm({})

const createTodo = () => {
  createForm.post('/todos', {
    onSuccess: () => {
      showCreateForm.value = false
      createForm.reset()
    },
  })
}

const toggleComplete = (todo: Todo) => {
  useForm({
    title: todo.title,
    description: todo.description,
    completed: !todo.completed,
  }).put(`/todos/${todo.id}`)
}

const deleteTodo = (todo: Todo) => {
  if (confirm(`Are you sure you want to delete "${todo.title}"?`)) {
    deleteForm.delete(`/todos/${todo.id}`)
  }
}
</script>

<style scoped>
.todos-list {
  margin-top: 2rem;
}

.todo-item {
  transition: opacity 0.2s;
}

.todo-item.is-completed {
  opacity: 0.6;
}

.checkbox {
  width: 20px;
  height: 20px;
  cursor: pointer;
}
</style>
