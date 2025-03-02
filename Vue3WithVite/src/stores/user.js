import {defineStore} from "pinia";

const useUserStore = defineStore('user', {
    state: () => ({
        user: null,
    }),

    actions: {
        login(user){
            this.user = user
        }
    },

    getters: {
        isLogin: (state) => state.user !== null
    }
})
export { useUserStore }