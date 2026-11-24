/**
 * Simple global event bus for synchronizing state across components and views
 * (e.g. sidebar, dashboard, tasks view, module details view).
 */

type Handler = (...args: any[]) => void

class EventBus {
  private events: Map<string, Set<Handler>> = new Map()

  on(event: string, handler: Handler): () => void {
    if (!this.events.has(event)) {
      this.events.set(event, new Set())
    }
    this.events.get(event)!.add(handler)

    return () => {
      this.off(event, handler)
    }
  }

  off(event: string, handler: Handler): void {
    const handlers = this.events.get(event)
    if (handlers) {
      handlers.delete(handler)
      if (handlers.size === 0) {
        this.events.delete(event)
      }
    }
  }

  emit(event: string, ...args: any[]): void {
    const handlers = this.events.get(event)
    if (handlers) {
      handlers.forEach(fn => {
        try {
          fn(...args)
        } catch (err) {
          console.error(`Error in event bus handler for "${event}":`, err)
        }
      })
    }
  }
}

export const eventBus = new EventBus()

export const SYNC_EVENTS = {
  TASKS_CHANGED: 'tasks:changed',
  ASSIGNMENTS_CHANGED: 'assignments:changed',
  MODULES_CHANGED: 'modules:changed'
} as const
