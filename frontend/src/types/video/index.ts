type VideoList = {
    url:string,
    site:string,
    title: string,
    type:string,
    streams:  stream[],
    caption:{
        [key:string]:caption
    },
}

type stream = {
    id:string,
    quality: string,
    size: number,
    ext: string
}

type caption = {
    url: string,
    size: number,
    ext: string
}

export type {VideoList}