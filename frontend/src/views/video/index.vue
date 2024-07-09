<script setup lang="ts">
import Search from '@/components/video/search.vue'
import List from '@/components/video/list.vue'

import{VideoList} from "@/types/video"
import {ref, watch} from "vue"
import Footer from "@/components/video/footer.vue";

const videoList = ref<VideoList[]>([
  {
    url:"https://www.bilibili.com/video/BV1VH4y1w7De",
    site:"哔哩哔哩 bilibili.com",
    title:"温庭筠：我拼尽了全力，还是过不好自己这一生",
    type:"video",
    streams:[
      {
        id:"32-12",
        quality:"清晰 480P hev1.1.6.L120.90",
        size:77059525,
        ext:"mp4",
      }
    ],
    caption:{
      "danmaku":{
        url:"https://comment.bilibili.com/1597996331.xml",
        size:0,
        ext:"xml"
      }
    }
  },
  {
    url:"https://www.bilibili.com/video/BV1VH4y1w7De",
    site:"哔哩哔哩 bilibili.com",
    title:"温庭筠：我拼尽了全力，还是过不好自己这一生",
    type:"video",
    streams:[
      {
        id:"32-12",
        quality:"清晰 480P hev1.1.6.L120.90",
        size:77059525,
        ext:"mp4",
      }
    ],
    caption:{
      "danmaku":{
        url:"https://comment.bilibili.com/1597996331.xml",
        size:0,
        ext:"xml"
      }
    }
  },
  {
    url:"https://www.bilibili.com/video/BV1VH4y1w7De",
    site:"哔哩哔哩 bilibili.com",
    title:"温庭筠：我拼尽了全力，还是过不好自己这一生",
    type:"video",
    streams:[
      {
        id:"32-12",
        quality:"清晰 480P hev1.1.6.L120.90",
        size:77059525,
        ext:"mp4",
      }
    ],
    caption:{
      "danmaku":{
        url:"https://comment.bilibili.com/1597996331.xml",
        size:0,
        ext:"xml"
      }
    }
  },
  {
    url:"https://www.bilibili.com/video/BV1VH4y1w7De",
    site:"哔哩哔哩 bilibili.com",
    title:"温庭筠：我拼尽了全力，还是过不好自己这一生",
    type:"video",
    streams:[
      {
        id:"32-12",
        quality:"清晰 480P hev1.1.6.L120.90",
        size:77059525,
        ext:"mp4",
      }
    ],
    caption:{
      "danmaku":{
        url:"https://comment.bilibili.com/1597996331.xml",
        size:0,
        ext:"xml"
      }
    }
  },
  {
    url:"https://www.bilibili.com/video/BV1VH4y1w7De",
    site:"哔哩哔哩 bilibili.com",
    title:"温庭筠：我拼尽了全力，还是过不好自己这一生",
    type:"video",
    streams:[
      {
        id:"32-12",
        quality:"清晰 480P hev1.1.6.L120.90",
        size:77059525,
        ext:"mp4",
      }
    ],
    caption:{
      "danmaku":{
        url:"https://comment.bilibili.com/1597996331.xml",
        size:0,
        ext:"xml"
      }
    }
  },
  {
    url:"https://www.bilibili.com/video/BV1VH4y1w7De",
    site:"哔哩哔哩 bilibili.com",
    title:"温庭筠：我拼尽了全力，还是过不好自己这一生",
    type:"video",
    streams:[
      {
        id:"32-12",
        quality:"清晰 480P hev1.1.6.L120.90",
        size:77059525,
        ext:"mp4",
      }
    ],
    caption:{
      "danmaku":{
        url:"https://comment.bilibili.com/1597996331.xml",
        size:0,
        ext:"xml"
      }
    }
  },
  {
    url:"https://www.bilibili.com/video/BV1VH4y1w7De",
    site:"哔哩哔哩 bilibili.com",
    title:"温庭筠：我拼尽了全力，还是过不好自己这一生",
    type:"video",
    streams:[
      {
        id:"32-12",
        quality:"清晰 480P hev1.1.6.L120.90",
        size:77059525,
        ext:"mp4",
      }
    ],
    caption:{
      "danmaku":{
        url:"https://comment.bilibili.com/1597996331.xml",
        size:0,
        ext:"xml"
      }
    }
  },
  {
    url:"https://www.bilibili.com/video/BV1VH4y1w7De",
    site:"哔哩哔哩 bilibili.com",
    title:"温庭筠：我拼尽了全力，还是过不好自己这一生",
    type:"video",
    streams:[
      {
        id:"32-12",
        quality:"清晰 480P hev1.1.6.L120.90",
        size:77059525,
        ext:"mp4",
      }
    ],
    caption:{
      "danmaku":{
        url:"https://comment.bilibili.com/1597996331.xml",
        size:0,
        ext:"xml"
      }
    }
  },
  {
    url:"https://www.bilibili.com/video/BV1VH4y1w7De",
    site:"哔哩哔哩 bilibili.com",
    title:"温庭筠：我拼尽了全力，还是过不好自己这一生",
    type:"video",
    streams:[
      {
        id:"32-12",
        quality:"清晰 480P hev1.1.6.L120.90",
        size:77059525,
        ext:"mp4",
      }
    ],
    caption:{
      "danmaku":{
        url:"https://comment.bilibili.com/1597996331.xml",
        size:0,
        ext:"xml"
      }
    }
  },
  {
    url:"https://www.bilibili.com/video/BV1VH4y1w7De",
    site:"哔哩哔哩 bilibili.com",
    title:"温庭筠：我拼尽了全力，还是过不好自己这一生",
    type:"video",
    streams:[
      {
        id:"32-12",
        quality:"清晰 480P hev1.1.6.L120.90",
        size:77059525,
        ext:"mp4",
      }
    ],
    caption:{
      "danmaku":{
        url:"https://comment.bilibili.com/1597996331.xml",
        size:0,
        ext:"xml"
      }
    }
  }
])

watch(videoList, (videoLists:VideoList[]) => {
  console.log(videoLists)
  if (videoLists !== videoList.value) {
    videoList.value = [...videoLists];
  }
})

import { getCurrentInstance } from 'vue'
import {MessageState} from "@/plugins/messagePlugin";

const { proxy } = getCurrentInstance()!

const clickTest = () => {
  proxy!.$message.error('This is an error message!')
}

</script>

<template>
  <Search v-model="videoList" />
  <List class="mt-3" v-model="videoList" />
</template>