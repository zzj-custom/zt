<template>
    <div class="d-flex justify-center align-center text-center h-100vh">
        <div>
            <img src="@/assets/images/background/errorimg.svg" width="500" alt="404" />
            <h1 class="text-h1 pt-3">Opps!!!</h1>
            <h4 class="text-h4 my-8">This page you are looking for could not be found.</h4>
            <v-btn flat color="primary" @click="goHome" class="mb-4">Go Back to Home</v-btn>
        </div>    
    </div>
</template>

<script setup>
 import {useRouter} from 'vue-router'
 const router = useRouter()

 const goHome = () => {
     router.push({
       name: 'Dashboard',
       path: '/',
     })
 }

</script>