<script setup lang="ts">
import { BrandBingIcon } from 'vue-tabler-icons';
import { Images } from '../../../wailsjs/go/app/App';
import {response} from '../../../wailsjs/go/models';
import {reactive, onMounted} from 'vue';

interface ImageList  {
  Id: number;
  Name: string;
  Copyright: string;
  CopyrightLink: string;
  Url: string;
  Start: string;
  End: string;
  ClickCount: number;
  DownloadCount: number;
}


const imagesList = reactive<ImageList[]>([]);

const list = async () => {
  await Images().then((res:response.Reply) => {
    res.result.forEach((item:ImageList) => {
      imagesList.push({
        Id: item.Id,
        Name: item.Name.length > 8 ? item.Name.substring(0, 8) + '...' : item.Name,
        Copyright: item.Copyright.length > 30 ? item.Copyright.substring(0, 30) + '...' : item.Copyright,
        CopyrightLink: item.CopyrightLink,
        Url: item.Url,
        Start: item.Start,
        End: item.End,
        ClickCount: item.ClickCount,
        DownloadCount: item.DownloadCount
      })
    })
  }).catch((error:any) => {
    console.log(error)
  });
}

onMounted(() => {
  list()
})

</script>
<template>
    <v-row>
        <v-col cols="12" lg="3" sm="6" v-for="images in imagesList" :key="images.Id">
            <v-card elevation="10" class="withbg" rounded="md" hover color="#26c6da" href="/">
                <v-img :cover="true" aspect-ratio="1" :src="images.Url" height="100%" class="rounded-t-md"></v-img>
                <div class="d-flex justify-end mr-4 mt-n5">
                    <v-btn size="40" icon class="bg-primary d-flex">
                        <v-avatar size="30" class="text-white">
                            <BrandBingIcon size="15" />
                        </v-avatar>
                    </v-btn>
                </div>
                <v-card-item class="pt-0 text-white">
                    <h6 class="text-h6" v-text="images.Copyright"></h6>
                    <div class="d-flex align-center justify-space-between mt-3">
                        <div>
                            <span class="text-h8" v-text="images.Name"></span>
                        </div>
                        <div class="justify-self-end">
                            <v-icon class="me-1" icon="mdi-heart"></v-icon>
                            <span class="subheading me-2">{{ images.ClickCount }}</span>
                            <v-icon class="me-1" icon="mdi-share-variant"></v-icon>
                            <span class="subheading">{{ images.DownloadCount }}</span>
          </div>
                    </div>
                </v-card-item>
            </v-card>
        </v-col>
    </v-row>
</template>
