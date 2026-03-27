import { ref } from 'vue'

export function useAsyncAction() {
  const running = ref(false)

  async function run<T>(action: () => Promise<T>) {
    running.value = true
    try {
      return await action()
    } finally {
      running.value = false
    }
  }

  return {
    running,
    run
  }
}
