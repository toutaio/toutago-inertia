# @toutaio/inertia-vue

[![npm version](https://badge.fury.io/js/@toutaio%2Finertia-vue.svg)](https://www.npmjs.com/package/@toutaio/inertia-vue)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Vue 3 adapter for Inertia.js with the Toutā Go framework.

## Installation

```bash
npm install @toutaio/inertia-vue
```

## Quick Start

```js
import { createApp, h } from 'vue'
import { createInertiaApp } from '@toutaio/inertia-vue'

createInertiaApp({
  resolve: name => {
    const pages = import.meta.glob('./Pages/**/*.vue', { eager: true })
    return pages[`./Pages/${name}.vue`]
  },
  setup({ el, App, props, plugin }) {
    createApp({ render: () => h(App, props) })
      .use(plugin)
      .mount(el)
  },
})
```

## Features

- 🚀 Modern SPA experience without REST API complexity
- 🔄 Automatic page transitions
- 📝 Forms with validation
- 🔗 Smart link handling
- 💾 Form state persistence
- 🎯 TypeScript support

## Core Components

### Link

Navigate between pages without full page reloads:

```vue
<template>
  <Link href="/users">Users</Link>
  <Link href="/posts" method="post" :data="{ title: 'New Post' }">
    Create Post
  </Link>
</template>
```

### Head

Manage document head:

```vue
<template>
  <Head title="Dashboard">
    <meta name="description" content="User dashboard">
  </Head>
</template>
```

## Composables

### useForm

Handle forms with validation:

```vue
<script setup>
import { useForm } from '@toutaio/inertia-vue'

const form = useForm({
  name: '',
  email: ''
})

const submit = () => {
  form.post('/users')
}
</script>

<template>
  <form @submit.prevent="submit">
    <input v-model="form.name" />
    <div v-if="form.errors.name">{{ form.errors.name }}</div>
    
    <button :disabled="form.processing">Submit</button>
  </form>
</template>
```

### usePage

Access page data:

```vue
<script setup>
import { usePage } from '@toutaio/inertia-vue'

const page = usePage()
console.log(page.props.user)
</script>
```

## Documentation

See the [main repository](https://github.com/toutaio/toutago-inertia) for complete documentation.

## License

MIT
