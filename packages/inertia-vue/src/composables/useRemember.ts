import { ref, watch, Ref } from 'vue'

type StorageType = 'session' | 'local'

/**
 * Remember form state across page visits
 * Stores data in sessionStorage or localStorage
 * 
 * @param initialValue - Initial value
 * @param key - Unique storage key
 * @param storage - Storage type ('session' or 'local')
 * @returns Reactive reference to the remembered value
 * 
 * @example
 * ```ts
 * const formData = useRemember({ name: '', email: '' }, 'contact-form')
 * // Values persist across page navigation
 * ```
 */
export function useRemember<T>(
  initialValue: T,
  key: string,
  storage: StorageType = 'session'
): Ref<T> {
  const storageKey = `inertia:remember:${key}`
  const storageObject = storage === 'local' ? localStorage : sessionStorage

  // Try to restore from storage
  let restoredValue = initialValue
  try {
    const stored = storageObject.getItem(storageKey)
    if (stored !== null) {
      restoredValue = JSON.parse(stored)
    }
  } catch (error) {
    // Ignore JSON parse errors, use initial value
    // Only warn in non-test environments
    if (typeof process === 'undefined' || process.env?.NODE_ENV !== 'test') {
      console.warn(`Failed to restore remembered value for key "${key}":`, error)
    }
  }

  const state = ref(restoredValue) as Ref<T>

  // Watch for changes and update storage
  watch(
    state,
    (newValue) => {
      if (newValue === undefined) {
        // Clear storage when set to undefined
        storageObject.removeItem(storageKey)
      } else {
        try {
          storageObject.setItem(storageKey, JSON.stringify(newValue))
        } catch (error) {
          console.error(`Failed to save remembered value for key "${key}":`, error)
        }
      }
    },
    { deep: true, immediate: true }
  )

  return state
}
