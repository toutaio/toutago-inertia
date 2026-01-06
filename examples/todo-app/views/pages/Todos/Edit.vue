<template>
  <div class="todos-edit">
    <h1 class="title">Edit Todo</h1>

    <div class="box">
      <form @submit.prevent="updateTodo">
        <div class="field">
          <label class="label">Title</label>
          <div class="control">
            <input
              v-model="form.title"
              type="text"
              class="input"
              :class="{ 'is-danger': form.errors.title }"
              placeholder="Enter todo title"
            />
          </div>
          <p v-if="form.errors.title" class="help is-danger">
            {{ form.errors.title }}
          </p>
        </div>

        <div class="field">
          <label class="label">Description</label>
          <div class="control">
            <textarea
              v-model="form.description"
              class="textarea"
              placeholder="Enter description (optional)"
            ></textarea>
          </div>
        </div>

        <div class="field">
          <div class="control">
            <label class="checkbox">
              <input v-model="form.completed" type="checkbox" />
              Completed
            </label>
          </div>
        </div>

        <div class="field is-grouped">
          <div class="control">
            <button
              type="submit"
              class="button is-primary"
              :class="{ 'is-loading': form.processing }"
              :disabled="form.processing"
            >
              Update Todo
            </button>
          </div>
          <div class="control">
            <Link href="/todos" class="button">Cancel</Link>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Link, useForm } from '@toutaio/inertia-vue'
import type { TodosEditPageProps } from '../../types'

const props = defineProps<TodosEditPageProps>()

const form = useForm({
  title: props.todo.title,
  description: props.todo.description,
  completed: props.todo.completed,
})

const updateTodo = () => {
  form.put(`/todos/${props.todo.id}`)
}
</script>
