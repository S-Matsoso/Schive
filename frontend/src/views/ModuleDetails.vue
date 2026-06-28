<template>
    <div class="pa-6">
        <h1 class="text-h4 mb-4">Module Dashboard</h1>

        <v-card class="pa-4" variant="tonal">
            <p>You are currently viewing the module with ID: <strong>{{ moduleId }}</strong></p>
        </v-card>
    </div>
</template>

<script setup>
    import { computed } from 'vue'
    import { useRoute } from 'vue-router'

    // useRoute gives us access to the current URL details
    const route = useRoute()

    // We use a computed property so that if the user clicks a different module 
    // in the sidebar, this ID updates automatically without refreshing the page
    const moduleId = computed(() => route.params.id)

</script>