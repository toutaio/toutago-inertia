<template>
  <div class="auth-login">
    <div class="columns is-centered">
      <div class="column is-5-tablet is-4-desktop">
        <div class="box">
          <h1 class="title has-text-centered">Login</h1>

          <form @submit.prevent="login">
            <div class="field">
              <label class="label">Email</label>
              <div class="control">
                <input
                  v-model="form.email"
                  type="email"
                  class="input"
                  :class="{ 'is-danger': form.errors.email }"
                  placeholder="you@example.com"
                />
              </div>
              <p v-if="form.errors.email" class="help is-danger">
                {{ form.errors.email }}
              </p>
            </div>

            <div class="field">
              <label class="label">Password</label>
              <div class="control">
                <input
                  v-model="form.password"
                  type="password"
                  class="input"
                  :class="{ 'is-danger': form.errors.password }"
                  placeholder="••••••••"
                />
              </div>
              <p v-if="form.errors.password" class="help is-danger">
                {{ form.errors.password }}
              </p>
            </div>

            <div class="field">
              <div class="control">
                <button
                  type="submit"
                  class="button is-primary is-fullwidth"
                  :class="{ 'is-loading': form.processing }"
                  :disabled="form.processing"
                >
                  Login
                </button>
              </div>
            </div>
          </form>

          <div class="notification is-info is-light mt-4">
            <p class="has-text-centered">
              <strong>Demo credentials:</strong><br />
              Email: demo@example.com<br />
              Password: password
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useForm } from '@toutaio/inertia-vue'
import type { LoginPageProps } from '../../types'

defineProps<LoginPageProps>()

const form = useForm({
  email: '',
  password: '',
})

const login = () => {
  form.post('/login')
}
</script>

<style scoped>
.auth-login {
  min-height: 50vh;
  padding-top: 3rem;
}

.mt-4 {
  margin-top: 1.5rem;
}
</style>
