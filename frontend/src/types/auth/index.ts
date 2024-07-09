type UserInfo = {
    Id:string,
    name:string ,
    email:string,
    avatar:string,
    mobile:string,
}

type LoginForm = {
    email:string,
    captcha:number,
}

export type {UserInfo,LoginForm}